package azuread

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

// Constants for Standard Column Descriptions
const (
	ColumnDescriptionTenant = "The Azure Tenant ID where the resource is located."
	ColumnDescriptionTags   = "A map of tags for the resource."
	ColumnDescriptionTitle  = "Title of the resource."
)

func isNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "Request_ResourceNotFound")
}

func isNotFoundErrorPredicate(notFoundErrors []string) plugin.ErrorPredicate {
	return func(err error) bool {
		if err != nil {
			for _, item := range notFoundErrors {
				if strings.Contains(err.Error(), item) {
					return true
				}
			}
		}
		return false
	}
}

func TagsToMap(tags []string) (*map[string]bool, error) {
	var turbotTagsMap map[string]bool
	if tags == nil {
		return nil, nil
	}

	turbotTagsMap = map[string]bool{}
	for _, i := range tags {
		turbotTagsMap[i] = true
	}

	return &turbotTagsMap, nil
}

func getListValues(listValue *proto.QualValueList) []string {
	values := make([]string, 0)
	for _, value := range listValue.Values {
		if strings.TrimRight(strings.TrimLeft(value.GetStringValue(), " "), " ") != "" {
			values = append(values, value.GetStringValue())
		}
	}
	return values
}

func getQualsValueByColumn(quals plugin.KeyColumnQualMap, columnName string, dataType string) interface{} {
	var value interface{}
	for _, q := range quals[columnName].Quals {
		if dataType == "string" {
			if q.Value.GetStringValue() != "" {
				value = q.Value.GetStringValue()
			} else {
				value = getListValues(q.Value.GetListValue())
			}
		}
		if dataType == "bool" {
			switch q.Operator {
			case "<>":
				value = !q.Value.GetBoolValue()
			case "=":
				value = q.Value.GetBoolValue()
			}
		}
		if dataType == "int64" {
			value = q.Value.GetInt64Value()
			if q.Value.GetInt64Value() == 0 {
				valueSlice := make([]*string, 0)
				for _, value := range q.Value.GetListValue().Values {
					val := strconv.FormatInt(value.GetInt64Value(), 10)
					valueSlice = append(valueSlice, &val)
				}
				value = valueSlice
			}
		}
		if dataType == "double" {
			value = q.Value.GetDoubleValue()
			if q.Value.GetDoubleValue() == 0 {
				valueSlice := make([]*string, 0)
				for _, value := range q.Value.GetListValue().Values {
					val := strconv.FormatFloat(value.GetDoubleValue(), 'f', 4, 64)
					valueSlice = append(valueSlice, &val)
				}
				value = valueSlice
			}

		}
		if dataType == "ipaddr" {
			value = q.Value.GetInetValue().Addr
			if q.Value.GetInetValue().Addr == "" {
				valueSlice := make([]*string, 0)
				for _, value := range q.Value.GetListValue().Values {
					val := value.GetInetValue().Addr
					valueSlice = append(valueSlice, &val)
				}
				value = valueSlice
			}
		}
		if dataType == "cidr" {
			value = q.Value.GetInetValue().Cidr
			if q.Value.GetInetValue().Addr == "" {
				valueSlice := make([]*string, 0)
				for _, value := range q.Value.GetListValue().Values {
					val := value.GetInetValue().Cidr
					valueSlice = append(valueSlice, &val)
				}
				value = valueSlice
			}
		}
		if dataType == "time" {
			value = getListValues(q.Value.GetListValue())
			if len(getListValues(q.Value.GetListValue())) == 0 {
				value = q.Value.GetTimestampValue().AsTime()
			}
		}
	}
	return value
}

type QualsColumn struct {
	ColumnName string
	ColumnType string
	FilterName string
}

//// Build common query filter

func buildCommaonQueryFilter(qualsColumns []QualsColumn, quals plugin.KeyColumnQualMap) []string {
	var filter []string
	for _, qualColumn := range qualsColumns {
		if quals[qualColumn.ColumnName] != nil {
			value := getQualsValueByColumn(quals, qualColumn.ColumnName, qualColumn.ColumnType)
			switch qualColumn.ColumnType {
			case "string":
				val, ok := value.([]string)
				if ok {
					var valueSlice []string
					for _, v := range val {
						valueSlice = append(valueSlice, fmt.Sprintf("%s eq '%s'", qualColumn.FilterName, v))
					}
					filter = append(filter, strings.Join(valueSlice, " or "))
				} else {
					val := value.(string)
					filter = append(filter, fmt.Sprintf("%s eq '%s'", qualColumn.FilterName, val))
				}
			case "bool":
				val := value.(bool)
				filter = append(filter, fmt.Sprintf("%s eq %t", qualColumn.FilterName, val))
			}
		}

	}
	return filter
}
