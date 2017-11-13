package main

func (e Enum) Equal(other Enum) bool {
	if len(e) != len(other) {
		return false
	}
	for i, field := range e {
		if field.Name != other[i].Name || !field.Type.Equal(other[i].Type) {
			return false
		}
	}
	return true
}

func (s Struct) Equal(other Struct) bool {
	if len(s) != len(other) {
		return false
	}
	for i, field := range s {
		if field.Name != other[i].Name || !field.Type.Equal(other[i].Type) {
			return false
		}
	}
	return true
}

func (t Tuple) Equal(other Tuple) bool {
	if len(t) != len(other) {
		return false
	}
	for i, typ := range t {
		if !typ.Equal(other[i]) {
			return false
		}
	}
	return true
}

func (p Pointer) Equal(other Pointer) bool {
	return p.Type.Equal(other.Type)
}

func (s Slice) Equal(other Slice) bool {
	return s.Type.Equal(other.Type)
}

func (t Type) Equal(other Type) bool {
	var result bool
	t.Match(
		func(i Ident) {
			other.Match(
				func(i2 Ident) { result = i == i2 },
				func(Enum) { result = false },
				func(Struct) { result = false },
				func(Tuple) { result = false },
				func(Pointer) { result = false },
				func(Slice) { result = false },
			)
		},
		func(e Enum) {
			other.Match(
				func(Ident) { result = false },
				func(e2 Enum) { result = e.Equal(e2) },
				func(Struct) { result = false },
				func(Tuple) { result = false },
				func(Pointer) { result = false },
				func(Slice) { result = false },
			)
		},
		func(s Struct) {
			other.Match(
				func(Ident) { result = false },
				func(Enum) { result = false },
				func(s2 Struct) { result = s.Equal(s2) },
				func(Tuple) { result = false },
				func(Pointer) { result = false },
				func(Slice) { result = false },
			)
		},
		func(t Tuple) {
			other.Match(
				func(Ident) { result = false },
				func(Enum) { result = false },
				func(Struct) { result = false },
				func(t2 Tuple) { result = t.Equal(t2) },
				func(Pointer) { result = false },
				func(Slice) { result = false },
			)
		},
		func(p Pointer) {
			other.Match(
				func(Ident) { result = false },
				func(Enum) { result = false },
				func(Struct) { result = false },
				func(Tuple) { result = false },
				func(p2 Pointer) { result = p.Equal(p2) },
				func(Slice) { result = false },
			)
		},
		func(s Slice) {
			other.Match(
				func(Ident) { result = false },
				func(Enum) { result = false },
				func(Struct) { result = false },
				func(Tuple) { result = false },
				func(Pointer) { result = false },
				func(s2 Slice) { result = s.Equal(s2) },
			)
		},
	)
	return result
}

func (td TypeDecl) Equal(other TypeDecl) bool {
	return td.Name == other.Name && td.Type.Equal(other.Type)
}

func (f File) Equal(other File) bool {
	if f.PackageName != other.PackageName {
		return false
	}
	if len(f.TypeDecls) != len(other.TypeDecls) {
		return false
	}

	for i, td := range f.TypeDecls {
		if !td.Equal(other.TypeDecls[i]) {
			return false
		}
	}

	return true
}
