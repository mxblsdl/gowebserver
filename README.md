# Go Web Server

This project is a simple web server built using Go. It demonstrates the use of handlers, middleware, and configuration management.

## Project Structure

```
go-webserver
├── cmd
│   └── main.go          # Entry point of the application
├── internal
│   ├── handlers
│   │   └── handlers.go  # HTTP handler functions
│   ├── middleware
│   │   └── middleware.go # Middleware functions
│   └── models
│       └── models.go    # Data models
├── pkg
│   └── config
│       └── config.go    # Configuration management
├── go.mod                # Module definition
├── go.sum                # Dependency checksums
└── README.md             # Project documentation
```

## Getting Started

### Prerequisites

- Go 1.16 or later

### Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd go-webserver
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

### Running the Server

To run the web server, execute the following command:

```bash
go run cmd/main.go
```

The server will start and listen on the specified port. You can access it in your web browser at `http://localhost:<port>`.

### API Endpoints

- `GET /` - Home page
- `GET /about` - About page

### License

This project is licensed under the MIT License. See the LICENSE file for details.