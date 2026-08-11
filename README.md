# Human Presence Verification for Go

**Presence SDK** — verify human presence before consequential actions execute.

This is the public, installable Go distribution of the Presence SDK and CLI.

> Applications send intent and context.
> Presence verifies the person.
> Decionis decides whether execution proceeds.

## Install

```sh
go get github.com/decionis/presence-go@v0.2.0
go install github.com/decionis/presence-go/cmd/presence@v0.2.0
```

```sh
export PRESENCE_API_KEY='presence_sk_…'
```

## Thirty seconds to a verdict

```go
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/decionis/presence-go/presence"
)

func main() {
	client := presence.NewClient("https://presence.decionis.com", os.Getenv("PRESENCE_API_KEY"))
	key := fmt.Sprintf("quickstart-%x", time.Now().UnixNano())
	result, _ := client.RunSandboxScenario(context.Background(), presence.ScenarioAllow, key)
	fmt.Println(result.PolicyVerdict)
}
```

The fixture is signed with `evidence_class: sandbox_fixture` and always runs in Shadow
Mode. Tenant credentials belong on trusted servers only.

## What happens

```text
Application

↓

Intent + Context

↓

Presence

↓

Biometrics

↓

Liveness

↓

Presence Record

↓

Decionis Verdict

↓

Continue
```

Every Presence Check produces a signed **Presence Record** — portable evidence of who
approved what, where, when, and on which device.

## Quickstart

```text
Install

↓

Configure

↓

Send Intent

↓

Receive Presence Result

↓

Verify Signature

↓

Done
```

Five minutes: [docs/quickstart.md](docs/quickstart.md).

## Supported scenarios

```text
Wire Transfer

Treasury

Password Reset

Executive Meeting

Vendor Change

AI Approval
```

## Examples

| Example                       | What it proves                                        |
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
- [presence.decionis.com](https://presence.decionis.com)

## License

Apache-2.0. Use of the hosted Presence service is governed separately.
