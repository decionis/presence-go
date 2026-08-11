package presence

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRunSandboxScenario(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/sandbox/scenarios/allow" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer presence_sk_test" {
			t.Fatalf("unexpected authorization %q", got)
		}
		if got := r.Header.Get("Idempotency-Key"); got != "quickstart-1234" {
			t.Fatalf("unexpected idempotency key %q", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"sandbox_fixture":   true,
			"scenario":          "allow",
			"policy_verdict":    "AUTHORIZED",
			"effective_verdict": "AUTHORIZED",
			"decision": map[string]any{
				"status": "success", "evidence_class": "sandbox_fixture", "verdict": "AUTHORIZED",
				"execution_token": "tok", "evaluated_at": "2026-08-04T00:00:00Z",
				"policy_evaluation": map[string]any{},
				"decision_dossier":  map[string]any{"dossier_id": "dos_1"},
			},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, "presence_sk_test")
	result, err := client.RunSandboxScenario(context.Background(), ScenarioAllow, "quickstart-1234")
	if err != nil {
		t.Fatal(err)
	}
	if !result.SandboxFixture || result.PolicyVerdict != VerdictAuthorized {
		t.Fatalf("unexpected result %#v", result)
	}
}

func TestAPIErrorReason(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(`{"status":"error","reason":"sandbox_fixture_requires_shadow"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "presence_sk_test")
	_, err := client.RunSandboxScenario(context.Background(), ScenarioAllow, "quickstart-1234")
	apiError, ok := err.(*APIError)
	if !ok || apiError.Reason != "sandbox_fixture_requires_shadow" {
		t.Fatalf("unexpected error %#v", err)
	}
}
