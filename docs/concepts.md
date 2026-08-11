# Concepts

Authentication proves a credential. Presence proves the person is participating before
the action executes.

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

Presence produces signals; Decionis decides. Presence never authorizes an action — it
proves the human, then Decionis evaluates Policy and returns the Verdict.

## The eight nouns

| Term                | Meaning                                                              | In this SDK                                       |
| ------------------- | -------------------------------------------------------------------- | ------------------------------------------------- |
| **Intent**          | What the application is about to do, named precisely.                | `ActionContext.Intent`                            |
| **Context**         | The facts around the Intent: actor, surface, target, amount.         | the `ActionContext` struct                        |
| **Presence Check**  | One verification of one person for one Intent.                       | `CreateSession` + `Decide`, or a sandbox scenario |
| **Presence Result** | The answer to a Presence Check: a Verdict plus its evidence.         | the `DecideResponse` returned by `Decide`         |
| **Presence Record** | The signed, portable evidence artifact behind every Presence Result. | `DecisionDossier` + `VerificationBundle`          |
| **Policy**          | The deterministic floors the evidence is compared against.           | `PolicyEvaluation`, `policy_thresholds`           |
| **Verdict**         | The decision: `AUTHORIZED`, `BLOCKED`, `ESCALATE`, or `RESTRAIN`.    | `Verdict`, `PolicyVerdict`, `EffectiveVerdict`    |
| **Shadow Mode**     | Presence evaluates and records; execution proceeds unchanged.        | `Enforcement.Mode == "SHADOW"`                    |

## Verdicts

- `AUTHORIZED` — execute. Autonomous terminal result.
- `BLOCKED` — reject. Autonomous terminal result.
- `ESCALATE` — held for human review. Only `ESCALATE` enters review.
- `RESTRAIN` — defensive hold.

A `BLOCKED` Verdict is a result, not an error. The SDK never returns an error for a
Verdict.

## The Presence Record

A Presence Record is the signed, portable evidence artifact produced by every Presence
Check.

```text
Presence Record

Who

What

Where

When

Why

Device

Challenge

Signature
```

In prose it is the Presence Record; on the wire it is the `dossier`. See
[Presence Checks](presence-check.md) for retrieval and verification.

## Credentials and tiers

- **Tenant Key** (`presence_sk_…`) — server-side secret. This SDK holds one.
- **Session Token** (`presence_st_…`) — narrow, short-lived, safe for browser and mobile.
- **Invitation** — one-time URL that hands a Presence Check to the subject's device.

Your server runs Presence Checks with this SDK. The subject-side ceremony — camera,
passkey, hosted flow — runs in the Browser, iOS, and Android SDKs, which never accept a
tenant key.

## Sandbox

Deterministic scenarios (`allow`, `block`, `escalate`) exercise the normal policy,
signing, and Presence Record path. Every sandbox response is permanently labelled
`evidence_class: sandbox_fixture` — it can never masquerade as production evidence.
Sandbox fixtures always run in Shadow Mode.
