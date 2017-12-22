package main

func (tr TypeRef) Substitute(ident Ident, typ Type) Type {
	if tr.Name == ident && len(tr.Params) < 1 {
		return typ
	}

	params := make([]Type, len(tr.Params))
	for i, param := range tr.Params {
		params[i] = param.Substitute(ident, typ)
	}
	return TypeRef_(TypeRef{Name: tr.Name, Params: params})
}

func (e Enum) Substitute(ident Ident, typ Type) Type {
	out := make(Enum, len(e))
	for i, field := range e {
		out[i] = Field{
			Name: field.Name,
			Type: field.Type.Substitute(ident, typ),
		}
	}
	return TypeEnum(out)
}

func (s Struct) Substitute(ident Ident, typ Type) Type {
	out := make(Struct, len(s))
	for i, field := range s {
		out[i] = Field{
			Name: field.Name,
			Type: field.Type.Substitute(ident, typ),
		}
	}
	return TypeStruct(out)
}

func (t Tuple) Substitute(ident Ident, typ Type) Type {
	out := make(Tuple, len(t))
	for i, t := range t {
		out[i] = t.Substitute(ident, typ)
	}
	return TypeTuple(out)
}

func (p Pointer) Substitute(ident Ident, typ Type) Type {
	return TypePointer(Pointer{p.Type.Substitute(ident, typ)})
}

func (s Slice) Substitute(ident Ident, typ Type) Type {
	return TypeSlice(Slice{s.Type.Substitute(ident, typ)})
}

func (t Type) Substitute(ident Ident, typ Type) Type {
	var out Type
	t.Match(
		func(tr TypeRef) { out = tr.Substitute(ident, typ) },
		func(e Enum) { out = e.Substitute(ident, typ) },
		func(s Struct) { out = s.Substitute(ident, typ) },
		func(t Tuple) { out = t.Substitute(ident, typ) },
		func(p Pointer) { out = p.Substitute(ident, typ) },
		func(s Slice) { out = s.Substitute(ident, typ) },
	)
	return out
}

func (td TypeDecl) Substitute(params []Type) (TypeDecl, error) {
	if len(params) < 1 {
		return td, nil
	}

	return TypeDecl.Substitute(
		TypeDecl{
			TypeCtor{td.Ctor.Name, td.Ctor.TypeVars[1:]},
			td.Type.Substitute(td.Ctor.TypeVars[0], params[0]),
		},
		params[1:],
	)
}
