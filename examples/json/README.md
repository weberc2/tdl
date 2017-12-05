JSON EXAMPLE
------------

This is a toy example that demonstrates what I consider to be a much nicer
API for specifying custom JSON marshaling behavior than Go's `encoding/json`
package currently provides. Specifically, I prefer to deal with a typed
representation of JSON than the raw-bytes interface that `encoding/json`
requires. It's also conceivably more performant than the standard reflection-
based serializing/deserializing, but I don't have any benchmarks (or indeed a
working non-toy implementation).

# USAGE

This assumes you have the tdl-go compiler binary installed somewhere in
your $PATH.

1. `cat ./json.tdl | tdl-go > ./types.go`
2. `go run *.go`

The type specification is simple in tdl, but complex in Go because many JSON
types are sum types. Here is the tdl specification:

```
type Field = {Name string; Value JSON}

type JSON = Array []JSON
    | Object []Field
    | String string
    | Number float64
    | Boolean bool
    | Null ()
```
