package object

import "fmt"

type ObjectType string

const(
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	NULL_OBJ = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
)

type Object interface{
	Type() ObjectType
	Inspect() string
}

// The monkey language has 3 types
type Integer struct{
	Value int64
}

func (i *Integer) Inspect() string{
	return fmt.Sprintf("%d",i.Value)
}
func (i *Integer) Type() ObjectType{
	return INTEGER_OBJ
}

type Boolean struct{
	Bool bool
}

func (b *Boolean) Inspect() string{
	return fmt.Sprintf("%t",b.Bool)
}
func (b *Boolean) Type() ObjectType{
	return BOOLEAN_OBJ
}

type Null struct{}

func (n *Null) Type() ObjectType{
    return NULL_OBJ
}
func (n *Null) Inspect() string{
	return "null"
}

type ReturnValue struct{
 	Value Object
}

func (rv *ReturnValue) Type() ObjectType{
	return RETURN_VALUE_OBJ
}
func (rv *ReturnValue) Inspect() string{
	return rv.Value.Inspect()
}
