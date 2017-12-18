// This file is generated. Do not modify.

package main

type BinExpr struct{
    Operator 5d38fbd0
    Left Column
    Right Column
}

type _ColumnVariant interface {
    _ColumnVariant()
}

type Column struct{
    variant _ColumnVariant
}

type _ColumnLiteral struct{
    π fa54d14c
}

func ColumnLiteral(π fa54d14c) Column {
    return Column {
        variant: _ColumnLiteral {
            π: π,
        },
    }
}

func (self _ColumnLiteral) _ColumnVariant() {}

type _ColumnColumnName struct{
    π string
}

func ColumnColumnName(π string) Column {
    return Column {
        variant: _ColumnColumnName {
            π: π,
        },
    }
}

func (self _ColumnColumnName) _ColumnVariant() {}

type _ColumnFunction struct{
    π be9035bd
}

func ColumnFunction(π be9035bd) Column {
    return Column {
        variant: _ColumnFunction {
            π: π,
        },
    }
}

func (self _ColumnFunction) _ColumnVariant() {}

func (self Column) Match(Literal func(π fa54d14c), ColumnName func(π string), Function func(π be9035bd)) {
    switch π := self.variant.(type) {
    case _ColumnLiteral:
        Literal(π.π)
        return 
    
    case _ColumnColumnName:
        ColumnName(π.π)
        return 
    
    case _ColumnFunction:
        Function(π.π)
        return 
    
    }
}

type Comparison struct{
    Operator 094fd78d
    Left Column
    Right Column
}

type _BoolExprVariant interface {
    _BoolExprVariant()
}

type BoolExpr struct{
    variant _BoolExprVariant
}

type _BoolExprComparison struct{
    π Comparison
}

func BoolExprComparison(π Comparison) BoolExpr {
    return BoolExpr {
        variant: _BoolExprComparison {
            π: π,
        },
    }
}

func (self _BoolExprComparison) _BoolExprVariant() {}

type _BoolExprAnd struct{
    π []BoolExpr
}

func BoolExprAnd(π []BoolExpr) BoolExpr {
    return BoolExpr {
        variant: _BoolExprAnd {
            π: π,
        },
    }
}

func (self _BoolExprAnd) _BoolExprVariant() {}

type _BoolExprOr struct{
    π []BoolExpr
}

func BoolExprOr(π []BoolExpr) BoolExpr {
    return BoolExpr {
        variant: _BoolExprOr {
            π: π,
        },
    }
}

func (self _BoolExprOr) _BoolExprVariant() {}

type _BoolExprBoolLit struct{
    π bool
}

func BoolExprBoolLit(π bool) BoolExpr {
    return BoolExpr {
        variant: _BoolExprBoolLit {
            π: π,
        },
    }
}

func (self _BoolExprBoolLit) _BoolExprVariant() {}

func (self BoolExpr) Match(Comparison func(π Comparison), And func(π []BoolExpr), Or func(π []BoolExpr), BoolLit func(π bool)) {
    switch π := self.variant.(type) {
    case _BoolExprComparison:
        Comparison(π.π)
        return 
    
    case _BoolExprAnd:
        And(π.π)
        return 
    
    case _BoolExprOr:
        Or(π.π)
        return 
    
    case _BoolExprBoolLit:
        BoolLit(π.π)
        return 
    
    }
}

type Query struct{
    Select []Column
    From b72dd03f
    Where BoolExpr
    OrderBy []struct{
        Column Column
        Descending bool
    }
    Limit 4b7d63ad
}
