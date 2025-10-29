@echo off
REM User CRUD API Test Script for Windows
REM Make sure the server is running on localhost:8080

set BASE_URL=http://localhost:8080/api/v1
set HEALTH_URL=http://localhost:8080/health

echo === User CRUD API Test Script ===
echo.

REM Test Health Check
echo 1. Testing Health Check...
curl -s -X GET %HEALTH_URL%
echo.
echo.

REM Test Create User
echo 2. Creating a new user...
curl -s -X POST %BASE_URL%/users -H "Content-Type: application/json" -d "{\"name\": \"John Doe\", \"email\": \"john@example.com\", \"age\": 30, \"phone\": \"+1234567890\", \"address\": \"123 Main St, City, State\"}"
echo.
echo.

REM Test Get All Users
echo 3. Getting all users...
curl -s -X GET %BASE_URL%/users
echo.
echo.

REM Test Get User by ID (you'll need to replace {user_id} with actual ID from previous response)
echo 4. Getting user by ID (replace {user_id} with actual ID)...
echo curl -s -X GET %BASE_URL%/users/{user_id}
echo.

REM Test Update User
echo 5. Updating user (replace {user_id} with actual ID)...
echo curl -s -X PUT %BASE_URL%/users/{user_id} -H "Content-Type: application/json" -d "{\"name\": \"John Smith\", \"age\": 31, \"address\": \"456 Oak Ave, New City, State\"}"
echo.

REM Test Delete User
echo 6. Deleting user (replace {user_id} with actual ID)...
echo curl -s -X DELETE %BASE_URL%/users/{user_id}
echo.

echo === Test completed ===
echo.
echo Note: For steps 4-6, you need to replace {user_id} with the actual user ID
echo returned from the create user request in step 2.
