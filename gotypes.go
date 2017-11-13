package main

import (
	"fmt"
	"strconv"
	"strings"
)

const indentString = "    "

type code interface {
	Render(indent int) string
}

func renderBlock(start, end, sep string, lines []code, indent int) string {
	if len(lines) < 1 {
		return start + end
	}
	baseIndent := strings.Repeat(indentString, indent)
	for _, line := range lines {
		start += "\n" + baseIndent + indentString + line.Render(indent+1) + sep
	}
	return start + "\n" + baseIndent + end
}

type GoType interface {
	code
	variantGoType()
}

type GoDecl interface {
	code
	variantGoDecl()
}

type GoExpr interface {
	GoStmt // all expressions are statements
	variantGoExpr()
}

type GoCase struct {
	Expr  GoExpr
	Stmts []GoStmt
}

func (gc GoCase) Render(indent int) string {
	stmts := make([]code, len(gc.Stmts))
	for i, stmt := range gc.Stmts {
		stmts[i] = stmt
	}
	return renderBlock(
		"case "+gc.Expr.Render(indent)+":",
		"",
		"",
		stmts,
		indent,
	)
}

type GoTypeSwitch struct {
	Assignment GoPair
	Cases      []GoCase
}

func (gts GoTypeSwitch) variantGoStmt() {}

func (gts GoTypeSwitch) Render(indent int) string {
	cases := make([]code, len(gts.Cases))
	for i, c := range gts.Cases {
		cases[i] = c
	}
	return renderBlock(
		fmt.Sprintf(
			"switch %s := %s.(type) {",
			gts.Assignment.Name.Render(indent-1),
			gts.Assignment.Value.Render(indent-1),
		),
		indentString+"}",
		"",
		cases,
		indent-1,
	)
}

type GoStmt interface {
	code
	variantGoStmt()
}

type GoIdent string

func (gi GoIdent) variantGoType() {}

func (gi GoIdent) variantGoExpr() {}

func (gi GoIdent) variantGoStmt() {}

func (gi GoIdent) Render(int) string { return string(gi) }

type GoField struct {
	Name GoIdent // Optional
	Type GoType
}

func (gf GoField) Render(indent int) string {
	if gf.Name == "" {
		return gf.Type.Render(indent)
	}
	return gf.Name.Render(indent) + " " + gf.Type.Render(indent)
}

type GoFuncType struct{ Args, Return []GoField }

func (gft GoFuncType) variantGoType() {}

func (gft GoFuncType) argsString(indent int) string {
	args := make([]string, len(gft.Args))
	for i, arg := range gft.Args {
		args[i] = arg.Render(indent)
	}
	return strings.Join(args, ", ")
}

func (gft GoFuncType) returnString(indent int) string {
	if len(gft.Return) > 1 ||
		(len(gft.Return) == 1 && gft.Return[0].Name != "") {
		// If there are multiple return values or if a single return value has
		// a label, then surround it with parentheses
		ret := make([]string, len(gft.Return))
		for i, arg := range gft.Return {
			ret[i] = arg.Render(indent)
		}
		return " (" + strings.Join(ret, ", ") + ")"
	} else if len(gft.Return) == 1 && gft.Return[0].Name == "" {
		// If there is exactly one return value and it is unlabeled, just
		// render the type
		return " " + gft.Return[0].Render(indent)
	}
	return ""
}

func (gft GoFuncType) Render(indent int) string {
	return "func(" + gft.argsString(indent) + ")" + gft.returnString(indent)
}

type GoInterfaceField struct {
	Name      GoIdent
	Signature GoFuncType
}

func (gif GoInterfaceField) Render(indent int) string {
	return fmt.Sprintf(
		"%s(%s)%s",
		gif.Name,
		gif.Signature.argsString(indent),
		gif.Signature.returnString(indent),
	)
}

type GoMethodSpec struct {
	Name GoIdent
	Recv GoType
	GoFuncType
}

func (gms GoMethodSpec) Render(indent int) string {
	return fmt.Sprintf(
		"(self %s) %s(%s)%s",
		gms.Recv.Render(indent),
		gms.Name,
		gms.argsString(indent),
		gms.returnString(indent),
	)
}

type GoStruct []GoField

func (gs GoStruct) variantGoType() {}

func (gs GoStruct) Render(indent int) string {
	fields := make([]code, len(gs))
	for i, field := range gs {
		fields[i] = field
	}
	return renderBlock("struct{", "}", "", fields, indent)
}

type GoInterface []GoInterfaceField

func (gi GoInterface) variantGoType() {}

func (gi GoInterface) Render(indent int) string {
	methods := make([]code, len(gi))
	for i, method := range gi {
		methods[i] = method
	}
	return renderBlock("interface {", "}", "", methods, indent)
}

type GoPointer struct {
	Type GoType
}

func (gp GoPointer) variantGoType() {}

func (gp GoPointer) Render(indent int) string {
	return "*" + gp.Type.Render(indent)
}

type GoSlice struct {
	Type GoType
}

func (gs GoSlice) variantGoType() {}

func (gs GoSlice) Render(indent int) string {
	return "[]" + gs.Type.Render(indent)
}

type GoBlock []GoStmt

func (gb GoBlock) Render(indent int) string {
	stmts := make([]code, len(gb))
	for i, stmt := range gb {
		stmts[i] = stmt
	}
	return renderBlock("{", "}", "", stmts, indent)
}

type GoPair struct {
	Name  GoIdent
	Value GoExpr
}

func (gp GoPair) Render(indent int) string {
	return gp.Name.Render(indent) + ": " + gp.Value.Render(indent)
}

type GoIntLit int

func (gil GoIntLit) variantGoExpr() {}

func (gil GoIntLit) variantGoStmt() {}

func (gil GoIntLit) Render(indent int) string { return strconv.Itoa(int(gil)) }

type GoStructLit struct {
	Type   GoType
	Fields []GoPair
}

func (gsl GoStructLit) variantGoExpr() {}

func (gsl GoStructLit) variantGoStmt() {}

func (gsl GoStructLit) Render(indent int) string {
	fs := make([]code, len(gsl.Fields))
	for i, field := range gsl.Fields {
		fs[i] = field
	}
	return gsl.Type.Render(indent) + renderBlock(" {", "}", ",", fs, indent)
}

type GoCallExpr struct {
	Fn   GoExpr
	Args []GoExpr
}

func (gce GoCallExpr) variantGoExpr() {}

func (gce GoCallExpr) variantGoStmt() {}

func (gce GoCallExpr) Render(indent int) string {
	args := make([]string, len(gce.Args))
	for i, arg := range gce.Args {
		args[i] = arg.Render(indent)
	}
	return gce.Fn.Render(indent) + "(" + strings.Join(args, ", ") + ")"
}

// type GoTypeCast struct {
// 	Type GoType
// 	Expr GoExpr
// }
//
// func (gtc GoTypeCast) variantGoExpr() {}
//
// func (gtc GoTypeCast) variantGoStmt() {}
//
// func (gtc GoTypeCast) Render(indent int) string {
// 	return GoCallExpr.Render(GoCallExpr{
// 		Fn:   GoIdent("(" + gtc.Type.Render(indent) + ")"),
// 		Args: []GoExpr{gtc.Expr},
// 	}, indent)
// }

type GoTypeDecl struct {
	Name GoIdent
	Type GoType
}

func (gtd GoTypeDecl) variantGoDecl() {}

func (gtd GoTypeDecl) Render(indent int) string {
	return "type " + gtd.Name.Render(indent) + " " + gtd.Type.Render(indent)
}

type GoFuncDecl struct {
	Name      GoIdent
	Signature GoFuncType
	Body      GoBlock
}

func (gfd GoFuncDecl) variantGoDecl() {}

func (gfd GoFuncDecl) Render(indent int) string {
	return fmt.Sprintf(
		"func %s(%s)%s %s",
		gfd.Name,
		gfd.Signature.argsString(indent),
		gfd.Signature.returnString(indent),
		gfd.Body.Render(indent),
	)
}

type GoMethodDecl struct {
	Signature GoMethodSpec
	Body      GoBlock
}

func (gmd GoMethodDecl) variantGoDecl() {}

func (gmd GoMethodDecl) Render(indent int) string {
	return fmt.Sprintf(
		"func %s %s",
		gmd.Signature.Render(indent),
		gmd.Body.Render(indent),
	)
}

type GoReturnStmt []GoExpr

func (grs GoReturnStmt) variantGoStmt() {}

func (grs GoReturnStmt) Render(indent int) string {
	exprs := make([]string, len(grs))
	for i, expr := range grs {
		exprs[i] = expr.Render(indent)
	}
	return "return " + strings.Join(exprs, ", ")
}

type GoFile struct {
	PackageName GoIdent
	Decls       []GoDecl
}

func (gf GoFile) Render() string {
	declStrings := make([]string, len(gf.Decls))
	for i, decl := range gf.Decls {
		declStrings[i] = decl.Render(0)
	}
	return strings.Join(
		append([]string{"package " + string(gf.PackageName)}, declStrings...),
		"\n\n",
	)
}
