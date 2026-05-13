# Design: Strict-Contracts Adoption — workflow-plugin-salesforce

**Date:** 2026-05-13
**Issue:** #5
**Author:** autonomous pipeline (brainstorming phase; revised cycle 2 — Option 1 adopted)

## Goal

Add `plugin.contracts.json` + proto-typed messages so the plugin passes `wfctl plugin validate --strict-contracts`. Mirror the worldsim PR #23 SDK pattern exactly. Closes issue #5.

## Context

The plugin currently has:
- 1 module type: `salesforce.provider`
- 72 step types across 14 Salesforce API categories
- No `plugin.contracts.json`, no `.proto`, no `.pb.go`
- CI validates `plugin.json` schema via `wfctl plugin validate --file plugin.json` at wfctl v0.3.56 (no `--strict-contracts` flag)

## Decision: Option 1 — Struct-for-all-inputs

Cycle-1 adversarial review found 6 Critical field-name drift bugs in the "12 shared Input message types" approach (field names in proto messages did not match the actual field names used in step Execute() implementations). The author adopted **Option 1** as the fix:

- Use `google.protobuf.Struct` for ALL step inputs (same as worldsim's `WorldsimCallInput.params`).
- This eliminates all 6 critical field-name drift findings in one stroke.
- Reduces 28+ message types (cycle-1 design had ~32 Input messages) to **3 messages total**.
- Steps already validate required fields at runtime — proto-level field-name enforcement adds nothing for this plugin's REST API call pattern.

## Proto Package Structure

```
proto/
  salesforce/v1/
    salesforce.proto        # all 3 messages in one file
gen/
  salesforce.pb.go          # generated, committed
```

Package: `workflow.plugin.salesforce.v1`
Go package: `github.com/GoCodeAlone/workflow-plugin-salesforce/gen;salesforcev1`

## Message Design (3 messages total)

```proto
syntax = "proto3";
package workflow.plugin.salesforce.v1;

import "google/protobuf/struct.proto";

option go_package = "github.com/GoCodeAlone/workflow-plugin-salesforce/gen;salesforcev1";

// SalesforceProviderConfig is the typed config for the salesforce.provider module.
// Fields mirror internal.salesforceModuleConfig.
message SalesforceProviderConfig {
  string login_url     = 1;  // OAuth login URL (default: https://login.salesforce.com)
  string client_id     = 2;  // OAuth client ID
  string client_secret = 3;  // OAuth client secret
  string access_token  = 4;  // Direct access token (alternative to OAuth)
  string instance_url  = 5;  // SF instance URL (required with access_token)
  string api_version   = 6;  // SF API version (default: v63.0)
}

// SalesforceStepInput carries dynamic runtime inputs for any salesforce step.
// All per-step parameters are passed as a free-form Struct to avoid requiring
// 72 separate proto message types. Steps validate required params at runtime.
message SalesforceStepInput {
  google.protobuf.Struct params = 1;
}

// SalesforceStepOutput holds the result of any salesforce step execution.
// All per-step outputs are returned in the data Struct (free-form SF REST response).
message SalesforceStepOutput {
  bool   success = 1;  // false if error field is non-empty
  string error   = 2;  // error description when success=false
  google.protobuf.Struct data = 3;  // SF REST API response body
}
```

## Step Config (no separate message — all 72 steps share StepInput/StepOutput)

Each step contract entry uses:
- Config: none (step config is part of the module config; steps have no extra per-step config fields in the current implementation)
- Input: `workflow.plugin.salesforce.v1.SalesforceStepInput`
- Output: `workflow.plugin.salesforce.v1.SalesforceStepOutput`

> Note: If `wfctl plugin validate --strict-contracts` requires a ConfigMessage for each step entry, use an empty `SalesforceStepConfig {}` message. This will be verified during implementation before committing contracts.go.

## plugin.contracts.json (73 entries)

All 73 contracts use `mode: strict_proto` and reference the 3 messages above:
- 1 module contract: `salesforce.provider` → config: `SalesforceProviderConfig`
- 72 step contracts: each step type → input: `SalesforceStepInput`, output: `SalesforceStepOutput`

The 72 step types (from `internal/step_registry.go`):
- SObject CRUD (7): salesforce_record_get, record_create, record_update, record_upsert, record_delete, record_describe, describe_global
- SOQL/SOSL (3): salesforce_query, query_all, search
- Collections (4): collection_insert, collection_update, collection_upsert, collection_delete
- Composite (2): composite_request, composite_tree
- Bulk API v2 (8): bulk_insert, bulk_update, bulk_upsert, bulk_delete, bulk_query, bulk_query_results, bulk_job_status, bulk_job_abort
- Tooling API (5): tooling_query, tooling_get, tooling_create, tooling_update, tooling_delete
- Apex (6): apex_execute, apex_get, apex_post, apex_patch, apex_put, apex_delete
- Reports/Dashboards (4): report_list, report_describe, report_run, dashboard_list, dashboard_describe, dashboard_refresh
- Approvals (4): approval_list, approval_submit, approval_approve, approval_reject
- Chatter (4): chatter_post, chatter_comment, chatter_like, chatter_feed_list
- Files (4): file_upload, file_download, content_version_create, content_document_get, content_document_delete
- Users (5): user_list, user_get, user_create, user_update, identity_get
- Flows/Events (2): flow_list, flow_run, event_publish
- Metadata (7): metadata_describe, metadata_list, metadata_read, metadata_create, metadata_update, metadata_delete, metadata_deploy, metadata_retrieve
- Misc (1): raw_request, org_limits

> Step type strings use the `step.salesforce_` prefix (from step_registry.go).

## ContractRegistry Implementation

Mirrors worldsim `internal/contracts.go` exactly:

```go
package internal

import (
    salesforcev1 "github.com/GoCodeAlone/workflow-plugin-salesforce/gen"
    pb "github.com/GoCodeAlone/workflow/plugin/external/proto"
    "google.golang.org/protobuf/reflect/protodesc"
    "google.golang.org/protobuf/types/descriptorpb"
    "google.golang.org/protobuf/types/known/structpb"
)

func (p *salesforcePlugin) ContractRegistry() *pb.ContractRegistry {
    return salesforceContractRegistry
}

var salesforceContractRegistry = &pb.ContractRegistry{
    FileDescriptorSet: &descriptorpb.FileDescriptorSet{
        File: []*descriptorpb.FileDescriptorProto{
            protodesc.ToFileDescriptorProto(structpb.File_google_protobuf_struct_proto),
            protodesc.ToFileDescriptorProto(salesforcev1.File_salesforce_proto),
        },
    },
    Contracts: []*pb.ContractDescriptor{
        // module
        {
            Kind:          pb.ContractKind_CONTRACT_KIND_MODULE,
            ModuleType:    "salesforce.provider",
            ConfigMessage: sfProtoPkg + "SalesforceProviderConfig",
            Mode:          pb.ContractMode_CONTRACT_MODE_STRICT_PROTO,
        },
        // 72 step entries follow (generated or hand-written loop)
        // Each entry: Kind=STEP, StepType="step.salesforce_X",
        // InputMessage=sfProtoPkg+"SalesforceStepInput",
        // OutputMessage=sfProtoPkg+"SalesforceStepOutput"
    },
}

const sfProtoPkg = "workflow.plugin.salesforce.v1."

var _ interface{ ContractRegistry() *pb.ContractRegistry } = (*salesforcePlugin)(nil)
```

## CI Update

Replace existing `wfctl plugin validate --file plugin.json` (wfctl v0.3.56) with worldsim-style job:
1. Check `plugin.json` and `plugin.contracts.json` both exist
2. Derive wfctl version from `go.mod`
3. Use `GoCodeAlone/setup-wfctl@bcd880980f5bbe8d192d0c20ff6279d25331f956`
4. Run `wfctl plugin validate --file plugin.json --strict-contracts`

> Note: The existing CI has a duplicate `test` + `test-lint` job pair that both do identical steps. This design does not change that — it only updates the `wfctl-strict-contracts` job.

## Files Changed

- `proto/salesforce/v1/salesforce.proto` — new (3 messages)
- `gen/salesforce.pb.go` — new (generated, committed)
- `internal/contracts.go` — new (ContractRegistry implementation)
- `internal/plugin.go` — add `ContractRegistry()` method
- `plugin.contracts.json` — new (73 entries)
- `Makefile` — add `proto-gen` target
- `.github/workflows/ci.yml` — update `wfctl-strict-contracts` job

## Assumptions

1. `google.protobuf.Struct` is accepted by the strict-contracts validator for input/output data fields. **Confirmed: worldsim uses this pattern and passes.**
2. All 72 step types retain their existing names and runtime behavior — no breaking changes to the step Execute() path.
3. The `salesforcePlugin` type in `internal/plugin.go` is the type that implements the sdk interface and can accept the `ContractRegistry()` method.
4. Committing the generated `.pb.go` file is correct (worldsim does this; protoc not required in CI).
5. `wfctl plugin validate --strict-contracts` passes when `plugin.contracts.json` + `gen/salesforce.pb.go` + `ContractRegistry()` method are present.
6. Steps do NOT have a per-step Config message in this plugin's implementation (steps get their config via the module). If `wfctl` requires a ConfigMessage per step contract, an empty `SalesforceStepConfig {}` message will be added. This is verified during implementation, not assumed.
7. The step count of 72 is authoritative from `internal/step_registry.go`. If `plugin.json`'s 72 `stepSchemas` entries and the registry entries diverge, the registry is the source of truth for contracts.

## Rollback

No runtime behavior changes — ContractRegistry is metadata-only. Rollback = revert the `ci.yml` change to the pre-strict job, then delete `plugin.contracts.json`, `proto/`, `gen/`, and `internal/contracts.go`. Plugin execution is unaffected because ContractRegistry is only called by the validator.

## Tradeoffs Accepted

- **No field-level enforcement at proto boundary**: Option 1 trades static field-name checking for zero drift risk. Accepted because: (a) 72 steps already validate fields at runtime, (b) cycle-1 design had 6 Critical bugs from drift, (c) worldsim uses the same pattern successfully.
- **73 nearly-identical contracts.json entries**: The file is mechanical/generated. Verbose but correct.
- **Step config not typed in proto**: Steps don't have meaningful per-step config in this plugin. Not typing it is YAGNI.
