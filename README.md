# My Go Project

This is a simple Go project that demonstrates a web server with user registration, authentication, and phone management functionality.

## Features

- User registration: Allows users to register by providing their login, password, name, and age.
- User authentication: Handles user authentication using JWT (JSON Web Token).
- Phone management: Supports adding, updating, searching, and deleting phone numbers associated with users.

## Prerequisites

- Go programming language (version 1.20.1)
- SQLite database (version 1.14.17)

## Getting Started

1. Clone the repository:

```
git clone https://github.com/beccoder/SimpleWebServer_2.0.git
```

2. Navigate to the project directory:

```
cd SimpleWebServer_2.0
```

3. Install the project dependencies:

```
go get -v ./...
```
- Or

```
 go mod download
 ```

4. Build the project:

```
go build -o myapp
```

- With CGO_ENABLED=1:

```
go build -tags cgo -o myapp
```

5. Run the application:

```
./myapp
```

6. Open your web browser and visit `http://localhost:8080` to access the application.

## Configuration

The application uses a SQLite database to store user and phone data. The database file is located at `./database.db` by default. You can modify the database configuration in the `main.go` file.

## API Endpoints

- `POST /user/register`: Register a new user by providing login, password, name, and age.
- `POST /user/auth`: Authenticate a user by providing login and password to obtain a JWT token.
- `GET /user/:name`: Get user information by name.
- `POST /user/phone`: Add a new phone number for the authenticated user.
- `GET /user/phone`: Search phone numbers associated with the authenticated user.
- `PUT /user/phone`: Update an existing phone number for the authenticated user.
- `DELETE /user/phone/:phone_id`: Delete a phone number by ID for the authenticated user.

## Sending Data Samples

### Endpoint: /user/register

- Method: POST

- Request Body:

```json
{
  "login": "john123",
  "password": "password123",
  "name": "John Doe",
  "age": 30
}
```

### Endpoint: /user/auth

- Method: POST

- Request Body:

```json
{
  "login": "john123",
  "password": "password123"
}
```

### Endpoint: /user/phone

- Method: POST

- Request Body:

```json
{
  "phone_number": "1234567890",
  "description": "Mobile",
  "is_fax": false
}
```

### Endpoint: /user/phone

- Method: PUT

- Request Body:

```json
{
  "id": 1,
  "phone_number": "9999999999",
  "description": "Home",
  "is_fax": false
}
```

These are just sample data samples in case you want to test them.