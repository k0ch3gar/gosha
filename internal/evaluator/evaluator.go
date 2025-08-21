package evaluator

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"kstmc.com/gosha/internal/analyzer"
	"kstmc.com/gosha/internal/ast"
	"kstmc.com/gosha/internal/object"
	"kstmc.com/gosha/internal/parser"
	"kstmc.com/gosha/internal/token"
)

var (
	NIL   = &object.Nil{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}

	isBashCommandInteractive map[string]bool = map[string]bool{
		"vi":    true,
		"vim":   true,
		"emacs": true,
		"nano":  true,
		"links": true,
	}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.ReadChanExpression:
		obj, _ := env.Get(node.Source.Value)
		chn := obj.(*object.ChanObject)
		val := <-chn.Chan
		return val
	case *ast.SendChanStatement:
		obj, _ := env.Get(node.Destination.Value)
		chn := obj.(*object.ChanObject)
		val := Eval(node.Source, env)
		chn.Chan <- val
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
	case *ast.DataTypeExpression:
		return &object.DataTypeObject{DataType: node.Type}
	case *ast.AssignStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}

		if !env.Contains(node.Name.Value) {
			return newError("unknown variable: %q", node.Name.Value)
		}

		env.Update(node.Name.Value, val)
	case *ast.GoStatement:
		return evalGoStatement(node.Expr, env)
	case *ast.VarStatement:
		var val object.Object
		if node.Value != nil {
			val = Eval(node.Value, env)
			if isError(val) {
				return val
			}
		}

		if env.Contains(node.Name.Value) {
			return newError("variable with such name already exists: %q", node.Name.Value)
		}

		if val == nil && node.Name.DataType != nil {
			env.Set(node.Name.Value, analyzer.NativeTypeToDefaultObj(*node.Name.DataType))
		} else if node.Name.DataType != nil && (*node.Name.DataType).Name() == parser.ANY.Name() {
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

		//switch fn := function.(type) {
		//case *object.Builtin:
		//	return applyBuiltin(fn, node.Arguments, env)
		//}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args, env)
	case *ast.IfStatement:
		return evalIfExpression(node, env)
	case *ast.ForStatement:
		return evalForStatement(node, env)
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
	case *ast.BashExpression:
		return evalBashExpression(node, env)
	case *ast.BashVarExpression:
		return evalBashVarExpression(node, env)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalPrefixExpression(node.Operator, right)
	case *ast.IndexExpression:
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}

		intIndex, ok := index.(*object.Integer)
		if !ok {
			return newError("expected integer type for index expression, got=%T", index)
		}

		obj := Eval(node.Left, env)
		if isError(index) {
			return obj
		}

		slice, ok := obj.(*object.SliceObject)
		if !ok {
			return newError("expected slice obect for index expression, got=%T", obj)
		}

		if intIndex.Value < 0 || int64(len(slice.Values)) <= intIndex.Value {
			return newError("out of bound error for slice '%s' at index %d", node.Left.String(), intIndex.Value)
		}

		return slice.Values[intIndex.Value]
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

func evalGoStatement(expr ast.Expression, env *object.Environment) object.Object {
	go Eval(expr, env)
	return NIL
}

func evalBashVarExpression(node *ast.BashVarExpression, env *object.Environment) object.Object {
	num, ok := strconv.Atoi(node.Value[1:])
	if ok == nil {
		if len(os.Args) <= num+1 {
			return newError("provided %d args but %d arg was called", len(os.Args)-1, num)
		}
		return &object.String{
			Value: os.Args[num+1],
		}
	}

	output, err := exec.Command("bash", "-c", "echo "+node.Value).Output()
	if err != nil {
		return newError("bash error %s", err.Error())
	}

	return &object.String{Value: string(output)}
}

func evalForStatement(stmt *ast.ForStatement, env *object.Environment) object.Object {
	for {
		condition := Eval(stmt.Condition, env)
		if isError(condition) {
			return condition
		}

		if condition == TRUE {
			result := Eval(stmt.Consequence, env)
			if result.Type().Name() == parser.RETURN.Name() {
				return result
			}
		} else {
			return NIL
		}
	}
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		env = object.NewEnclosedEnvironment(env)
		errors := analyzer.AnalyzeStatement(statement, parser.ANY, env)
		env = object.UnwrapEnvironment(env)
		if len(errors) != 0 {
			return newError("analyzer error %s", errors[0])
		}

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

func evalBashExpression(expr *ast.BashExpression, env *object.Environment) object.Object {
	var bashArgs []string
	for _, arg := range expr.Value {
		if arg[0] == '$' && env.Contains(arg[1:]) {
			arg, _ := env.Get(arg[1:])
			bashArgs = append(bashArgs, arg.Inspect())
		} else {
			bashArgs = append(bashArgs, arg)
		}
	}

	result, ok := isBashCommandInteractive[bashArgs[0]]
	if ok && result {
		cmd := exec.Command("bash", "-c", strings.Join(bashArgs, " "))
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setsid: true,
			Ctty:   int(os.Stdin.Fd()),
		}

		err := cmd.Run()

		if err != nil {
			return newError("bash error %s", err.Error())
		}
	} else {
		output, err := exec.Command("bash", "-c", strings.Join(bashArgs, " ")).Output()
		if err != nil {
			return newError("bash error %s", err.Error())
		}

		return &object.String{Value: string(output)}
	}

	return &object.String{Value: ""}
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

func applyFunction(fn object.Object, args []object.Object, env *object.Environment) object.Object {
	switch fn := fn.(type) {
	case *object.Builtin:
		return fn.Fn(args...)
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	default:
		return newError("not a function: %s", fn.Type().Name())
	}

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
	if ok {
		return val
	}

	builtin, ok := object.Builtins[node.Value]
	if ok {
		return builtin
	}

	return newError("unknown identifier: %s", node.Value)
}

func evalBlockStatement(bs *ast.BlockStatement, env *object.Environment) object.Object {
	env = object.NewEnclosedEnvironment(env)
	var result object.Object

	for _, statement := range bs.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == parser.RETURN || rt == parser.ERROR {
				env = object.UnwrapEnvironment(env)
				return result
			}
		}
	}

	env = object.UnwrapEnvironment(env)
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
	case left.Type() == parser.STRING && right.Type() == parser.STRING:
		return evalStringInfixExpression(operator, left, right)
	case operator == token.EQ:
		return rawBooleanToBooleanObject(left == right)
	case operator == token.NEQ:
		return rawBooleanToBooleanObject(left != right)
	case left.Type() == parser.BOOLEAN && right.Type() == parser.BOOLEAN:
		return evalBooleanInfixExpression(operator, left, right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type().Name(), operator, right.Type().Name())
	default:
		return newError("unknown operator: %s %s %s", left.Type().Name(), operator, right.Type().Name())
	}
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	case "==":
		if leftVal == rightVal {
			return TRUE
		} else {
			return FALSE
		}
	case "!=":
		if leftVal != rightVal {
			return TRUE
		} else {
			return FALSE
		}
	default:
		return newError("unsupported operator: %s %s %s", left.Type().Name(), operator, right.Type().Name())
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
		return newError("unknown operator: %s %s %s", left.Type().Name(), operator, right.Type().Name())
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
	case token.PERCENT:
		return &object.Integer{Value: leftVal % rightVal}
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
	case token.FOPER:
		return evalFoperPrefixOperatorExpression(right)
	case token.ASTERISK:
		return evalAsteriskPrefixOperatorExpression(right)
	case token.REF:
		return evalRefPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalRefPrefixOperatorExpression(right object.Object) object.Object {
	return &object.ReferenceObject{
		Value: &right,
	}
}

func evalAsteriskPrefixOperatorExpression(right object.Object) object.Object {
	switch right := right.(type) {
	case *object.ReferenceObject:
		return *right.Value
	default:
		return newError("expected reference: %s", right.Inspect())
	}
}

func evalFoperPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != parser.STRING {
		return newError("unknown operator: -f %s", right.Type())
	}

	value := right.(*object.String).Value
	_, err := os.Stat(value)
	if os.IsNotExist(err) {
		return FALSE
	} else {
		return TRUE
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
