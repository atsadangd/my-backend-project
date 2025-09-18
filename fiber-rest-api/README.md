# Fiber REST API

This project is a simple RESTful API built using the Fiber framework in Go. It responds to a GET request at the root endpoint (/) with a JSON object containing a greeting message.

## Project Structure

```
fiber-rest-api
├── cmd
│   └── server
│       └── main.go        # Entry point of the application
├── internal
│   ├── handlers
│   │   └── root.go        # Handler for the root endpoint
│   └── router
│       └── router.go      # Router setup for the application
├── go.mod                  # Module definition and dependencies
├── .gitignore              # Git ignore file
└── README.md               # Project documentation
```

## Getting Started

1. Clone the repository:
   ```
   git clone <repository-url>
   cd fiber-rest-api
   ```

2. Install the dependencies:
   ```
   go mod tidy
   ```

3. Run the application:
   ```
   go run cmd/server/main.go
   ```

4. Access the API:
   Open your browser or use a tool like `curl` or Postman to send a GET request to `http://localhost:3000/`. You should receive a response:
   ```json
   {
       "message": "hello world"
   }
   ```

## License

This project is licensed under the MIT License.