package util

import "encoding/json"

// JSONObject as an object
type JSONObject struct {
	data interface{}
}

// NewJSONObject creates a new JSONObject
func NewJSONObject(bytes []byte) JSONObject {
	var data interface{}
	json.Unmarshal(bytes, &data)
	return JSONObject{data: data}
}

func (json *JSONObject) getValue(key string) interface{} {
	return json.data.(map[string]interface{})[key]
}

// GetString retrives a string value from given key
func (json *JSONObject) GetString(key string) string {
	return json.getValue(key).(string)
}

// GetBoolean retrives a boolean value from given key
func (json *JSONObject) GetBoolean(key string) bool {
	return json.getValue(key).(bool)
}

// GetJSON retrives a JSON object from given key
func (json *JSONObject) GetJSON(key string) *JSONObject {
	return &JSONObject{data: json.getValue(key)}
}
