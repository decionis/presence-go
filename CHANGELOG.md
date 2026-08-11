# Changelog

Presence SDK versions advance in lockstep across supported languages. This public Go
mirror uses the standard `v{version}` tag required by Go modules.

## 0.2.0 — 2026-08-11

### Added

- Initial public release of the Presence SDK line: Presence Checks, sandbox scenarios,
  Presence Record retrieval and verification, and Shadow Mode enforcement reporting.
- `presence.Client` with `RunSandboxScenario`, `CreateSession`, `Decide`, `GetDossier`,
  `VerifyDossier`, and `Health`; every method takes a `context.Context` first.
- `presence` CLI: `quickstart`, `sandbox run`, `dossier verify`, `doctor`, `version`.
- Zero-dependency module (standard library only), versioned with the repository-root tag
  `v0.2.0`.

### Compatibility

- Requires Go 1.22 or later. API contract `presence-v1`. During `0.x`, breaking changes
  are stated explicitly in this file; pin an exact version.
