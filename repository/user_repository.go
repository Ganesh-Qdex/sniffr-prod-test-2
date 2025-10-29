package repository

import (
	"context"
	"fmt"
	"time"

	"user-crud/cache"
	"user-crud/database"
	"user-crud/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	collection *mongo.Collection
}

// NewUserRepository creates a new user repository
func NewUserRepository() *UserRepository {
	return &UserRepository{
		collection: database.Database.Collection("users"),
	}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	return user, nil
}

// GetByID retrieves a user by ID with caching
func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	// Try to get from cache first
	cacheKey := fmt.Sprintf("user:%s", id)
	var user models.User

	if err := cache.Get(ctx, cacheKey, &user); err == nil {
		return &user, nil
	}

	// If not in cache, get from database
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %v", err)
	}

	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	// Cache the result for 5 minutes
	cache.Set(ctx, cacheKey, &user, 5*time.Minute)

	return &user, nil
}

// GetAll retrieves all users with pagination and caching
func (r *UserRepository) GetAll(ctx context.Context, page, limit int) ([]*models.User, int64, error) {
	// Create cache key based on pagination
	cacheKey := fmt.Sprintf("users:page:%d:limit:%d", page, limit)

	// Try to get from cache first
	var cachedResult struct {
		Users []*models.User `json:"users"`
		Total int64          `json:"total"`
	}

	if err := cache.Get(ctx, cacheKey, &cachedResult); err == nil {
		return cachedResult.Users, cachedResult.Total, nil
	}

	// Calculate skip value for pagination
	skip := (page - 1) * limit

	// Get total count
	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %v", err)
	}

	// Set up pagination options
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}}) // Sort by creation date descending

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get users: %v", err)
	}
	defer cursor.Close(ctx)

	var users []*models.User
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, 0, fmt.Errorf("failed to decode user: %v", err)
		}
		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, fmt.Errorf("cursor error: %v", err)
	}

	// Cache the result for 2 minutes
	cachedResult = struct {
		Users []*models.User `json:"users"`
		Total int64          `json:"total"`
	}{
		Users: users,
		Total: total,
	}
	cache.Set(ctx, cacheKey, &cachedResult, 2*time.Minute)

	return users, total, nil
}

// Update updates a user by ID
func (r *UserRepository) Update(ctx context.Context, id string, updateData *models.UpdateUserRequest) (*models.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %v", err)
	}

	// Build update document with only non-zero values
	update := bson.M{
		"updated_at": time.Now(),
	}

	if updateData.Name != "" {
		update["name"] = updateData.Name
	}
	if updateData.Email != "" {
		update["email"] = updateData.Email
	}
	if updateData.Age > 0 {
		update["age"] = updateData.Age
	}
	if updateData.Phone != "" {
		update["phone"] = updateData.Phone
	}
	if updateData.Address != "" {
		update["address"] = updateData.Address
	}

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": update},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("user not found")
	}

	// Invalidate cache for this user
	cacheKey := fmt.Sprintf("user:%s", id)
	cache.Delete(ctx, cacheKey)

	// Invalidate users list cache
	cache.DeletePattern(ctx, "users:*")

	// Return the updated user
	return r.GetByID(ctx, id)
}

// Delete deletes a user by ID
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid user ID: %v", err)
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("user not found")
	}

	// Invalidate cache for this user
	cacheKey := fmt.Sprintf("user:%s", id)
	cache.Delete(ctx, cacheKey)

	// Invalidate users list cache
	cache.DeletePattern(ctx, "users:*")

	return nil
}
