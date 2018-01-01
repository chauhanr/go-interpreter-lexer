package evaluator

import (
	"testing"
	"go-interpreter-lexer/object"
	"go-interpreter-lexer/lexer"
	"go-interpreter-lexer/parser"
)

func TestEvaluateIntegerExpression(t *testing.T){
	tests := []struct{
		input string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-10", -10},
		{"-39", -39},
		{"5 + 5 + 5", 15},
		{"3 * 5 + 9", 24},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-5 + 10 + -5", 0},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10 ", 60},
		{"2 * (5 + 10)", 30},
	}

	for _, tt := range tests{
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvaluateBoolExpressions(t *testing.T){
	tests := []struct{
		input string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 1", false},
		{"1 > 1 ", false},
		{"1 > 2", false},
		{" 1 == 1", true},
		{"1 != 2", true},
		{"1 != 1", false},
		{"false == false", true},
		{"true == true", true},
		{"false != true", true},
		{"false != false", false},
		{" (1 < 2) == true", true},
		{" (1 < 2) == false", false},
		{" (1 > 2) == false", true},
		{" (1 > 2) == true", false},
		{"true == true", true},
		{"false != true", true},
		{"false != false", false},
	}

	for _, bb := range tests{
		evaluated := testEval(bb.input)
		testBooleanObject(t, evaluated, bb.expected)
	}
}

func testEval (input string) object.Object{
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	return Eval(program)
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool{
	result, ok := obj.(*object.Boolean)

	if !ok{
		t.Errorf("expected value was object.Boolean but got %T",obj)
	}
	if result.Bool != expected{
		t.Errorf("Expected %t but got %t", expected, result.Bool)
		return false
	}
	return true
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool{
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("Expected value is suppose to be of object.Integer Type by found %T",obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("Expected Value %d was not equal to the evaluated value %d", expected, result.Value)
		return false
	}
	return true
}

func TestBangOperator(t *testing.T){
	tests := []struct{
		input string
		expected bool
	}{
		{"!true", false},
		{ "!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!5", true},
	}

	for _, bang := range tests{
		evaluated := testEval(bang.input)
		testBooleanObject(t, evaluated,bang.expected)
	}
}

func TestIfElseExpression(t *testing.T){
	tests := []struct{
		input string
		expected interface{}
	}{
		{"if (true) { 10 }", 10 },
		{"if (false) { 10 }", nil},
		{"if ( 1 ) { 10 }", 10 },
		{"if ( 1 < 2) { 10 }", 10 },
		{"if ( 1 > 2) { 10 }", nil },
		{"if ( 1 < 2) { 10 } else { 20 } ", 10 },
		{"if ( 1 > 2) { 10 } else { 20 } ", 20 },
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		}else {
			testNullObject(t, evaluated)
		}
	}
}

func testNullObject(t *testing.T, evaluated object.Object) bool {
	if evaluated != NULL {
		t.Errorf("Object expected to be null but got %T", evaluated)
		return false
	}
	return true
}

func TestReturnStatements(t *testing.T){
	tests := []struct{
		input string
		expected int64
	}{
		{ "return 10;", 10},
		{ "return 10; 9; ", 10},
		{ "return 2 * 5; 9;", 10},
		{"9; return 2*5; 9;", 10},
		{`if (10>1) {
		if ( 10>1) {
			return 10;
		}
			return 1;
		`, 10},
	}

	for _,tt := range tests{
		evaluated := testEval(tt.input)
		testIntegerObject(t,evaluated,tt.expected)
	}
}