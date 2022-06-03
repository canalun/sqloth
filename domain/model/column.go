package model

import (
	"errors"
)

// TODO: reorder and reorganize source code

type ColumnTypeBase string
type ColumnTypeParam int

type ColumnType struct {
	Base  ColumnTypeBase
	Param ColumnTypeParam
}

const (
	Varchar    ColumnTypeBase = "varchar"
	Varbinary  ColumnTypeBase = "varbinary"
	Mediumblob ColumnTypeBase = "mediumblob"
	Text       ColumnTypeBase = "text"
	Int        ColumnTypeBase = "int"
	Tinyint    ColumnTypeBase = "tinyint"
	Timestamp  ColumnTypeBase = "timestamp"
	Datetime   ColumnTypeBase = "datetime"
	Json       ColumnTypeBase = "json"
)

func StrToColumnTypeBase(str string) (ColumnTypeBase, error) {
	switch str {
	case string(Varchar):
		return Varchar, nil
	case string(Varbinary):
		return Varbinary, nil
	case string(Mediumblob):
		return Mediumblob, nil
	case string(Text):
		return Text, nil
	case string(Int):
		return Int, nil
	case string(Tinyint):
		return Tinyint, nil
	case string(Timestamp):
		return Timestamp, nil
	case string(Datetime):
		return Datetime, nil
	case string(Json):
		return Json, nil
	default:
		return "", errors.New("unregistered type")
	}
}

// TODO: rethink if these types are needed
type ColumnData string
type ColumnName string
type Constraint struct {
	TableName
	ColumnName
}

func NewConstraint(tn TableName, cn ColumnName) Constraint {
	return Constraint{
		TableName:  tn,
		ColumnName: cn,
	}
}

type Column struct {
	Name         ColumnName
	Type         ColumnType
	Constraint   Constraint
	DataToInsert []ColumnData
}

func NewColumn(name ColumnName, ct ColumnType) Column {
	return Column{
		Name: ColumnName(name),
		Type: ct,
	}
}

func (c Column) HasConstraint() bool {
	return string(c.Constraint.TableName)+string(c.Constraint.ColumnName) != ""
}

func (c *Column) SetConstraint(constraint Constraint) {
	c.Constraint = constraint
}

func (c Column) GenerateRandomData() ColumnData {
	var data string
	switch c.Type.Base {
	case Varchar, Text, Varbinary, Mediumblob:
		data = generateRandomString(int(c.Type.Param))
	case Int:
		data = generateRandomInt(int(c.Type.Param))
	case Tinyint:
		data = generateRandomTinyint()
	case Timestamp, Datetime:
		data = generateRandomDate()
	case Json:
		data = generateRandomJson()
	}
	return ColumnData(data)
}
