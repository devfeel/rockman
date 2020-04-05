package _json

import (
	"encoding/json"
)

// GetJsonString marshals the object as string
func GetJsonString(obj interface{}) string {
	resByte, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(resByte)
}

// Marshal marshals the value as string
func Marshal(v interface{}) (string, error) {
	resByte, err := json.Marshal(v)
	if err != nil {
		return "", err
	} else {
		return string(resByte), nil
	}
}

// Unmarshal parses the JSON-encoded data and stores the result
func Unmarshal(value string, v interface{}) error {
	return json.Unmarshal([]byte(value), v)
}
