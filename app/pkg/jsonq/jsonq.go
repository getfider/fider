package jsonq

import (
	"encoding/json"
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
	var str string
	err := json.Unmarshal(*q.m[key], &str)
	return str, err
}

func (q *Query) Contains(key string) bool {
	return q.m[key] != nil
}
