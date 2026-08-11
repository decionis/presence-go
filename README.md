<!--
  Maintainers: this file is published VERBATIM as the README of
  github.com/decionis/presence-go by scripts/publish-public-sdk-mirrors.sh.
  Keep it public-facing: no monorepo paths, no docs/NN references, and only
  links that resolve from the public repository.
-->

# Human Presence Verification for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/decionis/presence-go.svg)](https://pkg.go.dev/github.com/decionis/presence-go)
[![License: Apache-2.0](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](LICENSE)

**Presence SDK** — verify human presence before consequential actions execute.

> Applications send intent and context.
> Presence verifies the person.
> Decionis decides whether execution proceeds.

This module ships two things:

- **`presence`** — the server-side Go client for the Presence verification API. The tenant credential stays on your servers.
- **`cmd/presence`** — a CLI for running sandbox scenarios, re-verifying Presence Records, and checking connectivity.

## Install

Requires Go 1.22+.

```sh
go get github.com/decionis/presence-go@v0.2.0
```

The CLI, if you want it:

```sh
go install github.com/decionis/presence-go/cmd/presence@v0.2.0
```

Set your tenant credential — server-side only, never in a client build:

```sh
export PRESENCE_API_KEY='presence_sk_…'
```

## Thirty seconds to a verdict

Run a deterministic sandbox fixture — no real subject involved, and never mistakable for production evidence:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/decionis/presence-go/presence"
)

func main() {
	client := presence.NewClient("https://presence.decionis.com", os.Getenv("PRESENCE_API_KEY"))

	key := fmt.Sprintf("quickstart-%x", time.Now().UnixNano())
	result, err := client.RunSandboxScenario(context.Background(), presence.ScenarioAllow, key)
	if err != nil {
		// Fail closed: no verdict means the action does not proceed.
		log.Fatal(err)
	}
	fmt.Println(result.EffectiveVerdict)
}
```

Every sandbox response is signed with `evidence_class: sandbox_fixture` and runs in Shadow Mode.

The same flow from the CLI:

```sh
presence quickstart                   # the allow scenario, end to end
presence sandbox run escalate         # deterministic fixtures: allow | block | escalate
presence dossier verify <dossier-id>  # re-verify a Presence Record's signature and chain
presence doctor                       # confirm the API is reachable
```

## How a check flows

```text
application → intent + context → Presence (biometrics · liveness) → signed Presence Record → Decionis verdict → continue
```

Every Presence Check produces a signed **Presence Record** — portable evidence of who approved what, where, when, and on which device. Verdicts are `AUTHORIZED`, `RESTRAIN`, `ESCALATE`, or `BLOCKED`; treat anything other than `AUTHORIZED` — including an error — as "do not proceed".

## Client surface

| Method                                              | What it does                                                                                         |
| --------------------------------------------------- | ---------------------------------------------------------------------------------------------------- |
| `RunSandboxScenario(ctx, scenario, idempotencyKey)` | Run a deterministic fixture (`allow`, `block`, `escalate`) — the fastest way to see a full decision. |
| `CreateSession(ctx, action)`                        | Open a verification session for one action; returns the narrow session token the browser gets.       |
| `Decide(ctx, request, idempotencyKey)`              | Evaluate an attested verification and seal the verdict into a signed Presence Record.                |
| `GetDossier(ctx, dossierID)`                        | Retrieve a sealed Presence Record (wire name: dossier).                                              |
| `VerifyDossier(ctx, dossierID)`                     | Re-verify a record's Ed25519 signature and chain integrity, independent of the issuance path.        |
| `Health(ctx)`                                       | Confirm the API is reachable.                                                                        |

Full reference: [pkg.go.dev/github.com/decionis/presence-go](https://pkg.go.dev/github.com/decionis/presence-go).

## Trust, then re-verify

A Presence Record is evidence precisely because anyone can check it again later:

```go
verification, err := client.VerifyDossier(context.Background(), dossierID)
if err != nil {
	log.Fatal(err)
}
fmt.Println(verification.SignatureValid, verification.ChainValid)
```

Records also verify offline against the published keys: [presence.decionis.com/.well-known/presence-dossier-jwks.json](https://presence.decionis.com/.well-known/presence-dossier-jwks.json).

## Examples

| Example                       | What it shows                                         |
| ----------------------------- | ----------------------------------------------------- |
| [simple](examples/simple)     | Your first Presence Check.                            |
| [banking](examples/banking)   | A wire transfer with a verified Presence Record.      |
| [zoom](examples/zoom)         | An executive meeting instruction held for review.     |
| [treasury](examples/treasury) | A treasury wire read through Shadow Mode enforcement. |

## Documentation

- [Installation](docs/installation.md)
- [Quickstart](docs/quickstart.md)
- [Concepts](docs/concepts.md)
- [Presence Checks](docs/presence-check.md)
- [Shadow Mode](docs/shadow-mode.md)
- [API Reference](docs/api.md)
- [Troubleshooting](docs/troubleshooting.md)

Product documentation lives at [presence.decionis.com](https://presence.decionis.com).

## Support

- [GitHub Issues](https://github.com/decionis/presence-go/issues)
- Security reports: [security.txt](https://presence.decionis.com/.well-known/security.txt)

## License

Apache-2.0. Use of the hosted Presence service is governed separately.
