# Copilot Instructions for Go User Management API

## Project Overview
This is a Go-based REST API for user management using Clean Architecture principles. The project implements a user registration system with membership levels and points tracking.

## Architecture
The project follows Clean Architecture with the following layers:

### 1. **Entity Layer** (`/entity`)
- Contains domain models and business entities
- **Key Files:**
  - `user.go` - User domain model with business logic methods

### 2. **Repository Layer** (`/repository`)
- Defines data access interfaces and implementations
- **Key Files:**
  - `user_repository.go` - Interface definition for user data operations
  - `memory_user_repository.go` - In-memory implementation of user repository

### 3. **Use Case Layer** (`/usecase`)
- Contains business logic and application services
- **Key Files:**
  - `user_usecase.go` - User business logic, request/response DTOs, and validation

### 4. **Handler Layer** (`/handler`)
- Handles HTTP requests and responses
- **Key Files:**
  - `http_handler.go` - HTTP handlers using Fiber framework

### 5. **Main Application** (`/`)
- Application entry point and dependency injection
- **Key Files:**
  - `main.go` - Application bootstrap and server configuration

### 6. **Scripts** (`/scripts`)
- Helper scripts for project management and development
- **Key Files:**
  - `kill-port.sh` - Utility script to kill processes running on specific ports

## Technology Stack
- **Language:** Go 1.23.1
- **Web Framework:** Fiber v2.52.0
- **UUID Generation:** google/uuid v1.5.0
- **Architecture:** Clean Architecture with Dependency Injection

## Coding Standards

### 1. **Package Organization**
- Use clear, descriptive package names
- Keep packages focused on single responsibility
- Follow Go naming conventions (lowercase, no underscores)

### 2. **Error Handling**
- Always handle errors explicitly
- Return meaningful error messages
- Use custom error types when appropriate
- Log errors at appropriate levels

### 3. **Interface Design**
- Define interfaces in the layer that uses them
- Keep interfaces small and focused
- Use descriptive method names

### 4. **Struct Design**
- Use JSON tags for API serialization
- Include validation tags where appropriate
- Provide constructor functions for complex structs

### 5. **HTTP Handler Patterns**
- Use consistent response structures
- Implement proper HTTP status codes
- Handle request parsing errors gracefully
- Validate input data before processing

## API Endpoints

### Health Check
- `GET /health` - Service health status

### User Management
- `POST /register` - Register new user
- `GET /user/:id` - Get user by ID
- `GET /users` - Get all users

## Request/Response Patterns

### Registration Request
```json
{
  "first_name": "string (required)",
  "last_name": "string (required)",
  "phone": "string (required)",
  "email": "string (required, valid email)"
}
```

### Standard Response Format
```json
{
  "success": "boolean",
  "message": "string",
  "user": "User object (optional)"
}
```

## Development Guidelines

### 1. **Adding New Features**
- Start with entity changes if needed
- Update repository interface and implementation
- Add use case logic with proper validation
- Implement HTTP handler
- Update main.go for dependency injection

### 2. **Testing**
- Write unit tests for each layer
- Use table-driven tests for multiple scenarios
- Mock dependencies in tests
- Test error conditions

### 3. **Validation**
- Validate input at the use case layer
- Use meaningful validation error messages
- Check business rules (e.g., email uniqueness)

### 4. **Error Responses**
- Use appropriate HTTP status codes
- Provide consistent error message format
- Log errors for debugging
- Don't expose internal implementation details

## Common Patterns

### 1. **Dependency Injection**
```go
// Initialize dependencies in main.go
userRepo := repository.NewMemoryUserRepository()
userUsecase := usecase.NewUserUsecase(userRepo)
httpHandler := handler.NewHTTPHandler(userUsecase)
```

### 2. **Repository Pattern**
```go
// Define interface in repository package
type UserRepository interface {
    Create(user *entity.User) error
    GetByID(id string) (*entity.User, error)
    // ... other methods
}
```

### 3. **Use Case Pattern**
```go
// Implement business logic in use case
func (u *userUsecase) Register(req RegisterRequest) (*RegisterResponse, error) {
    // Validation
    // Business logic
    // Repository calls
    // Response formatting
}
```

### 4. **HTTP Handler Pattern**
```go
// Handle HTTP requests
func (h *HTTPHandler) Register(c *fiber.Ctx) error {
    // Parse request
    // Call use case
    // Return appropriate response
}
```

## Best Practices

1. **Keep functions small and focused**
2. **Use meaningful variable and function names**
3. **Add comments for complex business logic**
4. **Follow Go idioms and conventions**
5. **Use context for request-scoped data**
6. **Implement proper logging**
7. **Handle panics gracefully**
8. **Use interfaces for testability**

## File Naming Conventions
- Use snake_case for file names
- Include package purpose in filename (e.g., `user_repository.go`)
- Keep file names descriptive and concise

## Import Organization
1. Standard library imports
2. Third-party library imports
3. Local package imports

## Memory Management
- Use pointers for large structs
- Be mindful of memory leaks in long-running applications
- Use appropriate data structures for performance

## Security Considerations
- Validate all input data
- Sanitize user inputs
- Use HTTPS in production
- Implement proper authentication/authorization when needed
- Don't expose sensitive information in error messages

## Performance Tips
- Use appropriate data structures
- Avoid unnecessary allocations
- Use connection pooling for databases
- Implement caching where appropriate
- Profile and benchmark critical paths

## Common Gotchas
- Remember to handle nil pointers
- Check for empty strings and zero values
- Validate email format properly
- Handle concurrent access to shared data
- Use proper error wrapping for context

## Future Enhancements
- Add database persistence layer
- Implement authentication and authorization
- Add input validation middleware
- Implement rate limiting
- Add comprehensive logging
- Add metrics and monitoring
- Implement graceful shutdown
- Add configuration management
