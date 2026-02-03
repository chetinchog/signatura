package signatura

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty string", "", ""},
		{"normal string", "test", "test"},
		{"unicode string", "测试", "测试"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := String(tt.input)
			if result == nil {
				t.Fatal("String() returned nil")
			}
			if *result != tt.want {
				t.Errorf("String() = %v, want %v", *result, tt.want)
			}
		})
	}
}

func TestEncodeFileToBase64(t *testing.T) {
	// Create temp file for testing
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	testContent := []byte("Hello, World!")

	if err := os.WriteFile(testFile, testContent, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	t.Run("valid file", func(t *testing.T) {
		result, err := EncodeFileToBase64(testFile)
		if err != nil {
			t.Errorf("EncodeFileToBase64() error = %v", err)
			return
		}
		if result == "" {
			t.Error("EncodeFileToBase64() returned empty string")
		}

		// Verify we can decode it back
		decoded, err := DecodeBase64(result)
		if err != nil {
			t.Errorf("Failed to decode result: %v", err)
		}
		if string(decoded) != string(testContent) {
			t.Errorf("Decoded content = %v, want %v", string(decoded), string(testContent))
		}
	})

	t.Run("file not found", func(t *testing.T) {
		_, err := EncodeFileToBase64("/nonexistent/file.txt")
		if err == nil {
			t.Error("EncodeFileToBase64() expected error for nonexistent file")
		}
	})
}

func TestEncodeBytesToBase64(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{"empty bytes", []byte{}},
		{"simple text", []byte("test")},
		{"binary data", []byte{0x00, 0xFF, 0xAA, 0x55}},
		{"large data", make([]byte, 1024)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EncodeBytesToBase64(tt.input)

			// Decode and verify
			decoded, err := DecodeBase64(result)
			if err != nil {
				t.Errorf("Failed to decode result: %v", err)
			}
			if string(decoded) != string(tt.input) {
				t.Errorf("Decoded content doesn't match original")
			}
		})
	}
}

func TestDecodeBase64(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []byte
		wantErr bool
	}{
		{"valid base64", "SGVsbG8sIFdvcmxkIQ==", []byte("Hello, World!"), false},
		{"empty string", "", []byte{}, false},
		{"invalid base64", "not-valid-base64!", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeBase64(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeBase64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && string(got) != string(tt.want) {
				t.Errorf("DecodeBase64() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func TestNewEmailValidation(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		inviteByEmail bool
	}{
		{"with invite", "test@example.com", true},
		{"without invite", "test@example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sig := NewEmailValidation(tt.email, tt.inviteByEmail)

			if sig.Validations.Email == nil {
				t.Error("Email validation is nil")
			} else if *sig.Validations.Email != tt.email {
				t.Errorf("Email = %v, want %v", *sig.Validations.Email, tt.email)
			}

			if tt.inviteByEmail {
				if len(sig.InviteChannel) == 0 {
					t.Error("InviteChannel not set")
				} else if sig.InviteChannel[0] != InviteChannelEmail {
					t.Errorf("InviteChannel = %v, want %v", sig.InviteChannel[0], InviteChannelEmail)
				}
			} else {
				if len(sig.InviteChannel) != 0 {
					t.Error("InviteChannel should be empty")
				}
			}
		})
	}
}

func TestNewEmailPhoneValidation(t *testing.T) {
	email := "test@example.com"
	phone := "+1234567890"
	sig := NewEmailPhoneValidation(email, String(phone), true)

	if sig.Validations.Email == nil {
		t.Error("Email validation is nil")
	} else if *sig.Validations.Email != email {
		t.Errorf("Email = %v, want %v", *sig.Validations.Email, email)
	}

	if sig.Validations.Phone == nil {
		t.Error("Phone should not be nil")
	} else if *sig.Validations.Phone != phone {
		t.Errorf("Phone = %v, want %v", *sig.Validations.Phone, phone)
	}

	if len(sig.InviteChannel) == 0 {
		t.Error("InviteChannel not set")
	} else if sig.InviteChannel[0] != InviteChannelEmail {
		t.Errorf("InviteChannel = %v, want %v", sig.InviteChannel[0], InviteChannelEmail)
	}
}

func TestNewBiometricValidation(t *testing.T) {
	sig := NewBiometricValidation()

	if sig.Validations.Biometric == nil {
		t.Error("Biometric validation is nil")
	} else if !*sig.Validations.Biometric {
		t.Error("Biometric should be true")
	}

	if len(sig.InviteChannel) != 0 {
		t.Error("InviteChannel should be empty for biometric")
	}
}

func TestNewAFIPValidation(t *testing.T) {
	cuit := "20-12345678-9"
	sig := NewAFIPValidation(cuit)

	if sig.Validations.AFIP == nil {
		t.Error("AFIP validation is nil")
	} else if *sig.Validations.AFIP != cuit {
		t.Errorf("AFIP = %v, want %v", *sig.Validations.AFIP, cuit)
	}
}

func TestIsDocumentPending(t *testing.T) {
	tests := []struct {
		name   string
		status string
		want   bool
	}{
		{"pending", DocumentStatusPending, true},
		{"completed", DocumentStatusCompleted, false},
		{"canceled", DocumentStatusCanceled, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &GetDocumentResponse{Status: tt.status}
			if got := IsDocumentPending(doc); got != tt.want {
				t.Errorf("IsDocumentPending() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsDocumentCompleted(t *testing.T) {
	tests := []struct {
		name   string
		status string
		want   bool
	}{
		{"pending", DocumentStatusPending, false},
		{"completed", DocumentStatusCompleted, true},
		{"canceled", DocumentStatusCanceled, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &GetDocumentResponse{Status: tt.status}
			if got := IsDocumentCompleted(doc); got != tt.want {
				t.Errorf("IsDocumentCompleted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsDocumentCanceled(t *testing.T) {
	tests := []struct {
		name   string
		status string
		want   bool
	}{
		{"pending", DocumentStatusPending, false},
		{"completed", DocumentStatusCompleted, false},
		{"canceled", DocumentStatusCanceled, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &GetDocumentResponse{Status: tt.status}
			if got := IsDocumentCanceled(doc); got != tt.want {
				t.Errorf("IsDocumentCanceled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSignedSignatures(t *testing.T) {
	now := time.Now()
	doc := &GetDocumentResponse{
		Signatures: []SignatureResponse{
			{ID: "1", SignedAt: &now},
			{ID: "2", SignedAt: nil},
			{ID: "3", SignedAt: &now},
			{ID: "4", DeclinedAt: &now},
		},
	}

	signed := GetSignedSignatures(doc)
	if len(signed) != 2 {
		t.Errorf("GetSignedSignatures() returned %d signatures, want 2", len(signed))
	}

	for _, sig := range signed {
		if sig.SignedAt == nil {
			t.Errorf("Expected all signatures to have SignedAt set")
		}
	}
}

func TestGetPendingSignatures(t *testing.T) {
	now := time.Now()
	doc := &GetDocumentResponse{
		Signatures: []SignatureResponse{
			{ID: "1", SignedAt: &now},
			{ID: "2", SignedAt: nil, DeclinedAt: nil},
			{ID: "3", SignedAt: nil, DeclinedAt: nil},
			{ID: "4", DeclinedAt: &now},
		},
	}

	pending := GetPendingSignatures(doc)
	if len(pending) != 2 {
		t.Errorf("GetPendingSignatures() returned %d signatures, want 2", len(pending))
	}

	for _, sig := range pending {
		if sig.SignedAt != nil || sig.DeclinedAt != nil {
			t.Errorf("Expected pending signatures to have nil SignedAt and DeclinedAt")
		}
	}
}

func TestGetDeclinedSignatures(t *testing.T) {
	now := time.Now()
	doc := &GetDocumentResponse{
		Signatures: []SignatureResponse{
			{ID: "1", SignedAt: &now},
			{ID: "2", SignedAt: nil, DeclinedAt: nil},
			{ID: "3", DeclinedAt: &now},
			{ID: "4", DeclinedAt: &now},
		},
	}

	declined := GetDeclinedSignatures(doc)
	if len(declined) != 2 {
		t.Errorf("GetDeclinedSignatures() returned %d signatures, want 2", len(declined))
	}

	for _, sig := range declined {
		if sig.DeclinedAt == nil {
			t.Errorf("Expected all declined signatures to have DeclinedAt set")
		}
	}
}
