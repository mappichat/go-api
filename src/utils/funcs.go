package utils

import (
	"github.com/iancoleman/strcase"
	"github.com/mitchellh/mapstructure"
)

func DecodeSnakeCase(input interface{}) (map[string]interface{}, error) {
	output := map[string]interface{}{}
	if err := mapstructure.Decode(input, &output); err != nil {
		return nil, err
	}
	newOut := map[string]interface{}{}
	for k, v := range output {
		newOut[strcase.ToSnake(k)] = v
	}
	return newOut, nil
}
