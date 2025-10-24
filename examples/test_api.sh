#!/bin/bash

# User CRUD API Test Script with Authentication
# Make sure the server is running on localhost:8080

BASE_URL="http://localhost:8080/api/v1"
HEALTH_URL="http://localhost:8080/health"

echo "=== User CRUD API Test Script with Authentication ==="
echo

# Test Health Check
echo "1. Testing Health Check..."
curl -s -X GET $HEALTH_URL | jq .
echo -e "\n"

# Test User Registration
echo "2. Registering a new user..."
REGISTER_RESPONSE=$(curl -s -X POST $BASE_URL/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123",
    "age": 30,
    "phone": "+1234567890",
    "address": "123 Main St, City, State"
  }')

echo $REGISTER_RESPONSE | jq .
TOKEN=$(echo $REGISTER_RESPONSE | jq -r '.token')
USER_ID=$(echo $REGISTER_RESPONSE | jq -r '.user.id')
echo -e "\n"

# Test User Login
echo "3. Testing user login..."
LOGIN_RESPONSE=$(curl -s -X POST $BASE_URL/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }')

echo $LOGIN_RESPONSE | jq .
TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token')
echo -e "\n"

# Test Get Profile
echo "4. Getting user profile..."
curl -s -X GET $BASE_URL/auth/profile \
  -H "Authorization: Bearer $TOKEN" | jq .
echo -e "\n"

# Test Get All Users (Protected)
echo "5. Getting all users (protected route)..."
curl -s -X GET $BASE_URL/users \
  -H "Authorization: Bearer $TOKEN" | jq .
echo -e "\n"

# Test Get User by ID (Protected)
echo "6. Getting user by ID: $USER_ID (protected route)..."
curl -s -X GET $BASE_URL/users/$USER_ID \
  -H "Authorization: Bearer $TOKEN" | jq .
echo -e "\n"

# Test Update Profile
echo "7. Updating user profile..."
curl -s -X PUT $BASE_URL/auth/profile \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Smith",
    "age": 31,
    "address": "456 Oak Ave, New City, State"
  }' | jq .
echo -e "\n"

# Test Get Updated Profile
echo "8. Getting updated profile..."
curl -s -X GET $BASE_URL/auth/profile \
  -H "Authorization: Bearer $TOKEN" | jq .
echo -e "\n"

# Test Change Password
echo "9. Changing password..."
curl -s -X PUT $BASE_URL/auth/change-password \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "current_password": "password123",
    "new_password": "newpassword456"
  }' | jq .
echo -e "\n"

# Test Login with New Password
echo "10. Testing login with new password..."
curl -s -X POST $BASE_URL/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "newpassword456"
  }' | jq .
echo -e "\n"

# Test Create Another User (Protected)
echo "11. Creating another user (protected route)..."
curl -s -X POST $BASE_URL/users \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Doe",
    "email": "jane@example.com",
    "password": "password123",
    "age": 25,
    "phone": "+0987654321",
    "address": "789 Pine St, Another City, State"
  }' | jq .
echo -e "\n"

# Test Get All Users (should show both users)
echo "12. Getting all users after creating second user..."
curl -s -X GET $BASE_URL/users \
  -H "Authorization: Bearer $TOKEN" | jq .
echo -e "\n"

# Test Unauthorized Access (should fail)
echo "13. Testing unauthorized access (should fail)..."
curl -s -X GET $BASE_URL/users | jq .
echo -e "\n"

echo "=== Test completed ==="
