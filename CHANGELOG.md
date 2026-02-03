# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.1] - 2026-02-03

### Added
- Coverage and test badges in README (90.5%, 50+ tests)
- Documentation of `ValidationValue` and `ValidationResponse` structures
- Signature status codes documentation (IN, SI, DE, PE)
- CHANGELOG.md with comprehensive release notes
- SECURITY.md with security policy and best practices
- References to all documentation files in README

### Fixed
- README field name errors: `SigningURL` → `URL`
- Removed references to non-existent `CompletedAt` field
- ListDocuments example now uses correct fields (`CreationDate`, `Status`)
- All documentation examples validated against real API

### Changed
- Updated README badges to reflect accurate metrics
- Improved documentation organization and structure
- Makefile now excludes examples/ directory from build targets

## [1.0.0] - 2026-02-03

### Added
- Initial production-ready release of Signatura Go Client
- Complete API client implementation with 7 core methods:
  - `CreateDocument` - Create documents for electronic signature
  - `GetDocument` - Retrieve document details
  - `ListDocuments` - List documents with filtering and pagination
  - `CancelDocument` - Cancel pending documents
  - `DownloadDocument` - Download signed document PDFs
  - `ResendInvitation` - Resend signature invitations
- Comprehensive helper functions (14 total):
  - String pointer helper
  - Base64 encoding/decoding utilities
  - Signature builder functions for email, phone, biometric, and AFIP validations
  - Document status check helpers
  - Signature filtering helpers
- Full validation response support with `ValidationValue` and `ValidationResponse` types
- Complete test suite with 90.3% code coverage
- GitHub Actions CI/CD pipeline
- Development tooling: Makefile, golangci-lint configuration
- AI assistant guidance in CLAUDE.md
- Comprehensive documentation in README.md and PROYECTO_COMPLETO.md

### Changed
- **BREAKING**: Updated all API response structures to match actual Signatura API v2
- **BREAKING**: `ListDocumentsResponse` now uses `results` and `count` fields instead of `documents` and `total_count`
- **BREAKING**: Validation fields now use nested `ValidationValue` structure with `validated` and `value` properties
- **BREAKING**: `SignatureResponse` uses `url` field instead of `signing_url`
- **BREAKING**: Signature status uses new constants: `SignatureStatusInvited` ("IN"), `SignatureStatusSigned` ("SI"), `SignatureStatusDeclined` ("DE"), `SignatureStatusPending` ("PE")
- **BREAKING**: `Biometric` validation field changed from `*interface{}` to `*bool`
- All helper functions now use `Status` field instead of timestamp fields for signature filtering

### Fixed
- Comprehensive input validation for all API methods
- Error handling consistency across all methods
- Nil pointer safety in response handlers
- Proper context cancellation handling
- All golangci-lint errcheck violations resolved
- Corrected all documentation examples to be runnable
- Fixed import paths throughout documentation

### Security
- Input validation prevents empty or invalid API requests
- Proper error wrapping provides context without exposing sensitive data
- All examples use environment variables for API keys

## [Unreleased]

### Planned
- Webhook signature verification helper
- Audit trail download method
- Rate limiting with exponential backoff
- Request ID tracking for debugging

---

[1.0.1]: https://github.com/chetinchog/signatura/releases/tag/v1.0.1
[1.0.0]: https://github.com/chetinchog/signatura/releases/tag/v1.0.0
