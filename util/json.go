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

func (json *JSONObject) GetString(key string) interface{} {
	return json.data.(map[string]interface{})[key]
}
