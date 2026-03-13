package internal

import (
	"fmt"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// stepConstructor is a function that creates a StepInstance.
type stepConstructor func(name string, config map[string]any) (sdk.StepInstance, error)

// stepRegistry maps step type strings to constructor functions.
var stepRegistry = map[string]stepConstructor{
	// SObject CRUD
	"step.salesforce_record_get":      func(n string, c map[string]any) (sdk.StepInstance, error) { return newRecordGetStep(n, c) },
	"step.salesforce_record_create":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newRecordCreateStep(n, c) },
	"step.salesforce_record_update":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newRecordUpdateStep(n, c) },
	"step.salesforce_record_upsert":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newRecordUpsertStep(n, c) },
	"step.salesforce_record_delete":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newRecordDeleteStep(n, c) },
	"step.salesforce_record_describe": func(n string, c map[string]any) (sdk.StepInstance, error) { return newRecordDescribeStep(n, c) },
	"step.salesforce_describe_global": func(n string, c map[string]any) (sdk.StepInstance, error) { return newDescribeGlobalStep(n, c) },

	// SOQL / SOSL
	"step.salesforce_query":     func(n string, c map[string]any) (sdk.StepInstance, error) { return newQueryStep(n, c) },
	"step.salesforce_query_all": func(n string, c map[string]any) (sdk.StepInstance, error) { return newQueryAllStep(n, c) },
	"step.salesforce_search":    func(n string, c map[string]any) (sdk.StepInstance, error) { return newSearchStep(n, c) },

	// Collections
	"step.salesforce_collection_insert": func(n string, c map[string]any) (sdk.StepInstance, error) { return newCollectionInsertStep(n, c) },
	"step.salesforce_collection_update": func(n string, c map[string]any) (sdk.StepInstance, error) { return newCollectionUpdateStep(n, c) },
	"step.salesforce_collection_upsert": func(n string, c map[string]any) (sdk.StepInstance, error) { return newCollectionUpsertStep(n, c) },
	"step.salesforce_collection_delete": func(n string, c map[string]any) (sdk.StepInstance, error) { return newCollectionDeleteStep(n, c) },

	// Composite
	"step.salesforce_composite_request": func(n string, c map[string]any) (sdk.StepInstance, error) { return newCompositeRequestStep(n, c) },
	"step.salesforce_composite_tree":    func(n string, c map[string]any) (sdk.StepInstance, error) { return newCompositeTreeStep(n, c) },

	// Bulk API v2
	"step.salesforce_bulk_insert":       func(n string, c map[string]any) (sdk.StepInstance, error) { return newBulkInsertStep(n, c) },
	"step.salesforce_bulk_update":       func(n string, c map[string]any) (sdk.StepInstance, error) { return newBulkUpdateStep(n, c) },
	"step.salesforce_bulk_upsert":       func(n string, c map[string]any) (sdk.StepInstance, error) { return newBulkUpsertStep(n, c) },
	"step.salesforce_bulk_delete":       func(n string, c map[string]any) (sdk.StepInstance, error) { return newBulkDeleteStep(n, c) },
	"step.salesforce_bulk_query":        func(n string, c map[string]any) (sdk.StepInstance, error) { return newBulkQueryStep(n, c) },
	"step.salesforce_bulk_query_results": func(n string, c map[string]any) (sdk.StepInstance, error) { return newBulkQueryResultsStep(n, c) },
	"step.salesforce_bulk_job_status":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newBulkJobStatusStep(n, c) },
	"step.salesforce_bulk_job_abort":    func(n string, c map[string]any) (sdk.StepInstance, error) { return newBulkJobAbortStep(n, c) },

	// Tooling API
	"step.salesforce_tooling_query":  func(n string, c map[string]any) (sdk.StepInstance, error) { return newToolingQueryStep(n, c) },
	"step.salesforce_tooling_get":    func(n string, c map[string]any) (sdk.StepInstance, error) { return newToolingGetStep(n, c) },
	"step.salesforce_tooling_create": func(n string, c map[string]any) (sdk.StepInstance, error) { return newToolingCreateStep(n, c) },
	"step.salesforce_tooling_update": func(n string, c map[string]any) (sdk.StepInstance, error) { return newToolingUpdateStep(n, c) },
	"step.salesforce_tooling_delete": func(n string, c map[string]any) (sdk.StepInstance, error) { return newToolingDeleteStep(n, c) },
	"step.salesforce_apex_execute":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newApexExecuteStep(n, c) },

	// Apex REST
	"step.salesforce_apex_get":    func(n string, c map[string]any) (sdk.StepInstance, error) { return newApexGetStep(n, c) },
	"step.salesforce_apex_post":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newApexPostStep(n, c) },
	"step.salesforce_apex_patch":  func(n string, c map[string]any) (sdk.StepInstance, error) { return newApexPatchStep(n, c) },
	"step.salesforce_apex_put":    func(n string, c map[string]any) (sdk.StepInstance, error) { return newApexPutStep(n, c) },
	"step.salesforce_apex_delete": func(n string, c map[string]any) (sdk.StepInstance, error) { return newApexDeleteStep(n, c) },

	// Reports & Dashboards
	"step.salesforce_report_list":         func(n string, c map[string]any) (sdk.StepInstance, error) { return newReportListStep(n, c) },
	"step.salesforce_report_describe":     func(n string, c map[string]any) (sdk.StepInstance, error) { return newReportDescribeStep(n, c) },
	"step.salesforce_report_run":          func(n string, c map[string]any) (sdk.StepInstance, error) { return newReportRunStep(n, c) },
	"step.salesforce_dashboard_list":      func(n string, c map[string]any) (sdk.StepInstance, error) { return newDashboardListStep(n, c) },
	"step.salesforce_dashboard_describe":  func(n string, c map[string]any) (sdk.StepInstance, error) { return newDashboardDescribeStep(n, c) },
	"step.salesforce_dashboard_refresh":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newDashboardRefreshStep(n, c) },

	// Approvals
	"step.salesforce_approval_list":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newApprovalListStep(n, c) },
	"step.salesforce_approval_submit": func(n string, c map[string]any) (sdk.StepInstance, error) { return newApprovalSubmitStep(n, c) },
	"step.salesforce_approval_approve": func(n string, c map[string]any) (sdk.StepInstance, error) { return newApprovalApproveStep(n, c) },
	"step.salesforce_approval_reject": func(n string, c map[string]any) (sdk.StepInstance, error) { return newApprovalRejectStep(n, c) },

	// Chatter
	"step.salesforce_chatter_post":      func(n string, c map[string]any) (sdk.StepInstance, error) { return newChatterPostStep(n, c) },
	"step.salesforce_chatter_comment":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newChatterCommentStep(n, c) },
	"step.salesforce_chatter_like":      func(n string, c map[string]any) (sdk.StepInstance, error) { return newChatterLikeStep(n, c) },
	"step.salesforce_chatter_feed_list": func(n string, c map[string]any) (sdk.StepInstance, error) { return newChatterFeedListStep(n, c) },

	// Files & Content
	"step.salesforce_file_upload":              func(n string, c map[string]any) (sdk.StepInstance, error) { return newFileUploadStep(n, c) },
	"step.salesforce_file_download":            func(n string, c map[string]any) (sdk.StepInstance, error) { return newFileDownloadStep(n, c) },
	"step.salesforce_content_version_create":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newContentVersionCreateStep(n, c) },
	"step.salesforce_content_document_get":     func(n string, c map[string]any) (sdk.StepInstance, error) { return newContentDocumentGetStep(n, c) },
	"step.salesforce_content_document_delete":  func(n string, c map[string]any) (sdk.StepInstance, error) { return newContentDocumentDeleteStep(n, c) },

	// Users & Identity
	"step.salesforce_user_get":    func(n string, c map[string]any) (sdk.StepInstance, error) { return newUserGetStep(n, c) },
	"step.salesforce_user_list":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newUserListStep(n, c) },
	"step.salesforce_user_create": func(n string, c map[string]any) (sdk.StepInstance, error) { return newUserCreateStep(n, c) },
	"step.salesforce_user_update": func(n string, c map[string]any) (sdk.StepInstance, error) { return newUserUpdateStep(n, c) },
	"step.salesforce_identity_get": func(n string, c map[string]any) (sdk.StepInstance, error) { return newIdentityGetStep(n, c) },
	"step.salesforce_org_limits":  func(n string, c map[string]any) (sdk.StepInstance, error) { return newOrgLimitsStep(n, c) },

	// Flows
	"step.salesforce_flow_list": func(n string, c map[string]any) (sdk.StepInstance, error) { return newFlowListStep(n, c) },
	"step.salesforce_flow_run":  func(n string, c map[string]any) (sdk.StepInstance, error) { return newFlowRunStep(n, c) },

	// Platform Events
	"step.salesforce_event_publish": func(n string, c map[string]any) (sdk.StepInstance, error) { return newEventPublishStep(n, c) },

	// Metadata
	"step.salesforce_metadata_describe": func(n string, c map[string]any) (sdk.StepInstance, error) { return newMetadataDescribeStep(n, c) },
	"step.salesforce_metadata_list":     func(n string, c map[string]any) (sdk.StepInstance, error) { return newMetadataListStep(n, c) },
	"step.salesforce_metadata_read":     func(n string, c map[string]any) (sdk.StepInstance, error) { return newMetadataReadStep(n, c) },
	"step.salesforce_metadata_create":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newMetadataCreateStep(n, c) },
	"step.salesforce_metadata_update":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newMetadataUpdateStep(n, c) },
	"step.salesforce_metadata_delete":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newMetadataDeleteStep(n, c) },
	"step.salesforce_metadata_deploy":   func(n string, c map[string]any) (sdk.StepInstance, error) { return newMetadataDeployStep(n, c) },
	"step.salesforce_metadata_retrieve": func(n string, c map[string]any) (sdk.StepInstance, error) { return newMetadataRetrieveStep(n, c) },

	// Generic
	"step.salesforce_raw_request": func(n string, c map[string]any) (sdk.StepInstance, error) { return newRawRequestStep(n, c) },
}

// createStep dispatches to the appropriate step constructor.
func createStep(typeName, name string, config map[string]any) (sdk.StepInstance, error) {
	constructor, ok := stepRegistry[typeName]
	if !ok {
		return nil, fmt.Errorf("salesforce plugin: unknown step type %q", typeName)
	}
	return constructor(name, config)
}

// allStepTypes returns all registered step type strings.
func allStepTypes() []string {
	types := make([]string, 0, len(stepRegistry))
	for k := range stepRegistry {
		types = append(types, k)
	}
	return types
}
