package parser

import (
	"testing"
	"go-interpreter-lexer/lexer"
	"go-interpreter-lexer/ast"
	"fmt"
	"strconv"
)

type identifiers struct{
	identifier string
}

var testParserCases = []struct{
	input string
	statementCount int
	expectedIdentifiers []identifiers
	errorsCount int
}{
	 {
	 	`let x = 5;
let y = 10;
let foobar = 83883;
`,
	 	3,
	 	[]identifiers{
	 	   {"x"},
	 	   {"y"},
	 	   {"foobar"},
		},
		0,
	 },
	/*{
		`
let x = 5;
let  = 10;
let 832323;
`,
		0,
		[]identifiers{
		},
		2,
	},*/
}

func TestLetStatements(t *testing.T){
	for i, statementCase := range testParserCases {
		t.Logf("Test Case index: %d\n", i)
		l := lexer.New(statementCase.input)
		p := New(l)

		program := p.ParseProgram()
		errCount := ParserErrorsCount(t, p)

		if errCount != 0 && errCount == statementCase.errorsCount{
			t.Logf("Expected errors match the errorCount")
			return
		}

		if program == nil {
			t.Fatalf("ParseProgram returned nil therefore exiting")
		}

		if len(program.Statements) != statementCase.statementCount {
			t.Fatalf("Program statements does not container correct number of statements expected : %d but got: %d", statementCase.statementCount, len(program.Statements))
		}else{
			t.Logf("Statement count is : %d\n",len(program.Statements))
		}

		for i, tt := range statementCase.expectedIdentifiers {
			stmt := program.Statements[i]
			if !testLetStatement(t, stmt, tt.identifier){
				return
			}
		}

	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool{
	if s.TokenLiteral() != "let"{
		t.Errorf("Statement token literal should be let but got %q \n", s.TokenLiteral())
		return false
	}
	letStmt, ok := s.(*ast.LetStatement)
	if !ok{
		t.Errorf("Statement is not an ast.Statement Type got %T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf ( "letStmt.Name.Value not '%s' got: %s ", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("let statement name is not %s got %s", name, letStmt.Name)
		return false
	}

	return true
}

func ParserErrorsCount(t *testing.T, p *Parser) int{
	errors := p.Errors()
	if len(errors) == 0 {
		return 0
	}
	t.Logf("parser has %d errors\n", len(p.Errors()))
	for _, msg := range errors{
		t.Logf("parser error: %q \n" ,msg)
	}
	return len(p.errors)
}

var testReturnStatements = []struct{
	input string
	errorCount int
	statementCount int
}{
	{
		`
return 5;
return 10;
return 832323;
`,
0,
3,
	},
}

func TestReturnStatements(t *testing.T){
	for _, returnCase :=range testReturnStatements{
		l := lexer.New(returnCase.input)
		p := New(l)

		program := p.ParseProgram()
		errCount := ParserErrorsCount(t,p)

		if errCount != 0 && errCount == returnCase.errorCount{
			t.Logf("Expected errors match the errorCount")
			return
		}

		if len(program.Statements) != 3{
			t.Fatalf("program.Statements does not contain %d statements got: %d", returnCase.statementCount, len(program.Statements))
		}

		for _, stmt := range program.Statements {
			returnStmt, ok := stmt.(*ast.ReturnStatement)
			if !ok{
				t.Errorf("stmt not a *ast.ReturnStatement. got= %T\n", stmt)
				continue
			}
			if returnStmt.TokenLiteral() != "return"{
				t.Errorf("Statement token literal is not 'return' got: %s \n", returnStmt.TokenLiteral())
			}
		}
	}
}


func TestIdentifierExpression(t *testing.T){
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	count := ParserErrorsCount(t,p)

	if count != 0 {
		t.Errorf("Expected 0 errors but found %d\n", count)
	}
	if len(program.Statements) != 1{
		t.Fatalf("program does not have enough statements got: %d",len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok{
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatements. got: %T\n",program.Statements[0])
	}
	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expression is not ast.Identifier got: %T\n",stmt.Expression)
	}

	if ident.Value != "foobar"{
		t.Errorf("identifier value not %s; got: %s", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar"{
		t.Errorf("identifier Literal not %s; got: %s", "foobar", ident.TokenLiteral())
	}
}

func TestIntegerLiterals(t *testing.T){
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	count := ParserErrorsCount(t,p)

	if count != 0 {
		t.Errorf("Expected 0 errors but found %d\n", count)
	}

	if len(program.Statements) != 1 {
		t.Errorf("Parser must have 1 statement but got: %d\n", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok{
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatements. got: %T\n",program.Statements[0])
	}

	in, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Expression is not ast.Identifier got: %T\n",stmt.Expression)
	}

	if in.Value != 5 {
		t.Errorf("identifier value not %s; got: %d", "5", in.Value)
	}

	if in.TokenLiteral() != "5"{
		t.Errorf("identifier Literal not %s; got: %s", "5", in.TokenLiteral())
	}

}

// testing the prefix operator support
func TestParsingPrefixExpression(t *testing.T){
	prefixTest := []struct{
		input string
		operator string
		value interface{}
	}{
		{"-15;", "-", 15},
		{"!5", "!", 5},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, tt := range prefixTest{
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		count := ParserErrorsCount(t,p)
		if count != 0{
			t.Errorf("There are parsing errors in the statements")
		}
		if len(program.Statements) != 1{
			t.Fatalf("Program statements does not contain %d statements got: %d\n", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got: %T\n", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)

		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got: %T\n", stmt.Expression)
		}
		if exp.Operator != tt.operator{
			t.Fatalf("exp.Operator is not %s. got: %s\n", tt.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.value){
			return
		}
	}
}


func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool{
	intLit , ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il is not a integer literal got: %T\n", il)
		return false
	}
	if intLit.Value != value {
		t.Errorf("Integer Literal value does not match expected %d but got %d\n", value, intLit.Value)
		return false
	}
	if intLit.TokenLiteral() != fmt.Sprintf("%d", value){
		t.Errorf("intLit.TokenLiteral is not %d got: %s",value, intLit.TokenLiteral() )
		return false
	}
	return true
}

func TestParsingInfixExpression(t *testing.T){
	infixTests := []struct{
		input string
		leftValue interface{}
		operator string
		rightValue interface{}
	}{
		{
			"5 + 5;",
			 5,
			 "+",
			 5,
		},
		{
			"5 == 5;",
			5,
			"==",
			5,
		},
		{
			"5 != 5;",
			5,
			"!=",
			5,
		},
		{
			"5 - 5;",
			5,
			"-",
			5,
		},
		{
			"5 * 5;",
			5,
			"*",
			5,
		},
		{
			"5 > 5;",
			5,
			">",
			5,
		},
		{
			"true == false", true, "==", false,
		},
	}


	for _, tc := range infixTests {
		l := lexer.New(tc.input)
		p := New(l)
		program := p.ParseProgram()

		count := ParserErrorsCount(t,p)
		if count != 0 {
			t.Fatalf("Statement should not a error as all are valid statements")
		}
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements got : %d\n", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not a ast.ExpressionStatement. got: %T\n", program.Statements[0])
		}

		testInfixExpression(t, stmt.Expression, tc.leftValue, tc.operator, tc.rightValue)
	}

}

func TestOperatorPrecendenceParsing(t *testing.T){
	tests := []struct{
		input string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b * c",
			"(a + (b * c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 3",
			"((5 + 5) * 3)",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
	}

	for _, tt := range tests{
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		count := ParserErrorsCount(t,p)
		if count != 0{
			t.Fatalf("Expression must all be valid")
		}
		actual := program.String()
		if actual != tt.expected{
			t.Errorf("expected: %q, got: %q", tt.expected, actual)
		}
	}
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool{
	ident, ok := exp.(*ast.Identifier)
	if !ok{
		t.Errorf("exp not of Identifier type got: %T", exp)
		return false
	}
	if ident.Value != value{
		t.Errorf("ident.Value not %s. got: %T", value, ident.Value)
		return false
	}
	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got: %T", value, ident.TokenLiteral())
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {

	switch v:= expected.(type){
		case int:
				return testIntegerLiteral(t,exp, int64(v))
		case int64:
				return testIntegerLiteral(t,exp, v)
		case string:
				return testIdentifier(t,exp,v)
		case bool:
				return testBooleanLiterals(t,exp,v)
	}

	t.Errorf("type of expression not handled. got: %T",exp)
	return false;
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("Expression is not an infix expression. got= %T",opExp)
		return false
	}
	if !testLiteralExpression(t, opExp.Left, left){
		return false
	}
	if opExp.Operator != operator{
		t.Errorf("Operator for the infix expression expected: %s got: %s ",operator, opExp.Operator)
		return false
	}
	if ! testLiteralExpression(t, opExp.Right, right){
		return false
	}
	return true
}

func TestBooleanExpression(t *testing.T){
	boolTests := []struct{
		input string
		expected string
	}{
		{"true", "true"},
		{"false", "false"},
	}

	for _, tc := range boolTests {
		l := lexer.New(tc.input)
		p := New(l)

		program := p.ParseProgram()
		count := ParserErrorsCount(t,p)

		if count != 0 {
			t.Fatalf("Expected 0 errors but got %d ", count)
		}

		if len(program.Statements) != 1 {
			t.Errorf("program.Statements is on length %d expected 1")
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok{
			t.Fatalf("expected the staement to be ExpressionStatement but got %T",stmt)
		}
		exp, ok := stmt.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf("Expression was not ast.Boolean type got: %T", exp)
		}
		gotValue, _ := strconv.ParseBool(tc.expected)
		if  gotValue != exp.Value{
			t.Errorf("Expected value %s but got= %s", gotValue,exp.Value)
		}
	}
}

func testBooleanLiterals(t *testing.T, exp ast.Expression, value bool) bool{
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean type got: %T", exp)
		return false
	}
	if bo.Value != value{
		t.Errorf("Boolean expression value is %t got: %t", value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value){
		t.Errorf("bo.TokenLiteral not %t got: %t",value, bo.TokenLiteral())
		return false
	}
	return true
}

func TestIfExpression(t *testing.T){
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	count := ParserErrorsCount(t, p)

	if count != 0 {
		t.Errorf("Expression should not have any error statements but got %d", count)
	}

	if len(program.Statements) != 1{
		t.Fatalf("program Body does not contain %d statements got %d instead", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got %T instead ", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Errorf("stmt.Expression is not an if expression got %T ", stmt)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y"){
		return
	}
	if len(exp.Consequence.Statements) != 1 {
		t.Errorf ( "consequence is not 1 statement got %d ", len(exp.Consequence.Statements))
	}
	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements in the consequence is not ExpressionStatement got %T", exp.Consequence.Statements[0])
	}
	if !testIdentifier(t, consequence.Expression,"x") { return }

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative was nil but got statement for alternative got: %+v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T){
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	count := ParserErrorsCount(t, p)

	if count != 0 {
		t.Errorf("Expression should not have any error statements but got %d", count)
	}

	if len(program.Statements) != 1{
		t.Fatalf("program Body does not contain %d statements got %d instead", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got %T instead ", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Errorf("stmt.Expression is not an if expression got %T ", stmt)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y"){
		return
	}
	if len(exp.Consequence.Statements) != 1 {
		t.Errorf ( "consequence is not 1 statement got %d ", len(exp.Consequence.Statements))
	}
	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements in the consequence is not ExpressionStatement got %T", exp.Consequence.Statements[0])
	}
	if !testIdentifier(t, consequence.Expression,"x") { return }

	if exp.Alternative == nil {
		t.Errorf("exp.Alternative was not nil but got nil")
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf ( "Alternative is not 1 statement got %d ", len(exp.Alternative.Statements))
	}

	alt, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements in the alternative is not ExpressionStatement got %T", exp.Alternative.Statements[0])
	}
	if !testIdentifier(t, alt.Expression,"y") { return }

}