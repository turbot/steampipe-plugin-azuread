package azuread

import (
	"context"
	"encoding/json"
	"strings"
	"strconv"

	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type RequestError struct {
	Code    string
	Message string
}

func (m *RequestError) Error() string {
	errStr, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(errStr)
}

func getErrorObject(err error) *RequestError {
	switch err := err.(type) {
	case *odataerrors.ODataError:
		terr := err.GetErrorEscaped()
		return &RequestError{
			Code:    *terr.GetCode(),
			Message: *terr.GetMessage(),
		}
	case *azidentity.AuthenticationFailedError:
		return &RequestError{
			Code:    strconv.Itoa(err.RawResponse.StatusCode),
			Message: err.Error(),
		}
	default:
		// If the error type is unknown
		// return the exact error
		return &RequestError{
			Message: err.Error(),
		}
	}
}

func isIgnorableErrorPredicate(ignoreErrorCodes []string) plugin.ErrorPredicateWithContext {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, err error) bool {
		if err != nil {
			if terr, ok := err.(*RequestError); ok {
				for _, item := range ignoreErrorCodes {
					if terr != nil && (terr.Code == item || strings.Contains(terr.Message, item)) {
						return true
					}
				}
			}
		}
		return false
	}
}