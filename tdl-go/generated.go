// This file is generated. Do not modify.

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

type TypeRef struct{
    Name Ident
    Params []Type
}

type _anonEnum0Variant interface {
    _anonEnum0Variant()
}

type anonEnum0 struct{
    variant _anonEnum0Variant
}

type _anonEnum0Ref_ struct{
    π TypeRef
}

func anonEnum0Ref_(π TypeRef) anonEnum0 {
    return anonEnum0 {
        variant: _anonEnum0Ref_ {
            π: π,
        },
    }
}

func (self _anonEnum0Ref_) _anonEnum0Variant() {}

type _anonEnum0Enum struct{
    π Enum
}

func anonEnum0Enum(π Enum) anonEnum0 {
    return anonEnum0 {
        variant: _anonEnum0Enum {
            π: π,
        },
    }
}

func (self _anonEnum0Enum) _anonEnum0Variant() {}

type _anonEnum0Struct struct{
    π Struct
}

func anonEnum0Struct(π Struct) anonEnum0 {
    return anonEnum0 {
        variant: _anonEnum0Struct {
            π: π,
        },
    }
}

func (self _anonEnum0Struct) _anonEnum0Variant() {}

type _anonEnum0Tuple struct{
    π Tuple
}

func anonEnum0Tuple(π Tuple) anonEnum0 {
    return anonEnum0 {
        variant: _anonEnum0Tuple {
            π: π,
        },
    }
}

func (self _anonEnum0Tuple) _anonEnum0Variant() {}

type _anonEnum0Pointer struct{
    π Pointer
}

func anonEnum0Pointer(π Pointer) anonEnum0 {
    return anonEnum0 {
        variant: _anonEnum0Pointer {
            π: π,
        },
    }
}

func (self _anonEnum0Pointer) _anonEnum0Variant() {}

type _anonEnum0Slice struct{
    π Slice
}

func anonEnum0Slice(π Slice) anonEnum0 {
    return anonEnum0 {
        variant: _anonEnum0Slice {
            π: π,
        },
    }
}

func (self _anonEnum0Slice) _anonEnum0Variant() {}

func (self anonEnum0) Match(Ref_ func(π TypeRef), Enum func(π Enum), Struct func(π Struct), Tuple func(π Tuple), Pointer func(π Pointer), Slice func(π Slice)) {
    switch π := self.variant.(type) {
    case _anonEnum0Ref_:
        Ref_(π.π)
        return 
    
    case _anonEnum0Enum:
        Enum(π.π)
        return 
    
    case _anonEnum0Struct:
        Struct(π.π)
        return 
    
    case _anonEnum0Tuple:
        Tuple(π.π)
        return 
    
    case _anonEnum0Pointer:
        Pointer(π.π)
        return 
    
    case _anonEnum0Slice:
        Slice(π.π)
        return 
    
    }
}

type _TypeVariant interface {
    _TypeVariant()
}

type Type struct{
    variant _TypeVariant
}

type _TypeRef_ struct{
    π TypeRef
}

func TypeRef_(π TypeRef) Type {
    return Type {
        variant: _TypeRef_ {
            π: π,
        },
    }
}

func (self _TypeRef_) _TypeVariant() {}

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

func (self Type) Match(Ref_ func(π TypeRef), Enum func(π Enum), Struct func(π Struct), Tuple func(π Tuple), Pointer func(π Pointer), Slice func(π Slice)) {
    switch π := self.variant.(type) {
    case _TypeRef_:
        Ref_(π.π)
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

type TypeCtor struct{
    Name Ident
    TypeVars []Ident
}

type TypeDecl struct{
    Ctor TypeCtor
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

type _anonEnum1Variant interface {
    _anonEnum1Variant()
}

type anonEnum1 struct{
    variant _anonEnum1Variant
}

type _anonEnum1Ident struct{
    π GoIdent
}

func anonEnum1Ident(π GoIdent) anonEnum1 {
    return anonEnum1 {
        variant: _anonEnum1Ident {
            π: π,
        },
    }
}

func (self _anonEnum1Ident) _anonEnum1Variant() {}

type _anonEnum1Func struct{
    π GoFuncType
}

func anonEnum1Func(π GoFuncType) anonEnum1 {
    return anonEnum1 {
        variant: _anonEnum1Func {
            π: π,
        },
    }
}

func (self _anonEnum1Func) _anonEnum1Variant() {}

type _anonEnum1Struct struct{
    π GoStruct
}

func anonEnum1Struct(π GoStruct) anonEnum1 {
    return anonEnum1 {
        variant: _anonEnum1Struct {
            π: π,
        },
    }
}

func (self _anonEnum1Struct) _anonEnum1Variant() {}

type _anonEnum1Interface struct{
    π GoInterface
}

func anonEnum1Interface(π GoInterface) anonEnum1 {
    return anonEnum1 {
        variant: _anonEnum1Interface {
            π: π,
        },
    }
}

func (self _anonEnum1Interface) _anonEnum1Variant() {}

type _anonEnum1Pointer struct{
    π GoPointer
}

func anonEnum1Pointer(π GoPointer) anonEnum1 {
    return anonEnum1 {
        variant: _anonEnum1Pointer {
            π: π,
        },
    }
}

func (self _anonEnum1Pointer) _anonEnum1Variant() {}

type _anonEnum1Slice struct{
    π GoSlice
}

func anonEnum1Slice(π GoSlice) anonEnum1 {
    return anonEnum1 {
        variant: _anonEnum1Slice {
            π: π,
        },
    }
}

func (self _anonEnum1Slice) _anonEnum1Variant() {}

func (self anonEnum1) Match(Ident func(π GoIdent), Func func(π GoFuncType), Struct func(π GoStruct), Interface func(π GoInterface), Pointer func(π GoPointer), Slice func(π GoSlice)) {
    switch π := self.variant.(type) {
    case _anonEnum1Ident:
        Ident(π.π)
        return 
    
    case _anonEnum1Func:
        Func(π.π)
        return 
    
    case _anonEnum1Struct:
        Struct(π.π)
        return 
    
    case _anonEnum1Interface:
        Interface(π.π)
        return 
    
    case _anonEnum1Pointer:
        Pointer(π.π)
        return 
    
    case _anonEnum1Slice:
        Slice(π.π)
        return 
    
    }
}

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

type _anonEnum2Variant interface {
    _anonEnum2Variant()
}

type anonEnum2 struct{
    variant _anonEnum2Variant
}

type _anonEnum2Ident struct{
    π GoIdent
}

func anonEnum2Ident(π GoIdent) anonEnum2 {
    return anonEnum2 {
        variant: _anonEnum2Ident {
            π: π,
        },
    }
}

func (self _anonEnum2Ident) _anonEnum2Variant() {}

type _anonEnum2IntLit struct{
    π GoIntLit
}

func anonEnum2IntLit(π GoIntLit) anonEnum2 {
    return anonEnum2 {
        variant: _anonEnum2IntLit {
            π: π,
        },
    }
}

func (self _anonEnum2IntLit) _anonEnum2Variant() {}

type _anonEnum2StructLit struct{
    π GoStructLit
}

func anonEnum2StructLit(π GoStructLit) anonEnum2 {
    return anonEnum2 {
        variant: _anonEnum2StructLit {
            π: π,
        },
    }
}

func (self _anonEnum2StructLit) _anonEnum2Variant() {}

type _anonEnum2Call struct{
    π GoCallExpr
}

func anonEnum2Call(π GoCallExpr) anonEnum2 {
    return anonEnum2 {
        variant: _anonEnum2Call {
            π: π,
        },
    }
}

func (self _anonEnum2Call) _anonEnum2Variant() {}

func (self anonEnum2) Match(Ident func(π GoIdent), IntLit func(π GoIntLit), StructLit func(π GoStructLit), Call func(π GoCallExpr)) {
    switch π := self.variant.(type) {
    case _anonEnum2Ident:
        Ident(π.π)
        return 
    
    case _anonEnum2IntLit:
        IntLit(π.π)
        return 
    
    case _anonEnum2StructLit:
        StructLit(π.π)
        return 
    
    case _anonEnum2Call:
        Call(π.π)
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

type _anonEnum3Variant interface {
    _anonEnum3Variant()
}

type anonEnum3 struct{
    variant _anonEnum3Variant
}

type _anonEnum3Expr struct{
    π GoExpr
}

func anonEnum3Expr(π GoExpr) anonEnum3 {
    return anonEnum3 {
        variant: _anonEnum3Expr {
            π: π,
        },
    }
}

func (self _anonEnum3Expr) _anonEnum3Variant() {}

type _anonEnum3Return struct{
    π GoReturnStmt
}

func anonEnum3Return(π GoReturnStmt) anonEnum3 {
    return anonEnum3 {
        variant: _anonEnum3Return {
            π: π,
        },
    }
}

func (self _anonEnum3Return) _anonEnum3Variant() {}

type _anonEnum3TypeSwitch struct{
    π GoTypeSwitch
}

func anonEnum3TypeSwitch(π GoTypeSwitch) anonEnum3 {
    return anonEnum3 {
        variant: _anonEnum3TypeSwitch {
            π: π,
        },
    }
}

func (self _anonEnum3TypeSwitch) _anonEnum3Variant() {}

func (self anonEnum3) Match(Expr func(π GoExpr), Return func(π GoReturnStmt), TypeSwitch func(π GoTypeSwitch)) {
    switch π := self.variant.(type) {
    case _anonEnum3Expr:
        Expr(π.π)
        return 
    
    case _anonEnum3Return:
        Return(π.π)
        return 
    
    case _anonEnum3TypeSwitch:
        TypeSwitch(π.π)
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

type _anonEnum4Variant interface {
    _anonEnum4Variant()
}

type anonEnum4 struct{
    variant _anonEnum4Variant
}

type _anonEnum4Type struct{
    π GoTypeDecl
}

func anonEnum4Type(π GoTypeDecl) anonEnum4 {
    return anonEnum4 {
        variant: _anonEnum4Type {
            π: π,
        },
    }
}

func (self _anonEnum4Type) _anonEnum4Variant() {}

type _anonEnum4Func struct{
    π GoFuncDecl
}

func anonEnum4Func(π GoFuncDecl) anonEnum4 {
    return anonEnum4 {
        variant: _anonEnum4Func {
            π: π,
        },
    }
}

func (self _anonEnum4Func) _anonEnum4Variant() {}

type _anonEnum4Method struct{
    π GoMethodDecl
}

func anonEnum4Method(π GoMethodDecl) anonEnum4 {
    return anonEnum4 {
        variant: _anonEnum4Method {
            π: π,
        },
    }
}

func (self _anonEnum4Method) _anonEnum4Variant() {}

func (self anonEnum4) Match(Type func(π GoTypeDecl), Func func(π GoFuncDecl), Method func(π GoMethodDecl)) {
    switch π := self.variant.(type) {
    case _anonEnum4Type:
        Type(π.π)
        return 
    
    case _anonEnum4Func:
        Func(π.π)
        return 
    
    case _anonEnum4Method:
        Method(π.π)
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

type GoComment string

type GoFile struct{
    FileComment GoComment
    PackageName GoIdent
    Decls []GoDecl
}

type Foo OptionalᐸListᐸintᐳᐳ

type Bar int

type _anonEnum5Variant interface {
    _anonEnum5Variant()
}

type anonEnum5 struct{
    variant _anonEnum5Variant
}

type _anonEnum5Cons struct{
    π struct{
        _0 int
        _1 Listᐸintᐳ
    }
}

func anonEnum5Cons(π struct{
    _0 int
    _1 Listᐸintᐳ
}) anonEnum5 {
    return anonEnum5 {
        variant: _anonEnum5Cons {
            π: π,
        },
    }
}

func (self _anonEnum5Cons) _anonEnum5Variant() {}

type _anonEnum5Nil struct{
    π struct{}
}

func anonEnum5Nil(π struct{}) anonEnum5 {
    return anonEnum5 {
        variant: _anonEnum5Nil {
            π: π,
        },
    }
}

func (self _anonEnum5Nil) _anonEnum5Variant() {}

func (self anonEnum5) Match(Cons func(π struct{
    _0 int
    _1 Listᐸintᐳ
}), Nil func(π struct{})) {
    switch π := self.variant.(type) {
    case _anonEnum5Cons:
        Cons(π.π)
        return 
    
    case _anonEnum5Nil:
        Nil(π.π)
        return 
    
    }
}

type _ListᐸintᐳVariant interface {
    _ListᐸintᐳVariant()
}

type Listᐸintᐳ struct{
    variant _ListᐸintᐳVariant
}

type _ListᐸintᐳCons struct{
    π struct{
        _0 int
        _1 Listᐸintᐳ
    }
}

func ListᐸintᐳCons(π struct{
    _0 int
    _1 Listᐸintᐳ
}) Listᐸintᐳ {
    return Listᐸintᐳ {
        variant: _ListᐸintᐳCons {
            π: π,
        },
    }
}

func (self _ListᐸintᐳCons) _ListᐸintᐳVariant() {}

type _ListᐸintᐳNil struct{
    π struct{}
}

func ListᐸintᐳNil(π struct{}) Listᐸintᐳ {
    return Listᐸintᐳ {
        variant: _ListᐸintᐳNil {
            π: π,
        },
    }
}

func (self _ListᐸintᐳNil) _ListᐸintᐳVariant() {}

func (self Listᐸintᐳ) Match(Cons func(π struct{
    _0 int
    _1 Listᐸintᐳ
}), Nil func(π struct{})) {
    switch π := self.variant.(type) {
    case _ListᐸintᐳCons:
        Cons(π.π)
        return 
    
    case _ListᐸintᐳNil:
        Nil(π.π)
        return 
    
    }
}

type _anonEnum6Variant interface {
    _anonEnum6Variant()
}

type anonEnum6 struct{
    variant _anonEnum6Variant
}

type _anonEnum6Some struct{
    π Listᐸintᐳ
}

func anonEnum6Some(π Listᐸintᐳ) anonEnum6 {
    return anonEnum6 {
        variant: _anonEnum6Some {
            π: π,
        },
    }
}

func (self _anonEnum6Some) _anonEnum6Variant() {}

type _anonEnum6None struct{
    π struct{}
}

func anonEnum6None(π struct{}) anonEnum6 {
    return anonEnum6 {
        variant: _anonEnum6None {
            π: π,
        },
    }
}

func (self _anonEnum6None) _anonEnum6Variant() {}

func (self anonEnum6) Match(Some func(π Listᐸintᐳ), None func(π struct{})) {
    switch π := self.variant.(type) {
    case _anonEnum6Some:
        Some(π.π)
        return 
    
    case _anonEnum6None:
        None(π.π)
        return 
    
    }
}

type _OptionalᐸListᐸintᐳᐳVariant interface {
    _OptionalᐸListᐸintᐳᐳVariant()
}

type OptionalᐸListᐸintᐳᐳ struct{
    variant _OptionalᐸListᐸintᐳᐳVariant
}

type _OptionalᐸListᐸintᐳᐳSome struct{
    π Listᐸintᐳ
}

func OptionalᐸListᐸintᐳᐳSome(π Listᐸintᐳ) OptionalᐸListᐸintᐳᐳ {
    return OptionalᐸListᐸintᐳᐳ {
        variant: _OptionalᐸListᐸintᐳᐳSome {
            π: π,
        },
    }
}

func (self _OptionalᐸListᐸintᐳᐳSome) _OptionalᐸListᐸintᐳᐳVariant() {}

type _OptionalᐸListᐸintᐳᐳNone struct{
    π struct{}
}

func OptionalᐸListᐸintᐳᐳNone(π struct{}) OptionalᐸListᐸintᐳᐳ {
    return OptionalᐸListᐸintᐳᐳ {
        variant: _OptionalᐸListᐸintᐳᐳNone {
            π: π,
        },
    }
}

func (self _OptionalᐸListᐸintᐳᐳNone) _OptionalᐸListᐸintᐳᐳVariant() {}

func (self OptionalᐸListᐸintᐳᐳ) Match(Some func(π Listᐸintᐳ), None func(π struct{})) {
    switch π := self.variant.(type) {
    case _OptionalᐸListᐸintᐳᐳSome:
        Some(π.π)
        return 
    
    case _OptionalᐸListᐸintᐳᐳNone:
        None(π.π)
        return 
    
    }
}
