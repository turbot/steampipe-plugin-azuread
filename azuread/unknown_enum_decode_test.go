package azuread

import (
	"testing"

	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

// Graph returns enum members newer than the SDK's generated enum (e.g. operationResult
// "clientError"). Those must decode to a nil field, not fail the whole page.
func TestDirectoryAuditPageSurvivesUnknownOperationResult(t *testing.T) {
	payload := `{"value":[
		{"id":"a1","activityDisplayName":"Update user","result":"clientError","resultReason":"Bad request"},
		{"id":"a2","activityDisplayName":"Add member to group","result":"success"}
	]}`
	node, err := jsonserialization.NewJsonParseNode([]byte(payload))
	if err != nil {
		t.Fatalf("parse node: %v", err)
	}
	parsed, err := node.GetObjectValue(models.CreateDirectoryAuditCollectionResponseFromDiscriminatorValue)
	if err != nil {
		t.Fatalf("decode directoryAudit page: %v", err)
	}
	audits := parsed.(models.DirectoryAuditCollectionResponseable).GetValue()
	if len(audits) != 2 {
		t.Fatalf("decoded %d audits, want 2", len(audits))
	}

	unknown := &ADDirectoryAuditReportInfo{audits[0]}
	if got := unknown.DirectoryAuditResult(); got != "" {
		t.Errorf("unknown result = %q, want empty", got)
	}
	if got := unknown.GetResultReason(); got == nil || *got != "Bad request" {
		t.Errorf("resultReason = %v, want Bad request", got)
	}
	if got := (&ADDirectoryAuditReportInfo{audits[1]}).DirectoryAuditResult(); got != "success" {
		t.Errorf("known result = %q, want success", got)
	}
}
