package main

import (
	"testing"
)

func TestCombinator(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		parser      Parser
		wantedRest  string
		wantedError bool
	}{
		{
			name:       "string-lit-no-remainder",
			input:      "foo",
			parser:     StringLit("foo"),
			wantedRest: "",
		},
		{
			name:       "string-lit-w-remainder",
			input:      "foobar",
			parser:     StringLit("foo"),
			wantedRest: "bar",
		},
		{
			name:       "unicode-class-letter",
			input:      "a",
			parser:     UnicodeClassLetter,
			wantedRest: "",
		},
		{
			name:        "unicode-class-letter-fails-for-digits",
			input:       "1",
			parser:      UnicodeClassLetter,
			wantedError: true,
		},
		{
			name:        "unicode-class-letter-fails-for-punctionation",
			input:       "*",
			parser:      UnicodeClassLetter,
			wantedError: true,
		},
		{
			name:       "seq-empty",
			input:      "foo",
			parser:     Seq{},
			wantedRest: "foo",
		},
		{
			name:       "seq-one-elt",
			input:      "foo",
			parser:     Seq{StringLit("foo")},
			wantedRest: "",
		},
		{
			name:  "seq-multi-elts",
			input: "foobarbaz",
			parser: Seq{
				StringLit("foo"),
				StringLit("bar"),
				StringLit("baz"),
			},
			wantedRest: "",
		},
		{
			name:       "repeat",
			input:      "abc123",
			parser:     Repeat{UnicodeClassLetter},
			wantedRest: "123",
		},
		{
			name:       "repeat-no-matches",
			input:      "123",
			parser:     Repeat{UnicodeClassLetter},
			wantedRest: "123",
		},
		{
			name:       "one-or-more-one-match",
			input:      "foobar",
			parser:     OneOrMore{StringLit("foo")},
			wantedRest: "bar",
		},
		{
			name:       "one-or-more-multi-match",
			input:      "foofoofoobar",
			parser:     OneOrMore{StringLit("foo")},
			wantedRest: "bar",
		},
		{
			name:        "one-or-more-fails-if-no-matches",
			input:       "bar",
			parser:      OneOrMore{StringLit("foo")},
			wantedError: true,
		},
		{
			name:       "any",
			input:      "barbaz",
			parser:     Any{StringLit("foo"), StringLit("bar")},
			wantedRest: "baz",
		},
		{
			name:       "opt-match",
			input:      "abc",
			parser:     Opt{RuneLit('a')},
			wantedRest: "bc",
		},
		{
			name:       "opt-no-match",
			input:      "123",
			parser:     Opt{RuneLit('a')},
			wantedRest: "123",
		},
		// This test is documentation. UnicodeClassSpace is horizontal
		// whitespace characters *except* tabs (for some reason). We should
		// take care not to confuse it with *all* horizontal spaces.
		{
			name:        "unicode-class-space-fails-for-tab",
			input:       "\t",
			parser:      CanSpace,
			wantedRest:  "\t",
			wantedError: true,
		},
		{
			name:   "can-horiz-ws",
			input:  "\t",
			parser: CanHorizWS,
		},
		{
			name:   "eos-newline-simple",
			input:  "x\ny",
			parser: Seq{RuneLit('x'), EOS, RuneLit('y')},
		},
		{
			name:   "eos-semicolon-simple",
			input:  "x;y",
			parser: Seq{RuneLit('x'), EOS, RuneLit('y')},
		},
		{
			name:   "eos-newline-padded-ws",
			input:  "x\t\n\ty",
			parser: Seq{RuneLit('x'), EOS, RuneLit('y')},
		},
		{
			name:   "eos-semicolon-padded-ws",
			input:  "x\t;\ty",
			parser: Seq{RuneLit('x'), EOS, RuneLit('y')},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			rest, err := testCase.parser.Parse(
				Input{Source: []byte(testCase.input)},
			)
			if err != nil {
				if testCase.wantedError {
					return
				}
				t.Fatal("Unexpected error:", err)
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
