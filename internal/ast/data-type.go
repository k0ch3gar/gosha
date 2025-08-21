package ast

import (
	"bytes"
	"strings"

	"kstmc.com/gosha/internal/token"
)

type DataType interface {
	Name() string
}

type DataTypeExpression struct {
	Token token.Token
	Type  DataType
}

func (dte *DataTypeExpression) expressionNode() {

}

func (dte *DataTypeExpression) String() string {
	return dte.Type.Name()
}

func (dte *DataTypeExpression) TokenLiteral() string {
	return dte.Token.Literal
}

type IntegerDataType struct {
}

func (idt *IntegerDataType) Name() string {
	return "int"
}

type StringDataType struct {
}

func (sdt *StringDataType) Name() string {
	return "string"
}

type BooleanDataType struct {
}

func (bdt *BooleanDataType) Name() string {
	return "bool"
}

type NilDataType struct {
}

func (ndt *NilDataType) Name() string {
	return "nil"
}

type AnyDataType struct {
}

func (adt *AnyDataType) Name() string {
	return "any"
}

type FunctionDataType struct {
	Parameters []DataType
	ReturnType DataType
}

func (fdt *FunctionDataType) Name() string {
	var out bytes.Buffer

	out.WriteString("func(")
	var paramsTemp []string
	for _, param := range fdt.Parameters {
		paramsTemp = append(paramsTemp, param.Name())
	}

	out.WriteString(strings.Join(paramsTemp, ", "))
	out.WriteString(") ")
	out.WriteString(fdt.ReturnType.Name())

	return out.String()
}

type ReturnDataType struct {
}

func (rdt *ReturnDataType) Name() string {
	return "return"
}

type ErrorDataType struct {
}

func (edt *ErrorDataType) Name() string {
	return "error"
}

type BuiltinDataType struct {
	Parameters []DataType
	ReturnType DataType
}

func (bdt *BuiltinDataType) Name() string {
	return "builtin"
}

type SliceDataType struct {
	Type DataType
}

func (sdt *SliceDataType) Name() string {
	return "[]" + sdt.Type.Name()
}

type ReferenceDataType struct {
	ValueType DataType
}

func (pdt *ReferenceDataType) Name() string {
	return "*" + pdt.ValueType.Name()
}

type PointerDataType struct {
	ValueType DataType
}

func (pdt *PointerDataType) Name() string {
	return "*" + pdt.ValueType.Name()
}

type ChanDataType struct {
	ValueType DataType
}

func (cdt *ChanDataType) Name() string {
	return "chan " + cdt.ValueType.Name()
}
