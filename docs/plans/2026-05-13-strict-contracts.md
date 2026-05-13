# Strict-Contracts Adoption Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add proto messages + `plugin.contracts.json` so `wfctl plugin validate --strict-contracts` passes for all 72 step types and 1 module type.

**Architecture:** Mirror worldsim PR #23 exactly. 3 proto messages (SalesforceProviderConfig, SalesforceStepInput with google.protobuf.Struct params, SalesforceStepOutput with Struct data). All 72 steps share the same generic input/output — no per-step message types. ContractRegistry() method added to salesforcePlugin. CI updated from v0.3.56 hardcoded to setup-wfctl derive-from-go.mod pattern.

**Tech Stack:** Go 1.26, google.golang.org/protobuf, worldsim pattern (workflow-plugin-worldsim), protoc + protoc-gen-go (developer only; generated file committed)

**Base branch:** feat/issue-5-strict-contracts

---

## Scope Manifest

**PR Count:** 1
**Tasks:** 6
**Estimated Lines of Change:** ~350 (1 proto file ~30 lines, 1 generated pb.go ~200 lines, contracts.go ~120 lines, plugin.go +5 lines, plugin.contracts.json ~160 lines, ci.yml +25 lines)

**Out of scope:**
- Changing any step Execute() logic or runtime behavior
- Adding typed per-step input messages (rejected in design; zero drift risk with Struct approach)
- Removing or adding step types
- Updating plugin.json schema fields

**PR Grouping:**

| PR # | Title | Tasks | Branch |
|------|-------|-------|--------|
| 1 | feat: add strict-contracts (proto + ContractRegistry + CI) | Task 1, Task 2, Task 3, Task 4, Task 5, Task 6 | feat/issue-5-strict-contracts |

**Status:** Draft

---

## Task 1: Write proto file

**Files:**
- Create: `proto/salesforce/v1/salesforce.proto`
- Create: `Makefile` (modify, add proto-gen target)

**Step 1: Write the proto file**

Create `proto/salesforce/v1/salesforce.proto`:

```proto
syntax = "proto3";
package workflow.plugin.salesforce.v1;

import "google/protobuf/struct.proto";

option go_package = "github.com/GoCodeAlone/workflow-plugin-salesforce/gen;salesforcev1";

// SalesforceProviderConfig is the typed config for the salesforce.provider module.
// Fields mirror internal.salesforceModule config keys (camelCase in YAML; proto uses snake_case).
message SalesforceProviderConfig {
  // login_url is the OAuth login URL (default: https://login.salesforce.com).
  string login_url = 1;
  // client_id is the OAuth client ID.
  string client_id = 2;
  // client_secret is the OAuth client secret.
  string client_secret = 3;
  // access_token is the direct access token (alternative to OAuth).
  string access_token = 4;
  // instance_url is the SF instance URL (required when using access_token).
  string instance_url = 5;
  // api_version is the SF API version (default: v58.0).
  string api_version = 6;
}

// SalesforceStepInput carries dynamic runtime inputs for any salesforce step.
// All per-step parameters are passed as a free-form Struct to avoid requiring
// 72 separate proto message types. Steps validate required params at runtime.
message SalesforceStepInput {
  // params holds the step-specific input parameters.
  google.protobuf.Struct params = 1;
}

// SalesforceStepOutput holds the result of any salesforce step execution.
// All per-step outputs are returned in the data Struct (free-form SF REST response).
message SalesforceStepOutput {
  // success indicates whether the step completed without error.
  bool success = 1;
  // error holds the error description when success is false.
  string error = 2;
  // data holds the step-specific output fields.
  google.protobuf.Struct data = 3;
}
```

**Step 2: Add proto-gen target to Makefile**

Add after the `clean:` target in `Makefile`:

```makefile
proto-gen:
	protoc \
		--proto_path=proto/salesforce/v1 \
		--go_out=gen \
		--go_opt=paths=source_relative \
		proto/salesforce/v1/salesforce.proto
```

Also add `proto-gen` to the `.PHONY` line at the top:
```makefile
.PHONY: build test install cross-build clean proto-gen
```

**Step 3: Verify proto syntax by checking it matches worldsim pattern**

Run: `diff proto/salesforce/v1/salesforce.proto /dev/null; echo "File exists: $?"`
Expected: `File exists: 0` (file exists)

**Step 4: Commit**

```bash
git add proto/salesforce/v1/salesforce.proto Makefile
git commit -m "feat: add salesforce proto messages for strict-contracts"
```

---

## Task 2: Generate and commit pb.go

**Files:**
- Create: `gen/salesforce.pb.go`

**Step 1: Check if protoc is available; if not, generate manually**

Run: `which protoc 2>/dev/null && echo "available" || echo "not available"`

If protoc is available:
```bash
mkdir -p gen
make proto-gen
```

If protoc is NOT available, generate the pb.go manually by adapting from worldsim:

The `gen/salesforce.pb.go` file must:
- Declare `package salesforcev1`
- Declare `var File_salesforce_proto` as a `protodesc.File` (accessible via `protodesc.ToFileDescriptorProto`)
- Define `SalesforceProviderConfig`, `SalesforceStepInput`, `SalesforceStepOutput` structs with proper protobuf reflection

The simplest approach when protoc is unavailable: copy `workflow-plugin-worldsim/gen/worldsim.pb.go` and adapt it. The pattern is identical — replace all `Worldsim*` with `Salesforce*` and field definitions accordingly.

**Step 2: Verify the generated file compiles**

Run: `GOWORK=off go build ./gen/... 2>&1`
Expected: no output (exit 0)

If build fails with missing imports, run `go mod tidy` first.

**Step 3: Verify File_salesforce_proto is exported**

Run: `grep 'File_salesforce_proto' gen/salesforce.pb.go | head -3`
Expected: Line containing `var File_salesforce_proto`

**Step 4: Commit**

```bash
git add gen/salesforce.pb.go
git commit -m "feat: add generated salesforce.pb.go"
```

---

## Task 3: Write contracts.go

**Files:**
- Create: `internal/contracts.go`

**Step 1: Write contracts.go**

Create `internal/contracts.go` — mirrors worldsim `internal/contracts.go` exactly with all 72 step entries:

```go
package internal

import (
	salesforcev1 "github.com/GoCodeAlone/workflow-plugin-salesforce/gen"
	pb "github.com/GoCodeAlone/workflow/plugin/external/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/structpb"
)

// ContractRegistry returns the typed contract descriptors for all salesforce
// module and step types. The workflow engine calls this via the
// sdk.ContractProvider interface for strict validation.
func (p *salesforcePlugin) ContractRegistry() *pb.ContractRegistry {
	return salesforceContractRegistry
}

// sfProtoPkg is the proto package prefix for all salesforce typed messages.
const sfProtoPkg = "workflow.plugin.salesforce.v1."

// salesforceContractRegistry declares STRICT_PROTO contracts for all salesforce
// module and step types. The FileDescriptorSet includes google.protobuf.Struct
// (used in SalesforceStepInput.params and SalesforceStepOutput.data fields).
var salesforceContractRegistry = &pb.ContractRegistry{
	FileDescriptorSet: &descriptorpb.FileDescriptorSet{
		File: []*descriptorpb.FileDescriptorProto{
			protodesc.ToFileDescriptorProto(structpb.File_google_protobuf_struct_proto),
			protodesc.ToFileDescriptorProto(salesforcev1.File_salesforce_proto),
		},
	},
	Contracts: []*pb.ContractDescriptor{
		// ── module ───────────────────────────────────────────────────────────
		{
			Kind:          pb.ContractKind_CONTRACT_KIND_MODULE,
			ModuleType:    "salesforce.provider",
			ConfigMessage: sfProtoPkg + "SalesforceProviderConfig",
			Mode:          pb.ContractMode_CONTRACT_MODE_STRICT_PROTO,
		},
		// ── SObject CRUD steps ───────────────────────────────────────────────
		sfStep("step.salesforce_record_get"),
		sfStep("step.salesforce_record_create"),
		sfStep("step.salesforce_record_update"),
		sfStep("step.salesforce_record_upsert"),
		sfStep("step.salesforce_record_delete"),
		sfStep("step.salesforce_record_describe"),
		sfStep("step.salesforce_describe_global"),
		// ── SOQL / SOSL ──────────────────────────────────────────────────────
		sfStep("step.salesforce_query"),
		sfStep("step.salesforce_query_all"),
		sfStep("step.salesforce_search"),
		// ── Collections ──────────────────────────────────────────────────────
		sfStep("step.salesforce_collection_insert"),
		sfStep("step.salesforce_collection_update"),
		sfStep("step.salesforce_collection_upsert"),
		sfStep("step.salesforce_collection_delete"),
		// ── Composite ────────────────────────────────────────────────────────
		sfStep("step.salesforce_composite_request"),
		sfStep("step.salesforce_composite_tree"),
		// ── Bulk API v2 ──────────────────────────────────────────────────────
		sfStep("step.salesforce_bulk_insert"),
		sfStep("step.salesforce_bulk_update"),
		sfStep("step.salesforce_bulk_upsert"),
		sfStep("step.salesforce_bulk_delete"),
		sfStep("step.salesforce_bulk_query"),
		sfStep("step.salesforce_bulk_query_results"),
		sfStep("step.salesforce_bulk_job_status"),
		sfStep("step.salesforce_bulk_job_abort"),
		// ── Tooling API ──────────────────────────────────────────────────────
		sfStep("step.salesforce_tooling_query"),
		sfStep("step.salesforce_tooling_get"),
		sfStep("step.salesforce_tooling_create"),
		sfStep("step.salesforce_tooling_update"),
		sfStep("step.salesforce_tooling_delete"),
		// ── Apex ─────────────────────────────────────────────────────────────
		sfStep("step.salesforce_apex_execute"),
		sfStep("step.salesforce_apex_get"),
		sfStep("step.salesforce_apex_post"),
		sfStep("step.salesforce_apex_patch"),
		sfStep("step.salesforce_apex_put"),
		sfStep("step.salesforce_apex_delete"),
		// ── Reports & Dashboards ─────────────────────────────────────────────
		sfStep("step.salesforce_report_list"),
		sfStep("step.salesforce_report_describe"),
		sfStep("step.salesforce_report_run"),
		sfStep("step.salesforce_dashboard_list"),
		sfStep("step.salesforce_dashboard_describe"),
		sfStep("step.salesforce_dashboard_refresh"),
		// ── Approvals ────────────────────────────────────────────────────────
		sfStep("step.salesforce_approval_list"),
		sfStep("step.salesforce_approval_submit"),
		sfStep("step.salesforce_approval_approve"),
		sfStep("step.salesforce_approval_reject"),
		// ── Chatter ──────────────────────────────────────────────────────────
		sfStep("step.salesforce_chatter_post"),
		sfStep("step.salesforce_chatter_comment"),
		sfStep("step.salesforce_chatter_like"),
		sfStep("step.salesforce_chatter_feed_list"),
		// ── Files & Content ──────────────────────────────────────────────────
		sfStep("step.salesforce_file_upload"),
		sfStep("step.salesforce_file_download"),
		sfStep("step.salesforce_content_version_create"),
		sfStep("step.salesforce_content_document_get"),
		sfStep("step.salesforce_content_document_delete"),
		// ── Users & Identity ─────────────────────────────────────────────────
		sfStep("step.salesforce_user_get"),
		sfStep("step.salesforce_user_list"),
		sfStep("step.salesforce_user_create"),
		sfStep("step.salesforce_user_update"),
		sfStep("step.salesforce_identity_get"),
		sfStep("step.salesforce_org_limits"),
		// ── Flows ────────────────────────────────────────────────────────────
		sfStep("step.salesforce_flow_list"),
		sfStep("step.salesforce_flow_run"),
		// ── Platform Events ──────────────────────────────────────────────────
		sfStep("step.salesforce_event_publish"),
		// ── Metadata ─────────────────────────────────────────────────────────
		sfStep("step.salesforce_metadata_describe"),
		sfStep("step.salesforce_metadata_list"),
		sfStep("step.salesforce_metadata_read"),
		sfStep("step.salesforce_metadata_create"),
		sfStep("step.salesforce_metadata_update"),
		sfStep("step.salesforce_metadata_delete"),
		sfStep("step.salesforce_metadata_deploy"),
		sfStep("step.salesforce_metadata_retrieve"),
		// ── Generic ──────────────────────────────────────────────────────────
		sfStep("step.salesforce_raw_request"),
	},
}

// sfStep returns a ContractDescriptor for a salesforce step using the shared
// SalesforceStepInput and SalesforceStepOutput messages.
func sfStep(stepType string) *pb.ContractDescriptor {
	return &pb.ContractDescriptor{
		Kind:          pb.ContractKind_CONTRACT_KIND_STEP,
		StepType:      stepType,
		InputMessage:  sfProtoPkg + "SalesforceStepInput",
		OutputMessage: sfProtoPkg + "SalesforceStepOutput",
		Mode:          pb.ContractMode_CONTRACT_MODE_STRICT_PROTO,
	}
}

// Compile-time assertion: salesforcePlugin implements sdk.ContractProvider.
var _ interface{ ContractRegistry() *pb.ContractRegistry } = (*salesforcePlugin)(nil)
```

**Step 2: Verify compile-time assertion compiles**

Run: `GOWORK=off go build ./internal/... 2>&1`
Expected: no output (exit 0)

If it fails with `pb.ContractDescriptor has no field StepType`, check the actual field name in the workflow proto — it may be `Type` rather than `StepType`. Grep the worldsim contracts.go or the proto package to confirm: `grep -r "StepType\|\.Type\b" /Users/jon/workspace/workflow-plugin-worldsim/internal/contracts.go`

**Step 3: Verify contract count**

Run: `grep -c 'sfStep(' internal/contracts.go`
Expected: `72`

**Step 4: Commit**

```bash
git add internal/contracts.go
git commit -m "feat: add ContractRegistry with 73 contracts (1 module + 72 steps)"
```

---

## Task 4: Add ContractRegistry method to plugin.go

**Files:**
- Modify: `internal/plugin.go`

**Step 1: Verify that contracts.go already attaches ContractRegistry to salesforcePlugin**

Run: `grep 'ContractRegistry' internal/contracts.go`
Expected: `func (p *salesforcePlugin) ContractRegistry() *pb.ContractRegistry {`

The compile-time assertion in contracts.go is sufficient — no changes needed to `plugin.go` unless the workflow SDK's ContractProvider interface requires registration at startup. Check:

Run: `grep -r 'ContractProvider\|ContractRegistry' /Users/jon/workspace/workflow-plugin-worldsim/internal/plugin.go 2>/dev/null || echo "not in worldsim plugin.go"`

If worldsim's `plugin.go` has no ContractProvider-related code (i.e., the SDK discovers it via interface assertion), no change is needed to `internal/plugin.go`.

If the SDK requires explicit registration (e.g., `sdk.RegisterContractProvider(p)`), add it to `NewSalesforcePlugin()` in `internal/plugin.go`.

**Step 2: Build the entire project to confirm no import cycles**

Run: `GOWORK=off go build ./... 2>&1`
Expected: no output (exit 0)

**Step 3: Run tests**

Run: `GOWORK=off go test ./... -count=1 2>&1 | tail -20`
Expected: All `PASS` or `ok` lines; no `FAIL`

**Step 4: Commit (only if plugin.go was modified)**

```bash
# Only run this if plugin.go was changed
git add internal/plugin.go
git commit -m "feat: wire ContractProvider interface in plugin registration"
```

---

## Task 5: Write plugin.contracts.json

**Files:**
- Create: `plugin.contracts.json`

**Step 1: Write plugin.contracts.json**

Create `plugin.contracts.json` (mirrors worldsim structure exactly):

```json
{
  "version": "1",
  "contracts": [
    {
      "kind": "module",
      "type": "salesforce.provider",
      "mode": "strict_proto",
      "config": "workflow.plugin.salesforce.v1.SalesforceProviderConfig"
    },
    {
      "kind": "step",
      "type": "step.salesforce_record_get",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_record_create",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_record_update",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_record_upsert",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_record_delete",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_record_describe",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_describe_global",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_query",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_query_all",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_search",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_collection_insert",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_collection_update",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_collection_upsert",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_collection_delete",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_composite_request",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_composite_tree",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_bulk_insert",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_bulk_update",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_bulk_upsert",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_bulk_delete",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_bulk_query",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_bulk_query_results",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_bulk_job_status",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_bulk_job_abort",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_tooling_query",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_tooling_get",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_tooling_create",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_tooling_update",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_tooling_delete",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_apex_execute",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_apex_get",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_apex_post",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_apex_patch",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_apex_put",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_apex_delete",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_report_list",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_report_describe",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_report_run",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_dashboard_list",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_dashboard_describe",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_dashboard_refresh",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_approval_list",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_approval_submit",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_approval_approve",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_approval_reject",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_chatter_post",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_chatter_comment",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_chatter_like",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_chatter_feed_list",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_file_upload",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_file_download",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_content_version_create",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_content_document_get",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_content_document_delete",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_user_get",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_user_list",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_user_create",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_user_update",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_identity_get",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_org_limits",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_flow_list",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_flow_run",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_event_publish",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_metadata_describe",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_metadata_list",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_metadata_read",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_metadata_create",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_metadata_update",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_metadata_delete",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_metadata_deploy",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_metadata_retrieve",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    },
    {
      "kind": "step",
      "type": "step.salesforce_raw_request",
      "mode": "strict_proto",
      "input": "workflow.plugin.salesforce.v1.SalesforceStepInput",
      "output": "workflow.plugin.salesforce.v1.SalesforceStepOutput"
    }
  ]
}
```

**Step 2: Verify entry count**

Run: `python3 -c "import json; d=json.load(open('plugin.contracts.json')); print(len(d['contracts']), 'contracts')"`
Expected: `73 contracts`

**Step 3: Verify JSON is valid**

Run: `python3 -m json.tool plugin.contracts.json > /dev/null && echo "valid JSON"`
Expected: `valid JSON`

**Step 4: Commit**

```bash
git add plugin.contracts.json
git commit -m "feat: add plugin.contracts.json (73 entries: 1 module + 72 steps)"
```

---

## Task 6: Update CI wfctl-strict-contracts job

**Files:**
- Modify: `.github/workflows/ci.yml`

**Step 1: Replace the `wfctl-strict-contracts` job in ci.yml**

Replace the existing `wfctl-strict-contracts` job block with the worldsim pattern:

```yaml
  wfctl-strict-contracts:
    name: Validate strict plugin contracts
    runs-on: ubuntu-latest
    permissions:
      contents: read
    steps:
      - uses: actions/checkout@v4

      - name: Check plugin.json and plugin.contracts.json exist
        run: |
          if [ ! -f plugin.json ]; then
            echo "::error::plugin.json is missing"
            exit 1
          fi
          if [ ! -f plugin.contracts.json ]; then
            echo "::error::plugin.contracts.json is missing"
            exit 1
          fi

      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod

      - name: Derive wfctl version from go.mod
        id: wfctl-ver
        run: |
          ver=$(go list -m -f '{{.Version}}' github.com/GoCodeAlone/workflow)
          echo "version=$ver" >> "$GITHUB_OUTPUT"

      - uses: GoCodeAlone/setup-wfctl@bcd880980f5bbe8d192d0c20ff6279d25331f956
        with:
          version: ${{ steps.wfctl-ver.outputs.version }}

      - name: Validate strict plugin contracts
        run: wfctl plugin validate --file plugin.json --strict-contracts
```

Rollback: revert ci.yml to the `go run github.com/GoCodeAlone/workflow/cmd/wfctl@v0.3.56 plugin validate --file plugin.json` form. This restores the old non-strict validation gate.

**Step 2: Verify ci.yml is valid YAML**

Run: `python3 -c "import yaml; yaml.safe_load(open('.github/workflows/ci.yml'))" && echo "valid YAML"`
Expected: `valid YAML`

**Step 3: Verify the hardcoded v0.3.56 version is gone**

Run: `grep 'v0.3.56' .github/workflows/ci.yml`
Expected: no output (empty)

**Step 4: Final build + test pass**

Run: `GOWORK=off go build ./... && GOWORK=off go vet ./... && echo "BUILD+VET OK"`
Expected: `BUILD+VET OK`

Run: `GOWORK=off go test ./... -count=1 2>&1 | tail -5`
Expected: All `ok` lines; no `FAIL`

**Step 5: Commit**

```bash
git add .github/workflows/ci.yml
git commit -m "ci: update wfctl-strict-contracts job to derive version from go.mod"
```
