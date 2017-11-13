package main

import (
	"fmt"
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

type OneOrMore struct {
	Parser
}

func (oom OneOrMore) Parse(input Input) (Input, error) {
	return Seq{oom.Parser, Repeat{oom.Parser}}.Parse(input)
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

type ParseError struct {
	Cause   error
	Message string
}

func (p ParseError) Error() string {
	return fmt.Sprintf("%s\n%s", p.Cause.Error(), p.Message)
}

var (
	CanSpace = Repeat{UnicodeClassSpace}
	CanWS    = Repeat{UnicodeClassWhiteSpace}
	EOF      = RuneLit(eof)
	EOS      = Seq{
		CanSpace,
		Any{RuneLit(';'), RuneLit('\n')},
		CanSpace,
	}
	Space = OneOrMore{UnicodeClassSpace}
	WS    = OneOrMore{UnicodeClassWhiteSpace}
)
