package main

import (
	"fmt"
)

func (i *Ident) Parser() Parser { return i }

func (i *Ident) Parse(input Input) (Input, error) {
	rest, err := Seq{
		Any{RuneLit('_'), UnicodeClassLetter},
		Repeat{Any{RuneLit('_'), UnicodeClassLetter, UnicodeClassDigit}},
	}.Parse(input)
	if err != nil {
		return Input{}, ParseError{
			Cause:   err,
			Message: fmt.Sprintf("Wanted IDENT; got '%s...'", input.Take(20)),
		}
	}
	*i = Ident(input.Source[input.Offset:rest.Offset])
	return rest, nil
}

func (i *Ident) EqualASTNode(other ASTNode) bool {
	otherIdent, ok := other.(*Ident)
	return ok && *i == *otherIdent
}

func (i *Ident) Name() string { return "IDENT" }

func (f *Field) Parser() Parser {
	// create a temp for the field name, since we don't want to modify the
	// field in the case where we successfully parse the field name but then
	// fail to parse the WS or the type.
	var name Ident
	return sideEffectParser{
		Seq{name.Parser(), WS, f.Type.Parser()},
		func() { f.Name = name },
	}
}

func (f *Field) EqualASTNode(other ASTNode) bool {
	otherField, ok := other.(*Field)
	return ok && f.Name == otherField.Name && f.Type.Equal(otherField.Type)
}

type enumField Field

func (ef *enumField) Parser() Parser {
	var name Ident
	return sideEffectParser{
		Seq{name.Parser(), WS, ef.Type.Parser()},
		func() { ef.Name = name },
	}
}

func (e *Enum) Parser() Parser {
	var field enumField
	var fields []Field

	// IDENT TYPE | ... IDENT TYPE
	return sideEffectParser{
		Seq{
			// There has to be at least one `IDENT TYPE |`; after we parse each
			// field, we need to append it to the list of parsed fields and
			// then clear it before the next field is parsed.
			OneOrMore{
				sideEffectParser{
					Seq{field.Parser(), CanWS, RuneLit('|'), CanWS},
					func() {
						fields = append(fields, Field(field))
						field = enumField{}
					},
				},
			},
			// Parse the last field; no need to reset `field` after we append
			// it to `fields`.
			sideEffectParser{
				field.Parser(),
				func() { fields = append(fields, Field(field)) },
			},
		},
		// Only modify `e` once we've completed parsing the Enum
		func() { *e = Enum(fields) },
	}
}

func (e *Enum) EqualASTNode(other ASTNode) bool {
	otherEnum, ok := other.(*Enum)
	return ok && e.Equal(*otherEnum)
}

func (s *Struct) Parser() Parser {
	var field Field
	var fields []Field
	return sideEffectParser{
		Seq{
			RuneLit('{'),
			CanWS,
			Opt{Seq{
				Repeat{sideEffectParser{
					Seq{field.Parser(), EOS},
					func() {
						fields = append(fields, field)
						field = Field{}
					},
				}},
				Opt{sideEffectParser{
					field.Parser(),
					func() { fields = append(fields, field) },
				}},
			}},
			CanWS,
			RuneLit('}'),
		},
		func() { *s = Struct(fields) },
	}
}

func (s *Struct) EqualASTNode(other ASTNode) bool {
	otherStruct, ok := other.(*Struct)
	return ok && s.Equal(*otherStruct)
}

func (t *Tuple) Parser() Parser {
	// `()` or `(TYPE, ... TYPE)`
	var typ Type
	var types []Type
	return sideEffectParser{
		Seq{
			RuneLit('('),
			Opt{Seq{
				// A non-empty tuple must have 2 or more types. A single-type
				// tuple is syntactically indistinguishable from a paren group.
				// E.g., `foo (bar baz)` could be a type constructor with a
				// tuple argument, or it could be a type constructor whose
				// argument is the result of another type constructor. Since
				// there's no use case (I think) for a single-element tuple, it
				// seems better to reserve this syntax for paren groups. Also,
				// Haskell made the same design decision.
				OneOrMore{sideEffectParser{
					Seq{CanWS, typ.Parser(), CanWS, RuneLit(',')},
					// Each time we successfully parse the "TYPE," sequence,
					// add the temp `typ` variable to the temp `types` slice
					// and clear `typ`.
					func() { types = append(types, typ); typ = Type{} },
				}},
				CanWS,
				// The last element in the tuple. Just like the previous
				// element(s), it needs to be added to the temp slice after it
				// is parsed, but because it is the last element, there is no
				// need to clear `typ`.
				sideEffectParser{
					typ.Parser(),
					func() { types = append(types, typ) },
				},
			}},
			RuneLit(')'),
		},
		// Once we've successfully parsed the tuple, mutate `t`.
		func() { *t = Tuple(types) },
	}
}

func (t *Tuple) EqualASTNode(other ASTNode) bool {
	otherTuple, ok := other.(*Tuple)
	return ok && t.Equal(*otherTuple)
}

func (p *Pointer) Parser() Parser {
	// No need for a temporary Type parser b/c p.Type.Parser() is the last in
	// the sequence. I.e., if we parse it correctly, then the whole Pointer is
	// good to go.
	return Seq{RuneLit('*'), CanWS, p.Type.Parser()}
}

func (p *Pointer) EqualASTNode(other ASTNode) bool {
	otherPointer, ok := other.(*Pointer)
	return ok && p.Equal(*otherPointer)
}

func (s *Slice) Parser() Parser {
	// No need for a temporary Type parser b/c s.Type.Parser() is the last in
	// the sequence. I.e., if we parse it correctly, then the whole Slice is
	// good to go.
	return Seq{StringLit("[]"), CanWS, s.Type.Parser()}
}

func (s *Slice) EqualASTNode(other ASTNode) bool {
	otherSlice, ok := other.(*Slice)
	return ok && s.Equal(*otherSlice)
}

type ASTNode interface {
	Parser() Parser
	EqualASTNode(other ASTNode) bool
}

type sideEffectNodeParser struct {
	node       ASTNode
	sideEffect func()
}

func (senp sideEffectNodeParser) Parse(input Input) (Input, error) {
	rest, err := senp.node.Parser().Parse(input)
	if err != nil {
		return Input{}, err
	}
	senp.sideEffect()
	return rest, nil
}

func (senp sideEffectNodeParser) Name() string {
	return senp.node.Parser().Name()
}

// TODO: Find a better way to handle recursive types. Without this type,
// Type.Parser() would return `Any{..., Seq{'(', Type.Parser(), ')'}}` which
// creates an infinite recursion loop that overflows the stack. This type gets
// around that by deferring the construction of the inner Type.Parser() until
// parse time. There's probably a general-purpose type that could solve this
// problem so we don't need to create a special type each time.
type typeParens struct {
	typ *Type
}

func (tp typeParens) Parse(input Input) (Input, error) {
	return Seq{
		RuneLit('('),
		CanWS,
		tp.typ.Parser(),
		CanWS,
		RuneLit(')'),
	}.Parse(input)
}

func (tp typeParens) Name() string {
	return tp.typ.Parser().Name()
}

func (tr *TypeRef) Parser() Parser {
	// Provision scratch variables so we only update `tr` if parsing succeeds
	var name Ident
	var params []Type
	var currentParam Type

	return sideEffectParser{
		Seq{
			&name,
			Opt{Seq{
				CanWS,
				RuneLit('['),
				Repeat{sideEffectParser{
					Seq{
						CanWS,
						currentParam.Parser(),
						CanWS,
						RuneLit(','),
					},
					func() { params = append(params, currentParam) },
				}},
				CanWS,
				sideEffectNodeParser{
					&currentParam,
					func() { params = append(params, currentParam) },
				},
				CanWS,
				RuneLit(']'),
			}},
		},

		// Only update `tr` if parsing succeeded
		func() {
			*tr = TypeRef{Name: name, Params: params}
			name = ""
			params = []Type{}
			currentParam = Type{}
		},
	}
}

func (tr *TypeRef) EqualASTNode(other ASTNode) bool {
	otherTypeRef, ok := other.(*TypeRef)
	return ok && tr.Equal(*otherTypeRef)
}

func (t *Type) Parser() Parser {
	var typeRef TypeRef
	var enum Enum
	var struct_ Struct
	var tuple Tuple
	var pointer Pointer
	var slice Slice

	return Any{
		sideEffectNodeParser{
			&enum,
			func() { *t = TypeEnum(enum) },
		},
		Any{
			sideEffectNodeParser{
				&struct_,
				func() { *t = TypeStruct(struct_) },
			},
			sideEffectNodeParser{
				&tuple,
				func() { *t = TypeTuple(tuple) },
			},
			sideEffectNodeParser{
				&pointer,
				func() { *t = TypePointer(pointer) },
			},
			sideEffectNodeParser{
				&slice,
				func() { *t = TypeSlice(slice) },
			},
			sideEffectNodeParser{
				&typeRef,
				func() { *t = TypeRef_(typeRef) },
			},
			typeParens{t},
		},
	}
}

func (t *Type) EqualASTNode(other ASTNode) bool {
	otherType, ok := other.(*Type)
	return ok && t.Equal(*otherType)
}

func (tc *TypeCtor) Parser() Parser {
	var typeName Ident
	var typeArgs []Ident
	var currentTypeArg Ident
	return sideEffectParser{
		Seq{
			typeName.Parser(),
			CanWS,
			Opt{Seq{
				RuneLit('['),
				Repeat{sideEffectParser{
					Seq{
						CanWS,
						&currentTypeArg,
						CanWS,
						RuneLit(','),
					},
					// Only update typeArgs if we get the full `T,`
					func() { typeArgs = append(typeArgs, currentTypeArg) },
				}},
				CanWS,
				sideEffectNodeParser{
					&currentTypeArg,
					func() { typeArgs = append(typeArgs, currentTypeArg) },
				},
				RuneLit(']'),
			}},
		},
		func() {
			*tc = TypeCtor{Name: typeName, TypeVars: typeArgs}
			typeName = ""
			typeArgs = nil
			currentTypeArg = ""
		},
	}
}

func (td *TypeDecl) Parser() Parser {
	var ctor TypeCtor
	var typ Type
	// type IDENT[A, B, ..., Z] = TYPE
	return sideEffectParser{
		Seq{
			StringLit("type"),
			WS,
			ctor.Parser(),
			CanWS,
			RuneLit('='),
			CanWS,
			typ.Parser(),
		},
		func() { *td = TypeDecl{Ctor: ctor, Type: typ} },
	}
}

func (td *TypeDecl) EqualASTNode(other ASTNode) bool {
	otherTypeDecl, ok := other.(*TypeDecl)
	return ok && td.Equal(*otherTypeDecl)
}

func (f *File) Parser() Parser {
	// These are all temporary buffer parsers; we don't want to modify `f`
	// until we know the parse was successful.
	var td TypeDecl
	var tds []TypeDecl
	var packageName Ident

	return sideEffectParser{
		// package IDENT
		//
		// type IDENT = TYPE
		// ...
		//
		// EOF
		Seq{
			CanWS,
			StringLit("package"),
			WS,
			packageName.Parser(),
			WS,
			// Every time we successfully parse a type decl, append it to the
			// list of type decls and clear the type decl for the next attempt.
			Repeat{sideEffectParser{
				Seq{td.Parser(), CanWS},
				func() { tds = append(tds, td); td = TypeDecl{} },
			}},
			EOF,
		},

		// If the parse was successful, copy the temporary buffer values into
		// `f`.
		func() { *f = File{PackageName: packageName, TypeDecls: tds} },
	}
}
