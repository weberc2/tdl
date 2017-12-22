package main

import (
	"fmt"
	"strings"
)

func identf(format string, args ...interface{}) Ident {
	return Ident(fmt.Sprintf(format, args...))
}

func (tr TypeRef) Mangle() Ident {
	if len(tr.Params) < 1 {
		return tr.Name
	}
	typeStrings := make([]string, len(tr.Params))
	for i, typ := range tr.Params {
		typeStrings[i] = string(typ.Mangle())
	}
	return identf(
		"%sᐸ%sᐳ",
		string(tr.Name),
		strings.Join(typeStrings, "ˇ"),
	)
}

func (e Enum) Mangle() Ident {
	fieldStrings := make([]string, len(e))
	for i, field := range e {
		fieldStrings[i] = fmt.Sprintf(
			"%sπ%s",
			field.Name,
			field.Type.Mangle(),
		)
	}
	return identf("__enum__%s", strings.Join(fieldStrings, "ꟾ"))
}

func (s Struct) Mangle() Ident {
	fieldStrings := make([]string, len(s))
	for i, field := range s {
		fieldStrings[i] = fmt.Sprintf(
			"%sπ%s",
			field.Name,
			field.Type.Mangle(),
		)
	}
	return identf("__struct__%s", strings.Join(fieldStrings, "ˇ"))
}

func (t Tuple) Mangle() Ident {
	typeStrings := make([]string, len(t))
	for i, typ := range t {
		typeStrings[i] = string(typ.Mangle())
	}
	return identf("__tuple__%s", strings.Join(typeStrings, "ˇ"))
}

func (p Pointer) Mangle() Ident {
	return identf("__ptr__%s", p.Type.Mangle())
}

func (s Slice) Mangle() Ident {
	return identf("__slice__%s", s.Type.Mangle())
}

func (typ Type) Mangle() Ident {
	var out Ident
	typ.Match(
		func(tr TypeRef) { out = tr.Mangle() },
		func(e Enum) { out = e.Mangle() },
		func(s Struct) { out = s.Mangle() },
		func(t Tuple) { out = t.Mangle() },
		func(p Pointer) { out = p.Mangle() },
		func(s Slice) { out = s.Mangle() },
	)
	return out
}
