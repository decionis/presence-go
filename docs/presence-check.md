# Presence Checks

A Presence Check is one verification of one person for one Intent. This page walks the
server-side surface: sessions, decisions, Verdict handling, and the Presence Record.

## Open a Session

A Session binds one Presence Check to one action. Create it with the Intent and Context:

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/decionis/presence-go/presence"
)

var wireTransfer = presence.ActionContext{
	Intent:           "banking.wire_transfer.execute",
	Surface:          "web_banking_portal",
	ActorID:          "usr_treasurer_014",
	TargetResourceID: "acc_operating_01",
	Amount:           250000,
	Currency:         "USD",
	RecipientIBAN:    "DE89370400440532013000",
}

func main() {
	client := presence.NewClient("https://presence.decionis.com", os.Getenv("PRESENCE_API_KEY"))

	session, err := client.CreateSession(context.Background(), wireTransfer)
	if err != nil {
		fmt.Fprintln(os.Stderr, "presence:", err)
		os.Exit(1)
	}

	fmt.Println(session.SessionID, session.Status)
}
```

The response carries a narrow `SessionToken` (`presence_st_…`) that is safe to hand to
the browser for the subject-side ceremony, and a `Nonce` that binds the hardware
challenge to this exact action. The tenant key never leaves your server.

## Submit a decision

After the subject-side ceremony attaches evidence to the session, submit the decision
request. The idempotency key is required — it is what makes sealing a Presence Record
safe to retry:

```go
result, err := client.Decide(context.Background(), request, "wire-"+transferID)
if err != nil {
	return err
}
verdict := result.Verdict
```

`request` is the schema-valid decision payload for the session — `Decide` accepts `any`
JSON-marshalable value (see the [API contract](../../../openapi/presence-v1.yaml)). The
hardware anchor is server-derived: Presence overwrites the WebAuthn block with what it
verified for that session, so a forged client block cannot pass.

## Handle the Verdict

```go
switch result.Verdict {
case "AUTHORIZED":
	execute()
case "ESCALATE":
	holdForReview()
case "BLOCKED", "RESTRAIN":
	reject()
}
```

Only `ESCALATE` enters human review. `AUTHORIZED` and `BLOCKED` are autonomous terminal
results, and a `BLOCKED` Verdict is a result, not an error — the SDK returns no error
for it.

## Retrieve the Presence Record

Every Presence Result carries the Presence Record (wire name: `dossier`) — the signed
evidence artifact behind the Verdict:

```go
bundle := result.VerificationBundle
record, err := client.GetDossier(context.Background(), bundle.DossierID)
```

`VerificationBundle` is a pointer and may be nil on deployments that do not attach
bundles — guard it before use. `record` is the decoded dossier as `map[string]any`.

## Verify the Presence Record

```go
verification, err := client.VerifyDossier(context.Background(), bundle.DossierID)
if err != nil {
	return err
}

fmt.Println(verification.SignatureValid)
fmt.Println(verification.ChainValid)
fmt.Println(verification.Valid)
```

Verification recomputes the Ed25519 signature and the hash-chain link server-side and
returns `signature_valid`, `chain_valid`, and `valid`. Independent verification uses the
public keys at `/.well-known/presence-dossier-jwks.json` and the bundle's
`public_keys_url`.

Verifying the Presence Record is what turns "Presence said yes" into evidence you can
hold independently.

## Sandbox scenarios

Before wiring the full ceremony, run the deterministic path end to end:

```go
key := fmt.Sprintf("quickstart-%x", time.Now().UnixNano())
result, err := client.RunSandboxScenario(context.Background(), presence.ScenarioAllow, key)
```

| Scenario   | Produces     | Use it to                       |
| ---------- | ------------ | ------------------------------- |
| `allow`    | `AUTHORIZED` | prove the integration works     |
| `block`    | `BLOCKED`    | exercise your rejection path    |
| `escalate` | `ESCALATE`   | exercise your human-review path |

Scenario names are lowercase inputs; Verdicts are uppercase outputs. Every fixture is
signed with `evidence_class: sandbox_fixture` and always runs in Shadow Mode. The
sandbox is available only on deployments that explicitly enable fixtures, with a
non-production tenant in Shadow Mode.
