# Installation

## Install Presence

```sh
go get github.com/decionis/presence-go@v0.2.0
```

The module path is `github.com/decionis/presence-go`; the import path is
`github.com/decionis/presence-go/presence`. The module depends on the standard
library only. The `presence` command-line tool is a separate binary in the same module:

```sh
go install github.com/decionis/presence-go/cmd/presence@v0.2.0
```

## Runtime baseline

| Requirement | Version                          |
| ----------- | -------------------------------- |
| Go          | 1.22 or later                    |
| Platform    | Any platform with TLS to the API |

## Credentials

```sh
export PRESENCE_API_KEY='presence_sk_…'
```

A tenant key (`presence_sk_…`) authorizes every server-side call. The optional
`PRESENCE_API_URL` overrides the default endpoint `https://presence.decionis.com` — set
it only when targeting a non-default sandbox deployment.

The library itself never reads the environment; pass the base URL and key into
`presence.NewClient`. The `presence` CLI reads both variables for you.

Tenant credentials belong on trusted servers only. Browser and mobile SDKs never accept
a tenant key.

## Version pinning

Every Presence SDK releases in lockstep — the same version number ships across every
language, currently `0.2.0`:

| Language | Package |
| --- | --- |
| Go | [`github.com/decionis/presence-go`](https://pkg.go.dev/github.com/decionis/presence-go) |
| Swift / iOS | [`decionis/presence-swift`](https://github.com/decionis/presence-swift) |
| Python | [`decionis-presence`](https://pypi.org/project/decionis-presence/) |
| Node.js | [`@decionis/presence-node`](https://www.npmjs.com/package/@decionis/presence-node) |
| Browser | [`@decionis/presence-web`](https://www.npmjs.com/package/@decionis/presence-web) · [widget](https://www.npmjs.com/package/@decionis/presence-widget) · [React](https://www.npmjs.com/package/@decionis/presence-react) · [Vue](https://www.npmjs.com/package/@decionis/presence-vue) |
| .NET | [`Decionis.Presence`](https://www.nuget.org/packages/Decionis.Presence/) · [CLI](https://www.nuget.org/packages/Decionis.Presence.Cli/) |
| Java | [`com.decionis.presence:presence-java`](https://central.sonatype.com/artifact/com.decionis.presence/presence-java) |
| Android | [`com.decionis.presence:presence-android`](https://central.sonatype.com/artifact/com.decionis.presence/presence-android) |

This public mirror publishes repository-root tags such as `v0.2.0`. During `0.x`, pin an
exact version and read the [changelog](../CHANGELOG.md) before upgrading:

```sh
go get github.com/decionis/presence-go@v0.2.0
```

`go.mod` records the exact version you install — avoid `@latest` during `0.x` so an
upgrade is always a deliberate edit.

## Next

Continue to the [Quickstart](quickstart.md) — five minutes to your first Verdict.
