package object

import "fmt"

type ObjectType string

const(
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	NULL_OBJ = "NULL"
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