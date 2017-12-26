package parser

import (
	"go-interpreter-lexer/lexer"
	"go-interpreter-lexer/token"
	"go-interpreter-lexer/ast"
	"fmt"
	"strconv"
)


/**
	The iota starts of with the value of the constant as 0 for _ and 1 to 7 for subsequent values.
	we will use this for the order or precedence of the operators as well. therefore the PREFIX will holder
	higher precedence than say an equal etc.
*/
const(
	_ int = iota
	LOWEST
	EQUALS // ==
	LESSGREATER // > or <
	SUM // +
	PRODUCT // *
	PREFIX  // -X or !X
	CALL    // myFunc(X)
)

type Parser struct{
	l *lexer.Lexer

	curToken token.Token
	peekToken token.Token
	errors []string
	// adding a series of infix and prefix func holder.
	prefixParseFuncs map[token.TokenType]prefixParseFn
	infixParseFuncs map[token.TokenType]infixParseFn
}

/*
	Each node in the ast can have two parsing functions associated with it
	1. prefix parsing function
	2. infix parsing function.
	When a token is found in a prefix prosition ++a  or i-- then we call the
	prefix function in the case of a infix position for a token we call the
	infix parse function
*/
type (
	prefixParseFn func() ast.Expression
	/*
		the infix parse function will take an input to the function which
		acts as the left side of the expression that are trying to eval.
	*/
	infixParseFn func(ast.Expression) ast.Expression
)

/*
	functions to register the infix and prefix functions.
*/
func (p *Parser) registerPrefixFn(tokenType token.TokenType, fn prefixParseFn){
	p.prefixParseFuncs[tokenType] = fn
}

func (p *Parser) registerInfixFn(tokenType token.TokenType, fn infixParseFn){
	p.infixParseFuncs[tokenType] = fn
}

func New(l *lexer.Lexer) *Parser{
	p := &Parser{
		l:l,
		errors: []string{},
	}

	p.prefixParseFuncs = make(map[token.TokenType]prefixParseFn)
	p.registerPrefixFn(token.IDENT, p.parseIdentifier)
	p.registerPrefixFn(token.INT, p.parseInteger)
	p.registerPrefixFn(token.BANG, p.parsePrefixExpression)
	p.registerPrefixFn(token.MINUS, p.parsePrefixExpression)


	// calling next token twice to get current and peek token populated.
    p.nextToken()
    p.nextToken()
	return p
}

// function to handle the prefix func for identifiers.
func (p *Parser) parseIdentifier() ast.Expression{
	return &ast.Identifier{ Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseInteger() ast.Expression{
	il := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0,64)
	if err != nil{
		msg := fmt.Sprintf("could not parse integer literal %q as integer", p.curToken.Literal)
		p.errors = append(p.errors,msg)
		return nil
	}
	il.Value = value
	return il
}

func (p *Parser) parsePrefixExpression() ast.Expression{
	pe := &ast.PrefixExpression{
		Token : p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	pe.Right = p.parseExpression(PREFIX)
	return pe
}

func (p *Parser) Errors() []string{
	return p.errors
}

func (p *Parser) peekError(t token.TokenType){
	msg := fmt.Sprintf("expected next token to be %s got: %s", t, p.peekToken.Type)
	p.errors = append(p.errors,msg)
}

func (p *Parser) nextToken(){
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program{
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF{
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement{
	//log.Printf("Current Token Type is : %s \n", p.curToken.Type)
	switch p.curToken.Type{
		case token.LET:
			return p.parseLetStatement()
		case token.RETURN:
			return p.parseReturnStatement()
		default:
			// handle expression statements
			return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement{
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT){
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(token.ASSIGN){
		return nil
	}
	//TODO: skip unto we get semicolon add more functionality later.
	for !p.curTokenIs(token.SEMICOLON){
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement{
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()
	// TODO: skipping the expression until we get a semicolon.
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement{
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON){
		p.nextToken()
	}

	return stmt
}

func (p *Parser) curTokenIs( t token.TokenType) bool{
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool{
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool{
	if p.peekTokenIs(t){
		p.nextToken()
		return true
	}else{
		p.peekError(t)
		return false
	}
}

func (p *Parser) parseExpression(precedence int) ast.Expression{
	prefix := p.prefixParseFuncs[p.curToken.Type]
	if prefix == nil{
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	return leftExp
}

func (p *Parser) noPrefixParseFnError(t token.TokenType){
	msg := fmt.Sprintf("no prefix parse function for %s found\n", t)
	p.errors = append(p.errors,msg)
}

