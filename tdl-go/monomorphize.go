package main

import (
	"fmt"
	"strings"
)

type Environment map[Ident]TypeDecl

func (e Environment) addPrimitive(name Ident) {
	e[name] = TypeDecl{TypeCtor{name, nil}, TypeRef_(TypeRef{name, nil})}
}

func NewEnvironment(typeDecls []TypeDecl) Environment {
	env := make(Environment, len(typeDecls))
	env.addPrimitive("string")
	env.addPrimitive("int")
	for _, td := range typeDecls {
		env[td.Ctor.Name] = td
	}
	return env
}

type TemplateErr struct {
	Name   Ident
	Wanted []Ident
	Got    []Type
}

func (te TemplateErr) Error() string {
	wantedStrings := make([]string, len(te.Wanted))
	for i, wanted := range te.Wanted {
		wantedStrings[i] = string(wanted)
	}

	gotStrings := make([]string, len(te.Got))
	for i, got := range te.Got {
		gotStrings[i] = got.String()
	}

	return fmt.Sprintf(
		"Type constructor `%s` wanted arguments [%s]; got [%s]",
		te.Name,
		strings.Join(wantedStrings, ", "),
		strings.Join(gotStrings, ", "),
	)
}

// Collect all type references from the type
func (t Type) TypeRefs() []TypeRef {
	var typeRefs []TypeRef
	t.Match(
		func(tr TypeRef) {
			if len(tr.Params) > 0 {
				for _, param := range tr.Params {
					typeRefs = append(typeRefs, param.TypeRefs()...)
				}
			}
			typeRefs = append(typeRefs, tr)
		},
		func(e Enum) {
			for _, field := range e {
				typeRefs = append(typeRefs, field.Type.TypeRefs()...)
			}
		},
		func(s Struct) {
			for _, field := range s {
				typeRefs = append(typeRefs, field.Type.TypeRefs()...)
			}
		},
		func(t Tuple) {
			for _, typ := range t {
				typeRefs = append(typeRefs, typ.TypeRefs()...)
			}
		},
		func(p Pointer) { typeRefs = p.Type.TypeRefs() },
		func(s Slice) { typeRefs = s.Type.TypeRefs() },
	)
	return typeRefs
}

func (t Type) MangleTypeRefs() Type {
	var type_ Type
	t.Match(
		func(tr TypeRef) { type_ = TypeRef_(TypeRef{Name: tr.Mangle()}) },
		func(e Enum) {
			out := make(Enum, len(e))
			for i, field := range e {
				out[i] = Field{field.Name, field.Type.MangleTypeRefs()}
			}
			type_ = TypeEnum(out)
		},
		func(s Struct) {
			out := make(Struct, len(s))
			for i, field := range s {
				out[i] = Field{field.Name, field.Type.MangleTypeRefs()}
			}
			type_ = TypeStruct(out)
		},
		func(t Tuple) {
			out := make(Tuple, len(t))
			for i, typ := range t {
				out[i] = typ.MangleTypeRefs()
			}
			type_ = TypeTuple(out)
		},
		func(p Pointer) {
			type_ = TypePointer(Pointer{p.Type.MangleTypeRefs()})
		},
		func(s Slice) {
			type_ = TypeSlice(Slice{s.Type.MangleTypeRefs()})
		},
	)
	return type_
}

func (td TypeDecl) Apply(params []Type) (Type, error) {
	// Validate that the parameters correspond to the type constructor. For
	// now, that just means checking the length, but this will likely entail
	// constraint checks at some point.
	if len(td.Ctor.TypeVars) != len(params) {
		return Type{}, TemplateErr{td.Ctor.Name, td.Ctor.TypeVars, params}
	}

	// Return a concrete type (all type variables have been substituted for
	// concrete types).
	typ := td.Type
	for i, param := range params {
		typ = typ.Substitute(td.Ctor.TypeVars[i], param)
	}
	return typ, nil
}

type UnknownIdentErr Ident

func (uie UnknownIdentErr) Error() string {
	return fmt.Sprintf("Unknown identifier: '%s'", string(uie))
}

func dedup(genericTypeRefs []TypeRef) []TypeRef {
	var out []TypeRef
OUTER:
	for _, typeRef := range genericTypeRefs {
		for _, existing := range out {
			// If the declaration has already been added to the output, move on
			// to the next declaration
			if existing.Equal(typeRef) {
				continue OUTER
			}
		}
		out = append(out, typeRef)
	}

	return out
}

func (f File) Monomorphize() (File, error) {
	env := NewEnvironment(f.TypeDecls)
	var concreteDecls []TypeDecl
	var typeRefs []TypeRef
	for _, typeDecl := range f.TypeDecls {
		if len(typeDecl.Ctor.TypeVars) > 0 {
			continue
		}

		typeRefs = append(typeRefs, typeDecl.Type.TypeRefs()...)
		concreteDecls = append(
			concreteDecls,
			TypeDecl{
				Ctor: typeDecl.Ctor,
				Type: typeDecl.Type.MangleTypeRefs(),
			},
		)
	}

	// For each unique concrete type reference, look up its type declaration
	// and create a new concrete type. Replace concrete parameterized type
	// references (e.g., `List[int]`) with mangled identifiers (e.g.,
	// `Listᐸintᐳ`). Stick these new type decls in the list of declarations
	// which will be compiled into Go.
	for _, typeRef := range dedup(typeRefs) {
		if len(typeRef.Params) > 0 {
			if typeDecl, found := env[typeRef.Name]; found {
				concreteType, err := typeDecl.Apply(typeRef.Params)
				if err != nil {
					return File{}, err
				}
				concreteDecls = append(
					concreteDecls,
					TypeDecl{
						Ctor: TypeCtor{Name: typeRef.Mangle()},
						Type: concreteType.MangleTypeRefs(),
					},
				)
				continue
			}
			return File{}, UnknownIdentErr(typeRef.Name)
		}
	}

	return File{PackageName: f.PackageName, TypeDecls: concreteDecls}, nil
}
