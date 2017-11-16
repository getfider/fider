package jsonq

import (
	"encoding/json"
	"strings"
)

type Query struct {
	m map[string]*json.RawMessage
}

func New(content string) *Query {
	var m map[string]*json.RawMessage
	err := json.Unmarshal([]byte(content), &m)
	if err != nil {
		panic(err)
	}
	return &Query{m: m}
}

func (q *Query) String(key string) (string, error) {
	data := q.get(key)
	if data != nil {
		var str string
		err := json.Unmarshal(*data, &str)
		return str, err
	}
	return "", nil
}

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
