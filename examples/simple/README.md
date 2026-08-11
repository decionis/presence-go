# simple

## What it proves

Your first Presence Check.

## Run it

```sh
export PRESENCE_API_KEY='presence_sk_…'
go run ./examples/simple
```

## What you'll see

```text
Policy verdict: AUTHORIZED
Effective verdict: AUTHORIZED (sandbox is SHADOW)
Evidence class: sandbox_fixture
Dossier: dos_…
```

The fixture is signed with `evidence_class: sandbox_fixture` and always runs in Shadow
Mode.
