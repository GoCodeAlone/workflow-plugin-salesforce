package main

import (
	"github.com/GoCodeAlone/workflow-plugin-salesforce/internal"
	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

var version = "dev"

func main() {
	sdk.Serve(internal.NewSalesforcePlugin())
}
