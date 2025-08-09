package analyzer

import (
	"fmt"

	"kstmc.com/gosha/internal/ast"
	"kstmc.com/gosha/internal/object"
	"kstmc.com/gosha/internal/parser"
)

func AnalyzeProgram(node *ast.Program, env *object.Environment) []string {
	env = object.NewEnclosedEnvironment(env)
	var errors []string
	for _, stmt := range node.Statements {
		errors = append(errors, AnalyzeStatement(stmt, parser.ANY, env)...)
	}

	return errors
}

func analyzeBlockStatement(node *ast.BlockStatement, returnType ast.DataType, env *object.Environment) []string {
	var errors []string
	for _, stmt := range node.Statements {
		errors = append(errors, AnalyzeStatement(stmt, returnType, env)...)
	}

	return errors
}

func AnalyzeStatement(stmt ast.Statement, returnType ast.DataType, env *object.Environment) []string {
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

	if ident.Type() != parser.ANY && ident.Type().Name() != exprType.Name() {
		msg := fmt.Sprintf("Analyzer error. type mismatch. expected %s, got %s", ident.Type().Name(), exprType.Name())
		errors = append(errors, msg)
	}

	return errors
}

func analyzeInitAssignStatement(stmt *ast.InitAssignStatement, env *object.Environment) []string {
	_, errors := AnalyzeExpression(stmt.Value, env)
	return errors
}

func analyzeVarStatement(stmt *ast.VarStatement, env *object.Environment) []string {
	var identType ast.DataType
	if stmt.Name.DataType != nil {
		identType = *stmt.Name.DataType
	} else {
		identType = parser.ANY
	}

	exprType, errors := AnalyzeExpression(stmt.Value, env)
	if len(errors) != 0 {
		return errors
	}

	if identType != parser.ANY && identType.Name() != exprType.Name() {
		msg := fmt.Sprintf("Analyzer error. type mismatch. expected %s, got %s", identType.Name(), exprType.Name())
		errors = append(errors, msg)
	}

	return errors
}

func analyzeReturnStatement(stmt *ast.ReturnStatement, returnType ast.DataType, env *object.Environment) []string {
	stmtReturnType, errors := AnalyzeExpression(stmt.ReturnValue, env)
	if len(errors) != 0 {
		return errors
	}

	if returnType != parser.ANY && stmtReturnType.Name() != returnType.Name() {
		msg := fmt.Sprintf("analyzer error. function returns %s, got=%s", returnType.Name(), stmtReturnType.Name())
		errors = append(errors, msg)
	}

	return errors
}

func analyzeExpressionStatement(expr *ast.ExpressionStatement, env *object.Environment) []string {
	_, errors := AnalyzeExpression(expr.Expression, env)
	return errors
}

func AnalyzeExpression(expr ast.Expression, env *object.Environment) (ast.DataType, []string) {
	var errors []string
	switch expr := expr.(type) {
	case *ast.IntegerLiteral:
		return parser.INT, errors
	case *ast.Boolean:
		return parser.BOOLEAN, errors
	case *ast.BashExpression:
		return parser.STRING, errors
	case *ast.Identifier:
		obj, ok := env.Get(expr.Value)
		if !ok {
			msg := fmt.Sprintf("analyzer error. unknown identifier %s", expr.Value)
			errors = append(errors, msg)
			return nil, errors
		}

		return obj.Type(), nil
	case *ast.CallExpression:
		return analyzeCallExpression(expr, env)
	case *ast.PrefixExpression:
		return analyzePrefixExpression(expr, env)
	case *ast.StringLiteral:
		return parser.STRING, errors
	case *ast.InfixExpression:
		return analyzeInfixExpression(expr, env)
	case *ast.FunctionLiteral:
		return analyzeFunctionLiteral(expr, env)
	default:
		msg := fmt.Sprintf("analyzer error. unexpected expression type %T", expr)
		errors = append(errors, msg)
		return nil, errors
	}
}

func analyzeIfStatement(expr *ast.IfStatement, returnType ast.DataType, env *object.Environment) []string {
	conditionType, errors := AnalyzeExpression(expr.Condition, env)
	if len(errors) != 0 {
		return errors
	}

	if conditionType != parser.BOOLEAN {
		msg := fmt.Sprintf("Analyzer error. expected boolean type for if expression, got %s", conditionType.Name())
		return []string{msg}
	}

	errors = append(errors, analyzeBlockStatement(expr.Consequence, returnType, env)...)
	errors = append(errors, analyzeBlockStatement(expr.Alternative, returnType, env)...)
	return errors
}

func analyzeCallExpression(expr *ast.CallExpression, env *object.Environment) (ast.DataType, []string) {
	dType, errors := AnalyzeExpression(expr.Function, env)
	if len(errors) != 0 {
		return nil, errors
	}

	fnType, ok := dType.(*ast.FunctionDataType)
	if !ok {
		msg := fmt.Sprintf("Analyzer error. Unsupported call type %T", dType)
		return nil, []string{msg}
	}

	for i, param := range expr.Arguments {
		var arg ast.DataType
		arg, tempErrors := AnalyzeExpression(param, env)
		if len(tempErrors) != 0 {
			return nil, append(errors, tempErrors...)
		}

		if fnType.Parameters[i].Name() != arg.Name() {
			msg := fmt.Sprintf("analyzer error. Incorrect type passed into function. expected %s, got=%s", fnType.Name(), arg.Name())
			errors = append(errors, msg)
		}
	}

	if len(errors) != 0 {
		return nil, errors
	}

	return fnType.ReturnValue, nil
}

func analyzeFunctionLiteral(expr *ast.FunctionLiteral, env *object.Environment) (ast.DataType, []string) {
	env = object.NewEnclosedEnvironment(env)
	fn := &ast.FunctionDataType{
		ReturnValue: expr.ReturnType,
	}

	for _, ident := range expr.Parameters {
		env.Set(ident.Value, RawTypeToObj(*ident.DataType))
		fn.Parameters = append(fn.Parameters, *ident.DataType)
	}

	errors := analyzeBlockStatement(expr.Body, expr.ReturnType, env)
	return fn, errors
}

func RawTypeToObj(rawType ast.DataType) object.Object {
	switch rawType.(type) {
	case *ast.IntegerDataType:
		return &object.Integer{}
	case *ast.BooleanDataType:
		return &object.Boolean{}
	case *ast.StringDataType:
		return &object.String{}
	case *ast.AnyDataType:
		return &object.Any{}
	default:
		return &object.Nil{}
	}
}

func analyzeInfixExpression(expr *ast.InfixExpression, env *object.Environment) (ast.DataType, []string) {
	leftType, errors := AnalyzeExpression(expr.Left, env)
	if len(errors) != 0 {
		return nil, errors
	}

	rightType, errors := AnalyzeExpression(expr.Right, env)
	if len(errors) != 0 {
		return nil, errors
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
		return nil, errors
	}
}

func analyzeNeqInfixOperator(leftType, rightType ast.DataType) (ast.DataType, []string) {
	if leftType == rightType {
		return parser.BOOLEAN, nil
	} else {
		msg := fmt.Sprintf("analyzer error. unsupported comparison for '!=' operator: %s and %s", leftType.Name(), rightType.Name())
		errors := []string{msg}
		return nil, errors
	}
}

func analyzeEqInfixOperator(leftType, rightType ast.DataType) (ast.DataType, []string) {
	if leftType == rightType {
		return parser.BOOLEAN, nil
	} else {
		msg := fmt.Sprintf("analyzer error. unsupported comparison for '==' operator: %s and %s", leftType.Name(), rightType.Name())
		errors := []string{msg}
		return nil, errors
	}
}

func analyzeLtInfixOperator(leftType, rightType ast.DataType) (ast.DataType, []string) {
	switch {
	case leftType == parser.INT && rightType == parser.INT:
		return parser.BOOLEAN, nil
	default:
		msg := fmt.Sprintf("analyzer error. unsupported expression type for '<' operator %s", rightType.Name())
		errors := []string{msg}
		return nil, errors
	}
}

func analyzeGtInfixOperator(leftType, rightType ast.DataType) (ast.DataType, []string) {
	switch {
	case leftType == parser.INT && rightType == parser.INT:
		return parser.BOOLEAN, nil
	default:
		msg := fmt.Sprintf("analyzer error. unsupported expression type for '>' operator %s", rightType.Name())
		errors := []string{msg}
		return nil, errors
	}
}

func analyzeAsteriksInfixOperator(leftType, rightType ast.DataType) (ast.DataType, []string) {
	switch {
	case leftType == parser.INT && rightType == parser.INT:
		return parser.INT, nil
	default:
		msg := fmt.Sprintf("analyzer error. unsupported expression type for '*' operator %s", rightType.Name())
		errors := []string{msg}
		return nil, errors
	}
}

func analyzeSlashInfixOperator(leftType, rightType ast.DataType) (ast.DataType, []string) {
	switch {
	case leftType == parser.INT && rightType == parser.INT:
		return parser.INT, nil
	default:
		msg := fmt.Sprintf("analyzer error. unsupported expression type for '/' operator %s", rightType.Name())
		errors := []string{msg}
		return nil, errors
	}
}

func analyzeMinusInfixOperator(leftType, rightType ast.DataType) (ast.DataType, []string) {
	switch {
	case leftType == parser.INT && rightType == parser.INT:
		return parser.INT, nil
	default:
		msg := fmt.Sprintf("analyzer error. unsupported expression type for '-' operator %s", rightType.Name())
		errors := []string{msg}
		return nil, errors
	}
}

func analyzePlusInfixOperator(leftType, rightType ast.DataType) (ast.DataType, []string) {
	switch {
	case leftType == parser.INT && rightType == parser.INT:
		return parser.INT, nil
	case leftType == parser.STRING && rightType == parser.STRING:
		return parser.STRING, nil
	default:
		msg := fmt.Sprintf("analyzer error. unsupported expression type for '+' operator %s", rightType.Name())
		errors := []string{msg}
		return nil, errors
	}
}

func analyzePrefixExpression(expr *ast.PrefixExpression, env *object.Environment) (ast.DataType, []string) {
	rightType, errors := AnalyzeExpression(expr.Right, env)
	if len(errors) != 0 {
		return parser.NIL, errors
	}

	switch expr.Operator {
	case "!":
		return analyzeBangPrefixOperator(rightType)
	case "-":
		return analyzeMinusPrefixOperator(rightType)
	default:
		msg := fmt.Sprintf("analyzer error. unsupportet prefix operator type %s", expr.Operator)
		errors = append(errors, msg)
		return nil, errors
	}
}

func analyzeMinusPrefixOperator(rightType ast.DataType) (ast.DataType, []string) {
	if rightType == parser.INT {
		return rightType, nil
	} else {
		msg := fmt.Sprintf("analyzer error. unsupported expression type for '-' operator %s", rightType.Name())
		errors := []string{msg}
		return nil, errors
	}
}

func analyzeBangPrefixOperator(rightType ast.DataType) (ast.DataType, []string) {
	if rightType == parser.BOOLEAN {
		return rightType, nil
	} else {
		msg := fmt.Sprintf("analyzer error. unsupported expression type for '!' operator %s", rightType.Name())
		errors := []string{msg}
		return nil, errors
	}
}
