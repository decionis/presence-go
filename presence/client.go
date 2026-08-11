// Package presence provides an idiomatic client for the Presence verification API.
package presence

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Verdict is the caller-facing result of a Presence decision.
type Verdict string

const (
	VerdictAuthorized Verdict = "AUTHORIZED"
	VerdictRestrain   Verdict = "RESTRAIN"
	VerdictEscalate   Verdict = "ESCALATE"
	VerdictBlocked    Verdict = "BLOCKED"
)

// SandboxScenario selects one deterministic developer fixture.
type SandboxScenario string

const (
	ScenarioAllow    SandboxScenario = "allow"
	ScenarioBlock    SandboxScenario = "block"
	ScenarioEscalate SandboxScenario = "escalate"
)

// ActionContext describes the consequential action being protected.
type ActionContext struct {
	Intent           string  `json:"intent"`
	Surface          string  `json:"surface"`
	ActorID          string  `json:"actor_id"`
	TargetResourceID string  `json:"target_resource_id"`
	ProviderEventID  string  `json:"provider_event_id,omitempty"`
	Amount           float64 `json:"amount,omitempty"`
	Currency         string  `json:"currency,omitempty"`
	RecipientIBAN    string  `json:"recipient_iban,omitempty"`
}

// CreatedSession contains the narrow credential that may be handed to a browser.
type CreatedSession struct {
	SessionID             string `json:"session_id"`
	SessionToken          string `json:"session_token"`
	SessionTokenExpiresIn int    `json:"session_token_expires_in"`
	Nonce                 string `json:"nonce"`
	CreatedAt             string `json:"created_at"`
	Status                string `json:"status"`
}

type DecisionDossier struct {
	DossierID         string `json:"dossier_id"`
	PreviousChainHash string `json:"previous_chain_hash"`
	ChainHash         string `json:"chain_hash"`
	Algorithm         string `json:"algorithm"`
	PublicKeyID       string `json:"public_key_id"`
	Signature         string `json:"signature"`
}

type EnforcementOutcome struct {
	Mode             string  `json:"mode"`
	EffectiveVerdict Verdict `json:"effective_verdict"`
	TrueVerdict      Verdict `json:"true_verdict"`
	Downgraded       bool    `json:"downgraded"`
	KillSwitchActive bool    `json:"kill_switch_active"`
}

type VerificationBundle struct {
	BundleVersion     string          `json:"bundle_version"`
	DossierID         string          `json:"dossier_id"`
	Kind              string          `json:"kind"`
	Algorithm         string          `json:"algorithm"`
	PublicKeyID       string          `json:"public_key_id"`
	PreviousChainHash string          `json:"previous_chain_hash"`
	ChainHash         string          `json:"chain_hash"`
	Signature         string          `json:"signature"`
	SealedPayload     json.RawMessage `json:"sealed_payload"`
	PublicKeysURL     string          `json:"public_keys_url"`
	VerifierURL       string          `json:"verifier_url"`
}

type DecideResponse struct {
	Status             string              `json:"status"`
	EvidenceClass      string              `json:"evidence_class,omitempty"`
	Verdict            Verdict             `json:"verdict"`
	ExecutionToken     *string             `json:"execution_token"`
	EvaluatedAt        string              `json:"evaluated_at"`
	LatencyMS          *float64            `json:"latency_ms,omitempty"`
	PolicyEvaluation   map[string]any      `json:"policy_evaluation"`
	EscalationContext  map[string]any      `json:"escalation_context,omitempty"`
	Enforcement        *EnforcementOutcome `json:"enforcement,omitempty"`
	DecisionDossier    DecisionDossier     `json:"decision_dossier"`
	VerificationBundle *VerificationBundle `json:"verification_bundle,omitempty"`
}

type SandboxFixtureResult struct {
	SandboxFixture   bool            `json:"sandbox_fixture"`
	Scenario         SandboxScenario `json:"scenario"`
	PolicyVerdict    Verdict         `json:"policy_verdict"`
	EffectiveVerdict Verdict         `json:"effective_verdict"`
	Decision         DecideResponse  `json:"decision"`
}

type DossierVerification struct {
	DossierID      string `json:"dossier_id"`
	SignatureValid bool   `json:"signature_valid"`
	ChainValid     bool   `json:"chain_valid"`
	Valid          bool   `json:"valid"`
}

// APIError preserves the status and machine-readable Presence reason.
type APIError struct {
	StatusCode int
	Reason     string
	Body       string
}

func (e *APIError) Error() string {
	if e.Reason != "" {
		return fmt.Sprintf("presence API returned %d: %s", e.StatusCode, e.Reason)
	}
	return fmt.Sprintf("presence API returned %d", e.StatusCode)
}

// Client is safe for concurrent use. APIKey must be a tenant-side credential.
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL: strings.TrimRight(baseURL, "/"),
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) RunSandboxScenario(ctx context.Context, scenario SandboxScenario, idempotencyKey string) (*SandboxFixtureResult, error) {
	if scenario != ScenarioAllow && scenario != ScenarioBlock && scenario != ScenarioEscalate {
		return nil, fmt.Errorf("invalid sandbox scenario %q", scenario)
	}
	var result SandboxFixtureResult
	if err := c.do(ctx, http.MethodPost, "/v1/sandbox/scenarios/"+string(scenario), nil, idempotencyKey, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) CreateSession(ctx context.Context, action ActionContext) (*CreatedSession, error) {
	var result CreatedSession
	if err := c.do(ctx, http.MethodPost, "/v1/sessions", action, "", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Decide sends a schema-valid request object. Generated models may be added
// without changing this stable transport boundary.
func (c *Client) Decide(ctx context.Context, request any, idempotencyKey string) (*DecideResponse, error) {
	var result DecideResponse
	if err := c.do(ctx, http.MethodPost, "/v1/decide", request, idempotencyKey, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetDossier(ctx context.Context, dossierID string) (map[string]any, error) {
	var result map[string]any
	if err := c.do(ctx, http.MethodGet, "/v1/dossiers/"+dossierID, nil, "", &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) VerifyDossier(ctx context.Context, dossierID string) (*DossierVerification, error) {
	var result DossierVerification
	if err := c.do(ctx, http.MethodGet, "/v1/dossiers/"+dossierID+"/verify", nil, "", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) Health(ctx context.Context) error {
	var result map[string]any
	return c.do(ctx, http.MethodGet, "/healthz", nil, "", &result)
}

func (c *Client) do(ctx context.Context, method, path string, body any, idempotencyKey string, result any) error {
	if c.BaseURL == "" {
		return errors.New("presence base URL is required")
	}
	if c.APIKey == "" && path != "/healthz" {
		return errors.New("presence API key is required")
	}

	var encoded io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("encode presence request: %w", err)
		}
		encoded = bytes.NewReader(payload)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+path, encoded)
	if err != nil {
		return fmt.Errorf("create presence request: %w", err)
	}
	if c.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.APIKey)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if idempotencyKey != "" {
		req.Header.Set("Idempotency-Key", idempotencyKey)
	}

	httpClient := c.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("call presence API: %w", err)
	}
	defer resp.Body.Close()
	payload, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return fmt.Errorf("read presence response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errorBody struct {
			Reason string `json:"reason"`
		}
		_ = json.Unmarshal(payload, &errorBody)
		return &APIError{StatusCode: resp.StatusCode, Reason: errorBody.Reason, Body: string(payload)}
	}
	if err := json.Unmarshal(payload, result); err != nil {
		return fmt.Errorf("decode presence response: %w", err)
	}
	return nil
}
