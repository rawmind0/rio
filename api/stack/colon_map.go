package stack

import (
	"fmt"

	"github.com/rancher/norman/types"
)

type MapToSlice struct {
	Field string
	Sep   string
}

func (d MapToSlice) FromInternal(data map[string]interface{}) {
}

func (d MapToSlice) ToInternal(data map[string]interface{}) error {
	v, ok := data[d.Field]
	if !ok {
		return nil
	}

	if m, ok := v.(map[string]interface{}); ok {
		var result []string
		for k, v := range m {
			result = append(result, fmt.Sprintf("%s%s%v", k, d.Sep, v))
		}
		data[d.Field] = result
	}

	return nil
}

func (MapToSlice) ModifySchema(schema *types.Schema, schemas *types.Schemas) error {
	return nil
}
