// Package signatura provides a comprehensive Go client for the Signatura API
// that enables electronic signatures with identity validation.
package signatura

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	// DefaultBaseURL is the production API endpoint
	DefaultBaseURL = "https://connect.signatura.co/api/v2"

	// DefaultTimeout for HTTP requests
	DefaultTimeout = 30 * time.Second

	// ValidationTypeEmail represents email validation
	ValidationTypeEmail = "EM"

	// ValidationTypePhone represents phone validation
	ValidationTypePhone = "PH"

	// ValidationTypeBiometric represents biometric validation
	ValidationTypeBiometric = "BI"

	// ValidationTypeAFIP represents AFIP key validation
	ValidationTypeAFIP = "AF"

	// InviteChannelEmail represents email invitation channel
	InviteChannelEmail = "EM"

	// InviteChannelSMS represents SMS invitation channel
	InviteChannelSMS = "SM"

	// DocumentStatusPending represents a pending document
	DocumentStatusPending = "PE"

	// DocumentStatusCompleted represents a completed document
	DocumentStatusCompleted = "CO"

	// DocumentStatusCanceled represents a canceled document
	DocumentStatusCanceled = "CA"

	// WebhookActionDocumentSigned represents a signed document webhook
	WebhookActionDocumentSigned = "DS"

	// WebhookActionSignatureDeclined represents a declined signature webhook
	WebhookActionSignatureDeclined = "SD"

	// WebhookActionDocumentChange represents a document status change webhook
	WebhookActionDocumentChange = "DC"
)

// Client represents the Signatura API client
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// Config holds configuration for the Signatura client
type Config struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
	Timeout    time.Duration
}

// New creates a new Signatura API client with the provided configuration
func New(config Config) *Client {
	if config.BaseURL == "" {
		config.BaseURL = DefaultBaseURL
	}

	httpClient := config.HTTPClient
	if httpClient == nil {
		timeout := config.Timeout
		if timeout == 0 {
			timeout = DefaultTimeout
		}
		httpClient = &http.Client{
			Timeout: timeout,
		}
	}

	return &Client{
		baseURL:    config.BaseURL,
		apiKey:     config.APIKey,
		httpClient: httpClient,
	}
}

// Validation represents identity validation requirements for a signature
type Validation struct {
	// Email validation (EM) - email address for validation
	Email *string `json:"EM,omitempty"`

	// Phone validation (PH) - phone number for validation
	Phone *string `json:"PH,omitempty"`

	// Biometric validation (BI) - set to true to enable biometric validation
	Biometric *bool `json:"BI,omitempty"`

	// AFIP key validation (AF) - AFIP CUIT/CUIL for validation
	AFIP *string `json:"AF,omitempty"`
}

// Signature represents a signer in a document
type Signature struct {
	// Validations defines the identity verification methods required
	Validations Validation `json:"validations"`

	// InviteChannel specifies how the signer will be invited (optional)
	// Can be "EM" for email or "SM" for SMS
	InviteChannel []string `json:"invite_channel,omitempty"`

	// SignerName is the name of the signer (optional)
	SignerName string `json:"signer_name,omitempty"`
}

// CreateDocumentRequest represents the request to create a new document
type CreateDocumentRequest struct {
	// Title of the document (visible to signers)
	Title string `json:"title"`

	// FileContent is the base64-encoded PDF content (RFC 4648)
	FileContent string `json:"file_content"`

	// Signatures defines the list of signers and their validation requirements
	Signatures []Signature `json:"signatures"`

	// Metadata is optional custom data associated with the document
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// ExpiresAt is the optional expiration date for the document
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// SignatureResponse represents a signature in the API response
type SignatureResponse struct {
	ID            string                 `json:"id"`
	Status        string                 `json:"status"`
	SignerName    string                 `json:"signer_name,omitempty"`
	SigningURL    string                 `json:"signing_url,omitempty"`
	Validations   Validation             `json:"validations"`
	InviteChannel []string               `json:"invite_channel,omitempty"`
	SignedAt      *time.Time             `json:"signed_at,omitempty"`
	DeclinedAt    *time.Time             `json:"declined_at,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// CreateDocumentResponse represents the response from creating a document
type CreateDocumentResponse struct {
	ID         string                 `json:"id"`
	Title      string                 `json:"title"`
	Status     string                 `json:"status"`
	Signatures []SignatureResponse    `json:"signatures"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
	ExpiresAt  *time.Time             `json:"expires_at,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// GetDocumentResponse represents the response from getting a document
type GetDocumentResponse struct {
	ID            string                 `json:"id"`
	Title         string                 `json:"title"`
	Status        string                 `json:"status"`
	Signatures    []SignatureResponse    `json:"signatures"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
	CompletedAt   *time.Time             `json:"completed_at,omitempty"`
	CanceledAt    *time.Time             `json:"canceled_at,omitempty"`
	ExpiresAt     *time.Time             `json:"expires_at,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	DownloadURL   string                 `json:"download_url,omitempty"`
	AuditTrailURL string                 `json:"audit_trail_url,omitempty"`
}

// ListDocumentsParams represents query parameters for listing documents
type ListDocumentsParams struct {
	// Status filter (PE, CO, CA)
	Status string

	// Limit the number of results (default 20, max 100)
	Limit int

	// Offset for pagination
	Offset int

	// CreatedAfter filters documents created after this date
	CreatedAfter *time.Time

	// CreatedBefore filters documents created before this date
	CreatedBefore *time.Time
}

// ListDocumentsResponse represents the response from listing documents
type ListDocumentsResponse struct {
	Documents  []GetDocumentResponse `json:"documents"`
	TotalCount int                   `json:"total_count"`
	Limit      int                   `json:"limit"`
	Offset     int                   `json:"offset"`
}

// WebhookEvent represents a webhook notification from Signatura
type WebhookEvent struct {
	// NotificationAction indicates the type of event (DS, SD, DC)
	NotificationAction string `json:"notification_action"`

	// DocumentID is the ID of the document
	DocumentID string `json:"document_id"`

	// SignatureID is the ID of the signature (for DS and SD events)
	SignatureID string `json:"signature_id,omitempty"`

	// NewStatus is the new document status (for DC events)
	NewStatus string `json:"new_status,omitempty"`

	// Timestamp of the event
	Timestamp time.Time `json:"timestamp"`

	// Metadata associated with the event
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// Error represents an API error response
type Error struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("signatura: %s - %s", e.Code, e.Message)
}

// doRequest performs an HTTP request to the Signatura API
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader

	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonData)
	}

	reqURL := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, method, reqURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "signatura-go-client/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

// handleResponse processes the API response and handles errors
func handleResponse(resp *http.Response, result interface{}) error {
	if resp == nil {
		return fmt.Errorf("response is nil")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		var apiErr Error
		if err := json.Unmarshal(body, &apiErr); err != nil {
			return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
		}
		return &apiErr
	}

	if result != nil && len(body) > 0 {
		if err := json.Unmarshal(body, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// CreateDocument creates a new document for signing
func (c *Client) CreateDocument(ctx context.Context, req CreateDocumentRequest) (*CreateDocumentResponse, error) {
	// Validate required fields
	if req.Title == "" {
		return nil, fmt.Errorf("title cannot be empty")
	}

	if req.FileContent == "" {
		return nil, fmt.Errorf("fileContent cannot be empty")
	}

	if len(req.Signatures) == 0 {
		return nil, fmt.Errorf("at least one signature is required")
	}

	// Validate each signature
	for i, sig := range req.Signatures {
		if sig.SignerName == "" && sig.Validations.Email == nil && sig.Validations.Phone == nil {
			return nil, fmt.Errorf("signature[%d]: signerName should be provided for email/phone validations", i)
		}

		// Validate that at least one validation method is specified
		if sig.Validations.Email == nil && sig.Validations.Phone == nil &&
			sig.Validations.Biometric == nil && sig.Validations.AFIP == nil {
			return nil, fmt.Errorf("signature[%d]: at least one validation method required", i)
		}
	}

	resp, err := c.doRequest(ctx, http.MethodPost, "/documents/create", req)
	if err != nil {
		return nil, err
	}

	var result CreateDocumentResponse
	if err := handleResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetDocument retrieves a document by ID
func (c *Client) GetDocument(ctx context.Context, documentID string) (*GetDocumentResponse, error) {
	if documentID == "" {
		return nil, fmt.Errorf("documentID cannot be empty")
	}

	path := fmt.Sprintf("/documents/%s", documentID)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var result GetDocumentResponse
	if err := handleResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ListDocuments retrieves a list of documents with optional filters
func (c *Client) ListDocuments(ctx context.Context, params ListDocumentsParams) (*ListDocumentsResponse, error) {
	query := url.Values{}

	if params.Status != "" {
		query.Set("status", params.Status)
	}

	if params.Limit > 0 {
		query.Set("limit", fmt.Sprintf("%d", params.Limit))
	}

	if params.Offset > 0 {
		query.Set("offset", fmt.Sprintf("%d", params.Offset))
	}

	if params.CreatedAfter != nil {
		query.Set("created_after", params.CreatedAfter.Format(time.RFC3339))
	}

	if params.CreatedBefore != nil {
		query.Set("created_before", params.CreatedBefore.Format(time.RFC3339))
	}

	path := "/documents"
	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var result ListDocumentsResponse
	if err := handleResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CancelDocument cancels a pending document
func (c *Client) CancelDocument(ctx context.Context, documentID string) error {
	if documentID == "" {
		return fmt.Errorf("documentID cannot be empty")
	}

	path := fmt.Sprintf("/documents/%s/cancel", documentID)
	resp, err := c.doRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return err
	}

	return handleResponse(resp, nil)
}

// DownloadDocument downloads the signed document PDF
func (c *Client) DownloadDocument(ctx context.Context, documentID string) ([]byte, error) {
	if documentID == "" {
		return nil, fmt.Errorf("documentID cannot be empty")
	}

	path := fmt.Sprintf("/documents/%s/download", documentID)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("download failed with status %d: %s", resp.StatusCode, string(body))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read document data: %w", err)
	}

	return data, nil
}

// ResendInvitation resends the invitation to a signer
func (c *Client) ResendInvitation(ctx context.Context, documentID, signatureID string) error {
	if documentID == "" {
		return fmt.Errorf("documentID cannot be empty")
	}

	if signatureID == "" {
		return fmt.Errorf("signatureID cannot be empty")
	}

	path := fmt.Sprintf("/documents/%s/signatures/%s/resend", documentID, signatureID)
	resp, err := c.doRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return err
	}

	return handleResponse(resp, nil)
}
