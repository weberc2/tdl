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

func (td TypeDecl) Apply(params []Type) (TypeDecl, error) {
	// Validate that the parameters correspond to the type constructor. For
	// now, that just means checking the length, but this will likely entail
	// constraint checks at some point.
	if len(td.Ctor.TypeVars) != len(params) {
		return TypeDecl{}, TemplateErr{td.Ctor.Name, td.Ctor.TypeVars, params}
	}

	// Return a concrete type; this specifically means the type constructor has
	// no arguments--if the type is generic, the type parameters will be
	// mangled into the name to prevent it from colliding with other concrete
	// implementations. This also means any type vars will be substituted for
	// their corresponding type parameters.
	substituted, err := td.Substitute(params)
	if err != nil {
		return TypeDecl{}, err
	}
	substituted.Ctor.Name = TypeRef{td.Ctor.Name, params}.Mangle()
	return substituted, nil
}
