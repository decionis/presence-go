# Shadow Mode

Shadow Mode observes every Presence Check and never blocks execution. Presence
evaluates, signs, and records the Verdict it would have enforced; your application
proceeds unchanged.

## Why it exists

You do not turn on enforcement for wire transfers on day one. Shadow Mode gives you the
full signed decision path — Policy evaluation, Verdicts, Presence Records — against real
traffic, while execution stays exactly as it is today.

## Reading enforcement

Shadow Mode is a tenant enforcement state, not an SDK flag. Every Presence Result
reports how enforcement treated it:

```go
enforcement := result.Decision.Enforcement

if enforcement != nil {
	fmt.Println("Mode:", enforcement.Mode)
	fmt.Println("Effective verdict:", enforcement.EffectiveVerdict)
	fmt.Println("True verdict:", enforcement.TrueVerdict)
	fmt.Println("Downgraded:", enforcement.Downgraded)
	fmt.Println("Kill switch:", enforcement.KillSwitchActive)
}
```

Three wire facts tell the story:

| Field               | Meaning                                   |
| ------------------- | ----------------------------------------- |
| `effective_verdict` | What was enforced.                        |
| `true_verdict`      | What Presence concluded.                  |
| `downgraded`        | Whether Shadow Mode softened the outcome. |

In Shadow Mode a `TrueVerdict` of `BLOCKED` arrives with an `EffectiveVerdict` of
`AUTHORIZED` and `Downgraded: true` — the block that would have happened, recorded and
signed, without happening. `KillSwitchActive` reports the tenant kill switch, which
forces observation regardless of mode.

## Enforcement modes

| Mode                  | Behavior                                        |
| --------------------- | ----------------------------------------------- |
| `SHADOW`              | Observe and record everything; enforce nothing. |
| `RESTRAIN_ONLY`       | Enforce only defensive holds.                   |
| `ENFORCE_HIGH_STAKES` | Enforce above the high-value floor.             |
| `ENFORCE`             | Enforce every Verdict.                          |

## The sandbox and Shadow Mode

Sandbox fixtures always run in Shadow Mode — you get the full signed decision path with
zero enforcement risk. That is why every quickstart line reads
`Effective verdict: AUTHORIZED (sandbox is SHADOW)`.

The `treasury` example reads all five enforcement fields from a `block` scenario:
[examples/treasury](../examples/treasury).

## Graduating to enforcement

Start in Shadow Mode, watch true verdicts against real traffic, then graduate surface by
surface. Enforcement changes are tenant rollout controls on the Presence control plane —
there is nothing to change in this SDK.
