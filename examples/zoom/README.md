# zoom

## What it proves

An executive meeting instruction held for review.

## Run it

```sh
export PRESENCE_API_KEY='presence_sk_…'
go run ./examples/zoom
```

## What you'll see

```text
Intent: meeting.executive_instruction.confirm
Policy verdict: ESCALATE
Hold: instruction awaits human review.
```

Only `ESCALATE` enters human review. `AUTHORIZED` and `BLOCKED` are autonomous terminal
results.
