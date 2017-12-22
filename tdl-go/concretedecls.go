package main

import (
	"fmt"
	"os"
)

type UnknownIdentErr Ident

func (uie UnknownIdentErr) Error() string {
	return fmt.Sprintf("Unknown identifier: '%s'", string(uie))
}

func (tr TypeRef) ConcreteDecls(env Environment) (Type, []TypeDecl, error) {
	if td, found := env[tr.Name]; found {
		// If the reference is already in the environment, just return the type
		if len(td.Ctor.TypeVars) < 1 {
			return TypeRef_(tr), nil, nil
		}

		var out []TypeDecl

		// handle nested generics (e.g., `List[Map[str, int]]`)
		for _, typ := range tr.Params {
			_, decls, err := typ.ConcreteDecls(env)
			if err != nil {
				return Type{}, nil, err
			}
			out = append(out, decls...)
		}

		concreteDecl, err := td.Apply(tr.Params)
		if err != nil {
			return Type{}, nil, err
		}
		return TypeRef_(TypeRef{
			Name: concreteDecl.Ctor.Name,
		}), append(out, concreteDecl), nil
	}
	return Type{}, nil, UnknownIdentErr(tr.Name)
}

func (e Enum) ConcreteDecls(env Environment) (Type, []TypeDecl, error) {
	var typeDecls []TypeDecl
	out := make(Enum, len(e))
	for i, field := range e {
		typ, decls, err := field.Type.ConcreteDecls(env)
		if err != nil {
			return Type{}, nil, err
		}
		out[i] = Field{field.Name, typ}
		typeDecls = append(typeDecls, decls...)
	}
	return TypeEnum(out), typeDecls, nil
}

func (s Struct) ConcreteDecls(env Environment) (Type, []TypeDecl, error) {
	var typeDecls []TypeDecl
	out := make(Struct, len(s))
	for i, field := range s {
		typ, decls, err := field.Type.ConcreteDecls(env)
		if err != nil {
			return Type{}, nil, err
		}
		out[i] = Field{field.Name, typ}
		typeDecls = append(typeDecls, decls...)
	}
	return TypeStruct(out), typeDecls, nil
}

func (t Tuple) ConcreteDecls(env Environment) (Type, []TypeDecl, error) {
	var typeDecls []TypeDecl
	out := make(Tuple, len(t))
	for i, typ := range t {
		typ, decls, err := typ.ConcreteDecls(env)
		if err != nil {
			return Type{}, nil, err
		}
		out[i] = typ
		typeDecls = append(typeDecls, decls...)
	}
	return TypeTuple(out), typeDecls, nil
}

func (p Pointer) ConcreteDecls(env Environment) (Type, []TypeDecl, error) {
	typ, decls, err := p.Type.ConcreteDecls(env)
	if err != nil {
		return Type{}, nil, err
	}
	return TypePointer(Pointer{typ}), decls, nil
}

func (s Slice) ConcreteDecls(env Environment) (Type, []TypeDecl, error) {
	typ, decls, err := s.Type.ConcreteDecls(env)
	if err != nil {
		return Type{}, nil, err
	}
	return TypeSlice(Slice{typ}), decls, nil
}

func (t Type) ConcreteDecls(env Environment) (Type, []TypeDecl, error) {
	var typ Type
	var typeDecls []TypeDecl
	var err error
	t.Match(
		func(tr TypeRef) { typ, typeDecls, err = tr.ConcreteDecls(env) },
		func(e Enum) { typ, typeDecls, err = e.ConcreteDecls(env) },
		func(s Struct) { typ, typeDecls, err = s.ConcreteDecls(env) },
		func(t Tuple) { typ, typeDecls, err = t.ConcreteDecls(env) },
		func(p Pointer) { typ, typeDecls, err = p.ConcreteDecls(env) },
		func(s Slice) { typ, typeDecls, err = s.ConcreteDecls(env) },
	)
	return typ, typeDecls, err
}

func applyMangledRefs(typ Type) Type {
	var out Type
	typ.Match(
		func(tr TypeRef) { out = TypeRef_(TypeRef{Name: tr.Mangle()}) },
		func(e Enum) {
			enum := make(Enum, len(e))
			for i, field := range e {
				enum[i] = Field{field.Name, applyMangledRefs(field.Type)}
			}
			out = TypeEnum(enum)
		},
		func(s Struct) {
			struct_ := make(Struct, len(s))
			for i, field := range s {
				struct_[i] = Field{field.Name, applyMangledRefs(field.Type)}
			}
			out = TypeStruct(struct_)
		},
		func(t Tuple) {
			tuple := make(Tuple, len(t))
			for i, typ := range t {
				tuple[i] = applyMangledRefs(typ)
			}
			out = TypeTuple(tuple)
		},
		func(p Pointer) { TypePointer(Pointer{applyMangledRefs(p.Type)}) },
		func(s Slice) { TypeSlice(Slice{applyMangledRefs(s.Type)}) },
	)
	return out
}

// ConcreteDecls returns all concrete declarations (including itself) if `td`
// is itself concrete, otherwise it returns no decls.
func (td TypeDecl) ConcreteDecls(env Environment) ([]TypeDecl, error) {
	if len(td.Ctor.TypeVars) < 1 {
		typ, decls, err := td.Type.ConcreteDecls(env)
		if err != nil {
			return nil, err
		}
		out := append(
			decls,
			TypeDecl{Ctor: td.Ctor, Type: typ},
		)
		return out, nil
	}
	return nil, nil
}

func dedup(typeDecls []TypeDecl) []TypeDecl {
	var out []TypeDecl
OUTER:
	for _, typeDecl := range typeDecls {
		for _, existing := range out {
			// If the declaration has already been added to the output, move on
			// to the next declaration
			if existing.Equal(typeDecl) {
				continue OUTER
			}
		}
		out = append(out, typeDecl)
	}

	return out
}

func (f File) ConcreteDecls() (File, error) {
	env := NewEnvironment(f.TypeDecls)
	var typeDecls []TypeDecl
	for _, typeDecl := range f.TypeDecls {
		concreteDecls, err := typeDecl.ConcreteDecls(env)
		if err != nil {
			return File{}, err
		}
		typeDecls = append(typeDecls, concreteDecls...)
	}

	// Replace nested concrete generics with concrete type references. e.g.,
	// (int, List[int]) => (int, List<int>)
	for i := range typeDecls {
		typeDecls[i].Type = applyMangledRefs(typeDecls[i].Type)
	}
	return File{PackageName: f.PackageName, TypeDecls: dedup(typeDecls)}, nil
}
