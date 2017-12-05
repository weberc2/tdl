package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Copied from encoding/json source code
func isValidTag(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		switch {
		case strings.ContainsRune("!#$%&()*+-./:<=>?@[]^_{|}~ ", c):
			// Backslash and quote chars are reserved, but
			// otherwise any punctuation chars are allowed
			// in a tag name.
		default:
			if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
				return false
			}
		}
	}
	return true
}

// FieldName is a JSON field name. It differs from a Go string in that a
// FieldName should be constrained to the subset of strings that are valid JSON
// field names. It should be created with NewFieldName() or MustFieldName() to
// guarantee that the string is valid. This is why FieldName is a struct
// instead of `type FieldName string`--the latter would make it too easy to
// create an invalid FieldName: `FieldName("")`.
//
// By making it a struct with a private Go string member, it's harder to
// construct an invalid String, at least from other packages (e.g.,
// FieldName{"\"} is a compile error). Obviously we can't control the
// zero-value case, but this is better than nothing.
type FieldName struct {
	s String
}

// MustFieldName constructs a new FieldName from a Go string like NewFieldName,
// but MustFieldName panics if `name` is invalid.
func MustFieldName(name string) FieldName {
	fname, err := NewFieldName(name)
	if err != nil {
		panic("Illegal field name: " + name)
	}
	return fname
}

// NewFieldName creates a new FieldName from a Go string. This is the preferred
// way to construct a FieldName as it performs the escaping and validation
// required to guarantee that the string is legal.
func NewFieldName(name string) (FieldName, error) {
	if !isValidTag(name) {
		return FieldName{}, fmt.Errorf(
			"Invalid JSON field name: '%s'",
			name,
		)
	}
	return FieldName{NewString(name)}, nil
}

func (fn FieldName) String() string { return fn.s.String() }

// String is a JSON string. It differs from a Go string in that a String should
// be constrained to the subset of strings that are valid JSON strings. It
// should be created with NewString() to guarantee that the string is valid.
// This is why String is a struct instead of `type String string`--the latter
// would make it too easy to create an invalid String: `String("\")`.
//
// By making it a struct with a private Go string member, it's harder to
// construct an invalid String, at least from other packages (e.g.,
// String{"\"} is a compile error). Obviously we can't control the zero-value
// case, but this is better than nothing.
type String struct {
	s string
}

func (s String) String() string { return s.s }

// NewString creates a new String from a Go string. This is the preferred way
// to construct a String as it performs escaping and guarantees that the string
// is legal.
func NewString(s string) String {
	// Escaping a string is complicated; we'll just piggyback off of the JSON
	// library while this is still a toy.
	data, err := json.Marshal(s)
	if err != nil {
		// I don't think marshalling a string can fail, but let's fail hard if
		// it does.
		panic("Unexpected string marshalling error: " + err.Error())
	}
	// `data[1:len(data)-1]` because we want to get rid of the quotes, which
	// we'll add back on later. This is an artifact of using json.Marshal() to
	// escape the string.
	return String{string(data[1 : len(data)-1])}
}

func marshalJSONArray(array []JSON) string {
	eltStrings := make([]string, len(array))
	for i, elt := range array {
		eltStrings[i] = elt.Marshal()
	}
	return "[" + strings.Join(eltStrings, ", ") + "]"
}

func marshalJSONObject(object []Field) string {
	fieldStrings := make([]string, len(object))
	for i, field := range object {
		fieldStrings[i] = "\"" + field.Name.String() + "\": " +
			field.Value.Marshal()
	}
	return "{" + strings.Join(fieldStrings, ", ") + "}"
}

// Marshal renders the JSON to a string
func (json JSON) Marshal() string {
	var result string
	json.Match(
		func(a []JSON) { result = marshalJSONArray(a) },
		func(o []Field) { result = marshalJSONObject(o) },
		func(s String) { result = "\"" + s.String() + "\"" },
		func(i float64) { result = strconv.FormatFloat(i, 'f', -1, 64) },
		func(b bool) { result = strconv.FormatBool(b) },
		func(struct{}) { result = "null" },
	)
	return result
}

type Person struct {
	Name string
	Age  int
}

var (
	personJSONNameKey FieldName
	personJSONAgeKey  FieldName
)

func init() {
	personJSONNameKey = MustFieldName("name")
	personJSONAgeKey = MustFieldName("age")
}

func (p Person) JSON() JSON {
	return JSONObject([]Field{{
		Name:  personJSONNameKey,
		Value: JSONString(NewString(p.Name)),
	}, {
		Name:  personJSONAgeKey,
		Value: JSONNumber(float64(p.Age)),
	}})
}

type People []Person

func (p People) JSON() JSON {
	elts := make([]JSON, len(p))
	for i, person := range p {
		elts[i] = person.JSON()
	}
	return JSONArray(elts)
}

func main() {
	people := People{
		Person{"bob", 54},
		Person{"jane", 45},
		Person{"sue", 63},
	}
	fmt.Println(people.JSON().Marshal())
}
