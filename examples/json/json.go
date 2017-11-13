package main

type Field struct{
    Name string
    Value JSON
}

type _JSONVariant interface {
    _JSONVariant()
}

type JSON struct{
    variant _JSONVariant
}

type _JSONArray struct{
    π []JSON
}

func JSONArray(π []JSON) JSON {
    return JSON {
        variant: _JSONArray {
            π: π,
        },
    }
}

func (self _JSONArray) _JSONVariant() {}

type _JSONObject struct{
    π []Field
}

func JSONObject(π []Field) JSON {
    return JSON {
        variant: _JSONObject {
            π: π,
        },
    }
}

func (self _JSONObject) _JSONVariant() {}

type _JSONString struct{
    π string
}

func JSONString(π string) JSON {
    return JSON {
        variant: _JSONString {
            π: π,
        },
    }
}

func (self _JSONString) _JSONVariant() {}

type _JSONNumber struct{
    π float64
}

func JSONNumber(π float64) JSON {
    return JSON {
        variant: _JSONNumber {
            π: π,
        },
    }
}

func (self _JSONNumber) _JSONVariant() {}

type _JSONBoolean struct{
    π bool
}

func JSONBoolean(π bool) JSON {
    return JSON {
        variant: _JSONBoolean {
            π: π,
        },
    }
}

func (self _JSONBoolean) _JSONVariant() {}

type _JSONNull struct{
    π struct{}
}

func JSONNull(π struct{}) JSON {
    return JSON {
        variant: _JSONNull {
            π: π,
        },
    }
}

func (self _JSONNull) _JSONVariant() {}

func (self JSON) Match(Array func(π []JSON), Object func(π []Field), String func(π string), Number func(π float64), Boolean func(π bool), Null func(π struct{})) {
    switch π := self.variant.(type) {
    case _JSONArray:
        Array(π.π)
        return 
    
    case _JSONObject:
        Object(π.π)
        return 
    
    case _JSONString:
        String(π.π)
        return 
    
    case _JSONNumber:
        Number(π.π)
        return 
    
    case _JSONBoolean:
        Boolean(π.π)
        return 
    
    case _JSONNull:
        Null(π.π)
        return 
    
    }
}
