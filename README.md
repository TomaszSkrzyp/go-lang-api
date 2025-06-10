
# Go/React/PostgreSQL Todo App

This project is a full-stack Todo application featuring a Go backend, React frontend, and PostgreSQL database. It supports JWT-based authentication, CRUD operations for tasks, and a modern CI/CD pipeline.

---

## Features

- **Backend:** Go (Golang) REST API with JWT authentication, task CRUD, and PostgreSQL integration.
- **Frontend:** React app with authentication, task management, and pagination.
- **Database:** PostgreSQL with schema initialization.
- **Authentication:** JWT-based login (hardcoded credentials for demo).
- **CI/CD:** Linting, testing, Docker builds, and deployment via GitLab CI.
- **Dockerized:** Multi-stage Dockerfile for efficient builds.
- **Dev Experience:** Hot reload, CORS enabled for local development.

---

## Getting Started

### Prerequisites

- Docker & Docker Compose
- Node.js (for local frontend dev)
- Go (for local backend dev)

### Environment Variables

Create a `.env` file with:

```env
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name
DB_HOST=todo-postgres-db
```

### Database Schema

Ensure the `schema_dump.sql` file exists in the root directory. It should define the `todos` table and any other required tables, for example:

```sql
CREATE TABLE IF NOT EXISTS todos (
    id SERIAL PRIMARY KEY,
    task TEXT NOT NULL,
    status TEXT NOT NULL,
    due DATE NOT NULL
);
```

**Note:** To populate the database with sample tasks on startup, uncomment the `storage.SeedSampleData()` line in your `main.go` file.

---

## Running Locally

### With Docker Compose

```sh
docker compose up --build
```

- Backend: [http://localhost:8090](http://localhost:8090)
- Frontend: Served statically by backend at `/`

### Without Docker

#### Backend

```sh
cd toDo
go run main.go
```

#### Frontend

```sh
cd todo-frontend
npm install
npm start
```

---

## API Endpoints

| Method | Endpoint            | Description                | Auth Required |
|--------|---------------------|----------------------------|--------------|
| POST   | `/login`            | Obtain JWT token           | No           |
| GET    | `/api/todos`        | List tasks (paginated)     | Yes          |
| POST   | `/api/todos`        | Add new task               | Yes          |
| GET    | `/api/todos/{id}`   | Get task by ID             | Yes          |
| PATCH  | `/api/todos/{id}`   | Update task or status up   | Yes          |
| DELETE | `/api/todos/{id}`   | Remove task                | Yes          |

**Note:** All `/api/todos*` endpoints require `Authorization: Bearer <token>` header.

---

## Authentication

- **Login:** POST `/login` with JSON body:
    ```json
    { "username": "admin", "password": "password123" }
    ```
- **Response:** `{ "token": "<JWT>" }`
- Use this token in the `Authorization` header for all protected endpoints.

---

## Frontend Usage

- Add, edit, remove, and update status of tasks.
- Pagination controls for task list.
- Login/logout functionality.

---

## CI/CD

- **Lint:** Go and React code linted in separate jobs.
- **Test:** Go and React tests run in pipeline.
- **Build:** Multi-stage Docker builds for backend and frontend.
- **Deploy:** Example deploy job (customize as needed).

---

## License

MIT

---

## Credits

- [Golang](https://golang.org/)
- [React](https://reactjs.org/)
- [PostgreSQL](https://www.postgresql.org/)
- [Gorilla Mux](https://github.com/gorilla/mux)
- [golang-jwt](https://github.com/golang-jwt/jwt)

---

```go