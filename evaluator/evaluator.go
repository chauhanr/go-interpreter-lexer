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
			return evalProgram(node)
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
		case *ast.BlockStatement:
			return evalBlockStatements(node)
		case *ast.IfExpression:
			return evalIfExpression(node)
		case *ast.ReturnStatement:
			val := Eval(node.ReturnValue)
			return &object.ReturnValue{Value: val}

	}
	return nil
}

func evalBlockStatements(block *ast.BlockStatement) object.Object{
	var result object.Object
	for _, stmt := range block.Statements{
		result := Eval(stmt)
		if result != nil && result.Type() == object.RETURN_VALUE_OBJ {
			return result
		}
	}
	return result
}

func evalProgram(program *ast.Program) object.Object{
	var result object.Object
	for _, stmt := range program.Statements{
		result := Eval(stmt)
		if returnValue, ok := result.(*object.ReturnValue); ok{
			return returnValue.Value
		}
	}
	return result
}


func evalIfExpression(ie *ast.IfExpression) object.Object{
	condition := Eval(ie.Condition)

	if isTruthy(condition) {
		return Eval(ie.Consequence)
	}else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	}else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool{
	switch obj {
		case NULL:
			return false
		case TRUE:
			return true
		case FALSE:
			return false
		default:
			return true
	}
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

func evalStatements(statements []ast.Statement) object.Object{
	var result object.Object
	for _, stmt := range statements{
		result = Eval(stmt)
		if returnValue, ok := result.(*object.ReturnValue); ok{
			return returnValue.Value
		}
	}
	return result
}

