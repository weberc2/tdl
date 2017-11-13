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

func (gi GoIdent) Render(indent int) string { return string(gi) }

func (gil GoIntLit) Render(indent int) string { return strconv.Itoa(int(gil)) }

func (gf GoField) Render(indent int) string {
	if gf.Name == "" {
		return gf.Type.Render(indent)
	}
	return gf.Name.Render(indent) + " " + gf.Type.Render(indent)
}

func (gp GoPair) Render(indent int) string {
	return gp.Name.Render(indent) + ": " + gp.Value.Render(indent)
}

func (gsl GoStructLit) Render(indent int) string {
	fs := make([]code, len(gsl.Fields))
	for i, field := range gsl.Fields {
		fs[i] = field
	}
	return gsl.Type.Render(indent) + renderBlock(" {", "}", ",", fs, indent)
}

func (gce GoCallExpr) Render(indent int) string {
	args := make([]string, len(gce.Args))
	for i, arg := range gce.Args {
		args[i] = arg.Render(indent)
	}
	return gce.Fn.Render(indent) + "(" + strings.Join(args, ", ") + ")"
}

func (expr GoExpr) Render(indent int) string {
	var result string
	expr.Match(
		func(gi GoIdent) { result = gi.Render(indent) },
		func(gil GoIntLit) { result = gil.Render(indent) },
		func(gsl GoStructLit) { result = gsl.Render(indent) },
		func(gce GoCallExpr) { result = gce.Render(indent) },
	)
	return result
}

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

func (gs GoStruct) Render(indent int) string {
	fields := make([]code, len(gs))
	for i, field := range gs {
		fields[i] = field
	}
	return renderBlock("struct{", "}", "", fields, indent)
}

func (gif GoInterfaceField) Render(indent int) string {
	return fmt.Sprintf(
		"%s(%s)%s",
		gif.Name,
		gif.Signature.argsString(indent),
		gif.Signature.returnString(indent),
	)
}

func (gi GoInterface) Render(indent int) string {
	interfaceFields := make([]code, len(gi))
	for i, interfaceField := range gi {
		interfaceFields[i] = interfaceField
	}
	return renderBlock("interface {", "}", "", interfaceFields, indent)
}

func (gp GoPointer) Render(indent int) string {
	return "*" + gp.Type.Render(indent)
}

func (gs GoSlice) Render(indent int) string {
	return "[]" + gs.Type.Render(indent)
}

func (gt GoType) Render(indent int) string {
	var result string
	gt.Match(
		func(gi GoIdent) { result = gi.Render(indent) },
		func(gft GoFuncType) { result = gft.Render(indent) },
		func(gst GoStruct) { result = gst.Render(indent) },
		func(gi GoInterface) { result = gi.Render(indent) },
		func(gp GoPointer) { result = gp.Render(indent) },
		func(gs GoSlice) { result = gs.Render(indent) },
	)
	return result
}

func (gtd GoTypeDecl) Render(indent int) string {
	return "type " + gtd.Name.Render(indent) + " " + gtd.Type.Render(indent)
}

func (grs GoReturnStmt) Render(indent int) string {
	exprs := make([]string, len(grs))
	for i, expr := range grs {
		exprs[i] = expr.Render(indent)
	}
	return "return " + strings.Join(exprs, ", ")
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

func (gs GoStmt) Render(indent int) string {
	var result string
	gs.Match(
		func(ge GoExpr) { result = ge.Render(indent) },
		func(grs GoReturnStmt) { result = grs.Render(indent) },
		func(gts GoTypeSwitch) { result = gts.Render(indent) },
	)
	return result
}

func (gb GoBlock) Render(indent int) string {
	stmts := make([]code, len(gb))
	for i, stmt := range gb {
		stmts[i] = stmt
	}
	return renderBlock("{", "}", "", stmts, indent)
}

func (gfd GoFuncDecl) Render(indent int) string {
	return fmt.Sprintf(
		"func %s(%s)%s %s",
		gfd.Name,
		gfd.Signature.argsString(indent),
		gfd.Signature.returnString(indent),
		gfd.Body.Render(indent),
	)
}

func (gms GoMethodSpec) Render(indent int) string {
	return fmt.Sprintf(
		"(self %s) %s(%s)%s",
		gms.Recv.Render(indent),
		gms.Name,
		gms.Signature.argsString(indent),
		gms.Signature.returnString(indent),
	)
}

func (gmd GoMethodDecl) Render(indent int) string {
	return fmt.Sprintf(
		"func %s %s",
		gmd.Signature.Render(indent),
		gmd.Body.Render(indent),
	)
}

func (gd GoDecl) Render(indent int) string {
	var result string
	gd.Match(
		func(gtd GoTypeDecl) { result = gtd.Render(indent) },
		func(gfd GoFuncDecl) { result = gfd.Render(indent) },
		func(gmd GoMethodDecl) { result = gmd.Render(indent) },
	)
	return result
}

func (gc GoComment) Render(indent int) string {
	ind := strings.Repeat(indentString, indent)
	commentPrefix := "// "
	// copied and modified from
	// https://gist.github.com/kennwhite/306317d81ab4a885a965e25aa835b8ef
	wordWrap := func(text string, lineWidth int) string {
		lineWidth = lineWidth - len([]rune(ind)) - len([]rune(commentPrefix))
		words := strings.Fields(strings.TrimSpace(text))
		if len(words) == 0 {
			return text
		}
		wrapped := words[0]
		spaceLeft := lineWidth - len(wrapped)
		for _, word := range words[1:] {
			if len(word)+1 > spaceLeft {
				wrapped += "\n" + commentPrefix + ind + word
				spaceLeft = lineWidth - len(word)
			} else {
				wrapped += " " + word
				spaceLeft -= 1 + len(word)
			}
		}

		return wrapped
	}

	return wordWrap(commentPrefix+string(gc), 80)
}

func (gf GoFile) Render() string {
	declStrings := make([]string, len(gf.Decls))
	for i, decl := range gf.Decls {
		declStrings[i] = decl.Render(0)
	}
	return strings.Join(
		append(
			[]string{
				gf.FileComment.Render(0),
				"package " + string(gf.PackageName),
			},
			declStrings...,
		),
		"\n\n",
	)
}
