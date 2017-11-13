package main

import (
	c "github.com/weberc2/gallium/combinator"
)

func parenGroup(input c.Input) c.Result {
	return c.Seq(
		c.Lit('('),
		c.CanWS,
		pType,
		c.CanWS,
		c.Lit(')'),
	).Get(2).Rename("ParenGroup")(input)
}

func atom(input c.Input) c.Result {
	return c.Any(
		pTuple,
		parenGroup,
		pStruct,
		pSlice,
		pIdent,
		pPointer,
	).Rename("Atom")(input)
}

var pIdent = c.Ident.Map(func(v interface{}) interface{} {
	return Ident(v.(string))
})

func pField(input c.Input) c.Result {
	return c.Seq(pIdent, c.WS, pType).MapSlice(
		func(vs []interface{}) interface{} {
			return Field{Name: vs[0].(Ident), Type: vs[2].(Type)}
		},
	).Rename("Field")(input)
}

func pFields(
	name string,
	field c.Parser,
	sep c.Parser,
	combine func(c.Parser) c.Parser,
) c.Parser {
	return func(input c.Input) c.Result {
		return c.Seq(
			combine(c.Seq(field, sep).Get(0)),
			field,
		).MapSlice(func(vs []interface{}) interface{} {
			return append(vs[0].([]interface{}), vs[1])
		}).Rename(name)(input)
	}
}

func pEnum(input c.Input) c.Result {
	return pFields(
		"Enum",
		pField,
		c.Seq(c.CanWS, c.Lit('|'), c.CanWS),
		c.OneOrMore,
	).MapSlice(func(vs []interface{}) interface{} {
		enum := make(Enum, len(vs))
		for i, v := range vs {
			enum[i] = v.(Field)
		}
		return enum
	}).Rename("Enum")(input)
}

func pTuple(input c.Input) c.Result {
	return c.Any(
		c.Seq(c.Lit('('), c.CanWS, c.Lit(')')).Map(
			func(v interface{}) interface{} { return Tuple{} },
		),
		c.Seq(
			c.Lit('('),
			c.CanWS,
			pFields(
				"Tuple",
				pType,
				c.Seq(c.CanWS, c.Lit(','), c.CanWS),
				c.OneOrMore,
			),
			c.CanWS,
			c.Lit(')'),
		).Get(2).MapSlice(func(vs []interface{}) interface{} {
			tuple := make(Tuple, len(vs))
			for i, v := range vs {
				tuple[i] = v.(Type)
			}
			return tuple
		}),
	).Rename("Tuple")(input)
}

func pPointer(input c.Input) c.Result {
	return c.Seq(c.Lit('*'), c.CanWS, atom).Get(2).Map(
		func(v interface{}) interface{} {
			return Pointer{v.(Type)}
		},
	).Rename("Pointer")(input)
}

func pStruct(input c.Input) c.Result {
	return c.Seq(
		c.Lit('{'),
		c.CanWS,
		c.Repeat(c.Seq(
			pField,
			c.Repeat(c.IsClass(c.UnicodeClassSpace)),
			c.Any(c.Lit(';'), c.Lit('\n')),
			c.CanWS,
		).Get(0)).MapSlice(func(vs []interface{}) interface{} {
			fields := make(Struct, len(vs))
			for i, v := range vs {
				fields[i] = v.(Field)
			}
			return fields
		}),
		c.Opt(pField),
		c.Lit('}'),
	).MapSlice(func(vs []interface{}) interface{} {
		if vs[3] != nil {
			return append(vs[2].(Struct), vs[3].(Field))
		}
		return vs[2].(Struct)
	}).Rename("Struct")(input)
}

func pType(input c.Input) c.Result {
	return c.Any(pEnum, pStruct, atom).Rename("Type")(input)
}

func pTypeDecl(input c.Input) c.Result {
	return c.Seq(
		c.StrLit("type"),
		c.WS,
		pIdent,
		c.CanWS,
		c.Lit('='),
		c.CanWS,
		pType,
	).MapSlice(func(vs []interface{}) interface{} {
		if e, ok := vs[6].(Enum); ok {
			return TypeDecl{vs[2].(Ident), e}
		}
		return TypeDecl{vs[2].(Ident), vs[6].(Type)}
	}).Rename("TypeDecl")(input)
}

func pSlice(input c.Input) c.Result {
	return c.Seq(c.StrLit("[]"), pType).Get(1).Map(
		func(v interface{}) interface{} { return Slice{v.(Type)} },
	).Rename("Slice")(input)
}

var ParseFile = c.Seq(
	c.Seq(c.StrLit("package"), c.WS, pIdent).Get(2),
	c.CanWS,
	c.Repeat(c.Seq(pTypeDecl, c.CanWS).Get(0)).MapSlice(
		func(vs []interface{}) interface{} {
			decls := make([]TypeDecl, len(vs))
			for i, v := range vs {
				decls[i] = v.(TypeDecl)
			}
			return decls
		},
	),
	c.EOF,
).MapSlice(func(vs []interface{}) interface{} {
	return File{
		PackageName: vs[0].(Ident),
		TypeDecls:   vs[2].([]TypeDecl),
	}
}).Rename("ParseFile")
