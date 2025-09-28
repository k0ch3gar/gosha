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
	env = object.NewEnclosedEnvironment(env)
	var errors []string
	for _, stmt := range node.Statements {
		errors = append(errors, AnalyzeStatement(stmt, returnType, env)...)
	}

	env = object.UnwrapEnvironment(env)

	return errors
}

func AnalyzeStatement(stmt ast.Statement, returnType ast.DataType, env *object.Environment) []string {
	switch stmt := stmt.(type) {
	case *ast.ExpressionStatement:
		return analyzeExpressionStatement(stmt, env)
	case *ast.SendChanStatement:
		return analyzeSendChanStatement(stmt, env)
	case *ast.ReturnStatement:
		return analyzeReturnStatement(stmt, returnType, env)
	case *ast.BreakStatement:
		return nil
	case *ast.AssignStatement:
		return analyzeAssignStatement(stmt, env)
	case *ast.IfStatement:
		return analyzeIfStatement(stmt, returnType, env)
	case *ast.GoStatement:
		_, errors := AnalyzeExpression(stmt.Expr, env)
		return errors
	case *ast.VarStatement:
		return analyzeVarStatement(stmt, env)
	case *ast.ForStatement:
		return analyzeForStatement(stmt, returnType, env)
	case *ast.InitAssignStatement:
		return analyzeInitAssignStatement(stmt, env)
	default:
		return []string{fmt.Sprintf("Analyzer error. Unsupported statement %T", stmt)}
	}
}

func analyzeSendChanStatement(stmt *ast.SendChanStatement, env *object.Environment) []string {
	obj, ok := env.Get(stmt.Destination.Value)
	if !ok {
		return []string{fmt.Sprintf("Analyzer error. Unknown identifier %s", stmt.Destination.Value)}
	}

	chn, ok := obj.(*object.ChanObject)
	if !ok {
		return []string{fmt.Sprintf("Analyzer error. Expected identifier to be *object.ChanObject type, got=%T", obj)}
	}

	exprType, errors := AnalyzeExpression(stmt.Source, env)
	if len(errors) > 0 {
		return errors
	}

	if exprType.Name() == parser.ANY.Name() || chn.ChanType.Name() == exprType.Name() {
		return nil
	} else {
		return []string{fmt.Sprintf("Analyzer error. Expression type and chan type mismatch. Chan type %T, expression type %T", chn.ChanType, exprType)}
	}
}

func analyzeForStatement(stmt *ast.ForStatement, returnType ast.DataType, env *object.Environment) []string {
	exprType, errors := AnalyzeExpression(stmt.Condition, env)
	if len(errors) != 0 {
		return errors
	}

	if exprType.Name() != parser.BOOLEAN.Name() {
		msg := fmt.Sprintf("Analyzer error. Expected boolean type for condition, got=%s", exprType.Name())
		return []string{msg}
	}

	errors = analyzeBlockStatement(stmt.Consequence, returnType, env)
	return errors
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

	if ident.Type() != parser.ANY && exprType != parser.ANY && ident.Type().Name() != exprType.Name() {
		msg := fmt.Sprintf("Analyzer error. type mismatch. expected %s, got %s", ident.Type().Name(), exprType.Name())
		errors = append(errors, msg)
	}

	return errors
}

func analyzeInitAssignStatement(stmt *ast.InitAssignStatement, env *object.Environment) []string {
	dType, errors := AnalyzeExpression(stmt.Value, env)
	env.Set(stmt.Name.Value, NativeTypeToDefaultObj(dType))
	return errors
}

func analyzeVarStatement(stmt *ast.VarStatement, env *object.Environment) []string {
	if stmt.Value == nil {
		env.Set(stmt.Name.Value, NativeTypeToDefaultObj(*stmt.Name.DataType))
		return nil
	}

	var identType ast.DataType
	exprType, errors := AnalyzeExpression(stmt.Value, env)
	if len(errors) != 0 {
		return errors
	}

	if stmt.Name.DataType != nil {
		identType = *stmt.Name.DataType
	} else {
		identType = exprType
	}

	if identType != parser.ANY && identType.Name() != exprType.Name() {
		msg := fmt.Sprintf("Analyzer error. type mismatch. expected %s, got %s", identType.Name(), exprType.Name())
		errors = append(errors, msg)
	}

	env.Set(stmt.Name.Value, NativeTypeToDefaultObj(identType))

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
		if ok {
			return obj.Type(), nil
		}

		fnObj, ok := object.Builtins[expr.Value]
		if ok {
			return fnObj.Type(), nil
		}

		msg := fmt.Sprintf("analyzer error. unknown identifier %s", expr.Value)
		errors = append(errors, msg)
		return nil, errors
	case *ast.CallExpression:
		return analyzeCallExpression(expr, env)
	case *ast.IndexExpression:
		return analyzeIndexExpression(expr, env)
	case *ast.ReadChanExpression:
		obj, ok := env.Get(expr.Source.Value)
		if !ok {
			msg := fmt.Sprintf("analyzer error. unknown identifier %s", expr.Source.Value)
			errors = append(errors, msg)
			return nil, errors
		}

		if chn, ok := obj.(*object.ChanObject); !ok {
			msg := fmt.Sprintf("analyzer error. expected identifier to be type *object.ChanObject, got %s", obj)
			errors = append(errors, msg)
			return nil, errors
		} else {
			return chn.ChanType, nil
		}
	case *ast.PrefixExpression:
		return analyzePrefixExpression(expr, env)
	case *ast.StringLiteral:
		return parser.STRING, errors
	case *ast.InfixExpression:
		return analyzeInfixExpression(expr, env)
	case *ast.BashVarExpression:
		return parser.STRING, errors
	case *ast.FunctionLiteral:
		return analyzeFunctionLiteral(expr, env)
	default:
		msg := fmt.Sprintf("analyzer error. unexpected expression type %T", expr)
		errors = append(errors, msg)
		return nil, errors
	}
}

func analyzeIndexExpression(expr *ast.IndexExpression, env *object.Environment) (ast.DataType, []string) {
	lType, errors := AnalyzeExpression(expr.Left, env)
	if len(errors) > 0 {
		return nil, errors
	}

	sliceType, ok := lType.(*ast.SliceDataType)
	if !ok {
		return nil, []string{fmt.Sprintf("Analyzer error. expected slice type for index expression, got=%T", lType)}
	}

	indexType, errors := AnalyzeExpression(expr.Index, env)
	if len(errors) > 0 {
		return nil, errors
	}

	_, ok = indexType.(*ast.IntegerDataType)
	if !ok {
		return nil, []string{fmt.Sprintf("Analyzer error. expected integer type for index expression, got=%T", lType)}
	}

	return sliceType.Type, nil
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
	if expr.Alternative != nil {
		errors = append(errors, analyzeBlockStatement(expr.Alternative, returnType, env)...)
	}
	return errors
}

func analyzeCallExpression(expr *ast.CallExpression, env *object.Environment) (ast.DataType, []string) {
	dType, errors := AnalyzeExpression(expr.Function, env)
	if len(errors) != 0 {
		return nil, errors
	}

	switch fnType := dType.(type) {
	case *ast.BuiltinDataType:
		return parser.ANY, nil
	case *ast.FunctionDataType:
		if len(expr.Arguments) != len(fnType.Parameters) {
			errors = append(errors, fmt.Sprintf("analyzer error. Incorrect parameter count. expected %d, got=%d", len(fnType.Parameters), len(expr.Arguments)))
			return nil, errors
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

		return fnType.ReturnType, nil
	default:
		msg := fmt.Sprintf("Analyzer error. Unsupported call type %T", dType)
		return nil, []string{msg}
	}
}

func analyzeFunctionLiteral(expr *ast.FunctionLiteral, env *object.Environment) (ast.DataType, []string) {
	env = object.NewEnclosedEnvironment(env)
	fn := &ast.FunctionDataType{
		ReturnType: expr.ReturnType,
	}

	for _, ident := range expr.Parameters {
		env.Set(ident.Value, NativeTypeToDefaultObj(*ident.DataType))
		fn.Parameters = append(fn.Parameters, *ident.DataType)
	}

	errors := analyzeBlockStatement(expr.Body, expr.ReturnType, env)
	return fn, errors
}

func NativeTypeToDefaultObj(rawType ast.DataType) object.Object {
	switch rawType := rawType.(type) {
	case *ast.IntegerDataType:
		return &object.Integer{}
	case *ast.BooleanDataType:
		return &object.Boolean{}
	case *ast.StringDataType:
		return &object.String{}
	case *ast.SliceDataType:
		return &object.SliceObject{ValueType: rawType.Type}
	case *ast.AnyDataType:
		return &object.Any{}
	case *ast.ReferenceDataType:
		val := NativeTypeToDefaultObj(rawType.ValueType)
		return &object.ReferenceObject{Value: &val}
	case *ast.ChanDataType:
		return &object.ChanObject{
			Chan:     nil,
			ChanType: rawType.ValueType,
		}
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
		if leftType == parser.ANY || rightType == parser.ANY {
			return parser.BOOLEAN, nil
		}
		return analyzeEqInfixOperator(leftType, rightType)
	case "!=":
		if leftType == parser.ANY || rightType == parser.ANY {
			return parser.BOOLEAN, nil
		}
		return analyzeNeqInfixOperator(leftType, rightType)
	case "<":
		if leftType == parser.ANY || rightType == parser.ANY {
			return parser.BOOLEAN, nil
		}
		return analyzeLtInfixOperator(leftType, rightType)
	case ">":
		if leftType == parser.ANY || rightType == parser.ANY {
			return parser.BOOLEAN, nil
		}
		return analyzeGtInfixOperator(leftType, rightType)
	case "%":
		return analyzePercentInfixOperator(leftType, rightType)
	default:
		msg := fmt.Sprintf("analyzer error. unsupported infix operator type %s", expr.Operator)
		errors = append(errors, msg)
		return nil, errors
	}
}

func analyzePercentInfixOperator(leftType ast.DataType, rightType ast.DataType) (ast.DataType, []string) {
	switch {
	case leftType.Name() == parser.INT.Name() && rightType.Name() == parser.INT.Name():
		return parser.INT, nil
	default:
		msg := fmt.Sprintf("analyzer error. unsupported expression for 'percent' operator: %s and %s", leftType.Name(), rightType.Name())
		errors := []string{msg}
		return nil, errors
	}
}

func analyzeNeqInfixOperator(leftType, rightType ast.DataType) (ast.DataType, []string) {
	if leftType == rightType {
		return parser.BOOLEAN, nil
	} else {
		msg := fmt.Sprintf("analyzer error. unsupported comparison for '==' operator: %s and %s", leftType.Name(), rightType.Name())
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
		msg := fmt.Sprintf("analyzer error. unsupported expression type for '<' operator: %s and %s", leftType.Name(), rightType.Name())
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
	case "-f":
		return analyzeFoperPrefixExpression(rightType)
	case "*":
		return analyzeAsteriksPrefixExpression(rightType)
	case "&":
		return analyzeRefPrefixExpression(rightType)
	default:
		msg := fmt.Sprintf("analyzer error. unsupportet prefix operator type %s", expr.Operator)
		errors = append(errors, msg)
		return nil, errors
	}
}

func analyzeRefPrefixExpression(rightType ast.DataType) (ast.DataType, []string) {
	return &ast.ReferenceDataType{
		ValueType: rightType,
	}, nil
}

func analyzeAsteriksPrefixExpression(rightType ast.DataType) (ast.DataType, []string) {
	switch rightType := rightType.(type) {
	case *ast.ReferenceDataType:
		return rightType.ValueType, nil
	default:
		msg := fmt.Sprintf("analyzer error. unsupported expression type for '*' operator %s", rightType.Name())
		errors := []string{msg}
		return nil, errors
	}
}

func analyzeFoperPrefixExpression(rightType ast.DataType) (ast.DataType, []string) {
	if rightType == parser.STRING {
		return parser.BOOLEAN, nil
	} else {
		msg := fmt.Sprintf("analyzer error. unsupported expression type for '-f' operator %s", rightType.Name())
		errors := []string{msg}
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
