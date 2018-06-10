package assert_test

import (
	"sync"
	"testing"
	"time"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/errors"
)

func TestBoolEquals(t *testing.T) {
	mockT := new(testing.T)
	RegisterT(mockT)

	var testCases = []struct {
		in1    bool
		in2    bool
		equals bool
	}{
		{true, true, true},
		{false, false, true},
		{false, true, false},
		{true, false, false},
	}

	for _, tt := range testCases {
		if Expect(tt.in1).Equals(tt.in2) == !tt.equals {
			t.Error("Equals assertion failed")
		}
	}
}

func TestBoolIsTrueFalse(t *testing.T) {
	mockT := new(testing.T)
	RegisterT(mockT)

	if Expect(false).IsTrue() {
		t.Error("IsTrue assertion failed")
	}

	if !Expect(true).IsTrue() {
		t.Error("IsTrue assertion failed")
	}

	if !Expect(false).IsFalse() {
		t.Error("IsFalse assertion failed")
	}

	if Expect(true).IsFalse() {
		t.Error("IsFalse assertion failed")
	}
}

func TestError(t *testing.T) {
	mockT := new(testing.T)
	RegisterT(mockT)

	notFound := errors.New("Not Found")
	noRecords := errors.New("No Records")

	if Expect(notFound).IsNil() {
		t.Error("IsNil assertion failed")
	}

	if !Expect(notFound).IsNotNil() {
		t.Error("IsNotNil assertion failed")
	}

	if !Expect(nil).IsNil() {
		t.Error("IsNil assertion failed")
	}

	if Expect(nil).IsNotNil() {
		t.Error("IsNotNil assertion failed")
	}

	var testCases = []struct {
		in1    error
		in2    error
		equals bool
	}{
		{notFound, notFound, true},
		{notFound, noRecords, false},
		{notFound, nil, false},
		{noRecords, nil, false},
		{nil, notFound, false},
		{nil, noRecords, false},
	}

	for _, tt := range testCases {
		if Expect(tt.in1).Equals(tt.in2) == !tt.equals {
			t.Error("Equals assertion failed")
		}
		if Expect(tt.in1).NotEquals(tt.in2) == tt.equals {
			t.Error("NotEquals assertion failed")
		}
	}
}

func TestFuncPanics(t *testing.T) {
	mockT := new(testing.T)
	RegisterT(mockT)

	if !Expect(func() {
		panic("Boom!")
	}).Panics() {
		t.Error("Panics assertion failed")
	}

	if Expect(func() {}).Panics() {
		t.Error("Panics assertion failed")
	}
}

func TestFuncEventuallyEquals(t *testing.T) {
	mockT := new(testing.T)
	RegisterT(mockT)
	mu := &sync.RWMutex{}

	completed := false
	go func() {
		time.Sleep(500 * time.Millisecond)
		mu.Lock()
		defer mu.Unlock()
		completed = true
	}()

	if !Expect(func() bool {
		mu.RLock()
		defer mu.RUnlock()
		return completed
	}).EventuallyEquals(true) {
		t.Error("EventuallyEquals assertion failed")
	}
}

func TestIntEquals(t *testing.T) {
	mockT := new(testing.T)
	RegisterT(mockT)

	var testCases = []struct {
		in1    int
		in2    int
		equals bool
	}{
		{3, 4, false},
		{243, 243, true},
		{214, 2140, false},
		{0, 0, true},
		{-9, -9, true},
	}

	for _, tt := range testCases {
		if Expect(tt.in1).Equals(tt.in2) == !tt.equals {
			t.Error("Equals assertion failed")
		}
		if Expect(tt.in1).NotEquals(tt.in2) == tt.equals {
			t.Error("NotEquals assertion failed")
		}
	}
}

func TestMapLength(t *testing.T) {
	mockT := new(testing.T)
	RegisterT(mockT)

	if !Expect(map[string]string{}).HasLen(0) {
		t.Error("HasLen assertion failed")
	}

	if !Expect(map[string]string{
		"Key": "Value",
	}).HasLen(1) {
		t.Error("HasLen assertion failed")
	}

	if Expect(map[string]string{
		"Key1": "Value1",
		"Key2": "Value2",
		"Key3": "Value3",
	}).HasLen(2) {
		t.Error("HasLen assertion failed")
	}
}

func TestMapEquals(t *testing.T) {
	mockT := new(testing.T)
	RegisterT(mockT)

	if !Expect(map[string]string{
		"Key1": "Value1",
		"Key2": "Value2",
		"Key3": "Value3",
	}).Equals(map[string]string{
		"Key1": "Value1",
		"Key2": "Value2",
		"Key3": "Value3",
	}) {
		t.Error("HasLen assertion failed")
	}

	if Expect(map[string]string{
		"Key1": "Value1",
		"Key2": "Value2",
		"Key3": "Value3",
	}).Equals(map[string]string{
		"Key1": "Value1",
		"Key2": "Value2222",
		"Key3": "Value3",
	}) {
		t.Error("HasLen assertion failed")
	}
}

func TestSliceLength(t *testing.T) {
	mockT := new(testing.T)
	RegisterT(mockT)

	if !Expect([]string{}).HasLen(0) {
		t.Error("HasLen assertion failed")
	}

	if Expect([]string{}).HasLen(1) {
		t.Error("HasLen assertion failed")
	}

	if !Expect([]string{"A", "B"}).HasLen(2) {
		t.Error("HasLen assertion failed")
	}

	if !Expect([4]int{}).HasLen(4) {
		t.Error("HasLen assertion failed")
	}
}

func TestSliceEquals(t *testing.T) {
	mockT := new(testing.T)
	RegisterT(mockT)

	if !Expect([]string{}).Equals([]string{}) {
		t.Error("Equals assertion failed")
	}

	if !Expect([]string{"A", "B"}).Equals([]string{"A", "B"}) {
		t.Error("Equals assertion failed")
	}

	if Expect([]string{"B", "A"}).Equals([]string{"A", "B"}) {
		t.Error("Equals assertion failed")
	}

	if Expect([]string{"B"}).Equals([]string{"A", "B"}) {
		t.Error("Equals assertion failed")
	}

	if Expect([]byte{23}).Equals([]string{"A", "B"}) {
		t.Error("Equals assertion failed")
	}

	if Expect([]byte{23}).Equals([]byte{43}) {
		t.Error("Equals assertion failed")
	}

	if !Expect([]byte{23}).Equals([]byte{23}) {
		t.Error("Equals assertion failed")
	}
}

func TestSliceIsNilOrNot(t *testing.T) {
	mockT := new(testing.T)
	RegisterT(mockT)

	type foo struct{}

	emptyFoo := []foo{}

	if Expect(emptyFoo).IsNil() {
		t.Error("IsNil assertion failed")
	}

	if !Expect(emptyFoo).IsNotNil() {
		t.Error("IsNotNil assertion failed")
	}

	var nilFoo []foo = nil

	if !Expect(nilFoo).IsNil() {
		t.Error("IsNil assertion failed")
	}

	if Expect(nilFoo).IsNotNil() {
		t.Error("IsNotNil assertion failed")
	}
}

func TestStringEquals(t *testing.T) {
	mockT := new(testing.T)
	RegisterT(mockT)

	var testCases = []struct {
		in1    string
		in2    string
		equals bool
	}{
		{"Hello", "World", false},
		{"Hello", "Hello", true},
		{"Hello", "Hello ", false},
		{"", "", true},
	}

	for _, tt := range testCases {
		if Expect(tt.in1).Equals(tt.in2) == !tt.equals {
			t.Error("Equals assertion failed")
		}
		if Expect(tt.in1).NotEquals(tt.in2) == tt.equals {
			t.Error("NotEquals assertion failed")
		}
	}
}

func TestStringEmpty(t *testing.T) {
	mockT := new(testing.T)
	RegisterT(mockT)

	var testCases = []struct {
		in1   string
		empty bool
	}{
		{"", true},
		{"Hello", false},
		{"  ", false},
	}

	for _, tt := range testCases {
		if Expect(tt.in1).IsEmpty() == !tt.empty {
			t.Error("IsEmpty assertion failed")
		}
	}
}

func TestStringNotEmpty(t *testing.T) {
	mockT := new(testing.T)
	RegisterT(mockT)

	var testCases = []struct {
		in1      string
		notempty bool
	}{
		{"", false},
		{"Hello", true},
		{"  ", true},
	}

	for _, tt := range testCases {
		if Expect(tt.in1).IsNotEmpty() == !tt.notempty {
			t.Error("IsNotEmpty assertion failed")
		}
	}
}

func TestStringLength(t *testing.T) {
	mockT := new(testing.T)
	RegisterT(mockT)

	if !Expect("Hello World").HasLen(11) {
		t.Error("HasLen assertion failed")
	}

	if Expect("Hello World").HasLen(5) {
		t.Error("HasLen assertion failed")
	}

	if Expect("").HasLen(5) {
		t.Error("HasLen assertion failed")
	}

	if !Expect("").HasLen(0) {
		t.Error("HasLen assertion failed")
	}
}

func TestStringContainsSubstring(t *testing.T) {
	mockT := new(testing.T)
	RegisterT(mockT)

	if !Expect("Hello World").ContainsSubstring("World") {
		t.Error("Contains assertion failed")
	}

	if Expect("Hello World").ContainsSubstring("John") {
		t.Error("Contains assertion failed")
	}
}

func TestStructEquals(t *testing.T) {
	mockT := new(testing.T)
	RegisterT(mockT)

	type foo struct {
		name string
		age  int
	}

	first := &foo{name: "Bob", age: 20}
	second := &foo{name: "John", age: 25}
	third := &foo{name: "Bob", age: 20}

	if Expect(first).Equals(second) {
		t.Error("Equals assertion failed")
	}

	if Expect(second).Equals(third) {
		t.Error("Equals assertion failed")
	}

	if !Expect(first).Equals(third) {
		t.Error("Equals assertion failed")
	}

	if Expect(map[string]string{
		"Key1": "Value1",
		"Key2": "Value2",
	}).Equals(third) {
		t.Error("Equals assertion failed")
	}

	if !Expect(map[string]string{
		"Key3": "Value3",
		"Key4": "Value4",
	}).Equals(map[string]string{
		"Key3": "Value3",
		"Key4": "Value4",
	}) {
		t.Error("Equals assertion failed")
	}
}

func TestStructIsNilOrNot(t *testing.T) {
	mockT := new(testing.T)
	RegisterT(mockT)

	type foo struct{}

	emptyFoo := &foo{}

	if Expect(emptyFoo).IsNil() {
		t.Error("IsNil assertion failed")
	}

	if !Expect(emptyFoo).IsNotNil() {
		t.Error("IsNotNil assertion failed")
	}

	var nilFoo *foo = nil

	if !Expect(nilFoo).IsNil() {
		t.Error("IsNil assertion failed")
	}

	if Expect(nilFoo).IsNotNil() {
		t.Error("IsNotNil assertion failed")
	}
}

func TestTimeSimilar(t *testing.T) {
	mockT := new(testing.T)
	RegisterT(mockT)

	base := time.Date(2010, 5, 3, 1, 10, 34, 0, time.UTC)
	time1 := base.Add(5 * time.Second)

	if !Expect(base).TemporarilySimilar(time1, 10*time.Second) {
		t.Error("Similar assertion failed")
	}

	if Expect(base).TemporarilySimilar(time1, time.Second) {
		t.Error("Similar assertion failed")
	}
}
