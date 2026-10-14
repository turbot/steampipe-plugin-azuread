package azuread

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/microsoft/kiota-abstractions-go/authentication"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

func TestIterateDirectoryAuditsKeepsRowsPastFirstPage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"value":[
			{"id":"page2-a","activityDisplayName":"Delete user","result":"success"},
			{"id":"page2-b","activityDisplayName":"Add app role assignment","result":"failure"}
		]}`)
	}))
	defer server.Close()

	firstPagePayload := fmt.Sprintf(`{"@odata.nextLink":"%s/v1.0/auditLogs/directoryAudits?$skiptoken=abc","value":[
		{"id":"page1-a","activityDisplayName":"Update user","result":"success"}
	]}`, server.URL)
	node, err := jsonserialization.NewJsonParseNode([]byte(firstPagePayload))
	if err != nil {
		t.Fatalf("parse node: %v", err)
	}
	parsed, err := node.GetObjectValue(models.CreateDirectoryAuditCollectionResponseFromDiscriminatorValue)
	if err != nil {
		t.Fatalf("decode first page: %v", err)
	}

	adapter, err := msgraphsdkgo.NewGraphRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient(
		&authentication.AnonymousAuthenticationProvider{},
		jsonserialization.NewJsonParseNodeFactory(),
		jsonserialization.NewJsonSerializationWriterFactory(),
		server.Client(),
	)
	if err != nil {
		t.Fatalf("adapter: %v", err)
	}

	var ids []string
	err = iterateDirectoryAudits(context.Background(), parsed.(models.DirectoryAuditCollectionResponseable), adapter, func(audit models.DirectoryAuditable) bool {
		ids = append(ids, *audit.GetId())
		return true
	})
	if err != nil {
		t.Fatalf("iterate: %v", err)
	}

	want := []string{"page1-a", "page2-a", "page2-b"}
	if fmt.Sprint(ids) != fmt.Sprint(want) {
		t.Errorf("iterated ids = %v, want %v", ids, want)
	}
}
