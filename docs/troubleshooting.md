# Troubleshooting

## The error taxonomy

A Verdict is never an error. `BLOCKED` returns normally; errors exist for transport,
authentication, contract, and server failure only.

| You see                                                      | Type        | What it means, what to do                                                            |
| ------------------------------------------------------------ | ----------- | ------------------------------------------------------------------------------------ |
| `presence API key is required`                               | plain error | No key reached the client. Export `PRESENCE_API_KEY` or set `APIKey`.                |
| `presence base URL is required`                              | plain error | The client was built with an empty `BaseURL`.                                        |
| `presence API returned 401: unauthorized`                    | `*APIError` | The key is wrong, revoked, or from another environment.                              |
| `presence API returned 409: sandbox_fixture_requires_shadow` | `*APIError` | The sandbox needs a non-production tenant in Shadow Mode.                            |
| `presence API returned 503: gate_unavailable`                | `*APIError` | The decision gate failed closed; the action remains held. Retry with the same key.   |
| `call presence API: …`                                       | wrapped     | DNS, TLS, connection, or context timeout before a response arrived.                  |
| `decode presence response: unexpected end of JSON input`     | wrapped     | The response was truncated at the 4 MiB safety cap, or an intermediary answered.     |
| `invalid sandbox scenario "…"`                               | plain error | Scenario must be `allow`, `block`, or `escalate` — returned before any network call. |

Inspect API failures with `errors.As`:

```go
var apiErr *presence.APIError
if errors.As(err, &apiErr) {
	log.Println(apiErr.StatusCode, apiErr.Reason)
}
```

Every `*APIError` carries `StatusCode`, the stable machine-readable `Reason`, and the
raw `Body` for support.

## Sandbox preconditions

`RunSandboxScenario` succeeds only when all of these hold:

- the deployment explicitly enables sandbox fixtures;
- the API uses the simulated gate;
- the credential belongs to a non-production tenant; and
- the tenant is in Shadow Mode, unless its kill switch is active.

Every fixture is labelled `evidence_class: sandbox_fixture` — by design it can never
pass as production evidence.

## Retries

Writes are retried only when they carry an idempotency key — retrying an unkeyed write
could seal a second Presence Record. Reads retry freely. Retry on 429 and 5xx; never on
other 4xx.

This SDK makes one attempt per call. When you retry `Decide` or `RunSandboxScenario`,
reuse the same idempotency key; the sealed result is returned instead of a duplicate.

## Timeouts

The default `http.Client` carries a 10-second timeout. Inject your own for slow links,
and bound calls with a context:

```go
client.HTTPClient = &http.Client{Timeout: 30 * time.Second}

ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
defer cancel()
```

## Diagnosing connectivity

```sh
presence doctor
```

prints `Presence API reachable at https://presence.decionis.com` when `/healthz` answers.
In code, `client.Health(ctx)` does the same probe without credentials.

## Still stuck

Open a [GitHub issue](https://github.com/decionis/presence-go/issues) with the `APIError`
fields (`StatusCode`, `Reason`) and the idempotency key you used — never include your
tenant key.
