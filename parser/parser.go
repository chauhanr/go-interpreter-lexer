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

var precedences = map[token.TokenType] int {
	token.EQ: EQUALS,
	token.NOT_EQ: EQUALS,
	token.LT: LESSGREATER,
	token.GT: LESSGREATER,
	token.PLUS: SUM,
	token.MINUS: SUM,
	token.SLASH: PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN: CALL,
}

// precedence function to determine the presedecne of the current operator with other operators.
func (p *Parser) peekPrecedence() int{
	if p, ok:= precedences[p.peekToken.Type]; ok{
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int{
  	if p, ok := precedences[p.curToken.Type]; ok{
  		return p
	}
	return LOWEST
}


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
	p.registerPrefixFn(token.TRUE, p.parseBoolean)
	p.registerPrefixFn(token.FALSE, p.parseBoolean)
	p.registerPrefixFn(token.LPAREN, p.parseGroupExpression)
	p.registerPrefixFn(token.IF, p.parseIfExpression)
	p.registerPrefixFn(token.FUNCTION, p.parseFunctionLiteral)

	// adding support for the INFIX parser operators.
    p.infixParseFuncs = make(map[token.TokenType]infixParseFn)
    p.registerInfixFn(token.PLUS, p.parseInfixExpression)
	p.registerInfixFn(token.MINUS, p.parseInfixExpression)
	p.registerInfixFn(token.GT, p.parseInfixExpression)
	p.registerInfixFn(token.LT, p.parseInfixExpression)
	p.registerInfixFn(token.SLASH, p.parseInfixExpression)
	p.registerInfixFn(token.ASTERISK, p.parseInfixExpression)
	p.registerInfixFn(token.EQ, p.parseInfixExpression)
	p.registerInfixFn(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfixFn(token.LPAREN, p.parseCallExpression)

	// calling next token twice to get current and peek token populated.
    p.nextToken()
    p.nextToken()
	return p
}

func (p *Parser) parseCallExpression(fn ast.Expression) ast.Expression {
	ce := &ast.CallExpression{Token: p.curToken, Function: fn}
	ce.Arguments = p.parseCallArguments()
	return ce
}

func (p *Parser) parseCallArguments() []ast.Expression{
	args := []ast.Expression{}

	// handle the situation where the call does not have any arguments.
	if p.peekTokenIs(token.LPAREN){
		p.nextToken()
		return args
	}
    // parse the first argument.
	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))
    // parse all subsequent arguments
	for p.peekTokenIs(token.COMMA){
		p.nextToken()
		p.nextToken()
		args = append(args,p.parseExpression(LOWEST))
	}
	if !p.expectPeek(token.RPAREN){
		return nil
	}

	return args
}

func (p *Parser) parseFunctionLiteral() ast.Expression{
	fLit := &ast.FunctionLiteral{ Token: p.curToken}

	if !p.expectPeek(token.LPAREN){
		return nil
	}
	fLit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE){
		return nil
	}

	fLit.Body = p.parseBlockStatement()

	return fLit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier{
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN){
		p.nextToken()
		return identifiers
	}
	p.nextToken()
	ident := &ast.Identifier{ Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers,ident)

	for p.peekTokenIs(token.COMMA){
		p.nextToken()
		p.nextToken()
		id := &ast.Identifier{ Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers,id)
	}
	if !p.expectPeek(token.RPAREN){
		return nil
	}
	return identifiers
}

func (p *Parser) parseIfExpression() ast.Expression{
	//p.nextToken()
    ifExp := &ast.IfExpression{Token: p.curToken}
	if !p.expectPeek(token.LPAREN){
		return nil
	}
	p.nextToken()
	ifExp.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN){
		return nil
	}
	if !p.expectPeek(token.LBRACE){
		return nil
	}

	ifExp.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE){
		p.nextToken()
		if !p.expectPeek(token.LBRACE){
			return nil
		}
		ifExp.Alternative = p.parseBlockStatement()
	}
	return ifExp
}


func (p *Parser) parseBlockStatement() *ast.BlockStatement{
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF){
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}



func (p *Parser) parseGroupExpression() ast.Expression{
	p.nextToken()
	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN){
		return nil
	}
	return exp
}

func (p *Parser) parseBoolean() ast.Expression{
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression{
	exp := &ast.InfixExpression{
		Token: p.curToken,
		Operator: p.curToken.Literal,
		Left: left,
	}
     precedence := p.curPrecedence()
     p.nextToken()
     exp.Right = p.parseExpression(precedence)
	 return exp
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
	switch p.curToken.Type{
		case token.LET:
			return p.parseLetStatement()
		case token.RETURN:
			return p.parseReturnStatement()
		default:
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

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	for p.peekTokenIs(token.SEMICOLON){
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement{
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

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

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence(){
		infix := p.infixParseFuncs[p.peekToken.Type]
		if infix == nil{
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) noPrefixParseFnError(t token.TokenType){
	msg := fmt.Sprintf("no prefix parse function for %s found\n", t)
	p.errors = append(p.errors,msg)
}

