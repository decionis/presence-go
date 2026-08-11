# Quickstart

Five minutes from install to a signed Verdict.

```text
1. Install Presence

↓

2. Create a Client

↓

3. Define Intent

↓

4. Start Verification

↓

5. Receive Presence Result

↓

6. Execute or Reject
```

## 1. Install Presence

```sh
go get github.com/decionis/presence-go@v0.2.0
export PRESENCE_API_KEY='presence_sk_…'
```

Your sandbox key is shown once during self-serve onboarding at
[presence.decionis.com](https://presence.decionis.com).

## 2. Create a Client

```go
import (
	"os"

	"github.com/decionis/presence-go/presence"
)

baseURL := os.Getenv("PRESENCE_API_URL")
if baseURL == "" {
	baseURL = "https://presence.decionis.com"
}
client := presence.NewClient(baseURL, os.Getenv("PRESENCE_API_KEY"))
```

The client has no default base URL — pass `https://presence.decionis.com` explicitly.
`NewClient` strips trailing slashes and installs a 10-second `HTTPClient` timeout;
replace the exported `HTTPClient` field to change it. The library never reads the
environment — the fallback above is your code.

## 3. Define Intent

An **Intent** names what your application is about to do. **Context** carries the facts
around it. Together they are the `ActionContext` (wire name: `action_context`):

```go
var wireTransfer = presence.ActionContext{
	Intent:           "banking.wire_transfer.execute",
	Surface:          "web_banking_portal",
	ActorID:          "usr_treasurer_014",
	TargetResourceID: "acc_operating_01",
	Amount:           250000,
	Currency:         "USD",
	RecipientIBAN:    "DE89370400440532013000",
}
```

## 4. Start Verification

Run your first Presence Check against the sandbox. Scenarios are deterministic inputs to
the normal policy, signing, and Presence Record path:

```go
key := fmt.Sprintf("quickstart-%x", time.Now().UnixNano())
result, err := client.RunSandboxScenario(context.Background(), presence.ScenarioAllow, key)
if err != nil {
	fmt.Fprintln(os.Stderr, "presence:", err)
	os.Exit(1)
}
```

The idempotency key makes the check safe to retry — reuse the same key and you get the
same sealed result, never a second one. Production ceremonies attach the subject's
biometrics and hardware evidence to the same flow; see
[Presence Checks](presence-check.md).

## 5. Receive Presence Result

```go
decision := result.Decision
fmt.Println("Policy verdict:", result.PolicyVerdict)
fmt.Println("Effective verdict:", result.EffectiveVerdict, "(sandbox is SHADOW)")
fmt.Println("Evidence class:", decision.EvidenceClass)
fmt.Println("Dossier:", decision.VerificationBundle.DossierID)
```

```text
Policy verdict: AUTHORIZED
Effective verdict: AUTHORIZED (sandbox is SHADOW)
Evidence class: sandbox_fixture
Dossier: dos_…
```

The scenario name is lowercase (`allow`); the Verdict it produces is uppercase
(`AUTHORIZED`).

## 6. Execute or Reject

Branch on the Verdict, and verify the Presence Record (wire name: `dossier`) before
trusting it:

```go
verification, err := client.VerifyDossier(context.Background(), decision.VerificationBundle.DossierID)
if err != nil {
	fmt.Fprintln(os.Stderr, "presence:", err)
	os.Exit(1)
}

if result.PolicyVerdict == "AUTHORIZED" && verification.Valid {
	executeWireTransfer()
} else {
	rejectWireTransfer()
}
```

Verifying the Presence Record is what turns "Presence said yes" into evidence you can
hold independently.

Done. Sandbox fixtures always run in Shadow Mode — you get the full signed decision path
with zero enforcement risk. Continue to [Presence Checks](presence-check.md) for the
production-shaped ceremony.
