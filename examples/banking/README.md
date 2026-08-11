# banking

## What it proves

A wire transfer with a verified Presence Record.

## Run it

```sh
export PRESENCE_API_KEY='presence_sk_…'
go run ./examples/banking
```

## What you'll see

```text
Intent: banking.wire_transfer.execute
Amount: 250000 USD to ACME
Policy verdict: AUTHORIZED
Dossier: dos_…
Signature valid: true
Chain valid: true
Valid: true
Execute: wire transfer may proceed.
```

Verifying the Presence Record is what turns "Presence said yes" into evidence you can
hold independently.
