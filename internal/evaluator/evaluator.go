package evaluator

import (
	"fmt"

	"kstmc.com/gosha/internal/ast"
	"kstmc.com/gosha/internal/object"
	"kstmc.com/gosha/internal/parser"
	"kstmc.com/gosha/internal/token"
)

var (
	NIL   = &object.Nil{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}

		return &object.ReturnValue{Value: val}
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.AssignStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}

		if !env.Contains(node.Name.Value) {
			return newError("unknown variable: %q", node.Name.Value)
		}

		env.Set(node.Name.Value, val)
	case *ast.VarStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}

		if env.Contains(node.Name.Value) {
			return newError("variable with such name already exists: %q", node.Name.Value)
		}

		if node.Name.DataType != nil && (*node.Name.DataType).Name() == parser.ANY.Name() {
			env.Set(node.Name.Value, &object.Any{Value: val.Inspect()})
		} else {
			env.Set(node.Name.Value, val)
		}
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.InitAssignStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}

		if env.Contains(node.Name.Value) {
			return newError("Identifier %q already exists", node.Name.Value)
		}

		env.Set(node.Name.Value, val)
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)
	case *ast.IfStatement:
		return evalIfExpression(node, env)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.Boolean:
		return rawBooleanToBooleanObject(node.Value)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		returnType := node.ReturnType
		name := node.Name
		function := &object.Function{Parameters: params, Env: env, Body: body, ReturnType: returnType, Name: name}
		if name != nil {
			if env.Contains(name.Value) {
				return newError("function %s already exists", name.Value)
			} else {
				env.Set(name.Value, function)
			}
		}

		return function
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)
	}

	return NIL
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	function, ok := fn.(*object.Function)
	if !ok {
		return newError("not a function: %s", fn.Type())
	}

	extendedEnv := extendFunctionEnv(function, args)
	evaluated := Eval(function.Body, extendedEnv)
	return unwrapReturnValue(evaluated)
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for argc, arg := range fn.Parameters {
		env.Set(arg.Value, args[argc])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)
	if !ok {
		return newError("unknown identifier: %s", node.Value)
	}

	return val
}

func evalBlockStatement(bs *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range bs.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == parser.RETURN || rt == parser.ERROR {
				return result
			}
		}
	}

	return result
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == parser.ERROR
	}

	return false
}

func evalIfExpression(ie *ast.IfStatement, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if condition == TRUE {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NIL
	}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == parser.INT && right.Type() == parser.INT:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == token.EQ:
		return rawBooleanToBooleanObject(left == right)
	case operator == token.NEQ:
		return rawBooleanToBooleanObject(left != right)
	case left.Type() == parser.BOOLEAN && right.Type() == parser.BOOLEAN:
		return evalBooleanInfixExpression(operator, left, right)
	case left.Type() == parser.STRING && right.Type() == parser.STRING:
		return evalStringInfixExpression(operator, left, right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	if operator != "+" {
		return newError("unsupported operator: %s %s %s", left.Type(), operator, right.Type())
	}

	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	return &object.String{Value: leftVal + rightVal}
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
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
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
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())

	}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case token.BANG:
		return evalBangOperatorExpression(right)
	case token.MINUS:
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != parser.INT {
		return newError("unknown operator: -%s", right.Type())
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
		return newError("unknown operator: !%s", right.Type())
	}
}

func rawBooleanToBooleanObject(rawBoolean bool) *object.Boolean {
	if rawBoolean {
		return TRUE
	}

	return FALSE
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
