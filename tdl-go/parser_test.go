package main

import (
	"strings"
	"testing"

	"github.com/kr/pretty"
)

func TestParser(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		node        ASTNode
		wantedNode  ASTNode
		wantedRest  string
		wantedError bool
	}{
		{
			name:       "ident-one-letter",
			input:      "x",
			node:       new(Ident),
			wantedNode: nodeIdent("x"),
		},
		{
			name:       "ident-one-underscore",
			input:      "_",
			node:       new(Ident),
			wantedNode: nodeIdent("_"),
		},
		{
			name:       "ident-multi-letter",
			input:      "abc",
			node:       new(Ident),
			wantedNode: nodeIdent("abc"),
		},
		{
			name:       "ident-w-nonleading-digits",
			input:      "a123",
			node:       new(Ident),
			wantedNode: nodeIdent("a123"),
		},
		{
			name:        "ident-cant-start-w-digit",
			input:       "1abc",
			node:        new(Ident),
			wantedError: true,
		},
		{
			name:       "ident-w-nonleading-underscores",
			input:      "abc_123_",
			node:       new(Ident),
			wantedNode: nodeIdent("abc_123_"),
		},
		{
			name:       "ident-breaks-on-ws",
			input:      "foo bar",
			node:       new(Ident),
			wantedNode: nodeIdent("foo"),
			wantedRest: " bar",
		},
		{
			name:       "type-ident",
			input:      "Foo",
			node:       &Type{},
			wantedNode: nodeType(TypeIdent("Foo")),
		},
		{
			name:  "enum-simple",
			input: "LeFoo Foo | LeBar Bar",
			node:  &Enum{},
			wantedNode: &Enum{
				Field{"LeFoo", TypeIdent("Foo")},
				Field{"LeBar", TypeIdent("Bar")},
			},
		},
		{
			name:  "enum-many-elts",
			input: "Yes () | No () | Maybe ()",
			node:  &Enum{},
			wantedNode: &Enum{
				Field{"Yes", TypeTuple(nil)},
				Field{"No", TypeTuple(nil)},
				Field{"Maybe", TypeTuple(nil)},
			},
		},
		{
			// TODO: Eventually we will want to permit "type-less" variants,
			// i.e., variants whose type is unit. For now, I'm running into too
			// many problems building that parser, so I'll just add this test
			// to assert that the code behaves as I think it does.
			name:        "enum-anon-fields-must-have-a-label-and-type",
			input:       "Nil | Cons List",
			node:        &Enum{},
			wantedError: true,
		},
		{
			name:  "enum-nested",
			input: "I int | C (I int | S string)",
			node:  &Enum{},
			wantedNode: &Enum{
				Field{Name: "I", Type: TypeIdent("int")},
				Field{
					Name: "C",
					Type: TypeEnum(Enum{
						Field{Name: "I", Type: TypeIdent("int")},
						Field{Name: "S", Type: TypeIdent("string")},
					}),
				},
			},
		},
		{
			name:       "tuple-empty",
			input:      "()",
			node:       &Tuple{},
			wantedNode: &Tuple{},
		},
		{
			name:        "tuple-one-elt-fails",
			input:       "(int)",
			node:        &Tuple{},
			wantedError: true,
		},
		{
			name:       "tuple-two-elts",
			input:      "(int, string)",
			node:       &Tuple{},
			wantedNode: &Tuple{TypeIdent("int"), TypeIdent("string")},
		},
		{
			name:       "tuple-three-elts",
			input:      "(x, y, z)",
			node:       &Tuple{},
			wantedNode: &Tuple{TypeIdent("x"), TypeIdent("y"), TypeIdent("z")},
		},
		{
			name:       "struct-empty",
			input:      "{}",
			node:       &Struct{},
			wantedNode: &Struct{},
		},
		{
			name:  "struct-one-field",
			input: "{Name Ident}",
			node:  &Struct{},
			wantedNode: &Struct{
				Field{Name: "Name", Type: TypeIdent("Ident")},
			},
		},
		{
			name:  "struct-two-fields",
			input: "{Name Ident; Type Type}",
			node:  &Struct{},
			wantedNode: &Struct{
				Field{Name: "Name", Type: TypeIdent("Ident")},
				Field{Name: "Type", Type: TypeIdent("Type")},
			},
		},
		{
			name:       "struct-ws-padding",
			input:      "{ Name Ident }",
			node:       &Struct{},
			wantedNode: &Struct{Field{Name: "Name", Type: TypeIdent("Ident")}},
		},
		{
			name:  "struct-multiline",
			input: "{\n\tName Ident\nType Type\n}",
			node:  &Struct{},
			wantedNode: &Struct{
				Field{Name: "Name", Type: TypeIdent("Ident")},
				Field{Name: "Type", Type: TypeIdent("Type")},
			},
		},
		{
			name:       "slice",
			input:      "[]string",
			node:       &Slice{},
			wantedNode: &Slice{TypeIdent("string")},
		},
		{
			name:  "type-parens",
			input: "(I int | S string)",
			node:  &Type{},
			wantedNode: nodeType(TypeEnum(Enum{
				Field{Name: "I", Type: TypeIdent("int")},
				Field{Name: "S", Type: TypeIdent("string")},
			})),
		},
		{
			name:  "type-parens-nested",
			input: "((I int | S string))",
			node:  &Type{},
			wantedNode: nodeType(TypeEnum(Enum{
				Field{Name: "I", Type: TypeIdent("int")},
				Field{Name: "S", Type: TypeIdent("string")},
			})),
		},
		{
			name:  "type-decl-simple",
			input: "type Field = {Name Ident; Type Type}",
			node:  &TypeDecl{},
			wantedNode: &TypeDecl{Name: "Field", Type: TypeStruct(Struct{
				Field{Name: "Name", Type: TypeIdent("Ident")},
				Field{Name: "Type", Type: TypeIdent("Type")},
			})},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			rest, err := testCase.node.Parser().Parse(Input{
				Source: []byte(testCase.input),
				Offset: 0,
			})
			if err != nil {
				if testCase.wantedError {
					return
				}
				t.Fatal("Unexpected error:", err)
			}
			if testCase.wantedError {
				t.Fatalf(
					"Expected an error; got type:\n%# v",
					pretty.Formatter(testCase.node),
				)
			}
			if !testCase.node.EqualASTNode(testCase.wantedNode) {
				t.Fatalf(
					"Wanted:\n%# v\n\nGot:\n%# v\n\n%s",
					pretty.Formatter(testCase.wantedNode),
					pretty.Formatter(testCase.node),
					strings.Join(
						pretty.Diff(testCase.wantedNode, testCase.node),
						"\n",
					),
				)
			}
			if rest.String() != testCase.wantedRest {
				t.Fatalf(
					"Wanted rest:\n'%s'\n\nGot rest:\n'%s'",
					testCase.wantedRest,
					rest,
				)
			}
		})
	}
}

// Because &Ident("foo") is not legal in Go...
func nodeIdent(ident Ident) ASTNode { return &ident }

// Because &TypeIdent("foo") is not legal in Go...
func nodeType(typ Type) ASTNode { return &typ }
