package signatura

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

// Helper functions for working with the Signatura API

// String returns a pointer to the provided string value.
// Useful for optional string fields in validation.
func String(s string) *string {
	return &s
}

// EncodeFileToBase64 reads a file and returns its base64-encoded content
func EncodeFileToBase64(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer func() { _ = file.Close() }()

	data, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return base64.StdEncoding.EncodeToString(data), nil
}

// EncodeBytesToBase64 encodes byte data to base64 string
func EncodeBytesToBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// DecodeBase64 decodes a base64 string to bytes
func DecodeBase64(encoded string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}
	return decoded, nil
}

// NewEmailValidation creates a signature with email validation
func NewEmailValidation(email string, inviteByEmail bool) Signature {
	sig := Signature{
		Validations: Validation{
			Email: String(email),
		},
	}

	if inviteByEmail {
		sig.InviteChannel = []string{InviteChannelEmail}
	}

	return sig
}

// NewEmailPhoneValidation creates a signature with email and phone validation
func NewEmailPhoneValidation(email string, phone *string, inviteByEmail bool) Signature {
	sig := Signature{
		Validations: Validation{
			Email: String(email),
			Phone: phone,
		},
	}

	if inviteByEmail {
		sig.InviteChannel = []string{InviteChannelEmail}
	}

	return sig
}

// NewBiometricValidation creates a signature with biometric validation
func NewBiometricValidation() Signature {
	biometric := true
	return Signature{
		Validations: Validation{
			Biometric: &biometric,
		},
	}
}

// NewAFIPValidation creates a signature with AFIP key validation
func NewAFIPValidation(cuitCuil string) Signature {
	return Signature{
		Validations: Validation{
			AFIP: String(cuitCuil),
		},
	}
}

// IsDocumentPending checks if a document is in pending status
func IsDocumentPending(doc *GetDocumentResponse) bool {
	return doc.Status == DocumentStatusPending
}

// IsDocumentCompleted checks if a document is completed
func IsDocumentCompleted(doc *GetDocumentResponse) bool {
	return doc.Status == DocumentStatusCompleted
}

// IsDocumentCanceled checks if a document is canceled
func IsDocumentCanceled(doc *GetDocumentResponse) bool {
	return doc.Status == DocumentStatusCanceled
}

// GetSignedSignatures returns all signed signatures from a document
func GetSignedSignatures(doc *GetDocumentResponse) []SignatureResponse {
	var signed []SignatureResponse
	for _, sig := range doc.Signatures {
		if sig.SignedAt != nil {
			signed = append(signed, sig)
		}
	}
	return signed
}

// GetPendingSignatures returns all pending signatures from a document
func GetPendingSignatures(doc *GetDocumentResponse) []SignatureResponse {
	var pending []SignatureResponse
	for _, sig := range doc.Signatures {
		if sig.SignedAt == nil && sig.DeclinedAt == nil {
			pending = append(pending, sig)
		}
	}
	return pending
}

// GetDeclinedSignatures returns all declined signatures from a document
func GetDeclinedSignatures(doc *GetDocumentResponse) []SignatureResponse {
	var declined []SignatureResponse
	for _, sig := range doc.Signatures {
		if sig.DeclinedAt != nil {
			declined = append(declined, sig)
		}
	}
	return declined
}
