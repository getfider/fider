package dto_test

import (
	"testing"

	"github.com/Spicy-Bush/fider-tarkov-community/app/models/dto"

	. "github.com/Spicy-Bush/fider-tarkov-community/app/pkg/assert"
)

func TestPropsMerge(t *testing.T) {
	RegisterT(t)

	p1 := dto.Props{
		"name": "Jon",
		"age":  26,
	}
	p2 := p1.Merge(dto.Props{
		"age":   30,
		"email": "john.snow@got.com",
	})
	Expect(p2).Equals(dto.Props{
		"name":  "Jon",
		"age":   30,
		"email": "john.snow@got.com",
	})
}
