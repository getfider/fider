package jsonq

import (
	"encoding/json"
	"strings"
)

//Query is a JSON query interface
type Query struct {
	json string
	m    map[string]*json.RawMessage
}

//New creates a new Query object based on given input as a JSON
func New(content string) *Query {
	var m map[string]*json.RawMessage
	if len(content) > 0 && string(content[0]) == "{" {
		err := json.Unmarshal([]byte(content), &m)
		if err != nil {
			panic(err)
		}
	}
	return &Query{
		m:    m,
		json: content,
	}
}

//String returns a string value from the json object based on its key
func (q *Query) String(key string) (string, error) {
	data := q.get(key)
	if data != nil {
		var str string
		err := json.Unmarshal(*data, &str)
		return str, err
	}
	return "", nil
}

//Int32 returns a integer value from the json object based on its key
func (q *Query) Int32(key string) (int, error) {
	data := q.get(key)
	if data != nil {
		var num int
		err := json.Unmarshal(*data, &num)
		return num, err
	}
	return 0, nil
}

//IsArray returns true if the json object is an array
func (q *Query) IsArray() bool {
	return q.m == nil
}

//ArrayLength returns number of elements in the array
func (q *Query) ArrayLength() int {
	if q.IsArray() {
		var arr []interface{}
		json.Unmarshal([]byte(q.json), &arr)
		return len(arr)
	}
	return 0
}

//Contains returns true if the json object has the key
func (q *Query) Contains(key string) bool {
	return q.get(key) != nil
}

func (q *Query) get(key string) *json.RawMessage {
	keys := strings.Split(key, ".")

	var message *json.RawMessage
	var m map[string]*json.RawMessage

	for _, key := range keys {
		if m != nil {
			message = m[key]
		} else {
			message = q.m[key]
		}

		if message != nil {
			bytes, _ := message.MarshalJSON()
			json.Unmarshal(bytes, &m)
		} else {
			return nil
		}
	}

	return message
}
