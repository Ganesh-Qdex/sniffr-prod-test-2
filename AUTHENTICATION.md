# Authentication System Documentation

This document describes the authentication system implemented in the User CRUD API.

## Overview

The authentication system uses JWT (JSON Web Tokens) for stateless authentication. Users can register, login, and access protected routes using JWT tokens.

## Features

- **User Registration**: Create new user accounts with password hashing
- **User Login**: Authenticate users and receive JWT tokens
- **Password Security**: Passwords are hashed using bcrypt
- **JWT Tokens**: Stateless authentication with configurable expiration
- **Protected Routes**: All user CRUD operations require authentication
- **Profile Management**: Users can view and update their own profiles
- **Password Change**: Secure password change functionality

## API Endpoints

### Public Endpoints

#### Register User
```
POST /api/v1/auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123",
  "age": 30,
  "phone": "+1234567890",
  "address": "123 Main St, City, State"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "64f8b1234567890abcdef123",
    "name": "John Doe",
    "email": "john@example.com",
    "age": 30,
    "phone": "+1234567890",
    "address": "123 Main St, City, State",
    "created_at": "2023-09-05T10:30:00Z",
    "updated_at": "2023-09-05T10:30:00Z"
  }
}
```

#### Login User
```
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "64f8b1234567890abcdef123",
    "name": "John Doe",
    "email": "john@example.com",
    "age": 30,
    "phone": "+1234567890",
    "address": "123 Main St, City, State",
    "created_at": "2023-09-05T10:30:00Z",
    "updated_at": "2023-09-05T10:30:00Z"
  }
}
```

### Protected Endpoints

All protected endpoints require the `Authorization: Bearer <token>` header.

#### Get User Profile
```
GET /api/v1/auth/profile
Authorization: Bearer <token>
```

#### Update User Profile
```
PUT /api/v1/auth/profile
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "John Smith",
  "age": 31,
  "address": "456 Oak Ave, New City, State"
}
```

#### Change Password
```
PUT /api/v1/auth/change-password
Authorization: Bearer <token>
Content-Type: application/json

{
  "current_password": "password123",
  "new_password": "newpassword456"
}
```

#### Protected User CRUD Operations
All existing user CRUD operations now require authentication:

- `GET /api/v1/users` - Get all users
- `POST /api/v1/users` - Create user
- `GET /api/v1/users/{id}` - Get user by ID
- `PUT /api/v1/users/{id}` - Update user
- `DELETE /api/v1/users/{id}` - Delete user

## Configuration

### Environment Variables

Create a `config.env` file with the following variables:

```env
# Database Configuration
MONGO_URI=mongodb://localhost:27017
DB_NAME=userdb

# Server Configuration
PORT=8080

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# Development/Production Environment
ENV=development
```

### JWT Secret

**Important**: Change the `JWT_SECRET` in production to a secure, random string. You can generate one using:

```bash
openssl rand -base64 32
```

## Security Features

1. **Password Hashing**: All passwords are hashed using bcrypt with default cost
2. **JWT Expiration**: Tokens expire after 24 hours by default
3. **Password Validation**: Minimum 6 characters required
4. **Email Validation**: Proper email format validation
5. **Token Verification**: All protected routes verify JWT tokens

## Testing

### Using the Test Scripts

#### Linux/macOS
```bash
chmod +x examples/test_api.sh
./examples/test_api.sh
```

#### Windows
```cmd
examples\test_api.bat
```

### Manual Testing

1. **Register a user**:
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123",
    "age": 30,
    "phone": "+1234567890",
    "address": "123 Main St, City, State"
  }'
```

2. **Login**:
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

3. **Access protected route** (replace `YOUR_TOKEN` with the token from login):
```bash
curl -X GET http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Error Handling

The API returns appropriate HTTP status codes:

- `200 OK` - Successful request
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Invalid or missing authentication
- `404 Not Found` - Resource not found
- `409 Conflict` - User already exists (registration)
- `500 Internal Server Error` - Server error

## Database Schema

The User model now includes a password field:

```go
type User struct {
    ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Name      string             `json:"name" bson:"name"`
    Email     string             `json:"email" bson:"email"`
    Password  string             `json:"-" bson:"password"` // Hidden from JSON responses
    Age       int                `json:"age" bson:"age"`
    Phone     string             `json:"phone" bson:"phone"`
    Address   string             `json:"address" bson:"address"`
    CreatedAt time.Time          `json:"created_at" bson:"created_at"`
    UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
```

## Dependencies Added

- `github.com/golang-jwt/jwt/v5` - JWT token handling
- `golang.org/x/crypto` - Password hashing (bcrypt)

## Running the Application

1. **Install dependencies**:
```bash
go mod tidy
```

2. **Set environment variables** (copy from `config.env`):
```bash
export JWT_SECRET="your-super-secret-jwt-key"
export MONGO_URI="mongodb://localhost:27017"
export DB_NAME="userdb"
export PORT="8080"
```

3. **Run the application**:
```bash
go run main.go
```

The server will start on `http://localhost:8080` with authentication enabled.
