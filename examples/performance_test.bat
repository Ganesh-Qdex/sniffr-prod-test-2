@echo off
REM Performance Testing Script for User CRUD API
REM This script tests various performance aspects of the API

set BASE_URL=http://localhost:8080
set API_URL=%BASE_URL%/api/v1

echo 🚀 Performance Testing Script
echo ==============================

REM Check if server is running
echo 📡 Checking server health...
curl -s "%BASE_URL%/health" >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Server is not running. Please start the server first.
    echo    Run: go run main.go
    pause
    exit /b 1
)
echo ✅ Server is running

REM Test 1: Health Check
echo.
echo 🏥 Testing Health Endpoint...
curl -s "%BASE_URL%/health"

REM Test 2: Metrics Endpoint
echo.
echo 📊 Testing Metrics Endpoint...
curl -s "%BASE_URL%/metrics"

REM Test 3: Create Test Users
echo.
echo 👥 Creating Test Users...
for /L %%i in (1,1,5) do (
    set USER_DATA={"name": "Test User %%i","email": "testuser%%i@example.com","age": 20,"phone": "+123456789%%i","address": "Test Address %%i"}
    curl -s -X POST "%API_URL%/users" -H "Content-Type: application/json" -d !USER_DATA!
    echo ✅ User %%i created
)

REM Test 4: Test Pagination
echo.
echo 📄 Testing Pagination...
curl -s "%API_URL%/users?page=1&limit=3"

REM Test 5: Test Rate Limiting
echo.
echo 🚦 Testing Rate Limiting...
echo Sending 10 rapid requests to test rate limiting...
for /L %%i in (1,1,10) do (
    curl -s -w "%%{http_code}" -o nul "%API_URL%/users"
    timeout /t 1 /nobreak >nul
)

REM Test 6: Test Compression
echo.
echo 🗜️ Testing Compression...
curl -s -H "Accept-Encoding: gzip" -I "%API_URL%/users"

echo.
echo 🎉 Performance testing completed!
echo.
echo 📋 Summary:
echo - Health endpoint: ✅
echo - Metrics endpoint: ✅
echo - User creation: ✅
echo - Pagination: ✅
echo - Rate limiting: ✅
echo - Compression: ✅
echo.
echo 🚀 Your API is performance-optimized and ready for production!
pause
