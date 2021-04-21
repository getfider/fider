package jsonq

import (
	"encoding/json"
	"strconv"
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
	if len(content) > 0 && content[:1] == "{" {
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

//String returns a string value from the json object based on its selector
func (q *Query) String(selector string) string {
	selectors := strings.Split(selector, ",")

	for _, s := range selectors {
		data := q.get(strings.TrimSpace(s))
		if data != nil {
			var str string
			if err := json.Unmarshal(*data, &str); err == nil && str != "" {
				return str
			}

			var num int
			if err := json.Unmarshal(*data, &num); err == nil {
				return strconv.Itoa(num)
			}
		}
	}

	return ""
}

//Int32 returns a integer value from the json object based on its selector
func (q *Query) Int32(selector string) int {
	data := q.get(selector)
	if data != nil {
		var num int
		err := json.Unmarshal(*data, &num)
		if err != nil {
			panic(err)
		}
		return num
	}
	return 0
}

//IsArray returns true if the json object is an array
func (q *Query) IsArray() bool {
	return q.m == nil
}

//ArrayLength returns number of elements in the array
func (q *Query) ArrayLength() int {
	if q.IsArray() {
		var arr []interface{}
		err := json.Unmarshal([]byte(q.json), &arr)
		if err != nil {
			panic(err)
		}
		return len(arr)
	}
	return 0
}

//Contains returns true if the json object has the key
func (q *Query) Contains(selector string) bool {
	return q.get(selector) != nil
}

func (q *Query) get(selector string) *json.RawMessage {
	if selector == "" {
		return nil
	}

	parts := strings.Split(selector, ".")

	var result *json.RawMessage
	var current map[string]*json.RawMessage

	for _, part := range parts {
		if part[len(part)-1:] == "]" {
			idx, _ := strconv.Atoi(part[len(part)-2 : len(part)-1])
			part = part[:len(part)-3]

			if current[part] != nil {
				result = current[part]
			} else {
				result = q.m[part]
			}
			if result != nil {
				var arr []*json.RawMessage
				bytes, _ := result.MarshalJSON()
				err := json.Unmarshal(bytes, &arr)
				if err != nil {
					panic(err)
				}

				if len(arr) > idx {
					bytes, _ = arr[idx].MarshalJSON()
					err := json.Unmarshal(bytes, &current)
					if err != nil {
						return arr[idx]
					}
				}
			}
		} else {

			if current[part] != nil {
				result = current[part]
			} else {
				result = q.m[part]
			}

			if result != nil {
				bytes, _ := result.MarshalJSON()
				_ = json.Unmarshal(bytes, &current)
			} else {
				return nil
			}

		}
	}

	return result
}
