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
