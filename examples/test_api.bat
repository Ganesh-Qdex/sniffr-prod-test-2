@echo off
REM User CRUD API Test Script with Authentication for Windows
REM Make sure the server is running on localhost:8080

set BASE_URL=http://localhost:8080/api/v1
set HEALTH_URL=http://localhost:8080/health

echo === User CRUD API Test Script with Authentication ===
echo.

REM Test Health Check
echo 1. Testing Health Check...
curl -s -X GET %HEALTH_URL%
echo.
echo.

REM Test User Registration
echo 2. Registering a new user...
curl -s -X POST %BASE_URL%/auth/register -H "Content-Type: application/json" -d "{\"name\": \"John Doe\", \"email\": \"john@example.com\", \"password\": \"password123\", \"age\": 30, \"phone\": \"+1234567890\", \"address\": \"123 Main St, City, State\"}"
echo.
echo.

REM Test User Login
echo 3. Testing user login...
curl -s -X POST %BASE_URL%/auth/login -H "Content-Type: application/json" -d "{\"email\": \"john@example.com\", \"password\": \"password123\"}"
echo.
echo.

REM Test Get Profile (replace {token} with actual token from login response)
echo 4. Getting user profile (replace {token} with actual token)...
echo curl -s -X GET %BASE_URL%/auth/profile -H "Authorization: Bearer {token}"
echo.

REM Test Get All Users (Protected)
echo 5. Getting all users (protected route, replace {token} with actual token)...
echo curl -s -X GET %BASE_URL%/users -H "Authorization: Bearer {token}"
echo.

REM Test Update Profile
echo 6. Updating user profile (replace {token} with actual token)...
echo curl -s -X PUT %BASE_URL%/auth/profile -H "Authorization: Bearer {token}" -H "Content-Type: application/json" -d "{\"name\": \"John Smith\", \"age\": 31, \"address\": \"456 Oak Ave, New City, State\"}"
echo.

REM Test Change Password
echo 7. Changing password (replace {token} with actual token)...
echo curl -s -X PUT %BASE_URL%/auth/change-password -H "Authorization: Bearer {token}" -H "Content-Type: application/json" -d "{\"current_password\": \"password123\", \"new_password\": \"newpassword456\"}"
echo.

REM Test Create Another User (Protected)
echo 8. Creating another user (protected route, replace {token} with actual token)...
echo curl -s -X POST %BASE_URL%/users -H "Authorization: Bearer {token}" -H "Content-Type: application/json" -d "{\"name\": \"Jane Doe\", \"email\": \"jane@example.com\", \"password\": \"password123\", \"age\": 25, \"phone\": \"+0987654321\", \"address\": \"789 Pine St, Another City, State\"}"
echo.

REM Test Unauthorized Access (should fail)
echo 9. Testing unauthorized access (should fail)...
curl -s -X GET %BASE_URL%/users
echo.
echo.

echo === Test completed ===
echo.
echo Note: For steps 4-8, you need to replace {token} with the actual JWT token
echo returned from the login request in step 3.

