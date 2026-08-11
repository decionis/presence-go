# API Reference

The complete public surface of package `presence`. The language-neutral contract is
[`openapi/presence-v1.yaml`](../../../openapi/presence-v1.yaml).

## Client

```go
import "github.com/decionis/presence-go/presence"

client := presence.NewClient("https://presence.decionis.com", apiKey)
```

```go
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

func NewClient(baseURL, apiKey string) *Client
```

Server-side Presence client, safe for concurrent use. `APIKey` must be a tenant-side
credential. Trailing slashes on `baseURL` are trimmed. The default `http.Client` carries
a 10-second timeout; set `HTTPClient` to inject your own transport or timeout. Every
method takes a `context.Context` first and honors cancellation.

## Methods

### `RunSandboxScenario(ctx, scenario, idempotencyKey) (*SandboxFixtureResult, error)`

`POST /v1/sandbox/scenarios/{scenario}` — run a deterministic, signed sandbox fixture.
`scenario` is `presence.ScenarioAllow`, `presence.ScenarioBlock`, or
`presence.ScenarioEscalate`; anything else returns an error before any network call.

### `CreateSession(ctx, action) (*CreatedSession, error)`

`POST /v1/sessions` — open a Session for one action. `action` is a
`presence.ActionContext`. Returns `SessionID`, `SessionToken`, `Nonce`, `CreatedAt`,
and `Status`.

### `Decide(ctx, request, idempotencyKey) (*DecideResponse, error)`

`POST /v1/decide` — submit a schema-valid decision request (`request any`) and receive
the Presence Result: `Verdict`, `ExecutionToken`, `PolicyEvaluation`, `Enforcement`,
`DecisionDossier`, and `VerificationBundle`.

### `GetDossier(ctx, dossierID) (map[string]any, error)`

`GET /v1/dossiers/{dossierID}` — retrieve one Presence Record (wire name: `dossier`).

### `VerifyDossier(ctx, dossierID) (*DossierVerification, error)`

`GET /v1/dossiers/{dossierID}/verify` — verify a Presence Record server-side. Returns
`DossierID`, `SignatureValid`, `ChainValid`, and `Valid`.

### `Health(ctx) error`

`GET /healthz` — unauthenticated reachability probe.

## Types

```go
type Verdict string            // VerdictAuthorized, VerdictRestrain, VerdictEscalate, VerdictBlocked
type SandboxScenario string    // ScenarioAllow, ScenarioBlock, ScenarioEscalate
```

`ActionContext` (json tags are the wire names): `Intent`, `Surface`, `ActorID`
(`actor_id`), `TargetResourceID` (`target_resource_id`), `ProviderEventID`, `Amount`,
`Currency`, `RecipientIBAN` (`recipient_iban`).

`SandboxFixtureResult`: `SandboxFixture`, `Scenario`, `PolicyVerdict`,
`EffectiveVerdict`, `Decision DecideResponse`.

`DecideResponse`: `Status`, `EvidenceClass`, `Verdict`, `ExecutionToken *string`,
`EvaluatedAt`, `LatencyMS *float64`, `PolicyEvaluation map[string]any`,
`EscalationContext map[string]any`, `Enforcement *EnforcementOutcome`,
`DecisionDossier DecisionDossier`, `VerificationBundle *VerificationBundle`.
`Enforcement` and `VerificationBundle` are nil-able pointers — guard them.

`EnforcementOutcome`: `Mode`, `EffectiveVerdict`, `TrueVerdict`, `Downgraded`,
`KillSwitchActive`.

`DecisionDossier` / `VerificationBundle`: `DossierID`, `ChainHash`,
`PreviousChainHash`, `Algorithm`, `PublicKeyID`, `Signature`; the bundle adds
`BundleVersion`, `Kind`, `SealedPayload json.RawMessage`, `PublicKeysURL`,
`VerifierURL`.

`DossierVerification`: `DossierID`, `SignatureValid`, `ChainValid`, `Valid`.

## Errors

```go
var apiErr *presence.APIError
if errors.As(err, &apiErr) {
	fmt.Println(apiErr.StatusCode, apiErr.Reason)
}
```

`APIError{StatusCode int, Reason string, Body string}` — a non-success response with its
stable machine-readable reason. Message: `presence API returned {status}: {reason}`
(lowercase per Go error convention). All other failures arrive as wrapped errors
(`call presence API: …`, `decode presence response: …`).

A `BLOCKED` Verdict is a result, not an error. See
[Troubleshooting](troubleshooting.md) for the reason-code taxonomy.

## CLI

```sh
go install github.com/decionis/presence-go/cmd/presence@v0.2.0
```

```text
presence quickstart [allow|block|escalate]
presence sandbox run <allow|block|escalate> [--idempotency-key KEY] [--json]
presence dossier verify <dossier-id>
presence doctor
presence version
```

Environment: `PRESENCE_API_KEY` (required for authenticated commands),
`PRESENCE_API_URL` (optional, defaults to `https://presence.decionis.com`). There is no
base-URL flag — the environment variable is the override. Errors print to stderr as
`presence: {message}` with exit code 1.

```text
Policy verdict: AUTHORIZED
Effective verdict: AUTHORIZED (sandbox is SHADOW)
Evidence class: sandbox_fixture
Dossier: dos_…
```

## Wire headers

| Header            | Value                                    |
| ----------------- | ---------------------------------------- |
| `Authorization`   | `Bearer presence_sk_…`                   |
| `Content-Type`    | `application/json` (when a body is sent) |
| `Idempotency-Key` | your stable retry key                    |
