# User CRUD API with Go and MongoDB

A RESTful API for user management built with Go and MongoDB. This application provides full CRUD (Create, Read, Update, Delete) operations for user entities.

## Features

- **Create User**: Add new users to the database
- **Read Users**: Get all users or a specific user by ID
- **Update User**: Modify existing user information
- **Delete User**: Remove users from the database
- **Health Check**: Monitor API health status
- **MongoDB Integration**: Persistent data storage
- **RESTful API**: Clean HTTP endpoints following REST conventions

## Project Structure

```
user-crud/
├── main.go                 # Application entry point
├── go.mod                  # Go module dependencies
├── models/
│   └── user.go            # User data models
├── database/
│   └── connection.go      # MongoDB connection setup
├── repository/
│   └── user_repository.go # Data access layer
├── handlers/
│   └── user_handler.go    # HTTP request handlers
└── routes/
    └── routes.go          # API route configuration
```

## Prerequisites

- Go 1.21 or higher
- MongoDB (local installation or MongoDB Atlas)
- Git

## Installation & Setup

### 1. Clone the Repository

```bash
git clone <repository-url>
cd user-crud
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Set Up MongoDB

#### Option A: Local MongoDB
1. Install MongoDB locally
2. Start MongoDB service:
   ```bash
   # On Windows
   net start MongoDB
   
   # On macOS with Homebrew
   brew services start mongodb-community
   
   # On Linux
   sudo systemctl start mongod
   ```

#### Option B: MongoDB Atlas (Cloud)
1. Create a free account at [MongoDB Atlas](https://www.mongodb.com/atlas)
2. Create a new cluster
3. Get your connection string

### 4. Environment Variables (Optional)

Create a `.env` file or set environment variables:

```bash
# MongoDB connection
export MONGO_URI="mongodb://localhost:27017"  # or your Atlas connection string
export DB_NAME="userdb"
export PORT="8080"
```

### 5. Run the Application

```bash
go run main.go
```

The server will start on `http://localhost:8080` (or your specified PORT).

## API Endpoints

### Health Check
- **GET** `/health` - Check API health status

### User Management
- **POST** `/api/v1/users` - Create a new user
- **GET** `/api/v1/users` - Get all users
- **GET** `/api/v1/users/{id}` - Get user by ID
- **PUT** `/api/v1/users/{id}` - Update user by ID
- **DELETE** `/api/v1/users/{id}` - Delete user by ID

## API Usage Examples

### 1. Create a User

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "age": 30,
    "phone": "+1234567890",
    "address": "123 Main St, City, State"
  }'
```

### 2. Get All Users

```bash
curl -X GET http://localhost:8080/api/v1/users
```

### 3. Get User by ID

```bash
curl -X GET http://localhost:8080/api/v1/users/{user_id}
```

### 4. Update User

```bash
curl -X PUT http://localhost:8080/api/v1/users/{user_id} \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Smith",
    "age": 31
  }'
```

### 5. Delete User

```bash
curl -X DELETE http://localhost:8080/api/v1/users/{user_id}
```

### 6. Health Check

```bash
curl -X GET http://localhost:8080/health
```

## Data Models

### User Model
```json
{
  "id": "ObjectId",
  "name": "string",
  "email": "string",
  "age": "number",
  "phone": "string",
  "address": "string",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

## Error Handling

The API returns appropriate HTTP status codes:

- `200 OK` - Successful GET, PUT requests
- `201 Created` - Successful POST requests
- `204 No Content` - Successful DELETE requests
- `400 Bad Request` - Invalid request data
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server errors

## Development

### Running in Development Mode

```bash
# Install air for hot reloading (optional)
go install github.com/cosmtrek/air@latest

# Run with hot reloading
air
```

### Building for Production

```bash
# Build the application
go build -o user-crud main.go

# Run the binary
./user-crud
```

## Testing

You can test the API using tools like:
- **curl** (command line)
- **Postman** (GUI)
- **Insomnia** (GUI)
- **HTTPie** (command line)

## Docker Support (Optional)

Create a `Dockerfile`:

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

Build and run with Docker:

```bash
docker build -t user-crud .
docker run -p 8080:8080 user-crud
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

This project is open source and available under the [MIT License](LICENSE).

## Support

For issues and questions, please create an issue in the repository or contact the development team.