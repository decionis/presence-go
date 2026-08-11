# treasury

## What it proves

A treasury wire read through Shadow Mode enforcement.

## Run it

```sh
export PRESENCE_API_KEY='presence_sk_…'
go run ./examples/treasury
```

## What you'll see

```text
Intent: corporate_treasury.wire_transfer.execute
Policy verdict: BLOCKED
Mode: SHADOW
Effective verdict: AUTHORIZED
True verdict: BLOCKED
Downgraded: true
Kill switch: false
Shadow Mode recorded the block; execution proceeded unchanged.
```

Shadow Mode observes every Presence Check and never blocks execution.
