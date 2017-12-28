package ast

import (
	"go-interpreter-lexer/token"
	"bytes"
)

/*
	Every node in the ast will implement the node interface and will return the literal value associated wit the node
    on calling the tokenLiteral method. Some of the nodes will implement the Statement interface and Expression interface.
*/

type Node interface{
	TokenLiteral() string
	String() string
}

type Statement interface{
	Node
	statementNode()
}

type Expression interface{
	Node
	expressionNode()
}


/**
	The program structure will hold on the slices of the statement and expression.
*/
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string{
	if len(p.Statements) >0 {
		return p.Statements[0].TokenLiteral()
	}else{
		return ""
	}
}

func (p *Program) String() string{
	var out bytes.Buffer

	for _,s := range p.Statements{
		out.WriteString(s.String())
	}

	return out.String()
}

type LetStatement struct{
	Token token.Token
	Name *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode(){}
func (ls *LetStatement) TokenLiteral() string{
	return ls.Token.Literal
}

func (ls LetStatement) String() string{
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() +" ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	out.WriteString(ls.Value.String())
    out.WriteString(";")
	return out.String()
}


/**
	the return statement is one where we have the following syntax each time
 	return <expression>
	e.g. return 39;
    	 return add(5,10);
*/
type ReturnStatement struct{
	Token token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode(){}
func (rs *ReturnStatement) TokenLiteral() string{
	return rs.Token.Literal
}
func (rs *ReturnStatement) String() string{
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

/*
	Structure to hold on to the name of the identifier that identifies a statement or expression.
*/
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode(){}
func (i *Identifier) TokenLiteral() string{
	return i.Token.Literal
}
func (i *Identifier) String() string{
	return i.Value
}

type IntegerLiteral struct{
	Token token.Token
	Value int64
}

func(i *IntegerLiteral) expressionNode(){}
func(i *IntegerLiteral) TokenLiteral() string{
	return i.Token.Literal
}
func (i *IntegerLiteral) String() string{
	return i.Token.Literal
}

/**
	Expression statements are statements that represent a single or combination of expressions.
*/
type ExpressionStatement struct{
	Token token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode(){}
func (es *ExpressionStatement) TokenLiteral() string{
	return es.Token.Literal
}
func (es *ExpressionStatement) String() string{
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type InfixExpression struct{
	Token token.Token
	Left Expression
	Operator string
	Right Expression
}

func (ie *InfixExpression) expressionNode(){}
func (ie *InfixExpression) TokenLiteral() string{
	return ie.Token.Literal
}

func (ie *InfixExpression) String() string{
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" "+ie.Operator+" ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

type PrefixExpression struct{
	Token token.Token
	Operator string
	Right Expression
}

func (pe *PrefixExpression) expressionNode(){}
func (pe *PrefixExpression) TokenLiteral() string{
	return pe.Token.Literal
}

func (pe *PrefixExpression) String() string{
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type Boolean struct{
	Token token.Token
	Value bool
}

func( b *Boolean) expressionNode(){}
func( b *Boolean) TokenLiteral() string{
	return b.Token.Literal
}
func (b *Boolean) String() string{
	return b.Token.Literal
}


type IfExpression struct{
	Token token.Token
	Condition Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ife *IfExpression) expressionNode(){}
func (ife *IfExpression) TokenLiteral() string{
	return ife.Token.Literal
}
func (ife *IfExpression) String() string{
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ife.Condition.String())
	out.WriteString(ife.Consequence.String())
	out.WriteString(" ")
	if ife.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ife.Alternative.String())
	}
	return out.String()
}

type BlockStatement struct{
	Token token.Token
	Statements []Statement
}

func(bs *BlockStatement) statementNode(){}
func(bs *BlockStatement) TokenLiteral() string{
	return bs.Token.Literal
}
func(bs *BlockStatement) String() string{
	var out bytes.Buffer
	for _, s := range bs.Statements{
		out.WriteString(s.String())
	}
	return out.String()
}