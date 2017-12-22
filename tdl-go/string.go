package main

import (
	"fmt"
	"strings"
)

func fieldStrings(fields []Field, delim string) string {
	fieldStrings := make([]string, len(fields))
	for i, field := range fields {
		fieldStrings[i] = string(field.Name) + " " + field.Type.String()
	}
	return strings.Join(fieldStrings, delim)
}

func (tr TypeRef) String() string {
	if len(tr.Params) < 1 {
		return string(tr.Name)
	}
	paramStrings := make([]string, len(tr.Params))
	for i, param := range tr.Params {
		paramStrings[i] = param.String()
	}
	return fmt.Sprintf("%s[%s]", tr.Name, strings.Join(paramStrings, ", "))
}

func (e Enum) String() string {
	return fmt.Sprintf("(%s)", fieldStrings([]Field(e), " | "))
}

func (s Struct) String() string {
	return fmt.Sprintf("{%s}", fieldStrings([]Field(s), "; "))
}

func (t Tuple) String() string {
	typeStrings := make([]string, len(t))
	for i, typ := range t {
		typeStrings[i] = typ.String()
	}
	return fmt.Sprintf("(%s)", strings.Join(typeStrings, ", "))
}

func (p Pointer) String() string { return "*" + p.Type.String() }

func (s Slice) String() string { return "[]" + s.Type.String() }

func (t Type) String() string {
	var out string
	t.Match(
		func(tr TypeRef) { out = tr.String() },
		func(e Enum) { out = e.String() },
		func(s Struct) { out = s.String() },
		func(t Tuple) { out = t.String() },
		func(p Pointer) { out = p.String() },
		func(s Slice) { out = s.String() },
	)
	return out
}

func (tc TypeCtor) String() string {
	if len(tc.TypeVars) < 1 {
		return string(tc.Name)
	}
	typeVarStrings := make([]string, len(tc.TypeVars))
	for i, typeVar := range tc.TypeVars {
		typeVarStrings[i] = string(typeVar)
	}
	return fmt.Sprintf("%s[%s]", tc.Name, strings.Join(typeVarStrings, ", "))
}

func (td TypeDecl) String() string {
	return fmt.Sprintf("type %s = %s", td.Ctor.String(), td.Type.String())
}
