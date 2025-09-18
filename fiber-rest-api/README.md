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

## New features

This project now includes profile management and avatar upload:

- GET /profile (protected) - return id, email, first_name, last_name, phone, avatar
- PUT /profile (protected) - update first_name, last_name, phone
- POST /profile/avatar (protected) - upload avatar (multipart/form-data)
- GET /profile/ui - minimal web UI to view/edit profile and upload avatar
- GET /docs and GET /docs/swagger.json - OpenAPI + Swagger UI

Authentication is JWT-based. Obtain a token via POST /auth/login then send it as Authorization: Bearer <token>.

## Database migration

If you have an existing `data.db` created before these changes, the users table may lack new columns. To migrate manually, run these SQL statements against the sqlite database file (example uses the bundled `sqlite3` tool):

```sh
sqlite3 data.db "ALTER TABLE users ADD COLUMN first_name TEXT;"
sqlite3 data.db "ALTER TABLE users ADD COLUMN last_name TEXT;"
sqlite3 data.db "ALTER TABLE users ADD COLUMN phone TEXT;"
sqlite3 data.db "ALTER TABLE users ADD COLUMN avatar TEXT;"
```

Note: I ran these ALTER TABLE commands against `fiber-rest-api/data.db` in the workspace and confirmed the columns were added.

## Curl examples

Register:
```sh
curl -X POST http://localhost:3000/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"email":"user@example.com","password":"secret"}'
```

Login (returns JWT):
```sh
curl -X POST http://localhost:3000/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"user@example.com","password":"secret"}'
```

Get profile:
```sh
curl -H "Authorization: Bearer <token>" http://localhost:3000/profile
```

Update profile:
```sh
curl -X PUT http://localhost:3000/profile \
  -H "Authorization: Bearer <token>" \
  -H 'Content-Type: application/json' \
  -d '{"first_name":"ชื่อ","last_name":"นามสกุล","phone":"0812345678"}'
```

Upload avatar:
```sh
curl -X POST http://localhost:3000/profile/avatar \
  -H "Authorization: Bearer <token>" \
  -F "avatar=@/path/to/avatar.jpg"
```

Open profile UI in a browser:

http://localhost:3000/profile/ui

Open Swagger UI:

http://localhost:3000/docs

## Notes
- Uploaded avatars are served from `/uploads` (stored in the `uploads/` directory).
- Avatar upload accepts common image extensions (.jpg/.jpeg/.png/.gif) and limits size to 5MB.
- Profile fields have basic server-side validation (max lengths and phone character checks).

## License

This project is licensed under the MIT License.