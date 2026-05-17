# AGENTS.md - Agentic Coding Guidelines for groupie-tracker

This file provides coding guidelines and instructions for AI agents working in this repository.

## Project Overview

- **Language**: Go 1.25
- **Framework**: gorilla/mux for HTTP routing
- **Type**: Web application (server-side HTML rendering)
- **Dependencies**: github.com/gorilla/mux v1.8.1

## Project Structure

```
groupie-tracker/
├── main.go              # Entry point, router setup
├── api/                 # API client layer
│   ├── fetch.go         # HTTP calls to external API
│   └── models.go        # Data structures
├── handlers/            # HTTP request handlers
│   ├── home.go          # Home page handler
│   └── artist.go        # Artist detail page handler
├── viewmodels/          # View data models
│   └── page.go          # PageData struct
├── views/               # HTML templates
│   ├── base.html
│   ├── header.html
│   ├── footer.html
│   ├── index.html
│   └── artist.html
├── static/              # Static assets (CSS, images)
├── go.mod
└── go.sum
```

## Build & Run Commands

### Build
```bash
go build -o groupie-tracker .
```

### Run
```bash
go run .
# Or after building:
./groupie-tracker
```
Server runs at http://localhost:8000

### Dependencies
```bash
go mod download
go mod tidy
```

### Testing
No tests currently exist in this project. To add tests:
```bash
# Run all tests
go test ./...

# Run tests in specific package
go test ./api
go test ./handlers

# Run a single test
go test -run TestFunctionName ./...

# Run with verbose output
go test -v ./...
```

### Linting
This project does not have a linter configured. To add linting, install golangci-lint:
```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
golangci-lint run
```

### Formatting
```bash
# Format code (gofmt)
gofmt -w .

# Format and organize imports
go fmt ./...

# Show differences (don't modify)
gofmt -d .
```

## Code Style Guidelines

### Imports
- Use Go's standard import organization:
  1. Standard library (fmt, net/http, etc.)
  2. Third-party packages (github.com/gorilla/mux)
  3. Internal packages (main/api, main/handlers)
- Use import aliases only when necessary (e.g., `mux "github.com/gorilla/mux"`)

### Formatting
- Use `go fmt` or gofmt for consistent formatting
- Maximum line length: 120 characters (soft limit)
- Use tabs for indentation, not spaces

### Types & Declarations
- Use explicit type annotations for function return types
- Prefer concrete types over interfaces unless polymorphism is needed
- Use struct tags for JSON serialization (e.g., `json:"name"`)

### Naming Conventions
- **Packages**: lowercase, short, no underscores (api, handlers, viewmodels)
- **Files**: lowercase with underscores (home.go, artist_page.go)
- **Types/Structs**: PascalCase (Artist, PageData)
- **Functions/Variables**: camelCase (GetArtists, handleHome)
- **Constants**: PascalCase or camelCase depending on scope
- **Acronyms**: preserve original casing (URL, ID, API)

### Error Handling
- Always handle errors explicitly - do not ignore with `_`
- Return meaningful error values from functions
- Usefmt.Println for logging errors in handlers (current pattern)
- Consider using structured logging for production code
- Return appropriate HTTP status codes (404, 500, etc.)

### HTTP Handlers
- Handler functions should have signature: `func(w http.ResponseWriter, r *http.Request)`
- Use mux.Vars(r) to extract path parameters
- Check r.Method for request type validation
- Use http.Error for error responses with appropriate status codes
- Close response bodies (already handled in api/fetch.go with defer)

### Template Usage
- Use template.Must for template parsing at package level
- Define template functions in template.FuncMap
- Use ExecuteTemplate with explicit template name (e.g., "base.html")

### Constants
- Define constants for magic numbers (e.g., itemsPerPage := 12)
- Consider constants for URLs if used across multiple files

## Best Practices

1. **Defer cleanup**: Always defer resp.Body.Close() after HTTP requests
2. **Input validation**: Validate URL parameters (artist IDs, page numbers)
3. **Error pages**: Return appropriate HTTP status codes on errors
4. **Graceful shutdown**: Not implemented but recommended for production
5. **Static files**: Serve from /static/ prefix as configured in main.go

## Common Patterns

### Making HTTP API calls
```go
resp, err := http.Get(url)
if err != nil {
    return nil, err
}
defer resp.Body.Close()
// decode resp.Body into target type
```

### Route definition (main.go)
```go
mux := mux.NewRouter()
mux.HandleFunc("/path", handlers.HandlerFunction)
mux.HandleFunc("/path/{id}", handlers.HandlerFunctionWithParam)
```

### Handler with path parameter
```go
func Handler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
    // convert id to appropriate type
}
```

## Notes for Agents

- This is a learning/project exercise codebase
- External API: https://groupietrackers.herokuapp.com/api
- No database - all data fetched from external API
- Templates use Go's html/template (auto-escapes HTML)
- Current templating approach uses ParseFiles for modular templates