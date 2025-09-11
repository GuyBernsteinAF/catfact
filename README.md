# Cat Facts API

A Go application that fetches and displays cat facts from the Cat Facts API. The application can run as either a CLI client or an HTTP server with API endpoints.

## Features

- **CLI Client**: Interactive command-line interface for fetching cat facts
- **HTTP Server**: RESTful API endpoints for programmatic access
- **Three Phases**: Different modes of operation for fetching facts
    - Phase 1: Fetch a single cat fact
    - Phase 2: Fetch 5 cat facts sequentially
    - Phase 3: Fetch 10 cat facts concurrently using goroutines

## Project Structure

```
catfacts/
├── cmd/
│   ├── client/
│   │   └── main.go     # CLI client application
│   └── server/
│       └── main.go     # HTTP server application
├── internal/
│   └── phases.go       # Core business logic
├── docs/
│   └── swagger.yaml    # API documentation
├── go.mod
└── README.md
```

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd catfacts
```

2. Install dependencies:
```bash
go mod download
```

## Usage

### Running the CLI Client

```bash
go run cmd/client/main.go
```

The client will prompt you to select a phase (1, 2, or 3) and display the corresponding cat facts.

### Running the HTTP Server

```bash
go run cmd/server/main.go
```

The server will start on port 8090 by default.

## API Endpoints

### Base URL
```
http://localhost:8090
```

### Endpoints

#### GET /phase-one
Fetches a single cat fact.

**Response:**
```
A cat fact string
```

#### GET /phase-two
Fetches 5 cat facts sequentially.

**Response:**
```
1. First cat fact
2. Second cat fact
3. Third cat fact
4. Fourth cat fact
5. Fifth cat fact
```

#### GET /phase-three
Fetches 10 cat facts concurrently using goroutines.

**Response:**
```
1. First cat fact
2. Second cat fact
...
10. Tenth cat fact
```

#### GET /headers
Debug endpoint that returns all request headers.

**Response:**
```
Header-Name: header-value
...
```

## API Documentation

For detailed API documentation, see the [Swagger/OpenAPI specification](docs/swagger.yaml).

To view the Swagger documentation in a browser, you can use tools like:
- [Swagger Editor](https://editor.swagger.io/) - Paste the swagger.yaml content
- [Swagger UI](https://swagger.io/tools/swagger-ui/) - For local viewing

## Development

### Building the Applications

Build the client:
```bash
go build -o bin/client cmd/client/main.go
```

Build the server:
```bash
go build -o bin/server cmd/server/main.go
```

### Testing

Run tests:
```bash
go test ./...
```

### Adding Swagger UI (Optional)

To add Swagger UI to your server, you can use the `swaggo/http-swagger` package:

1. Install the package:
```bash
go get -u github.com/swaggo/http-swagger
go get -u github.com/swaggo/files
```

2. Add the Swagger UI handler to your server (see server code for implementation)

## External Dependencies

- [Cat Facts API](https://catfact.ninja/) - External API for cat facts

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Acknowledgments

- Cat Facts API (https://catfact.ninja/) for providing the cat facts data