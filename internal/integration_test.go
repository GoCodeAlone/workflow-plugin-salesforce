package internal_test

import (
	"testing"

	"github.com/GoCodeAlone/workflow/wftest"
)

// TestIntegration_QueryRecords verifies that a pipeline using salesforce_query
// passes results through to downstream steps.
func TestIntegration_QueryRecords(t *testing.T) {
	h := wftest.New(t, wftest.WithYAML(`
pipelines:
  query:
    steps:
      - name: fetch
        type: step.salesforce_query
        config:
          soql: "SELECT Id, Name FROM Account LIMIT 1"
      - name: confirm
        type: step.set
        config:
          values:
            queried: true
`),
		wftest.MockStep("step.salesforce_query", wftest.Returns(map[string]any{
			"records":   []any{map[string]any{"Id": "001xx", "Name": "Acme"}},
			"totalSize": 1,
		})),
	)

	result := h.ExecutePipeline("query", nil)
	if result.Error != nil {
		t.Fatal(result.Error)
	}
	if result.Output["queried"] != true {
		t.Errorf("expected queried=true, got %v", result.Output["queried"])
	}
}

// TestIntegration_CreateRecord verifies that a pipeline using salesforce_record_create
// surfaces the returned record ID downstream.
func TestIntegration_CreateRecord(t *testing.T) {
	h := wftest.New(t, wftest.WithYAML(`
pipelines:
  create:
    steps:
      - name: create_account
        type: step.salesforce_record_create
        config:
          sobject: Account
          fields:
            Name: "Test Corp"
      - name: mark_done
        type: step.set
        config:
          values:
            created: true
`),
		wftest.MockStep("step.salesforce_record_create", wftest.Returns(map[string]any{
			"id":      "001xx000003GYn1AAG",
			"success": true,
		})),
	)

	result := h.ExecutePipeline("create", nil)
	if result.Error != nil {
		t.Fatal(result.Error)
	}
	if result.Output["created"] != true {
		t.Errorf("expected created=true, got %v", result.Output["created"])
	}
	// Verify the step output is accessible in StepResults.
	createOut := result.StepResults["create_account"]
	if createOut == nil {
		t.Fatal("expected step result for create_account")
	}
	if createOut["success"] != true {
		t.Errorf("expected success=true in step result, got %v", createOut["success"])
	}
}

// TestIntegration_BulkInsert verifies that a pipeline using salesforce_bulk_insert
// records call arguments via a Recorder and confirms pipeline completes successfully.
func TestIntegration_BulkInsert(t *testing.T) {
	rec := wftest.RecordStep("step.salesforce_bulk_insert")
	rec.WithOutput(map[string]any{
		"jobId":  "750xx000000blbXAAQ",
		"state":  "JobComplete",
		"failed": 0,
	})

	h := wftest.New(t, wftest.WithYAML(`
pipelines:
  bulk:
    steps:
      - name: bulk_load
        type: step.salesforce_bulk_insert
        config:
          sobject: Contact
          records:
            - FirstName: Alice
              LastName: Smith
            - FirstName: Bob
              LastName: Jones
      - name: summarize
        type: step.set
        config:
          values:
            bulk_done: true
`),
		rec,
	)

	result := h.ExecutePipeline("bulk", nil)
	if result.Error != nil {
		t.Fatal(result.Error)
	}
	if result.Output["bulk_done"] != true {
		t.Errorf("expected bulk_done=true, got %v", result.Output["bulk_done"])
	}
	if rec.CallCount() != 1 {
		t.Errorf("expected bulk_insert to be called once, got %d", rec.CallCount())
	}
	calls := rec.Calls()
	cfg := calls[0].Config
	if cfg["sobject"] != "Contact" {
		t.Errorf("expected sobject=Contact in step config, got %v", cfg["sobject"])
	}
}
