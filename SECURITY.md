# Security Policy

## Supported Versions

We release patches for security vulnerabilities for the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

If you discover a security vulnerability within this project, please send an email to the maintainers. All security vulnerabilities will be promptly addressed.

**Please do not report security vulnerabilities through public GitHub issues.**

### What to Include in Your Report

- A description of the vulnerability
- Steps to reproduce the issue
- Potential impact of the vulnerability
- Any possible mitigations you've identified

### Response Timeline

- **Initial Response**: Within 48 hours of report submission
- **Status Update**: Within 7 days with assessment and planned timeline
- **Resolution**: Depends on severity, but critical issues will be prioritized

### Disclosure Policy

- Security vulnerabilities will be disclosed publicly only after a fix is available
- We will credit researchers who responsibly disclose vulnerabilities (unless they prefer to remain anonymous)
- We aim to release security patches within 30 days of a verified report

## Security Best Practices

When using the Signatura Go Client:

### API Key Protection

- **Never commit API keys** to version control
- Store API keys in environment variables or secure secret management systems
- Rotate API keys regularly
- Use different API keys for development, staging, and production environments

```go
// Good - Use environment variables
apiKey := os.Getenv("SIGNATURA_API_KEY")

// Bad - Never hardcode
apiKey := "sk_live_1234567890abcdef" // DON'T DO THIS
```

### Input Validation

- Always validate user input before passing to API methods
- The client performs basic validation, but additional application-level validation is recommended
- Sanitize file names and titles before passing to `CreateDocument`

### Error Handling

- Handle all errors returned by the client
- Don't expose raw API error messages to end users
- Log errors securely without including sensitive data

```go
doc, err := client.GetDocument(ctx, docID)
if err != nil {
    log.Printf("Failed to retrieve document: %v", err)
    // Return generic error to user
    return fmt.Errorf("unable to retrieve document")
}
```

### Context Timeouts

- Always use context with timeouts for API calls
- Set reasonable timeout values based on your use case

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

result, err := client.CreateDocument(ctx, req)
```

### Webhook Security

When implementing webhooks:

- Verify webhook signatures (use HMAC verification)
- Use HTTPS endpoints only
- Validate webhook payload structure before processing
- Implement idempotency to handle duplicate webhook deliveries
- Rate limit webhook endpoints to prevent abuse

### TLS/HTTPS

- The client uses HTTPS by default for all API communication
- Never disable TLS certificate verification in production
- Keep Go and system CA certificates up to date

### Dependency Management

- This package has zero external dependencies (stdlib only)
- Keep your Go version up to date to receive security patches
- Regularly update the package to get security fixes: `go get -u github.com/chetinchog/signatura`

### Testing

- Never use production API keys in tests
- Use test mode or sandbox environments when available
- Don't commit test fixtures containing sensitive data

## Known Security Considerations

### Base64 Encoding

- This client uses base64 encoding for file uploads
- Base64 is encoding, not encryption - files are transmitted over HTTPS
- Never base64-encode secrets or passwords

### Context Cancellation

- All methods respect context cancellation
- Ensure contexts are properly canceled to prevent resource leaks

### Memory Management

- Large file downloads are loaded into memory
- For very large files, consider streaming or chunking strategies
- The client does not impose file size limits, but be aware of memory constraints

## Security Updates

Security updates will be released as patch versions (e.g., 1.0.1, 1.0.2) and documented in CHANGELOG.md with a `[SECURITY]` tag.

To receive security notifications:
- Watch this repository on GitHub
- Enable security alerts in your GitHub settings
- Subscribe to release notifications

## Third-Party Security

This package uses only Go standard library. However, be aware of:
- Go runtime security issues (keep Go updated)
- Operating system CA certificate security
- Network layer security (ensure TLS 1.2+ is used)

## Compliance

This client is designed to work with the Signatura electronic signature API. Compliance requirements (GDPR, CCPA, etc.) are primarily the responsibility of:
- The Signatura service provider
- Your application's use of the service

Consult legal counsel to ensure your use case meets regulatory requirements.

## Contact

For security concerns: Open a private security advisory on GitHub or contact the maintainers directly.

For general questions: Open a public issue on GitHub.
