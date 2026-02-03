package signatura

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	t.Run("with default values", func(t *testing.T) {
		client := New(Config{
			APIKey: "test-key",
		})

		if client.apiKey != "test-key" {
			t.Errorf("APIKey = %v, want test-key", client.apiKey)
		}
		if client.baseURL != DefaultBaseURL {
			t.Errorf("BaseURL = %v, want %v", client.baseURL, DefaultBaseURL)
		}
		if client.httpClient == nil {
			t.Error("httpClient is nil")
		}
	})

	t.Run("with custom values", func(t *testing.T) {
		customClient := &http.Client{Timeout: 60 * time.Second}
		customURL := "https://custom.api.com"

		client := New(Config{
			APIKey:     "custom-key",
			BaseURL:    customURL,
			HTTPClient: customClient,
		})

		if client.baseURL != customURL {
			t.Errorf("BaseURL = %v, want %v", client.baseURL, customURL)
		}
		if client.httpClient != customClient {
			t.Error("Custom HTTP client not used")
		}
	})

	t.Run("with empty APIKey", func(t *testing.T) {
		client := New(Config{
			APIKey: "",
		})

		// Should not panic, but apiKey should be empty
		if client.apiKey != "" {
			t.Errorf("APIKey = %v, want empty string", client.apiKey)
		}
	})
}

func TestCreateDocument_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.Header.Get("Authorization") == "" {
			t.Error("Missing Authorization header")
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Content-Type = %v, want application/json", r.Header.Get("Content-Type"))
		}

		// Verify request body
		var req CreateDocumentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request: %v", err)
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(CreateDocumentResponse{
			ID:     "doc-123",
			Title:  req.Title,
			Status: DocumentStatusPending,
		})
	}))
	defer server.Close()

	client := New(Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	doc, err := client.CreateDocument(context.Background(), CreateDocumentRequest{
		Title:       "Test Document",
		FileContent: "base64content",
		Signatures: []Signature{
			NewEmailValidation("test@example.com", true),
		},
	})

	if err != nil {
		t.Fatalf("CreateDocument() error = %v", err)
	}

	if doc.ID != "doc-123" {
		t.Errorf("Document ID = %v, want doc-123", doc.ID)
	}
	if doc.Status != DocumentStatusPending {
		t.Errorf("Document Status = %v, want %v", doc.Status, DocumentStatusPending)
	}
}

func TestCreateDocument_ValidationErrors(t *testing.T) {
	client := New(Config{APIKey: "test-key"})

	tests := []struct {
		name    string
		req     CreateDocumentRequest
		wantErr string
	}{
		{
			name: "empty title",
			req: CreateDocumentRequest{
				Title:       "",
				FileContent: "base64",
				Signatures:  []Signature{NewEmailValidation("test@example.com", true)},
			},
			wantErr: "title cannot be empty",
		},
		{
			name: "empty fileContent",
			req: CreateDocumentRequest{
				Title:       "Test",
				FileContent: "",
				Signatures:  []Signature{NewEmailValidation("test@example.com", true)},
			},
			wantErr: "fileContent cannot be empty",
		},
		{
			name: "no signatures",
			req: CreateDocumentRequest{
				Title:       "Test",
				FileContent: "base64",
				Signatures:  []Signature{},
			},
			wantErr: "at least one signature is required",
		},
		{
			name: "signature without validation method",
			req: CreateDocumentRequest{
				Title:       "Test",
				FileContent: "base64",
				Signatures: []Signature{
					{
						SignerName:  "John",
						Validations: Validation{},
					},
				},
			},
			wantErr: "signature[0]: at least one validation method required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.CreateDocument(context.Background(), tt.req)
			if err == nil {
				t.Error("Expected error, got nil")
				return
			}
			if err.Error() != tt.wantErr {
				t.Errorf("Error = %v, want %v", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestCreateDocument_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Error{
			Code:    "INVALID_REQUEST",
			Message: "Invalid file format",
		})
	}))
	defer server.Close()

	client := New(Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	_, err := client.CreateDocument(context.Background(), CreateDocumentRequest{
		Title:       "Test",
		FileContent: "invalid-base64",
		Signatures:  []Signature{NewEmailValidation("test@example.com", true)},
	})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	apiErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("Expected *Error, got %T", err)
	}

	if apiErr.Code != "INVALID_REQUEST" {
		t.Errorf("Error code = %v, want INVALID_REQUEST", apiErr.Code)
	}
}

func TestCreateDocument_ContextTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := New(Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := client.CreateDocument(ctx, CreateDocumentRequest{
		Title:       "Test",
		FileContent: "base64",
		Signatures:  []Signature{NewEmailValidation("test@example.com", true)},
	})

	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
}

func TestGetDocument_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(GetDocumentResponse{
			ID:     "doc-123",
			Title:  "Test Document",
			Status: DocumentStatusPending,
		})
	}))
	defer server.Close()

	client := New(Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	doc, err := client.GetDocument(context.Background(), "doc-123")
	if err != nil {
		t.Fatalf("GetDocument() error = %v", err)
	}

	if doc.ID != "doc-123" {
		t.Errorf("Document ID = %v, want doc-123", doc.ID)
	}
}

func TestGetDocument_EmptyID(t *testing.T) {
	client := New(Config{APIKey: "test-key"})

	_, err := client.GetDocument(context.Background(), "")
	if err == nil {
		t.Error("Expected error for empty documentID")
	}
	if err.Error() != "documentID cannot be empty" {
		t.Errorf("Error = %v, want 'documentID cannot be empty'", err)
	}
}

func TestGetDocument_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(Error{
			Code:    "NOT_FOUND",
			Message: "Document not found",
		})
	}))
	defer server.Close()

	client := New(Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	_, err := client.GetDocument(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	apiErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("Expected *Error, got %T", err)
	}
	if apiErr.Code != "NOT_FOUND" {
		t.Errorf("Error code = %v, want NOT_FOUND", apiErr.Code)
	}
}

func TestListDocuments_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		// Verify query parameters
		status := r.URL.Query().Get("status")
		if status != DocumentStatusCompleted {
			t.Errorf("Status = %v, want %v", status, DocumentStatusCompleted)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ListDocumentsResponse{
			Documents: []GetDocumentResponse{
				{ID: "doc-1", Status: DocumentStatusCompleted},
				{ID: "doc-2", Status: DocumentStatusCompleted},
			},
			TotalCount: 2,
			Limit:      10,
			Offset:     0,
		})
	}))
	defer server.Close()

	client := New(Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	result, err := client.ListDocuments(context.Background(), ListDocumentsParams{
		Status: DocumentStatusCompleted,
		Limit:  10,
	})

	if err != nil {
		t.Fatalf("ListDocuments() error = %v", err)
	}

	if len(result.Documents) != 2 {
		t.Errorf("Documents count = %d, want 2", len(result.Documents))
	}
	if result.TotalCount != 2 {
		t.Errorf("TotalCount = %d, want 2", result.TotalCount)
	}
}

func TestCancelDocument_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := New(Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	err := client.CancelDocument(context.Background(), "doc-123")
	if err != nil {
		t.Errorf("CancelDocument() error = %v", err)
	}
}

func TestCancelDocument_EmptyID(t *testing.T) {
	client := New(Config{APIKey: "test-key"})

	err := client.CancelDocument(context.Background(), "")
	if err == nil {
		t.Error("Expected error for empty documentID")
	}
	if err.Error() != "documentID cannot be empty" {
		t.Errorf("Error = %v, want 'documentID cannot be empty'", err)
	}
}

func TestDownloadDocument_Success(t *testing.T) {
	expectedData := []byte("PDF content here")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
		w.Write(expectedData)
	}))
	defer server.Close()

	client := New(Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	data, err := client.DownloadDocument(context.Background(), "doc-123")
	if err != nil {
		t.Fatalf("DownloadDocument() error = %v", err)
	}

	if string(data) != string(expectedData) {
		t.Errorf("Downloaded data doesn't match expected")
	}
}

func TestDownloadDocument_EmptyID(t *testing.T) {
	client := New(Config{APIKey: "test-key"})

	_, err := client.DownloadDocument(context.Background(), "")
	if err == nil {
		t.Error("Expected error for empty documentID")
	}
	if err.Error() != "documentID cannot be empty" {
		t.Errorf("Error = %v, want 'documentID cannot be empty'", err)
	}
}

func TestDownloadDocument_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Document not found"))
	}))
	defer server.Close()

	client := New(Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	_, err := client.DownloadDocument(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestResendInvitation_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := New(Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	err := client.ResendInvitation(context.Background(), "doc-123", "sig-456")
	if err != nil {
		t.Errorf("ResendInvitation() error = %v", err)
	}
}

func TestResendInvitation_EmptyIDs(t *testing.T) {
	client := New(Config{APIKey: "test-key"})

	tests := []struct {
		name        string
		documentID  string
		signatureID string
		wantErr     string
	}{
		{
			name:        "empty documentID",
			documentID:  "",
			signatureID: "sig-123",
			wantErr:     "documentID cannot be empty",
		},
		{
			name:        "empty signatureID",
			documentID:  "doc-123",
			signatureID: "",
			wantErr:     "signatureID cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.ResendInvitation(context.Background(), tt.documentID, tt.signatureID)
			if err == nil {
				t.Error("Expected error, got nil")
				return
			}
			if err.Error() != tt.wantErr {
				t.Errorf("Error = %v, want %v", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestHandleResponse_NilResponse(t *testing.T) {
	err := handleResponse(nil, nil)
	if err == nil {
		t.Error("Expected error for nil response")
	}
	if err.Error() != "response is nil" {
		t.Errorf("Error = %v, want 'response is nil'", err)
	}
}

func TestHandleResponse_ErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Error{
			Code:    "VALIDATION_ERROR",
			Message: "Invalid input",
		})
	}))
	defer server.Close()

	resp, _ := http.Get(server.URL)

	var result interface{}
	err := handleResponse(resp, &result)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	apiErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("Expected *Error, got %T", err)
	}

	if apiErr.Code != "VALIDATION_ERROR" {
		t.Errorf("Error code = %v, want VALIDATION_ERROR", apiErr.Code)
	}
}

func TestHandleResponse_MalformedJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json{"))
	}))
	defer server.Close()

	resp, _ := http.Get(server.URL)

	var result GetDocumentResponse
	err := handleResponse(resp, &result)

	if err == nil {
		t.Error("Expected error for malformed JSON")
	}
}

func TestError_Error(t *testing.T) {
	err := &Error{
		Code:    "TEST_ERROR",
		Message: "This is a test error",
	}

	expected := "signatura: TEST_ERROR - This is a test error"
	if err.Error() != expected {
		t.Errorf("Error() = %v, want %v", err.Error(), expected)
	}
}
