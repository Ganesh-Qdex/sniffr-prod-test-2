#!/bin/bash

# User CRUD API Test Script
# Make sure the server is running on localhost:8080

BASE_URL="http://localhost:8080/api/v1"
HEALTH_URL="http://localhost:8080/health"

echo "=== User CRUD API Test Script ==="
echo

# Test Health Check
echo "1. Testing Health Check..."
curl -s -X GET $HEALTH_URL | jq .
echo -e "\n"

# Test Create User
echo "2. Creating a new user..."
USER_RESPONSE=$(curl -s -X POST $BASE_URL/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "age": 30,
    "phone": "+1234567890",
    "address": "123 Main St, City, State"
  }')

echo $USER_RESPONSE | jq .
USER_ID=$(echo $USER_RESPONSE | jq -r '.id')
echo -e "\n"

# Test Get All Users
echo "3. Getting all users..."
curl -s -X GET $BASE_URL/users | jq .
echo -e "\n"

# Test Get User by ID
echo "4. Getting user by ID: $USER_ID"
curl -s -X GET $BASE_URL/users/$USER_ID | jq .
echo -e "\n"

# Test Update User
echo "5. Updating user..."
curl -s -X PUT $BASE_URL/users/$USER_ID \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Smith",
    "age": 31,
    "address": "456 Oak Ave, New City, State"
  }' | jq .
echo -e "\n"

# Test Get Updated User
echo "6. Getting updated user..."
curl -s -X GET $BASE_URL/users/$USER_ID | jq .
echo -e "\n"

# Test Delete User
echo "7. Deleting user..."
curl -s -X DELETE $BASE_URL/users/$USER_ID
echo -e "\n"

# Test Get All Users (should be empty or not contain deleted user)
echo "8. Getting all users after deletion..."
curl -s -X GET $BASE_URL/users | jq .
echo -e "\n"

echo "=== Test completed ==="
