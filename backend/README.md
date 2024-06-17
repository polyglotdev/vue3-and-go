# Golang Backend API for Vue.js Application

[Repo for frontend app](https://github.com/polyglotdev/vuejs-frontend)

## Getting Started

### Prerequisites

- Go 1.22.4 or higher
- Git

### Installation

1. Clone the repository:

```bash
git clone https://github.com/polyglotdev/golang-backend-api.git
```

2. Change into the project directory:

```bash
cd golang-backend-api
```

3. Install the dependencies:

```bash
go mod download
```

4. Run the application:

```bash
go run cmd/api/main.go
```

5. Open your web browser and navigate to `http://localhost:8081/users/login`.

## Features

- User login endpoint
- JSON response formatting
- Error handling

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Design Considerations

- Use the `chi` framework for routing and middleware.
- Use the `go-chi/cors` middleware for CORS support.
- Use the `go-chi/chi/v5/middleware` package for middleware.
- Use the `encoding/json` package for JSON encoding and decoding.
- Use the `log` package for logging.
- Use the `errors` package for error handling.
- Use the `io` package for reading and writing data.
- Use the `net/http` package for HTTP requests and responses.
- Use the `http.MaxBytesReader` function to limit the size of the request body.
- Use the `json.NewDecoder` function to decode JSON data.
- Use the `json.NewEncoder` function to encode JSON data.
- Use the `json.MarshalIndent` function to marshal JSON data with indentation.
- Use the `json.Unmarshal` function to unmarshal JSON data.
- Use the `http.Error` function to send HTTP errors.
- Use the `http.Header` type to set HTTP headers.
- Use the `http.ResponseWriter` interface to write HTTP responses.
- Use the `io.EOF` error to check if the request body contains only a single JSON value.
- Use the `errors.New` function to create custom error messages.
