# Design: Strict-Contracts Adoption — workflow-plugin-salesforce

**Date:** 2026-05-13
**Issue:** #5
**Author:** autonomous pipeline (brainstorming phase)

## Goal

Add `plugin.contracts.json` + proto-typed messages so the plugin passes `wfctl plugin validate --strict-contracts`. Mirror the worldsim PR #23 SDK pattern. Closes issue #5.

## Context

The plugin currently has:
- 1 module type: `salesforce.provider`
- 72 step types across 14 Salesforce API categories
- No `plugin.contracts.json`, no `.proto`, no `.pb.go`
- CI validates `plugin.json` schema via `wfctl plugin validate --file plugin.json` (no `--strict-contracts` flag)

The worldsim precedent (1 step type dispatching 70+ operations via an `action` string) used `google.protobuf.Struct` to carry dynamic parameters. Salesforce has 72 individual step types — they must remain as-is per `feedback_force_strict_contracts_no_compat`.

## Approach: Shared Message Groups (Single PR)

Rather than 72 unique Input message types, we group the 72 steps into ~12 shared Input message types based on their actual parameter shapes. All steps share a common Output shape (`SalesforceOutput` with `data: google.protobuf.Struct`), avoiding the structpb typed-slice blocker.

### Why single PR

All 72 steps share structurally identical Output shapes (free-form SF REST response). The Input messages cover 10-12 distinct parameter shapes across all 72 steps. No step has a unique input/output signature that requires its own message pair. Single PR is correct.

### Why not unique per-step messages

72 unique Input + 72 unique Output messages = 144 proto messages. Most would be identical or differ by 1-2 fields. This is maximum complexity for minimum benefit. The worldsim precedent confirms `Struct` is the right tool for dynamic API payloads.

## Proto Package Structure

```
proto/
  salesforce/v1/
    salesforce.proto        # all messages in one file
gen/
  salesforce.pb.go          # generated, committed
```

Package: `workflow.plugin.salesforce.v1`
Go package: `github.com/GoCodeAlone/workflow-plugin-salesforce/gen;salesforcev1`

## Message Design

### Module Config

```proto
message SalesforceProviderConfig {
  string login_url     = 1;  // OAuth login URL
  string client_id     = 2;  // OAuth client ID
  string client_secret = 3;  // OAuth client secret
  string access_token  = 4;  // Direct access token (alternative to OAuth)
  string instance_url  = 5;  // SF instance URL (required with access_token)
  string api_version   = 6;  // SF API version (default: v58.0)
}
```

### Shared Output (all 72 steps)

```proto
message SalesforceOutput {
  bool success              = 1;  // false if error field is non-empty
  string error              = 2;  // error description when success=false
  google.protobuf.Struct data = 3;  // free-form SF REST API response body
}
```

### Input Message Groups (12 types covering all 72 steps)

| Message | Steps |
|---|---|
| `RecordInput` (sobject_type, record_id, fields) | record_get, record_update, record_delete |
| `RecordCreateInput` (sobject_type, fields) | record_create |
| `RecordUpsertInput` (sobject_type, external_id_field, external_id_value, fields) | record_upsert |
| `RecordDescribeInput` (sobject_type) | record_describe |
| `EmptyInput` () | describe_global, approval_list, user_list, identity_get, org_limits, flow_list, report_list, dashboard_list, metadata_describe |
| `QueryInput` (soql) | query, query_all, bulk_query |
| `SearchInput` (sosl) | search |
| `CollectionInput` (records, all_or_none) | collection_insert, collection_update |
| `CollectionDeleteInput` (ids, all_or_none) | collection_delete |
| `CollectionUpsertInput` (sobject_type, external_id_field, records, all_or_none) | collection_upsert |
| `CompositeRequestInput` (composite_request, all_or_none) | composite_request |
| `CompositeTreeInput` (sobject_type, records) | composite_tree |
| `BulkJobInput` (sobject_type, external_id_field) | bulk_insert, bulk_update, bulk_upsert, bulk_delete |
| `BulkJobQueryInput` (job_id, max_records) | bulk_query_results |
| `BulkJobStatusInput` (job_id, job_type) | bulk_job_status, bulk_job_abort |
| `ToolingInput` (sobject_type, record_id, fields, soql) | tooling_query, tooling_get, tooling_create, tooling_update, tooling_delete |
| `ApexExecuteInput` (apex_body, log_levels) | apex_execute |
| `ApexRestInput` (apex_path, body) | apex_get, apex_post, apex_patch, apex_put, apex_delete |
| `ReportInput` (report_id) | report_describe, report_run |
| `DashboardInput` (dashboard_id) | dashboard_describe, dashboard_refresh |
| `ApprovalActionInput` (record_id, work_item_id, comments) | approval_submit, approval_approve, approval_reject |
| `ChatterPostInput` (subject_id, text) | chatter_post |
| `ChatterCommentInput` (feed_element_id, text) | chatter_comment, chatter_like, chatter_feed_list |
| `FileUploadInput` (title, path_on_client, version_data, linked_entity_id) | file_upload, content_version_create |
| `FileDownloadInput` (content_version_id) | file_download, content_document_get, content_document_delete |
| `UserGetInput` (user_id) | user_get |
| `UserCreateInput` (fields) | user_create |
| `UserUpdateInput` (user_id, fields) | user_update |
| `FlowRunInput` (flow_api_name, inputs) | flow_run |
| `EventPublishInput` (event_type, payload) | event_publish |
| `MetadataInput` (metadata_type, full_names, metadata) | metadata_list, metadata_read, metadata_create, metadata_update, metadata_delete, metadata_deploy, metadata_retrieve |
| `RawRequestInput` (method, path, body) | raw_request |

### Step Config (all steps share common module reference)

```proto
message SalesforceStepConfig {
  string module = 1;  // salesforce.provider module name (default: "salesforce")
}
```

Each step type uses `SalesforceStepConfig` as its Config message, plus its specific Input message and `SalesforceOutput` as Output.

## ContractRegistry Implementation

Follows worldsim pattern exactly:
- `internal/contracts.go` implements `ContractRegistry() *pb.ContractRegistry`
- `salesforcePlugin` gets `ContractRegistry()` method
- `plugin.contracts.json` declares all 73 contracts (1 module + 72 steps)

## CI Update

Replace existing `wfctl-strict-contracts` job with worldsim-style job that:
1. Checks both `plugin.json` and `plugin.contracts.json` exist
2. Derives wfctl version from `go.mod`
3. Uses `GoCodeAlone/setup-wfctl@bcd880980f5bbe8d192d0c20ff6279d25331f956`
4. Runs `wfctl plugin validate --file plugin.json --strict-contracts`

## Files Changed

- `proto/salesforce/v1/salesforce.proto` — new
- `gen/salesforce.pb.go` — new (generated)
- `internal/contracts.go` — new
- `internal/plugin.go` — add `ContractRegistry()` method
- `plugin.contracts.json` — new
- `Makefile` — add `proto-gen` target
- `.github/workflows/ci.yml` — update strict-contracts job

## Assumptions

1. `google.protobuf.Struct` is accepted by the strict-contracts validator for Output data fields (confirmed: worldsim uses this pattern).
2. All 72 step types retain their existing names and runtime behavior — no breaking changes to the step execute path.
3. The `SalesforceStepConfig.module` field (string, default "salesforce") correctly mirrors how `getModuleName(config)` works.
4. Committing the generated `.pb.go` file is correct (worldsim does this; no protoc required in CI).
5. `wfctl plugin validate --strict-contracts` passes when `plugin.contracts.json` + `plugin.pb.go` + `ContractRegistry()` are all present.

## Rollback

No runtime behavior changes — ContractRegistry is metadata-only. If the CI job fails, revert the `.github/workflows/ci.yml` change to restore the pre-strict job. The proto/gen/contracts files can be removed without affecting plugin execution.
