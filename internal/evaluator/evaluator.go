package evaluator

import (
	"kstmc.com/gosha/internal/ast"
	"kstmc.com/gosha/internal/object"
	"kstmc.com/gosha/internal/token"
)

var (
	NIL   = &object.Nil{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue)
		return &object.ReturnValue{Value: val}
	case *ast.Program:
		return evalProgram(node)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IfExpression:
		return evalIfExpression(node)
	case *ast.BlockStatement:
		return evalBlockStatement(node)
	case *ast.Boolean:
		return rawBooleanToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)
	}

	return NIL
}

func evalProgram(program *ast.Program) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement)

		returnValue, ok := result.(*object.ReturnValue)
		if ok {
			return returnValue.Value
		}
	}

	return result
}

func evalBlockStatement(bs *ast.BlockStatement) object.Object {
	var result object.Object

	for _, statement := range bs.Statements {
		result = Eval(statement)

		if result != nil && result.Type() == object.RETURN_VALUE_OBJ {
			return result
		}
	}

	return result
}

func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)
	if condition == TRUE {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	} else {
		return NIL
	}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == token.EQ:
		return rawBooleanToBooleanObject(left == right)
	case operator == token.NEQ:
		return rawBooleanToBooleanObject(left != right)
	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
		return evalBooleanInfixExpression(operator, left, right)
	default:
		return NIL
	}
}

func evalBooleanInfixExpression(operator string, left, right object.Object) object.Object {
	switch operator {
	case token.AND:
		return rawBooleanToBooleanObject(left == TRUE && right == TRUE)
	case token.OR:
		return rawBooleanToBooleanObject(left == TRUE || right == TRUE)
	case token.EQ:
		return rawBooleanToBooleanObject(left == right)
	case token.NEQ:
		return rawBooleanToBooleanObject(left != right)
	default:
		return NIL
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case token.PLUS:
		return &object.Integer{Value: leftVal + rightVal}
	case token.MINUS:
		return &object.Integer{Value: leftVal - rightVal}
	case token.ASTERISK:
		return &object.Integer{Value: leftVal * rightVal}
	case token.SLASH:
		return &object.Integer{Value: leftVal / rightVal}
	case token.EQ:
		return rawBooleanToBooleanObject(leftVal == rightVal)
	case token.NEQ:
		return rawBooleanToBooleanObject(leftVal != rightVal)
	case token.LT:
		return rawBooleanToBooleanObject(leftVal < rightVal)
	case token.GT:
		return rawBooleanToBooleanObject(leftVal > rightVal)
	default:
		return NIL
	}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case token.BANG:
		return evalBangOperatorExpression(right)
	case token.MINUS:
		return evalMinusPrefixOperatorExpression(right)
	default:
		return NIL
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return NIL
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case FALSE:
		return TRUE
	case TRUE:
		return FALSE
	default:
		return NIL
	}
}

func rawBooleanToBooleanObject(rawBoolean bool) *object.Boolean {
	if rawBoolean {
		return TRUE
	}

	return FALSE
}
