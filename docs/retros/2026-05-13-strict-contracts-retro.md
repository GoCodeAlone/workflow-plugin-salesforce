# Retro: Strict Contracts (Proto + ContractRegistry + CI)

**PR:** #6 — feat: add strict-contracts (proto + ContractRegistry + CI)
**Merged:** 2026-05-13
**Branch:** feat/issue-5-strict-contracts
**Design:** docs/plans/2026-05-13-strict-contracts-design.md
**Plan:** docs/plans/2026-05-13-strict-contracts.md
**Related ADRs:** none

## Adversarial-review findings, scored

### Cycle 1 (original design — 12 shared input message types)

| Phase | Finding | Severity | Outcome |
|---|---|---|---|
| design | Field-name drift: 12 shared message type design would silently diverge from plugin.json step names (no 1:1 mapping enforced) | Critical | Prescient — design was revised to Option 1 (3 messages), eliminating the entire drift class |
| design | Missing rollback story for proto schema changes (breaking change if message types renamed) | Critical | Resolved upfront — Option 1 has near-zero rollback surface (3 messages never change) |
| design | `SalesforceQueryInput.soql` naming inconsistent with existing step param key `query` | Critical | Resolved upfront — Option 1 uses Struct for all params, no per-field naming to drift |
| design | `SalesforceRecordCreateInput` missing `external_id_field` required for upsert steps | Critical | Resolved upfront — Option 1: all params in Struct, no missing fields possible |
| design | 12 shared message types still requires per-step field docs duplicated 72× | Important | Resolved upfront — Option 1 single `SalesforceStepInput{Struct params}` needs no per-step duplication |
| design | `SalesforceMetadataInput` too broad — covers both read and write metadata operations with incompatible required fields | Important | Resolved upfront — Option 1 makes this moot |

### Cycle 2 (revised design — Option 1, 3 messages)

| Phase | Finding | Severity | Outcome |
|---|---|---|---|
| design | Runtime validation discipline: without per-field proto types, missing required params produce runtime errors not proto errors | Important | Resolved upfront — design explicitly documents this trade-off and cites worldsim precedent |
| design | api_version proto comment says v58.0 but code defaults to v63.0 | Minor | Prescient — Copilot caught this in code review; proto comment and pb.go updated before merge |
| plan | make proto-gen invocation path concern (source_relative + proto_path combo) | Minor | False positive — verified stable: output correctly lands at gen/salesforce.pb.go |

## Gate misses

| Issue | Gate that missed | Why it slipped | Fix idea (optional) |
|---|---|---|---|
| `api_version` default comment `v58.0` vs actual `v63.0` in client.go | adversarial-design-review (design, cycle 2) | Design doc proto snippet was written before client.go was checked for actual default value; no explicit "verify constants match code" checklist item | Add "verify hardcoded defaults match code" to design-phase adversarial checklist for proto comment fields |

## Missed skill activations

| Gate | Fired? | Notes |
|---|---|---|
| brainstorming | yes | Session 1 (cycle 1 adversarial review context) |
| adversarial-design-review (design) | yes | Twice: cycle 1 (FAIL, 6 Critical) + cycle 2 after Option 1 revision (PASS) |
| writing-plans | yes | |
| adversarial-design-review (plan) | yes | PASS |
| alignment-check | yes | PASS |
| scope-lock | yes | Locked 2026-05-13T00:00:00Z |
| subagent-driven-development | yes | Sequential mode, 6 tasks |
| finishing-a-development-branch | yes | PR #6 created |
| pr-monitoring | yes | CI fix (v0.51.7 → v0.51.8 broken tag), Copilot review addressed, merge |
| post-merge-retrospective | yes | This document |

Full pipeline fired. No missed activations.

## What worked

- **Adversarial review cycle 1 was genuinely prescient.** The 6 Critical findings on the 12-message-type design were all real — field drift would have silently diverged plugin.json step names from proto field names. Choosing Option 1 (3 messages) eliminated the entire drift class before a single line of implementation code was written.
- **Option 1 design drastically reduced implementation surface.** 3 message types instead of 28+ means zero per-step proto maintenance. The ContractRegistry function is mechanical (72 identical `sfStep()` calls) rather than 72 uniquely-shaped message types.
- **setup-wfctl CI pattern works end-to-end.** Deriving wfctl version from go.mod + `GoCodeAlone/setup-wfctl@SHA` is robust and avoids hard-coding binary URLs in the workflow. `wfctl plugin validate --strict-contracts` passed first run after v0.51.8 bump.
- **Scope lock held.** 6 tasks, 1 PR — manifest honored exactly, no scope drift.

## What didn't

- **v0.51.7 broken tag cost one CI cycle.** The workflow `v0.51.7` tag had no `wfctl-linux-amd64` asset. This should have been caught by checking the GitHub release assets before committing the go.mod pin. Add a verification step: after bumping workflow dependency, confirm `wfctl-linux-amd64` exists on the corresponding release before pushing.
- **proto comment accuracy not verified against code.** The `api_version` default comment (`v58.0`) diverged from `defaultAPIVersion = "v63.0"` in client.go. This is a low-stakes miss but required an extra commit after Copilot caught it. A "verify proto comments match code constants" step in Task 1 (proto file authoring) would have prevented it.
- **Merge conflict with main required manual resolution.** Main had a concurrent `v0.51.7` bump (#4) while the feature branch was in flight with `v0.51.8`. The merge was clean (keep v0.51.8) but added a commit and re-triggered CI. Could have been avoided by rebasing onto main before final push.

## Plugin-level follow-ups

**One actionable pattern from this retro:**

The `api_version` proto comment miss (design doc + proto snippet listing wrong default) suggests adversarial-design-review's design-phase checklist should include:

> **Proto constant accuracy**: for any proto message field with a "default:" comment, verify the stated default matches the actual constant in the implementation code. Grep for the constant name and compare values.

This is a narrow but recurring risk for any plugin that writes proto before the code constants are finalized. Not broad enough for a new bug class, but worth a bullet under "Unstated assumptions" in the adversarial checklist — the assumption is "the default I'm documenting matches what the code does."

No other plugin-level changes warranted from this single PR. Pattern confirmation would require a second retro with the same miss.
