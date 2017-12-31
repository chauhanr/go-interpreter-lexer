package evaluator

import (
	"go-interpreter-lexer/ast"
	"go-interpreter-lexer/object"
)

var(
	TRUE = &object.Boolean{Bool: true}
	FALSE = &object.Boolean{Bool: false}
	NULL = &object.Null{}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type){
		case *ast.Program:
			return evalStatatments(node.Statements)
		case *ast.ExpressionStatement:
			return Eval(node.Expression)
		case *ast.IntegerLiteral:
			return &object.Integer{Value: node.Value}
		case *ast.Boolean:
			return nativeBoolToBooleanObject(node.Value)
		case *ast.PrefixExpression:
			right := Eval(node.Right)
			return evalPrefixExpression(node.Operator,right)
	case *ast.InfixExpression:
			right := Eval(node.Right)
			left := Eval(node.Left)
			return evalInfixExpression(node.Operator, left, right)
	}
	return nil
}

func evalPrefixExpression(operator string, right object.Object) object.Object{
	switch operator{
		case "!":
			return evalBangOperatorExpression(right)
		case "-":
			return evalMinusOperatorExpression(right)
		default:
			return NULL
	}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object{
	switch {
		case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		  return evalIntegerInfixExpression(operator, left, right)
		case operator == "==":
			return nativeBoolToBooleanObject(left == right)
		case operator == "!=" :
			return nativeBoolToBooleanObject(left != right)
		default :
			return NULL
	}
}

func evalIntegerInfixExpression(operator string, leftVal object.Object, rightVal object.Object) object.Object{
	left := leftVal.(*object.Integer).Value
	right := rightVal.(*object.Integer).Value

	switch operator{
		case "+":
			return &object.Integer{Value: (left + right)}
		case "-" :
			return &object.Integer{Value: (left - right)}
		case "*" :
			return &object.Integer{Value: (left * right)}
		case "/":
			if right == 0{
				return NULL
			}else {
				return &object.Integer{Value: (left / right)}
			}
		case ">":
			return nativeBoolToBooleanObject(left > right)
		case "<":
			return nativeBoolToBooleanObject(left < right)
		case "==":
			return nativeBoolToBooleanObject(left == right)
		case "!=":
			return nativeBoolToBooleanObject(left != right)
		default:
			return NULL
	}
}

func evalBangOperatorExpression(right object.Object) object.Object{
	switch right{
		case TRUE:
			return FALSE
		case FALSE:
			return TRUE
		case NULL:
			return TRUE
		default :
			return FALSE
	}
}

func evalMinusOperatorExpression(right object.Object) object.Object{
	if right.Type() != object.INTEGER_OBJ{
		return NULL
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func nativeBoolToBooleanObject(Bool bool) *object.Boolean{
	if Bool == true{
		return TRUE
	}else{
		return FALSE
	}
}

func evalStatatments(statements []ast.Statement) object.Object{
	var result object.Object
	for _, stmt := range statements{
		result = Eval(stmt)
	}
	return result
}

