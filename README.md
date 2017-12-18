README
------

Disclaimer: This is an experiment in code generation; it's largely a project
for my own amusement and utility. Don't use this in production.

`tdl` is a terse, expressive type description language. For now, I'm using it
to generate Go boilerplate, which is particularly useful for modeling sum
types. Besides generating type declarations, `tdl` could also be used to
generate serializers/deserializers/validators, and not just for Go but for any
language. This would make it usable as a simplified, expressive alternative for
protobufs. I'd really like to use it as the type grammar of a ML-like language
that compiles to Go.

If any of this interests you, drop me a line: weberc2@gmail.com

# TDL-GO

Right now, there is only one implementation (`tdl-go`)--a Go type definition
compiler (tdl goes in, Go type stubs come out). Eventually I want to build the
aforementioned applications as well as compilers for other languages
(`tdl-py` and `tdl-js` are high on my list).


# INSTALL

The repository contains `generated.go`, which is generated from `types.tdl`, so
all you need to install the compiler is the go toolchain: `go install
./tdl-go`.

# USAGE

`tdl-go` currently reads from stdin and writes to stdout, so you invoke it like
this: `cat ./types.tdl | tdl-go > ./generated.go`.

# EXAMPLES

`tdl-go` is self-hosted, so check out [its type description][0] and [the
generated output][1]. Also, ./examples/json contains a toy program that
serializes Go data types to JSON by first converting them to generated JSON
types.

Here is the type description for TDL's AST:

```
type Ident = string

type Field = {Name Ident; Type Type}

type Enum = []Field

type Struct = []Field

type Tuple = []Type

type Pointer = {Type Type}

type Slice = {Type Type}

type Type = Ident Ident
    | Enum Enum
    | Struct Struct
    | Tuple Tuple
    | Pointer Pointer
    | Slice Slice

type TypeDecl = {Name Ident; Type Type}
```

[0]:./tdl-go/types.tdl
[1]:./tdl-go/generated.go
