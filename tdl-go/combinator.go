package main

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const eof = rune(0)

type Input struct {
	Source []byte
	Offset int
}

func (i Input) Len() int { return len(i.Source) - i.Offset }

func (i Input) Cons() (rune, Input) {
	r, sz := utf8.DecodeRune([]byte(i.Source[i.Offset:]))
	if r == utf8.RuneError {
		if sz == 0 { // According to utf8 docs, (RuneError, 0) means eof
			return 0, Input{}
		}
		// According to utf8 docs, (RuneError, 1) means invalid utf-8
		panic("Invalid utf-8")
	}
	return r, Input{Source: i.Source, Offset: i.Offset + sz}
}

func (i Input) Take(n int) Input {
	if i.Len() <= n {
		return i
	}
	return Input{i.Source[:i.Offset+n], i.Offset}
}

func (i Input) String() string { return string(i.Source[i.Offset:]) }

type Parser interface {
	Parse(input Input) (Input, error)
	Name() string
}

type RuneLit rune

func (rl RuneLit) Parse(input Input) (Input, error) {
	head, rest := input.Cons()
	if head != rune(rl) {
		return Input{}, fmt.Errorf(
			"Wanted '%s'; got '%s...'",
			string(rl),
			input.Take(20),
		)
	}
	return rest, nil
}

func (rl RuneLit) Name() string {
	return "RuneLit['" + string(rl) + "']"
}

type StringLit string

func (sl StringLit) Parse(input Input) (Input, error) {
	start := input
	var head rune
	for _, r := range []rune(string(sl)) {
		if head, input = input.Cons(); head == r {
			continue
		}
		return Input{}, fmt.Errorf(
			"Wanted '%s'; got '%s...'",
			sl,
			start.Take(20),
		)
	}
	return input, nil
}

func (sl StringLit) Name() string {
	return "StringLit[\"" + string(sl) + "\"]"
}

type UnicodeClass struct {
	ClassName  string
	RangeTable *unicode.RangeTable
}

func (uc UnicodeClass) Parse(input Input) (Input, error) {
	head, tail := input.Cons()
	if unicode.Is(uc.RangeTable, head) {
		return tail, nil
	}
	return Input{}, fmt.Errorf(
		"Wanted unicode.%s; got '%s...'",
		uc.ClassName,
		input.Take(20),
	)
}

func (uc UnicodeClass) Name() string {
	return "UnicodeClass[" + uc.ClassName + "]"
}

var (
	UnicodeClassLetter     = UnicodeClass{"LETTER", unicode.Letter}
	UnicodeClassDigit      = UnicodeClass{"DIGIT", unicode.Digit}
	UnicodeClassSpace      = UnicodeClass{"SPACE", unicode.Space}
	UnicodeClassWhiteSpace = UnicodeClass{"WHITESPACE", unicode.White_Space}
)

type Seq []Parser

func (s Seq) Parse(input Input) (Input, error) {
	for _, p := range s {
		rest, err := p.Parse(input)
		if err != nil {
			return Input{}, err
		}
		input = rest
	}
	return input, nil
}

func (s Seq) Name() string {
	parserNames := make([]string, len(s))
	for i, parser := range s {
		parserNames[i] = parser.Name()
	}
	return "Seq[" + strings.Join(parserNames, ", ") + "]"
}

type Any []Parser

func (a Any) Parse(input Input) (Input, error) {
	for _, p := range a {
		rest, err := p.Parse(input)
		if err != nil {
			continue
		}
		return rest, nil
	}
	return Input{}, fmt.Errorf("Failed to match '%s...'", input.Take(20))
}

func (a Any) Name() string {
	parserNames := make([]string, len(a))
	for i, parser := range a {
		parserNames[i] = parser.Name()
	}
	return "Any[" + strings.Join(parserNames, ", ") + "]"
}

type Not struct {
	Parser
}

func (n Not) Parse(input Input) (Input, error) {
	if _, err := n.Parser.Parse(input); err != nil {
		return input, nil
	}
	return Input{}, fmt.Errorf("Matched %s", n.Name())
}

func (n Not) Name() string {
	return "Not[" + n.Parser.Name() + "]"
}

type Repeat struct {
	Parser
}

func (r Repeat) Parse(input Input) (Input, error) {
	for {
		rest, err := r.Parser.Parse(input)
		if err != nil {
			return input, nil
		}
		input = rest
	}
}

func (r Repeat) Name() string {
	return "Repeat[" + r.Parser.Name() + "]"
}

type OneOrMore struct {
	Parser
}

func (oom OneOrMore) Parse(input Input) (Input, error) {
	return Seq{oom.Parser, Repeat{oom.Parser}}.Parse(input)
}

func (oom OneOrMore) Name() string {
	return "OneOrMore[" + oom.Parser.Name() + "]"
}

type Opt struct {
	Parser
}

func (o Opt) Parse(input Input) (Input, error) {
	rest, err := o.Parser.Parse(input)
	if err != nil {
		return input, nil
	}
	return rest, nil
}

func (o Opt) Name() string { return "Opt[" + o.Parser.Name() + "]" }

type Rename struct {
	Parser
	NewName string
}

func (r Rename) Name() string {
	return r.NewName
}

type sideEffectParser struct {
	parser     Parser
	sideEffect func()
}

func (sep sideEffectParser) Parse(input Input) (Input, error) {
	rest, err := sep.parser.Parse(input)
	if err != nil {
		return Input{}, err
	}
	sep.sideEffect()
	return rest, nil
}

func (sep sideEffectParser) Name() string { return sep.parser.Name() }

type ParseError struct {
	Cause   error
	Message string
}

func (p ParseError) Error() string {
	return fmt.Sprintf("%s\n%s", p.Cause.Error(), p.Message)
}

var (
	CanSpace = Rename{Repeat{UnicodeClassSpace}, "CanSpace"}
	CanWS    = Rename{Repeat{UnicodeClassWhiteSpace}, "CanWS"}
	EOF      = Rename{RuneLit(eof), "EOF"}
	EOS      = Rename{
		Seq{CanSpace, Any{RuneLit(';'), RuneLit('\n')}, CanSpace},
		"EOS",
	}
	Space = Rename{OneOrMore{UnicodeClassSpace}, "Space"}
	WS    = Rename{OneOrMore{UnicodeClassWhiteSpace}, "WS"}
)
