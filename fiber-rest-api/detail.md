# API Detail

This document describes the request flow for user registration, authentication, profile retrieval and update, and the database schema. Diagrams are provided in Mermaid format.

## Sequence Diagram

```mermaid
sequenceDiagram
    participant C as Client (Browser / API Client)
    participant UI as Profile UI
    participant S as API Server (Fiber)
    participant DB as SQLite

    Note over C,S: Register
    C->>S: POST /auth/register { email, password }
    S->>DB: INSERT INTO users (email, password)
    DB-->>S: OK (new id)
    S-->>C: 201 Created

    Note over C,S: Login
    C->>S: POST /auth/login { email, password }
    S->>DB: SELECT id, password FROM users WHERE email = ?
    DB-->>S: row (id, hashed_password)
    S->>S: verify password, create JWT
    S-->>C: 200 { token }

    Note over C,S: View Profile (authenticated)
    C->>S: GET /profile (Authorization: Bearer <token>)
    S->>S: Auth middleware validate JWT -> extracts user id
    S->>DB: SELECT email, first_name, last_name, phone FROM users WHERE id = ?
    DB-->>S: row (email, first_name, last_name, phone)
    S-->>C: 200 { id, email, first_name, last_name, phone }

    Note over C,S: Update Profile (authenticated)
    C->>S: PUT /profile (Authorization: Bearer <token>) { first_name, last_name, phone }
    S->>S: Auth middleware validate JWT -> extracts user id
    S->>DB: UPDATE users SET first_name=?, last_name=?, phone=? WHERE id=?
    DB-->>S: OK
    S->>DB: SELECT email, first_name, last_name, phone FROM users WHERE id = ?
    DB-->>S: updated row
    S-->>C: 200 { id, email, first_name, last_name, phone }

    Note over C,UI: Profile UI
    C->>UI: Open /profile/ui
    UI-->>C: HTML page (fetches /profile with token)
```

## ER Diagram

```mermaid
erDiagram
    USERS {
        INTEGER id PK "autoincrement"
        TEXT email "unique, not null"
        TEXT password "hashed"
        TEXT first_name
        TEXT last_name
        TEXT phone
    }

    %% No other tables in current schema
```

## Notes
- JWT: The server issues a signed JWT on /auth/login. The token must be provided as an Authorization header: `Bearer <token>` for protected endpoints.
- Migration: If you already have an existing SQLite database without profile columns, add columns with:
  - ALTER TABLE users ADD COLUMN first_name TEXT;
  - ALTER TABLE users ADD COLUMN last_name TEXT;
  - ALTER TABLE users ADD COLUMN phone TEXT;

- Rendering: Many Markdown viewers (VS Code, GitHub) can render Mermaid diagrams with appropriate plugins or built-in support. Use a Mermaid live editor to preview if needed.
