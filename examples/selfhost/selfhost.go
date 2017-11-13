package main

type Ident string

type Field struct{
    Name Ident
    Type Type
}

type Enum []Field

type Struct []Field

type Tuple []Type

type Pointer struct{
    Type Type
}

type Slice struct{
    Type Type
}

type _TypeVariant interface {
    _TypeVariant()
}

type Type struct{
    variant _TypeVariant
}

type _TypeIdent struct{
    π Ident
}

func TypeIdent(π Ident) Type {
    return Type {
        variant: _TypeIdent {
            π: π,
        },
    }
}

func (self _TypeIdent) _TypeVariant() {}

type _TypeEnum struct{
    π Enum
}

func TypeEnum(π Enum) Type {
    return Type {
        variant: _TypeEnum {
            π: π,
        },
    }
}

func (self _TypeEnum) _TypeVariant() {}

type _TypeStruct struct{
    π Struct
}

func TypeStruct(π Struct) Type {
    return Type {
        variant: _TypeStruct {
            π: π,
        },
    }
}

func (self _TypeStruct) _TypeVariant() {}

type _TypeTuple struct{
    π Tuple
}

func TypeTuple(π Tuple) Type {
    return Type {
        variant: _TypeTuple {
            π: π,
        },
    }
}

func (self _TypeTuple) _TypeVariant() {}

type _TypePointer struct{
    π Pointer
}

func TypePointer(π Pointer) Type {
    return Type {
        variant: _TypePointer {
            π: π,
        },
    }
}

func (self _TypePointer) _TypeVariant() {}

type _TypeSlice struct{
    π Slice
}

func TypeSlice(π Slice) Type {
    return Type {
        variant: _TypeSlice {
            π: π,
        },
    }
}

func (self _TypeSlice) _TypeVariant() {}

func (self Type) Match(Ident func(π Ident), Enum func(π Enum), Struct func(π Struct), Tuple func(π Tuple), Pointer func(π Pointer), Slice func(π Slice)) {
    switch π := self.variant.(type) {
    case _TypeIdent:
        Ident(π.π)
        return 
    
    case _TypeEnum:
        Enum(π.π)
        return 
    
    case _TypeStruct:
        Struct(π.π)
        return 
    
    case _TypeTuple:
        Tuple(π.π)
        return 
    
    case _TypePointer:
        Pointer(π.π)
        return 
    
    case _TypeSlice:
        Slice(π.π)
        return 
    
    }
}

type TypeDecl struct{
    Name Ident
    Type Type
}

type File struct{
    PackageName Ident
    TypeDecls []TypeDecl
}

type GoIdent string

type GoPair struct{
    Name GoIdent
    Value GoExpr
}

type GoCase struct{
    Expr GoExpr
    Stmts []GoStmt
}

type GoTypeSwitch struct{
    Assignment GoPair
    Cases []GoCase
}

type GoField struct{
    Name GoIdent
    Type GoType
}

type GoFuncType struct{
    Args []GoField
    Return []GoField
}

type GoInterfaceField struct{
    Name GoIdent
    Signature GoFuncType
}

type GoMethodSpec struct{
    Name GoIdent
    Recv GoType
    Signature GoFuncType
}

type GoStruct []GoField

type GoInterface []GoInterfaceField

type GoPointer struct{
    Type GoType
}

type GoSlice struct{
    Type GoType
}

type GoBlock []GoStmt

type GoIntLit int

type GoStructLit struct{
    Type GoType
    Fields []GoPair
}

type GoCallExpr struct{
    Fn GoExpr
    Args []GoExpr
}

type GoTypeDecl struct{
    Name GoIdent
    Type GoType
}

type GoFuncDecl struct{
    Name GoIdent
    Signature GoFuncType
    Body GoBlock
}

type GoMethodDecl struct{
    Signature GoMethodSpec
    Body GoBlock
}

type GoReturnStmt []GoExpr

type _GoTypeVariant interface {
    _GoTypeVariant()
}

type GoType struct{
    variant _GoTypeVariant
}

type _GoTypeIdent struct{
    π GoIdent
}

func GoTypeIdent(π GoIdent) GoType {
    return GoType {
        variant: _GoTypeIdent {
            π: π,
        },
    }
}

func (self _GoTypeIdent) _GoTypeVariant() {}

type _GoTypeFunc struct{
    π GoFuncType
}

func GoTypeFunc(π GoFuncType) GoType {
    return GoType {
        variant: _GoTypeFunc {
            π: π,
        },
    }
}

func (self _GoTypeFunc) _GoTypeVariant() {}

type _GoTypeStruct struct{
    π GoStruct
}

func GoTypeStruct(π GoStruct) GoType {
    return GoType {
        variant: _GoTypeStruct {
            π: π,
        },
    }
}

func (self _GoTypeStruct) _GoTypeVariant() {}

type _GoTypeInterface struct{
    π GoInterface
}

func GoTypeInterface(π GoInterface) GoType {
    return GoType {
        variant: _GoTypeInterface {
            π: π,
        },
    }
}

func (self _GoTypeInterface) _GoTypeVariant() {}

type _GoTypePointer struct{
    π GoPointer
}

func GoTypePointer(π GoPointer) GoType {
    return GoType {
        variant: _GoTypePointer {
            π: π,
        },
    }
}

func (self _GoTypePointer) _GoTypeVariant() {}

type _GoTypeSlice struct{
    π GoSlice
}

func GoTypeSlice(π GoSlice) GoType {
    return GoType {
        variant: _GoTypeSlice {
            π: π,
        },
    }
}

func (self _GoTypeSlice) _GoTypeVariant() {}

func (self GoType) Match(Ident func(π GoIdent), Func func(π GoFuncType), Struct func(π GoStruct), Interface func(π GoInterface), Pointer func(π GoPointer), Slice func(π GoSlice)) {
    switch π := self.variant.(type) {
    case _GoTypeIdent:
        Ident(π.π)
        return 
    
    case _GoTypeFunc:
        Func(π.π)
        return 
    
    case _GoTypeStruct:
        Struct(π.π)
        return 
    
    case _GoTypeInterface:
        Interface(π.π)
        return 
    
    case _GoTypePointer:
        Pointer(π.π)
        return 
    
    case _GoTypeSlice:
        Slice(π.π)
        return 
    
    }
}

type _GoExprVariant interface {
    _GoExprVariant()
}

type GoExpr struct{
    variant _GoExprVariant
}

type _GoExprIdent struct{
    π GoIdent
}

func GoExprIdent(π GoIdent) GoExpr {
    return GoExpr {
        variant: _GoExprIdent {
            π: π,
        },
    }
}

func (self _GoExprIdent) _GoExprVariant() {}

type _GoExprIntLit struct{
    π GoIntLit
}

func GoExprIntLit(π GoIntLit) GoExpr {
    return GoExpr {
        variant: _GoExprIntLit {
            π: π,
        },
    }
}

func (self _GoExprIntLit) _GoExprVariant() {}

type _GoExprStructLit struct{
    π GoStructLit
}

func GoExprStructLit(π GoStructLit) GoExpr {
    return GoExpr {
        variant: _GoExprStructLit {
            π: π,
        },
    }
}

func (self _GoExprStructLit) _GoExprVariant() {}

type _GoExprCall struct{
    π GoCallExpr
}

func GoExprCall(π GoCallExpr) GoExpr {
    return GoExpr {
        variant: _GoExprCall {
            π: π,
        },
    }
}

func (self _GoExprCall) _GoExprVariant() {}

func (self GoExpr) Match(Ident func(π GoIdent), IntLit func(π GoIntLit), StructLit func(π GoStructLit), Call func(π GoCallExpr)) {
    switch π := self.variant.(type) {
    case _GoExprIdent:
        Ident(π.π)
        return 
    
    case _GoExprIntLit:
        IntLit(π.π)
        return 
    
    case _GoExprStructLit:
        StructLit(π.π)
        return 
    
    case _GoExprCall:
        Call(π.π)
        return 
    
    }
}

type _GoStmtVariant interface {
    _GoStmtVariant()
}

type GoStmt struct{
    variant _GoStmtVariant
}

type _GoStmtExpr struct{
    π GoExpr
}

func GoStmtExpr(π GoExpr) GoStmt {
    return GoStmt {
        variant: _GoStmtExpr {
            π: π,
        },
    }
}

func (self _GoStmtExpr) _GoStmtVariant() {}

type _GoStmtReturn struct{
    π GoReturnStmt
}

func GoStmtReturn(π GoReturnStmt) GoStmt {
    return GoStmt {
        variant: _GoStmtReturn {
            π: π,
        },
    }
}

func (self _GoStmtReturn) _GoStmtVariant() {}

type _GoStmtTypeSwitch struct{
    π GoTypeSwitch
}

func GoStmtTypeSwitch(π GoTypeSwitch) GoStmt {
    return GoStmt {
        variant: _GoStmtTypeSwitch {
            π: π,
        },
    }
}

func (self _GoStmtTypeSwitch) _GoStmtVariant() {}

func (self GoStmt) Match(Expr func(π GoExpr), Return func(π GoReturnStmt), TypeSwitch func(π GoTypeSwitch)) {
    switch π := self.variant.(type) {
    case _GoStmtExpr:
        Expr(π.π)
        return 
    
    case _GoStmtReturn:
        Return(π.π)
        return 
    
    case _GoStmtTypeSwitch:
        TypeSwitch(π.π)
        return 
    
    }
}

type _GoDeclVariant interface {
    _GoDeclVariant()
}

type GoDecl struct{
    variant _GoDeclVariant
}

type _GoDeclType struct{
    π GoTypeDecl
}

func GoDeclType(π GoTypeDecl) GoDecl {
    return GoDecl {
        variant: _GoDeclType {
            π: π,
        },
    }
}

func (self _GoDeclType) _GoDeclVariant() {}

type _GoDeclFunc struct{
    π GoFuncDecl
}

func GoDeclFunc(π GoFuncDecl) GoDecl {
    return GoDecl {
        variant: _GoDeclFunc {
            π: π,
        },
    }
}

func (self _GoDeclFunc) _GoDeclVariant() {}

type _GoDeclMethod struct{
    π GoMethodDecl
}

func GoDeclMethod(π GoMethodDecl) GoDecl {
    return GoDecl {
        variant: _GoDeclMethod {
            π: π,
        },
    }
}

func (self _GoDeclMethod) _GoDeclVariant() {}

func (self GoDecl) Match(Type func(π GoTypeDecl), Func func(π GoFuncDecl), Method func(π GoMethodDecl)) {
    switch π := self.variant.(type) {
    case _GoDeclType:
        Type(π.π)
        return 
    
    case _GoDeclFunc:
        Func(π.π)
        return 
    
    case _GoDeclMethod:
        Method(π.π)
        return 
    
    }
}

type GoFile struct{
    PackageName GoIdent
    Decls []GoDecl
}
