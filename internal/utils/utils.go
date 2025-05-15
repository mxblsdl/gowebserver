package utils

import (
	"encoding/json"
	"html/template"
)

func SafeJSON(data interface{}) (template.JS, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "{}", err
	}
	return template.JS(jsonBytes), nil
}
