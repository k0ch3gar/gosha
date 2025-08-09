package analyzer

import (
	"fmt"

	"kstmc.com/gosha/internal/ast"
	"kstmc.com/gosha/internal/object"
)

func AnalyzeProgram(node *ast.Program, env *object.Environment) []string {
	env = object.NewEnclosedEnvironment(env)
	var errors []string
	for _, stmt := range node.Statements {
		errors = append(errors, analyzeStatement(stmt, object.ANY_OBJ, env)...)
	}

	return errors
}

func analyzeBlockStatement(node *ast.BlockStatement, returnType object.ObjectType, env *object.Environment) []string {
	var errors []string
	for _, stmt := range node.Statements {
		errors = append(errors, analyzeStatement(stmt, returnType, env)...)
	}

	return errors
}

func analyzeStatement(stmt ast.Statement, returnType object.ObjectType, env *object.Environment) []string {
	switch stmt := stmt.(type) {
	case *ast.ExpressionStatement:
		return analyzeExpressionStatement(stmt, env)
	case *ast.ReturnStatement:
		return analyzeReturnStatement(stmt, returnType, env)
	case *ast.AssignStatement:
		return analyzeAssignStatement(stmt, env)
	case *ast.IfStatement:
		return analyzeIfStatement(stmt, returnType, env)
	case *ast.VarStatement:
		return analyzeVarStatement(stmt, env)
	case *ast.InitAssignStatement:
		return analyzeInitAssignStatement(stmt, env)
	default:
		return []string{fmt.Sprintf("Analyzer error. Unsupported statement %T", stmt)}
	}
}

func analyzeAssignStatement(stmt *ast.AssignStatement, env *object.Environment) []string {
	var errors []string
	if !env.Contains(stmt.Name.Value) {
		msg := fmt.Sprintf("Analyzer error. Unknown identifier %s", stmt.Name.Value)
		errors = append(errors, msg)
		return errors
	}

	ident, _ := env.Get(stmt.Name.Value)

	exprType, errors := AnalyzeExpression(stmt.Value, env)
	if len(errors) != 0 {
		return errors
	}

	if !env.Contains(stmt.Name.Value) {
		msg := fmt.Sprintf("Analyzer error. Unknown identifier %s", stmt.Name.Value)
		errors = append(errors, msg)
		return errors
	}

	if ident.Type() != object.ANY_OBJ && ident.Type() != exprType {
		msg := fmt.Sprintf("Analyzer error. type mismatch. expected %s, got %s", ident.Type(), exprType)
		errors = append(errors, msg)
	}

	return errors
}

func analyzeInitAssignStatement(stmt *ast.InitAssignStatement, env *object.Environment) []string {
	_, errors := AnalyzeExpression(stmt.Value, env)
	return errors
}

func analyzeVarStatement(stmt *ast.VarStatement, env *object.Environment) []string {
	var identType object.ObjectType
	if stmt.Name.DataType != nil {
		identType = RawTypeToObj(stmt.Name.DataType.Name).Type()
	} else {
		identType = object.ANY_OBJ
	}

	exprType, errors := AnalyzeExpression(stmt.Value, env)
	if len(errors) != 0 {
		return errors
	}

	if identType != object.ANY_OBJ && identType != exprType {
		msg := fmt.Sprintf("Analyzer error. type mismatch. expected %s, got %s", identType, exprType)
		errors = append(errors, msg)
	}

	return errors
}

func analyzeReturnStatement(stmt *ast.ReturnStatement, returnType object.ObjectType, env *object.Environment) []string {
	stmtReturnType, errors := AnalyzeExpression(stmt.ReturnValue, env)
	if len(errors) != 0 {
		return errors
	}

	if returnType != object.ANY_OBJ && stmtReturnType != returnType {
		msg := fmt.Sprintf("analyzer error. function returns %s, got=%s", returnType, stmtReturnType)
		errors = append(errors, msg)
	}

	return errors
}

func analyzeExpressionStatement(expr *ast.ExpressionStatement, env *object.Environment) []string {
	_, errors := AnalyzeExpression(expr.Expression, env)
	return errors
}

func AnalyzeExpression(expr ast.Expression, env *object.Environment) (object.ObjectType, []string) {
	var errors []string
	switch expr := expr.(type) {
	case *ast.IntegerLiteral:
		return object.INTEGER_OBJ, errors
	case *ast.Boolean:
		return object.BOOLEAN_OBJ, errors
	case *ast.Identifier:
		obj, ok := env.Get(expr.Value)
		if !ok {
			msg := fmt.Sprintf("analyzer error. unknown identifier %s", expr.Value)
			errors = append(errors, msg)
			return "", errors
		}

		return obj.Type(), nil
	case *ast.CallExpression:
		return analyzeCallExpression(expr, env)
	case *ast.PrefixExpression:
		return analyzePrefixExpression(expr, env)
	case *ast.StringLiteral:
		return object.STRING_OBJ, errors
	case *ast.InfixExpression:
		return analyzeInfixExpression(expr, env)
	case *ast.FunctionLiteral:
		return analyzeFunctionLiteral(expr, env)
	default:
		msg := fmt.Sprintf("analyzer error. unexpected expression type %T", expr)
		errors = append(errors, msg)
		return "", errors
	}
}

func analyzeIfStatement(expr *ast.IfStatement, returnType object.ObjectType, env *object.Environment) []string {
	conditionType, errors := AnalyzeExpression(expr.Condition, env)
	if len(errors) != 0 {
		return errors
	}

	if conditionType != object.BOOLEAN_OBJ {
		msg := fmt.Sprintf("Analyzer error. expected boolean type for if expression, got %s", conditionType)
		return []string{msg}
	}

	errors = append(errors, analyzeBlockStatement(expr.Consequence, returnType, env)...)
	errors = append(errors, analyzeBlockStatement(expr.Alternative, returnType, env)...)
	return errors
}

func analyzeCallExpressionThroughAnonymousFunc(expr *ast.CallExpression, env *object.Environment) (object.ObjectType, []string) {
	fn, ok := expr.Function.(*ast.FunctionLiteral)
	if !ok {
		msg := fmt.Sprintf("Analyze error. Unsupported %s function call", expr.Function.String())
		return "", []string{msg}
	}

	var errors []string
	for i, param := range fn.Parameters {
		var arg object.ObjectType
		arg, tempErrors := AnalyzeExpression(expr.Arguments[i], env)
		if len(tempErrors) != 0 {
			return "", append(errors, tempErrors...)
		}

		if RawTypeToObj(param.DataType.Name).Type() != arg {
			msg := fmt.Sprintf("analyzer error. Incorrect type passed into %s function. expected %s, got=%s", fn.Name.Value, RawTypeToObj(param.DataType.Name).Type(), arg)
			errors = append(errors, msg)
		}
	}

	if len(errors) != 0 {
		return "", errors
	}

	return RawTypeToObj(fn.ReturnType.Name).Type(), nil
}

func analyzeCallExpression(expr *ast.CallExpression, env *object.Environment) (object.ObjectType, []string) {
	fnName, ok := expr.Function.(*ast.Identifier)
	if !ok {
		return analyzeCallExpressionThroughAnonymousFunc(expr, env)
	}

	fnObj, ok := env.Get(fnName.Value)
	if !ok {
		msg := fmt.Sprintf("Analyzer error. Unknown function %s", fnName.Value)
		return "", []string{msg}
	}

	fn, ok := fnObj.(*object.Function)
	if !ok {
		msg := fmt.Sprintf("Analyzer error. Unsupported call for type %T", fnObj)
		return "", []string{msg}
	}

	var errors []string
	for i, param := range fn.Parameters {
		var arg object.ObjectType
		arg, tempErrors := AnalyzeExpression(expr.Arguments[i], env)
		if len(tempErrors) != 0 {
			return "", append(errors, tempErrors...)
		}

		if RawTypeToObj(param.DataType.Name).Type() != arg {
			msg := fmt.Sprintf("analyzer error. Incorrect type passed into %s function. expected %s, got=%s", fn.Name.Value, RawTypeToObj(param.DataType.Name).Type(), arg)
			errors = append(errors, msg)
		}
	}

	if len(errors) != 0 {
		return "", errors
	}

	return RawTypeToObj(fn.ReturnType.Name).Type(), nil
}

func analyzeFunctionLiteral(expr *ast.FunctionLiteral, env *object.Environment) (object.ObjectType, []string) {
	env = object.NewEnclosedEnvironment(env)
	for _, ident := range expr.Parameters {
		env.Set(ident.Value, RawTypeToObj(ident.DataType.Name))
	}

	errors := analyzeBlockStatement(expr.Body, RawTypeToObj(expr.ReturnType.Name).Type(), env)
	return object.FUNCTION_OBJ, errors
}

func RawTypeToObj(rawType string) object.Object {
	switch rawType {
	case "int":
		return &object.Integer{}
	case "bool":
		return &object.Boolean{}
	case "nil":
		return &object.Nil{}
	case "string":
		return &object.String{}
	case "any":
		return &object.Any{}
	default:
		return &object.Nil{}
	}
}

func analyzeInfixExpression(expr *ast.InfixExpression, env *object.Environment) (object.ObjectType, []string) {
	leftType, errors := AnalyzeExpression(expr.Left, env)
	if len(errors) != 0 {
		return "", errors
	}

	rightType, errors := AnalyzeExpression(expr.Right, env)
	if len(errors) != 0 {
		return "", errors
	}

	switch expr.Operator {
	case "+":
		return analyzePlusInfixOperator(leftType, rightType)
	case "-":
		return analyzeMinusInfixOperator(leftType, rightType)
	case "*":
		return analyzeAsteriksInfixOperator(leftType, rightType)
	case "/":
		return analyzeSlashInfixOperator(leftType, rightType)
	case "==":
		return analyzeEqInfixOperator(leftType, rightType)
	case "!=":
		return analyzeNeqInfixOperator(leftType, rightType)
	case "<":
		return analyzeLtInfixOperator(leftType, rightType)
	case ">":
		return analyzeGtInfixOperator(leftType, rightType)
	default:
		msg := fmt.Sprintf("analyzer error. unsupported infix operator type %s", expr.Operator)
		errors = append(errors, msg)
		return "", errors
	}
}

func analyzeNeqInfixOperator(leftType, rightType object.ObjectType) (object.ObjectType, []string) {
	if leftType == rightType {
		return object.BOOLEAN_OBJ, nil
	} else {
		msg := fmt.Sprintf("analyzer error. unsupported comparison for '!=' operator: %s and %s", leftType, rightType)
		errors := []string{msg}
		return "", errors
	}
}

func analyzeEqInfixOperator(leftType, rightType object.ObjectType) (object.ObjectType, []string) {
	if leftType == rightType {
		return object.BOOLEAN_OBJ, nil
	} else {
		msg := fmt.Sprintf("analyzer error. unsupported comparison for '==' operator: %s and %s", leftType, rightType)
		errors := []string{msg}
		return "", errors
	}
}

func analyzeLtInfixOperator(leftType, rightType object.ObjectType) (object.ObjectType, []string) {
	switch {
	case leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ:
		return object.BOOLEAN_OBJ, nil
	default:
		msg := fmt.Sprintf("analyzer error. unsupported expression type for '<' operator %s", rightType)
		errors := []string{msg}
		return "", errors
	}
}

func analyzeGtInfixOperator(leftType, rightType object.ObjectType) (object.ObjectType, []string) {
	switch {
	case leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ:
		return object.BOOLEAN_OBJ, nil
	default:
		msg := fmt.Sprintf("analyzer error. unsupported expression type for '>' operator %s", rightType)
		errors := []string{msg}
		return "", errors
	}
}

func analyzeAsteriksInfixOperator(leftType, rightType object.ObjectType) (object.ObjectType, []string) {
	switch {
	case leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ:
		return object.INTEGER_OBJ, nil
	default:
		msg := fmt.Sprintf("analyzer error. unsupported expression type for '*' operator %s", rightType)
		errors := []string{msg}
		return "", errors
	}
}

func analyzeSlashInfixOperator(leftType, rightType object.ObjectType) (object.ObjectType, []string) {
	switch {
	case leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ:
		return object.INTEGER_OBJ, nil
	default:
		msg := fmt.Sprintf("analyzer error. unsupported expression type for '/' operator %s", rightType)
		errors := []string{msg}
		return "", errors
	}
}

func analyzeMinusInfixOperator(leftType, rightType object.ObjectType) (object.ObjectType, []string) {
	switch {
	case leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ:
		return object.INTEGER_OBJ, nil
	default:
		msg := fmt.Sprintf("analyzer error. unsupported expression type for '-' operator %s", rightType)
		errors := []string{msg}
		return "", errors
	}
}

func analyzePlusInfixOperator(leftType, rightType object.ObjectType) (object.ObjectType, []string) {
	switch {
	case leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ:
		return object.INTEGER_OBJ, nil
	case leftType == object.STRING_OBJ && rightType == object.STRING_OBJ:
		return object.STRING_OBJ, nil
	default:
		msg := fmt.Sprintf("analyzer error. unsupported expression type for '+' operator %s", rightType)
		errors := []string{msg}
		return "", errors
	}
}

func analyzePrefixExpression(expr *ast.PrefixExpression, env *object.Environment) (object.ObjectType, []string) {
	rightType, errors := AnalyzeExpression(expr.Right, env)
	if len(errors) != 0 {
		return object.NIL_OBJ, errors
	}

	switch expr.Operator {
	case "!":
		return analyzeBangPrefixOperator(rightType)
	case "-":
		return analyzeMinusPrefixOperator(rightType)
	default:
		msg := fmt.Sprintf("analyzer error. unsupportet prefix operator type %s", expr.Operator)
		errors = append(errors, msg)
		return "", errors
	}
}

func analyzeMinusPrefixOperator(rightType object.ObjectType) (object.ObjectType, []string) {
	if rightType == object.INTEGER_OBJ {
		return rightType, nil
	} else {
		msg := fmt.Sprintf("analyzer error. unsupported expression type for '-' operator %s", rightType)
		errors := []string{msg}
		return "", errors
	}
}

func analyzeBangPrefixOperator(rightType object.ObjectType) (object.ObjectType, []string) {
	if rightType == object.BOOLEAN_OBJ {
		return rightType, nil
	} else {
		msg := fmt.Sprintf("analyzer error. unsupported expression type for '!' operator %s", rightType)
		errors := []string{msg}
		return "", errors
	}
}
