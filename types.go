package main

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
)

type Field struct {
	Name Ident
	Type Type
}

type Enum []Field

const pi = GoIdent("π")

func (e Enum) ToGo(name Ident) []GoDecl {
	decls := make([]GoDecl, 0, 3*len(e)+2)
	variantMethodIdent := GoIdent("_" + name + "Variant")
	variantInterfaceTypeIdent := GoIdent("_" + name + "Variant")

	// Variant type
	variantType := GoTypeDecl{
		variantInterfaceTypeIdent,
		GoInterface{GoInterfaceField{Name: variantMethodIdent}},
	}
	decls = append(decls, variantType)

	// Enum type
	decls = append(decls, GoTypeDecl{
		GoIdent(name),
		GoStruct{GoField{Name: "variant", Type: variantInterfaceTypeIdent}},
	})

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
			Type: GoFuncType{Args: []GoField{{
				Name: pi,
				Type: variant.Type.GoType(),
			}}},
		}

		// The case is:
		// case {{variantGoTypeIdent}}:
		//     {{variantMatchFuncIdent}}(pi.pi)
		//     return
		matchCases[i] = GoCase{
			Expr: variantGoTypeIdent,
			Stmts: []GoStmt{
				GoCallExpr{
					Fn:   variantMatchFuncIdent,
					Args: []GoExpr{pi + "." + pi},
				},
				GoReturnStmt{},
			},
		}

		decls = append(
			decls,
			// Type
			GoTypeDecl{
				Name: variantGoTypeIdent,
				Type: variantGoTypeDefinition,
			},
			// Constructor
			GoFuncDecl{
				Name: variantConstructorIdent,
				Signature: GoFuncType{
					Args:   []GoField{{Name: pi, Type: variant.Type.GoType()}},
					Return: []GoField{{Type: GoIdent(name)}},
				},
				Body: GoBlock{GoReturnStmt{GoStructLit{
					Type: GoIdent(name),
					Fields: []GoPair{{
						Name: GoIdent("variant"),
						Value: GoStructLit{
							Type:   variantGoTypeIdent,
							Fields: []GoPair{{Name: pi, Value: pi}},
						},
					}},
				}}},
			},
			// Method
			GoMethodDecl{
				Signature: GoMethodSpec{
					Name: variantMethodIdent,
					Recv: variantGoTypeIdent,
				},
			},
		)
	}

	// Declare the match method
	decls = append(decls, GoMethodDecl{
		Signature: GoMethodSpec{
			Name:       "Match",
			Recv:       GoIdent(name),
			GoFuncType: GoFuncType{Args: matchFuncArgs},
		},
		Body: GoBlock{
			GoTypeSwitch{
				Assignment: GoPair{Name: pi, Value: GoIdent("self.variant")},
				Cases:      matchCases,
			},
		},
	})

	return decls
}

func (e Enum) GoDecls(name Ident) []GoDecl { return e.ToGo(name) }

func (e Enum) id() []byte {
	out := []byte("enum")
	for _, v := range e {
		out = append(out, []byte(v.Name)...)
		out = append(out, v.Type.id()...)
	}
	return out
}

func (e Enum) anonName() Ident {
	s := md5.Sum(e.id())
	return Ident(hex.EncodeToString(s[:4]))
}

func (e Enum) GoType() GoType { return GoIdent(e.anonName()) }

func (e Enum) Constituents() []Type {
	var out []Type
	for _, v := range e {
		out = append(out, v.Type.Constituents()...)
	}
	return out
}

type Struct []Field

func (s Struct) GoDecls(name Ident) []GoDecl {
	return []GoDecl{GoTypeDecl{GoIdent(name), s.GoType()}}
}

func (s Struct) GoType() GoType {
	goStruct := make(GoStruct, len(s))
	for i, field := range s {
		goStruct[i] = GoField{
			Name: GoIdent(field.Name),
			Type: field.Type.GoType(),
		}
	}
	return goStruct
}

func (s Struct) Constituents() []Type {
	var out []Type
	for _, field := range s {
		out = append(out, field.Type.Constituents()...)
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

type Tuple []Type

func (t Tuple) GoDecls(name Ident) []GoDecl {
	return []GoDecl{GoTypeDecl{Name: GoIdent(name), Type: t.GoType()}}
}

func (t Tuple) GoType() GoType {
	out := make(GoStruct, len(t))
	for i, typ := range t {
		out[i] = GoField{
			Name: GoIdent("_" + strconv.Itoa(i)),
			Type: typ.GoType(),
		}
	}
	return out
}

func (t Tuple) Constituents() []Type {
	var out []Type
	for _, typ := range t {
		out = append(out, typ.Constituents()...)
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

type Pointer struct {
	Type Type
}

func (p Pointer) GoType() GoType { return GoPointer{p.Type.GoType()} }

func (p Pointer) GoDecls(name Ident) []GoDecl {
	return []GoDecl{GoTypeDecl{Name: GoIdent(name), Type: p.GoType()}}
}

func (p Pointer) Constituents() []Type { return p.Type.Constituents() }

func (p Pointer) id() []byte {
	return append([]byte("pointer"), p.Type.id()...)
}

type Slice struct {
	Type Type
}

func (s Slice) GoType() GoType { return GoSlice{s.Type.GoType()} }

func (s Slice) GoDecls(name Ident) []GoDecl {
	return []GoDecl{GoTypeDecl{Name: GoIdent(name), Type: s.GoType()}}
}

func (s Slice) Constituents() []Type { return s.Type.Constituents() }

func (s Slice) id() []byte {
	return append([]byte("slice"), s.Type.id()...)
}

type Ident string

func (i Ident) GoType() GoType { return GoIdent(i) }

func (i Ident) GoDecls(name Ident) []GoDecl {
	return []GoDecl{GoTypeDecl{GoIdent(name), GoIdent(i)}}
}

func (i Ident) Constituents() []Type { return nil }

func (i Ident) id() []byte { return append([]byte("ident"), []byte(i)...) }

type Type interface {
	GoType() GoType
	GoDecls(name Ident) []GoDecl
	Constituents() []Type
	id() []byte
}

type TypeDecl struct {
	Name Ident
	Type Type
}

func (td TypeDecl) GoDecls() []GoDecl {
	var decls []GoDecl
	for _, constituent := range td.Type.Constituents() {
		if e, ok := constituent.(Enum); ok {
			decls = append(decls, e.GoDecls(e.anonName())...)
		}
	}
	return append(decls, td.Type.GoDecls(td.Name)...)
}

type File struct {
	PackageName Ident
	TypeDecls   []TypeDecl
}

func (f File) ToGo() GoFile {
	var decls []GoDecl
	for _, typeDecl := range f.TypeDecls {
		decls = append(decls, typeDecl.GoDecls()...)
	}
	return GoFile{
		PackageName: GoIdent(f.PackageName),
		Decls:       decls,
	}
}

// func main() {
// 	s := bufio.NewScanner(os.Stdin)
// 	for {
// 		fmt.Print(" > ")
// 		if !s.Scan() {
// 			break
// 		}
// 		r := combinator.Any(pTypeDecl, pType)(combinator.Input(s.Text()))
// 		if r.Err != nil {
// 			fmt.Println(r.Err)
// 			continue
// 		}
// 		if r.Rest != "" {
// 			fmt.Printf("Unmatched: '%s'\n", r.Rest)
// 			continue
// 		}
// 		switch x := r.Value.(type) {
// 		case Type:
// 			for _, d := range x.GoDecls() {
// 				fmt.Println(d.Render(0) + "\n")
// 			}
// 			fmt.Println(x.GoType().Render(0))
// 			if err := s.Err(); err != nil {
// 				fmt.Println(err)
// 				return
// 			}
// 		case TypeDecl:
// 			for _, d := range x.Type.GoDecls() {
// 				fmt.Println(d.Render(0) + "\n")
// 			}
// 			fmt.Println(x.Type.GoType().Render(0))
// 			if err := s.Err(); err != nil {
// 				fmt.Println(err)
// 				return
// 			}
// 		}
// 	}
// }
