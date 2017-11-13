package main

import (
	"fmt"
	"strconv"
	"strings"
)

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
		fieldStrings[i] = fmt.Sprintf(
			"\"%s\": %s",
			field.Name,
			field.Value.Marshal(),
		)
	}
	return "{" + strings.Join(fieldStrings, ", ") + "}"
}

func (json JSON) Marshal() string {
	var result string
	json.Match(
		func(a []JSON) { result = marshalJSONArray(a) },
		func(o []Field) { result = marshalJSONObject(o) },
		func(s string) { result = "\"" + s + "\"" },
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

func (p Person) JSON() JSON {
	return JSONObject([]Field{{
		Name:  "name",
		Value: JSONString(p.Name),
	}, {
		Name:  "age",
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
