# Golang Gin Realworld Example App

This project is a Golang implementation of the RealWorld example app, which provides a full-stack application demonstrating best practices in API design and development. The application is built using the Gin web framework and GORM for ORM.

## Features

- User authentication and management
- Article creation, reading, updating, and deletion (CRUD)
- Article interactions (like, favorite, etc.)
- Integration tests for API endpoints

## Project Structure

```
golang-gin-realworld-example-app
├── .gitignore
├── .travis.yml
├── doc.go
├── go.mod
├── go.mod-generated_by_mod_init
├── go.sum
├── gorm.db
├── hello.go
├── integration_test.go
├── LICENSE
├── logo.png
├── readme.md
├── realworld-server
├── articles/
│   ├── doc.go
│   ├── models.go
│   ├── routers.go
│   ├── serializers.go
│   ├── unit_test.go
│   └── validators.go
├── common/
│   ├── database.go
│   ├── unit_test.go
│   └── utils.go
├── scripts/
│   ├── coverage.sh
│   └── gofmt.sh
└── users/
    ├── doc.go
    ├── middlewares.go
    ├── models.go
    ├── routers.go
    ├── serializers.go
    ├── unit_test.go
    └── validators.go
```

## Getting Started

1. Clone the repository:
   ```
   git clone <repository-url>
   cd golang-gin-realworld-example-app
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Run the application:
   ```
   go run realworld-server
   ```

4. Run integration tests:
   ```
   go test -v integration_test.go
   ```

## API Documentation

Refer to the API documentation in the `doc.go` files for detailed information on the available endpoints and their usage.

## License

This project is licensed under the MIT License. See the LICENSE file for more details.