package main

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
)

const pi = GoIdent("π")

func (e Enum) GoDecls(name Ident) []GoDecl {
	decls := make([]GoDecl, 0, 3*len(e)+2)
	variantMethodIdent := GoIdent("_" + name + "Variant")
	variantInterfaceTypeIdent := GoIdent("_" + name + "Variant")

	// Variant type
	variantType := GoTypeDecl{
		variantInterfaceTypeIdent,
		GoTypeInterface(GoInterface{GoInterfaceField{
			Name: variantMethodIdent,
		}}),
	}
	decls = append(decls, GoDeclType(variantType))

	// Enum type
	decls = append(decls, GoDeclType(GoTypeDecl{
		GoIdent(name),
		GoTypeStruct(GoStruct{GoField{
			Name: "variant",
			Type: GoTypeIdent(variantInterfaceTypeIdent),
		}}),
	}))

	matchFuncArgs := make([]GoField, len(e))
	matchCases := make([]GoCase, len(e))

	for i, variant := range e {
		variantConstructorIdent := GoIdent(name + variant.Name)
		variantGoTypeIdent := GoIdent("_" + name + variant.Name)

		// We can't simply create a new type off the definition type (e.g.,
		// `type VariantType GoType`) in case the definition type is an
		// interface (Go doesn't allow us to define new functions on an
		// interface) so instead we wrap it with a struct and name the field
		// `pi` for lulz.
		variantGoTypeDefinition := GoStruct{{
			Name: pi,
			Type: variant.Type.GoType(),
		}}

		variantMatchFuncIdent := GoIdent(variant.Name)
		matchFuncArgs[i] = GoField{
			Name: variantMatchFuncIdent,
			Type: GoTypeFunc(GoFuncType{Args: []GoField{{
				Name: pi,
				Type: variant.Type.GoType(),
			}}}),
		}

		// The case is:
		// case {{variantGoTypeIdent}}:
		//     {{variantMatchFuncIdent}}(pi.pi)
		//     return
		matchCases[i] = GoCase{
			Expr: GoExprIdent(variantGoTypeIdent),
			Stmts: []GoStmt{
				GoStmtExpr(GoExprCall(GoCallExpr{
					Fn:   GoExprIdent(variantMatchFuncIdent),
					Args: []GoExpr{GoExprIdent(pi + "." + pi)},
				})),
				GoStmtReturn(GoReturnStmt{}),
			},
		}

		decls = append(
			decls,
			// Type
			GoDeclType(GoTypeDecl{
				Name: variantGoTypeIdent,
				Type: GoTypeStruct(variantGoTypeDefinition),
			}),
			// Constructor
			GoDeclFunc(GoFuncDecl{
				Name: variantConstructorIdent,
				Signature: GoFuncType{
					Args:   []GoField{{Name: pi, Type: variant.Type.GoType()}},
					Return: []GoField{{Type: GoTypeIdent(GoIdent(name))}},
				},
				Body: GoBlock{GoStmtReturn(GoReturnStmt{
					GoExprStructLit(GoStructLit{
						Type: GoTypeIdent(GoIdent(name)),
						Fields: []GoPair{{
							Name: GoIdent("variant"),
							Value: GoExprStructLit(GoStructLit{
								Type: GoTypeIdent(variantGoTypeIdent),
								Fields: []GoPair{{
									Name:  pi,
									Value: GoExprIdent(pi),
								}},
							}),
						}},
					}),
				})},
			}),
			// Method
			GoDeclMethod(GoMethodDecl{
				Signature: GoMethodSpec{
					Name: variantMethodIdent,
					Recv: GoTypeIdent(variantGoTypeIdent),
				},
			}),
		)
	}

	// Declare the match method
	decls = append(decls, GoDeclMethod(GoMethodDecl{
		Signature: GoMethodSpec{
			Name:      "Match",
			Recv:      GoTypeIdent(GoIdent(name)),
			Signature: GoFuncType{Args: matchFuncArgs},
		},
		Body: GoBlock{
			GoStmtTypeSwitch(GoTypeSwitch{
				Assignment: GoPair{
					Name:  pi,
					Value: GoExprIdent("self.variant"),
				},
				Cases: matchCases,
			}),
		},
	}))

	return decls
}

var counter = 0
var anonEnums = map[string]int{}

func (e Enum) anonName() Ident {
	hash := md5.Sum(e.id())
	s := hex.EncodeToString(hash[:])

	var count int
	if i, found := anonEnums[s]; found {
		count = i
	} else {
		anonEnums[s] = counter
		count = counter
		counter++
	}
	return Ident("anonEnum" + strconv.Itoa(count))
	// s := md5.Sum(e.id())
	// return Ident("π" + hex.EncodeToString(s[:4]))
}

func (i Ident) GoType() GoType { return GoTypeIdent(GoIdent(i)) }

func (s Struct) GoType() GoType {
	goStruct := make(GoStruct, len(s))
	for i, field := range s {
		goStruct[i] = GoField{
			Name: GoIdent(field.Name),
			Type: field.Type.GoType(),
		}
	}
	return GoTypeStruct(goStruct)
}

func (t Tuple) GoType() GoType {
	goStruct := make(GoStruct, len(t))
	for i, typ := range t {
		goStruct[i] = GoField{
			Name: GoIdent("_" + strconv.Itoa(i)),
			Type: typ.GoType(),
		}
	}
	return GoTypeStruct(goStruct)
}

func (p Pointer) GoType() GoType {
	return GoTypePointer(GoPointer{p.Type.GoType()})
}

func (s Slice) GoType() GoType {
	return GoTypeSlice(GoSlice{s.Type.GoType()})
}

func (t Type) GoType() GoType {
	var result GoType
	t.Match(
		func(i Ident) { result = GoTypeIdent(GoIdent(i)) },
		func(e Enum) { result = GoTypeIdent(GoIdent(e.anonName())) },
		func(s Struct) { result = s.GoType() },
		func(t Tuple) { result = t.GoType() },
		func(p Pointer) { result = GoTypePointer(GoPointer{p.Type.GoType()}) },
		func(s Slice) { result = GoTypeSlice(GoSlice{s.Type.GoType()}) },
	)
	return result
}

func (t Type) GoDecls(name Ident) []GoDecl {
	var result []GoDecl
	goDecls := func(t interface {
		GoType() GoType
	}) []GoDecl {
		return []GoDecl{GoDeclType(GoTypeDecl{GoIdent(name), t.GoType()})}
	}
	t.Match(
		func(i Ident) { result = goDecls(i) },
		func(e Enum) { result = e.GoDecls(name) },
		func(s Struct) { result = goDecls(s) },
		func(t Tuple) { result = goDecls(t) },
		func(p Pointer) { result = goDecls(p) },
		func(s Slice) { result = goDecls(s) },
	)
	return result
}

func (e Enum) id() []byte {
	out := []byte("enum")
	for _, v := range e {
		out = append(out, []byte(v.Name)...)
		out = append(out, v.Type.id()...)
	}
	return out
}

func (s Struct) id() []byte {
	out := []byte("struct")
	for _, f := range s {
		out = append(out, []byte(f.Name)...)
		out = append(out, f.Type.id()...)
	}
	return out
}

func (t Tuple) id() []byte {
	out := []byte("tuple")
	for _, typ := range t {
		out = append(out, typ.id()...)
	}
	return out
}

func (t Type) id() []byte {
	var result []byte
	t.Match(
		func(i Ident) { result = append([]byte("ident"), []byte(i)...) },
		func(e Enum) { result = e.id() },
		func(s Struct) { result = s.id() },
		func(t Tuple) { result = t.id() },
		func(p Pointer) { result = append([]byte("pointer"), p.Type.id()...) },
		func(s Slice) { result = append([]byte("slice"), s.Type.id()...) },
	)
	return result
}

func (e Enum) Constituents() []Type {
	constituents := []Type{TypeEnum(e)}
	for _, field := range e {
		constituents = append(constituents, field.Type.Constituents()...)
	}
	return constituents
}

func (s Struct) Constituents() []Type {
	var constituents []Type
	for _, field := range s {
		constituents = append(constituents, field.Type.Constituents()...)
	}
	return constituents
}

func (t Tuple) Constituents() []Type {
	var constituents []Type
	for _, typ := range t {
		constituents = append(constituents, typ.Constituents()...)
	}
	return constituents
}

func (t Type) Constituents() []Type {
	var result []Type
	t.Match(
		func(i Ident) { result = nil },
		func(e Enum) { result = e.Constituents() },
		func(s Struct) { result = s.Constituents() },
		func(t Tuple) { result = t.Constituents() },
		func(p Pointer) { result = p.Type.Constituents() },
		func(s Slice) { result = s.Type.Constituents() },
	)
	return result
}

func (td TypeDecl) GoDecls() []GoDecl {
	var decls []GoDecl
	for _, constituent := range td.Type.Constituents() {
		constituent.Match(
			func(i Ident) {},
			func(e Enum) { decls = append(decls, e.GoDecls(e.anonName())...) },
			func(s Struct) {},
			func(t Tuple) {},
			func(p Pointer) {},
			func(s Slice) {},
		)
	}
	return append(decls, td.Type.GoDecls(td.Name)...)
}

func (f File) ToGo() GoFile {
	var decls []GoDecl
	for _, typeDecl := range f.TypeDecls {
		decls = append(decls, typeDecl.GoDecls()...)
	}
	return GoFile{
		FileComment: "This file is generated. Do not modify.",
		PackageName: GoIdent(f.PackageName),
		Decls:       decls,
	}
}
