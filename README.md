README
------

Disclaimer: This is all a pre-alpha prototype; it's largely a project for my
own amusement and utility. As such, I may stop supporting it at any time. Also,
everything will likely change very quickly. Don't use this for anything
serious.

`tdl` is a type description language--a language for describing data types.
By itself, this isn't very useful, but I think it may have many potential
metaprogramming applications:

* Generating type definitions for verbose programming languages
    - Most popular programming languages don't have first class support for sum
      types, and the workarounds tend to be verbose and/or not typesafe
* Generating serializers and deserializers
    - Have a multilingual system? Describe the types with tdl and then generate
      the implementations in each target language. Guarantee compatibility and
      avoid tediously keeping redundant type definitions in sync.
* Generating parsers
* Generating API server/client stubs (interface description language)
* Generating validators (e.g., jsonschema)
    - This is useful for validating the structure of data received over an API,
      stored in a schemaless database, etc
* TDL could be a subset of a type system for a new programming language or part
  of a hygienic macro system for an existing language

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
