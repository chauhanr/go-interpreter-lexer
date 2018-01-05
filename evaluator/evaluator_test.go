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
	env := object.NewEnvironment()

	return Eval(program, env)
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
		t.Errorf("Evaluated value is suppose to be of object.Integer Type by found %T",obj)
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

func TestErrorHandling (t *testing.T){
	tests := []struct{
		input string
		expectedMessage string
	}{
		{"5 + true;", "type mismatch: INTEGER + BOOLEAN"},
		{"5 + true; 5;", "type mismatch: INTEGER + BOOLEAN"},
		{"-true", "unknown operator: -BOOLEAN"},
		{"true + false;", "unknown operator: BOOLEAN + BOOLEAN"},
		{`if (10 > 1) {
		   if ( 10 > 1) {
				return true + true;
			}
		`, "unknown operator: BOOLEAN + BOOLEAN"},
		{"foobar", "identifier not found: foobar"},
		{`"Hello"-"World"`, "unknown operator: STRING - STRING"},
	}

	for _, tt :=range tests{
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if ! ok{
			t.Errorf("no error object returned. got= %T(%+v)", evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expectedMessage{
			t.Errorf("wrong error message.expected=%q, got=%q", tt.expectedMessage, errObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T){
	tests := []struct{
		input string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a; ", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests{
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}


func TestFunctionObject(t *testing.T){
	input := "fn(x) {x+2;};"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok{
		t.Fatalf("object is not a function. got: %T(%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong number of parameters %+v", fn.Parameters)
	}
	if fn.Parameters[0].String() != "x"{
		t.Fatalf("parameter is not 'x' got %q", fn.Parameters[0])
	}

	expectBody := "(x + 2)"
	if fn.Body.String() != expectBody{
		t.Fatalf("body of the function is not %q, got %q", expectBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T){
	tests := []struct{
		input string
		expected int64
	}{
		{ "let identity = fn(x) {x;}  identity(5);", 5},
		{"let identity = fn(x) { return x;}; identity(5)", 5},
		{"let double = fn(x) { x*2;}; double(5); ",10},
		{"let add = fn(x, y) { x + y; }; add(4, 6);", 10},
		{"let add = fn(x, y) { x + y; }; add(4 + 6, add(5,5));", 20},
		{ "fn(x) {x; }(5)", 5},
		{`fn( ) { 5;}()`, 5},
	}

	for _, tt := range tests{
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosures(t *testing.T){
	input := `
	let newAdder = fn(x) {
		fn(y) { x+y };
	};
	let addTwo = newAdder(2);
	addTwo(2);`

	testIntegerObject(t, testEval(input), 4)

}

func TestStringLiteralExpression(t *testing.T){
	input := `fn() {"hello world!"}();`
	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok{
		t.Fatalf("expected object.String got :%T", evaluated)
	}
	if str.Value != "hello world!" {
		t.Fatalf("expected value %s but got %s","hello world!" ,str.Value)
	}
}

func TestStringConcatenation(t *testing.T){
	input := `"Hello"+" "+"World!"`
	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("Object is not String. got= %T", evaluated)
	}
	if str.Value != "Hello World!" {
		t.Fatalf("The expected value of concatenated string %s but got %s", "Hello World!",str.Value)
	}
}

func TestBuiltInFunction(t *testing.T){
	tests := [] struct{
		input string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{"len(1)", "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got 2, want=1"},

	}

	for _, tt := range tests{
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type){
		case int:
			testIntegerObject(t,evaluated,int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok{
				t.Errorf("object is not error got.%T(+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected{
				t.Errorf("wrong error message. expected %q got %q", expected, errObj.Message)
			}
		}
	}
}