package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	APIKeyEnvironmentVariableName    = "TRUORA_API_KEY"
	APIServerEnvironmentVariableName = "TRUORA_API_SERVER"
)

var (
	ErrAPIKeyNotProvided    = fmt.Errorf("API key not provided")
	ErrAPIServerNotProvided = fmt.Errorf("API server not provided")
)

type IdentityFlowConfig struct {
	// RedirectUrls               *ConfigURL                  `json:"redirect_urls,omitempty"`
	// Messages                   *Messages                   `json:"messages,omitempty"`
	// ContinueFlowInNewDevice    bool                        `json:"continue_flow_in_new_device,omitempty"`
	EnableDesktopFlow bool `json:"enable_desktop_flow,omitempty"`
	// EnablePostponedWebProcess  bool                        `json:"enable_postponed_web_process,omitempty"`
	// IgnoreInitialMessage       bool                        `json:"ignore_initial_message,omitempty"`
	// TimeToLive                 int64                       `json:"time_to_live,omitempty"`
	Lang string `json:"lang,omitempty"`
	// EnableFollowUp             bool                        `json:"enable_follow_up,omitempty"`
	// FollowUpDelay              int64                       `json:"follow_up_delay,omitempty"`
	// FollowUpMessage            string                      `json:"follow_up_message,omitempty"`
	// StartBusinessHours         time.Time                   `json:"start_business_hours,omitempty"`
	// EndBusinessHours           time.Time                   `json:"end_business_hours,omitempty"`
}

type ResponseOption struct {
	Alias string `json:"alias,omitempty"`
	Value string `json:"value"`
}

type Input struct {
	Type string `json:"type"`
	Name string `json:"name"`
	// Placeholder      string                    `json:"placeholder"`
	// Description      string                    `json:"description"`
	// Options          []string                  `json:"options,omitempty"`
	// Length           int                       `json:"length"`
	// ReadOnly         bool                      `json:"read_only"`
	// MediaURL         string                    `json:"media_url,omitempty"`
	// MediaName        string                    `json:"media_name,omitempty"`
	// MessageMediaType media.MessageMediaType    `json:"message_media_type,omitempty"`
	// MediaType        string                    `json:"media_type,omitempty"`
	// FileUploadURL    string                    `json:"file_upload_url,omitempty"`
	ResponseOptions []*ResponseOption `json:"response_options,omitempty"`
	// MediaID          string                    `json:"media_id,omitempty"`
	// ProductCatalog   *whatsapp.ProductsCatalog `json:"product_catalog,omitempty"`
	// Product          *whatsapp.Product         `json:"product,omitempty"`
	// Contacts         []whatsapp.Contact        `json:"contacts,omitempty"`
	// Footer           *whatsapp.Footer          `json:"footer,omitempty"`
}

type Step struct {
	StepID      string `json:"step_id"`
	Type        string `json:"type"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	// Logo             string                    `json:"logo,omitempty"`
	// Config           *VerificationConfigValues `json:"config"`
	ExpectedInputs []*Input `json:"expected_inputs,omitempty"`
	// AsyncStep      *bool    `json:"async_step"`
	// VerificationID   string                    `json:"verification_id,omitempty"`
	// StartDate  *time.Time `json:"start_date,omitempty"`
	// FinishDate *time.Time `json:"finish_date,omitempty"`
}

type IdentityVerification struct {
	VerificationID string                 `json:"verification_id"`
	Name           string                 `json:"name"`
	Config         map[string]interface{} `json:"config,omitempty"`
	Steps          []*Step                `json:"steps,omitempty"`
	Logic          []string               `json:"if,omitempty"`
}

type IdentityProcessFlow struct {
	FlowID   string `json:"flow_id"`
	ClientID string `json:"client_id"`
	// RefFlowID string `json:"ref_flow_id,omitempty"`
	// HasDraft bool `json:"has_draft,omitempty"`
	// Username         string              `json:"username,omitempty"`
	Version          int64               `json:"version,omitempty"`
	Name             string              `json:"name"`
	Status           string              `json:"status"`
	Type             string              `json:"type,omitempty"`
	CreationDate     *time.Time          `json:"creation_date,omitempty"`
	UpdateDate       *time.Time          `json:"update_date,omitempty"`
	VersionStartDate *time.Time          `json:"version_start_date,omitempty"`
	VersionEndDate   *time.Time          `json:"version_end_date,omitempty"`
	Config           *IdentityFlowConfig `json:"config,omitempty"`
	// CreationSource        FlowCreationSource      `json:"creation_source,omitempty"`
	IdentityVerifications []*IdentityVerification `json:"identity_verifications"`
}

type IdentityProcessFlowResponse struct {
	FlowID   string `json:"flow_id"`
	ClientID string `json:"client_id"`
	// RefFlowID string `json:"ref_flow_id,omitempty"`
	// HasDraft bool `json:"has_draft,omitempty"`
	// Username         string              `json:"username,omitempty"`
	Version          int64               `json:"version,omitempty"`
	Name             string              `json:"name"`
	Status           string              `json:"status"`
	Type             string              `json:"type,omitempty"`
	CreationDate     *time.Time          `json:"creation_date,omitempty"`
	UpdateDate       *time.Time          `json:"update_date,omitempty"`
	VersionStartDate *time.Time          `json:"version_start_date,omitempty"`
	VersionEndDate   *time.Time          `json:"version_end_date,omitempty"`
	Config           *IdentityFlowConfig `json:"config,omitempty"`
	// CreationSource        FlowCreationSource      `json:"creation_source,omitempty"`
	IdentityVerifications []*IdentityVerification `json:"identity_verifications"`
}

type TruoraClientOption func(*TruoraClient)

type TruoraClient struct {
	APIKey     string
	APIServer  string
	HTTPClient *http.Client
}

func WithAPIKey(apiKey string) TruoraClientOption {
	return func(client *TruoraClient) {
		client.APIKey = apiKey
	}

}

func WithAPIServer(apiServer string) TruoraClientOption {
	return func(client *TruoraClient) {
		client.APIServer = apiServer
	}
}

func NewClient(opts ...TruoraClientOption) (*TruoraClient, error) {
	client := &TruoraClient{
		APIKey:     os.Getenv(APIKeyEnvironmentVariableName),
		APIServer:  os.Getenv(APIServerEnvironmentVariableName),
		HTTPClient: &http.Client{},
	}

	for _, o := range opts {
		o(client)
	}

	if client.APIKey == "" {
		return nil, ErrAPIKeyNotProvided
	}

	if client.APIServer == "" {
		return nil, ErrAPIServerNotProvided
	}

	return client, nil
}

func (c *TruoraClient) GetFlow(flowID string) (*IdentityProcessFlowResponse, error) {
	client := c.HTTPClient

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/flows/%s", c.APIServer, flowID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Truora-Api-Key", c.APIKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var flow IdentityProcessFlowResponse
	if err := json.NewDecoder(resp.Body).Decode(&flow); err != nil {
		return nil, err
	}

	return &flow, nil
}

func (c *TruoraClient) CreateFlow(ctx context.Context, flow *IdentityProcessFlow) (*IdentityProcessFlowResponse, error) {
	client := c.HTTPClient

	marshalledFlow, err := json.Marshal(flow)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewBuffer(marshalledFlow)

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/v1/flows", c.APIServer), reader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Truora-Api-Key", c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		stringBody := new(bytes.Buffer)
		stringBody.ReadFrom(resp.Body)
		return nil, fmt.Errorf("error creating flow: %s\n%s", resp.Status, stringBody)
	}

	defer resp.Body.Close()

	var flowResponse IdentityProcessFlowResponse
	if err := json.NewDecoder(resp.Body).Decode(&flowResponse); err != nil {
		return nil, err
	}

	return &flowResponse, nil
}

func (c *TruoraClient) UpdateFlow(flowID string, flow *IdentityProcessFlow) (*IdentityProcessFlowResponse, error) {
	client := c.HTTPClient

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/flows/%s", c.APIServer, flowID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Truora-Api-Key", c.APIKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		stringBody := new(bytes.Buffer)
		stringBody.ReadFrom(resp.Body)
		return nil, fmt.Errorf("error updating flow: %s\n%s", resp.Status, stringBody)
	}

	defer resp.Body.Close()

	var flowResponse IdentityProcessFlowResponse
	if err := json.NewDecoder(resp.Body).Decode(&flowResponse); err != nil {
		return nil, err
	}

	return &flowResponse, nil
}

func (c *TruoraClient) DeleteFlow(flowID string) error {
	client := c.HTTPClient

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v1/flows/%s", c.APIServer, flowID), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Truora-Api-Key", c.APIKey)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		stringBody := new(bytes.Buffer)
		stringBody.ReadFrom(resp.Body)
		return fmt.Errorf("error deleting flow: %s\n%s", resp.Status, stringBody)
	}

	defer resp.Body.Close()

	return nil
}
