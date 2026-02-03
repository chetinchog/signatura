# Claude Code Instructions for Signatura Go Client

## Project Overview

This is a production-ready Go client for the Signatura electronic signature API. The project follows Go best practices with comprehensive testing (90%+ coverage) and zero external dependencies.

**Module:** `github.com/chetinchog/signatura`
**Go Version:** 1.21+
**License:** MIT

## Development Guidelines

### Code Style
- Follow standard Go conventions (gofmt, go vet)
- Use meaningful variable names
- Keep functions focused and short (prefer < 50 lines)
- Add doc comments for all exported symbols
- Start doc comments with the symbol name

### Testing Requirements
- Maintain 70-90% test coverage (current: 90.3%)
- Write table-driven tests where applicable
- Mock external dependencies using httptest.Server
- Test both happy paths and error cases
- Test input validation thoroughly
- Test context cancellation where applicable

### Before Committing
1. Run `make test` - all tests must pass
2. Run `make fmt` - code must be formatted
3. Run `make vet` - no vet warnings
4. Run `make coverage` - verify coverage is maintained
5. Update tests if adding new functionality
6. Update README if adding new features or changing APIs

### Key Project Constraints
- **Zero external dependencies** - only use Go stdlib
- **Context-aware** - all API methods accept context.Context
- **Type-safe** - use strongly typed structs, avoid interface{} (except for metadata)
- **Error handling** - wrap errors with context, return descriptive messages
- **Input validation** - validate all user inputs before API calls
- **No comments inside functions** - per user preference in ~/.claude/CLAUDE.md

## Important Files

### Source Files
- `client.go` (416 lines) - Core API client with 7 methods
- `helpers.go` (145 lines) - Helper functions and builders (14 functions)

### Test Files
- `client_test.go` - Comprehensive unit tests for client (90%+ coverage)
- `helpers_test.go` - Unit tests for all helper functions (100% coverage)

### Configuration
- `go.mod` - Module definition
- `Makefile` - Development tasks
- `.gitignore` - Git ignore rules
- `.golangci.yml` - Linter configuration
- `.github/workflows/ci.yml` - CI/CD pipeline

### Documentation
- `README.md` - User-facing documentation with examples
- `PROYECTO_COMPLETO.md` - Project overview (Spanish)
- `CLAUDE.md` - This file (AI assistant guidance)

## API Structure

### Base URL
Production: `https://connect.signatura.co/api/v2`

### Authentication
Uses Bearer token authentication via `Authorization` header.

### API Methods (7 total)

1. **CreateDocument** - Creates a new document with signatures
   - Validates: title, fileContent, signatures array
   - Returns: CreateDocumentResponse

2. **GetDocument** - Retrieves document by ID
   - Validates: documentID not empty
   - Returns: GetDocumentResponse

3. **ListDocuments** - Lists documents with optional filters
   - Validates: none (all params optional)
   - Returns: ListDocumentsResponse

4. **CancelDocument** - Cancels a pending document
   - Validates: documentID not empty
   - Returns: error only

5. **DownloadDocument** - Downloads signed PDF
   - Validates: documentID not empty
   - Returns: []byte (PDF content)

6. **ResendInvitation** - Resends invitation to signer
   - Validates: documentID and signatureID not empty
   - Returns: error only

7. **New** - Creates a new client instance
   - Config fields: APIKey, BaseURL (optional), HTTPClient (optional)

### Validation Types

| Type | Field | JSON Key | Description |
|------|-------|----------|-------------|
| Email | *string | EM | Email address for validation |
| Phone | *string | PH | Phone number for validation |
| Biometric | *bool | BI | Biometric validation flag |
| AFIP | *string | AF | AFIP CUIT/CUIL |

**Important:** At least one validation method is required per signature.

### Helper Functions (14 total)

#### Builders
- `String(s string) *string` - Helper to create string pointer
- `NewEmailValidation(email, inviteByEmail)` - Create email signature
- `NewEmailPhoneValidation(email, phone, inviteByEmail)` - Email + phone signature
- `NewBiometricValidation()` - Create biometric signature
- `NewAFIPValidation(cuitCuil)` - Create AFIP signature

#### Encoding
- `EncodeFileToBase64(filePath)` - Encode file to base64
- `EncodeBytesToBase64(data)` - Encode bytes to base64
- `DecodeBase64(encoded)` - Decode base64 string

#### Status Checkers
- `IsDocumentPending(doc)` - Check if pending
- `IsDocumentCompleted(doc)` - Check if completed
- `IsDocumentCanceled(doc)` - Check if canceled

#### Signature Filters
- `GetSignedSignatures(doc)` - Filter signed signatures (SignedAt != nil)
- `GetPendingSignatures(doc)` - Filter pending (SignedAt == nil && DeclinedAt == nil)
- `GetDeclinedSignatures(doc)` - Filter declined (DeclinedAt != nil)

## Common Tasks

### Add New API Method
1. Read API documentation
2. Add method to `client.go`
3. Add input validation
4. Add comprehensive tests to `client_test.go`
5. Update README with example
6. Run `make test && make coverage`

### Add New Helper Function
1. Implement in `helpers.go`
2. Add doc comment
3. Add tests to `helpers_test.go`
4. Update README if user-facing
5. Run `make test && make coverage`

### Fix Bug
1. Write failing test first (TDD)
2. Fix the bug
3. Verify test passes
4. Check coverage maintained
5. Commit with descriptive message

### Add Validation
1. Add validation check at start of method
2. Return descriptive error message
3. Add test case for validation error
4. Update documentation if affects API

## Testing Strategy

### Mock HTTP Responses
Use `httptest.NewServer` for mocking API responses:

```go
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // Verify request
    if r.Method != http.MethodPost {
        t.Errorf("Expected POST, got %s", r.Method)
    }

    // Return mock response
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}))
defer server.Close()

client := New(Config{BaseURL: server.URL, APIKey: "test"})
```

### Test Context Cancellation
Test that long operations respect context:

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
defer cancel()

_, err := client.CreateDocument(ctx, request)
// Should return context deadline exceeded error
```

### Test Input Validation
Test all validation error cases:

```go
tests := []struct {
    name    string
    input   SomeRequest
    wantErr string
}{
    {"empty title", SomeRequest{Title: ""}, "title cannot be empty"},
    // ... more cases
}
```

### Test Error Responses
Test API error handling:

```go
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(Error{Code: "TEST", Message: "error"})
}))
```

## Module Information

**Import Path:**
```go
import "github.com/chetinchog/signatura"
```

**Installation:**
```bash
go get github.com/chetinchog/signatura
```

**Dependencies:**
None - uses only Go standard library

## CI/CD

### GitHub Actions
- Runs on push to `main` and `prod` branches
- Tests on Go 1.21, 1.22, 1.23
- Runs tests, vet, fmt, lint
- Uploads coverage to Codecov
- All checks must pass before merge

### Local Verification
Before pushing:
```bash
make test      # Run tests
make coverage  # Check coverage
make fmt       # Format code
make vet       # Run vet
make lint      # Run linter (if installed)
```

## Error Handling Patterns

### API Errors
API errors are returned as `*Error` type:

```go
if err != nil {
    if apiErr, ok := err.(*signatura.Error); ok {
        log.Printf("API Error: %s - %s", apiErr.Code, apiErr.Message)
    }
}
```

### Input Validation Errors
Validation errors are standard errors:

```go
if req.Title == "" {
    return nil, fmt.Errorf("title cannot be empty")
}
```

### Wrapped Errors
HTTP and JSON errors are wrapped for context:

```go
if err := json.Marshal(req); err != nil {
    return nil, fmt.Errorf("failed to marshal request: %w", err)
}
```

## Breaking Changes Policy

This is a v1 package. Breaking changes require:
1. Major version bump (v2.0.0)
2. Detailed migration guide
3. Deprecation warnings in v1.x

Non-breaking changes:
- Adding new methods (OK)
- Adding new optional fields (OK)
- Adding new validation methods (OK)
- Improving error messages (OK)

Breaking changes:
- Changing method signatures (NOT OK)
- Removing methods (NOT OK)
- Changing struct fields (NOT OK)
- Changing validation behavior (NOT OK)

## Performance Considerations

- Use context timeouts for all API calls
- Default HTTP client timeout: 30 seconds
- API rate limits apply (check Signatura docs)
- Large PDFs should be streamed, not loaded entirely in memory
- No connection pooling limits (uses Go defaults)

## Security Considerations

- API keys should be stored in environment variables
- Never commit API keys to version control
- Use HTTPS for all API communication (enforced)
- Validate all inputs before sending to API
- Handle webhook signatures if implementing webhooks

## Future Enhancements

Planned features (not yet implemented):
- Webhook signature verification helper
- Audit trail download method
- Rate limiting handling with backoff
- Request ID tracking for debugging

When implementing these, follow the patterns established in existing code.

## Support

- **Issues:** https://github.com/chetinchog/signatura/issues
- **API Docs:** https://docs.signatura.co
- **Signatura Support:** help@signatura.co

---

**Last Updated:** 2026-02-03
**Coverage:** 90.3%
**Go Version:** 1.21+
