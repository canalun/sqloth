package model

import (
	"errors"
	"strconv"
	"strings"
)

// Sqloth use the below models to grasp the given schema.
//	- Schema
// 	- Table
// 	- Column
//	- ForeignKey
// Schema has Table model(s).
// Each Table has its Column(s).
// Then, each Column has its constraint(s) as ForeignKey(s)

// In sqloth, it is Column that is responsible for data generation,
// because the format of random data is determined by column data type,
// and it seems to be a good choice to seen data generation as a behavior of Column.

////// Schema //////////////////////////

type Schema struct {
	Tables []Table
}

func (s *Schema) AddTable(table Table) {
	s.Tables = append(s.Tables, table)
}

func (s *Schema) LastTable() *Table {
	return &s.Tables[len(s.Tables)-1]
}

////// Table ///////////////////////////

type Table struct {
	Name    TableName
	Columns []Column
}
type TableName string

func NewTable(name TableName) Table {
	return Table{
		Name:    TableName(name),
		Columns: []Column{},
	}
}

func (t *Table) AddColumns(column Column) {
	t.Columns = append(t.Columns, column)
}

////// Column //////////////////////////

type Column struct {
	Name          ColumnName
	FullName      ColumnFullName
	Type          ColumnType
	ForeignKeys   []ForeignKey
	AutoIncrement bool
	Unsigned      bool
}
type ColumnName string
type ColumnFullName string

func NewColumn(fullName ColumnFullName, ct ColumnType) Column {
	name := strings.Split(string(fullName), ".")[1]
	return Column{
		Name:     ColumnName(name),
		FullName: fullName,
		Type:     ct,
	}
}

func NewColumnFullName(tn TableName, cn ColumnName) ColumnFullName {
	return ColumnFullName(string(tn) + "." + string(cn))
}

// e.g. INT(10) converts to ColumnType{Base: Int, Param: ColumnTypeParam(10)}
type ColumnType struct {
	Base  ColumnTypeBase
	Param ColumnTypeParam
}
type ColumnTypeBase string
type ColumnTypeParam int

const (
	Varchar    ColumnTypeBase = "varchar"
	Varbinary  ColumnTypeBase = "varbinary"
	Mediumblob ColumnTypeBase = "mediumblob"
	Text       ColumnTypeBase = "text"
	Tinyint    ColumnTypeBase = "tinyint"
	Smallint   ColumnTypeBase = "smallint"
	Mediumint  ColumnTypeBase = "mediumint"
	Int        ColumnTypeBase = "int"
	Bigint     ColumnTypeBase = "bigint"
	Timestamp  ColumnTypeBase = "timestamp"
	Datetime   ColumnTypeBase = "datetime"
	Json       ColumnTypeBase = "json"
)

func NewColumnType(str string) (ct ColumnType, err error) {
	f := func(r rune) bool {
		return r == '(' || r == ')' || r == ','
	}
	l := strings.FieldsFunc(str, f)

	base, err := newColumnTypeBase(strings.ToLower(l[0]))
	if err != nil {
		return ColumnType{}, err
	}

	var param int
	if len(l) > 1 {
		param, err = strconv.Atoi(l[1])
		if err != nil {
			return ColumnType{}, err
		}
	}
	if base == Text {
		param = 100
	}
	// TODO: handle varbinary appropriately
	if base == Varbinary {
		param = 100
	}
	// TODO: handle mediumblob appropriately
	if base == Mediumblob {
		param = 100
	}

	return ColumnType{
		Base:  base,
		Param: ColumnTypeParam(param),
	}, nil
}

func newColumnTypeBase(str string) (ColumnTypeBase, error) {
	switch str {
	case string(Varchar):
		return Varchar, nil
	case string(Varbinary):
		return Varbinary, nil
	case string(Mediumblob):
		return Mediumblob, nil
	case string(Text):
		return Text, nil
	case string(Smallint), string(Int), string(Mediumint), string(Bigint):
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

func (c Column) HasForeignKey() bool {
	return len(c.ForeignKeys) > 0
}

func (c *Column) SetForeignKey(foreignKey ForeignKey) {
	c.ForeignKeys = append(c.ForeignKeys, foreignKey)
}

func (c *Column) SetAutoIncrement() {
	c.AutoIncrement = true
}

func (c *Column) SetUnsigned(b bool) {
	c.Unsigned = b
}

////// Foreign Key //////////////////////////

type ForeignKey struct {
	TableName
	ColumnName
}

func NewForeignKey(tn TableName, cn ColumnName) ForeignKey {
	return ForeignKey{
		TableName:  tn,
		ColumnName: cn,
	}
}

////// Data Generation //////////////////////////

func (c Column) GenerateData(n int) []Value {
	d := []Value{}
	switch c.AutoIncrement {
	case true:
		for i := 0; i < n; i++ {
			d = append(d, Value("NULL"))
		}
	default:
		for i := 0; i < n; i++ {
			d = append(d, Value(c.GenerateRandomData()))
		}
	}
	return d
}

func (c Column) GenerateRandomData() string {
	var data string
	switch c.Type.Base {
	case Varchar, Text, Varbinary, Mediumblob:
		data = generateRandomString(int(c.Type.Param))
	case Int:
		data = generateRandomInt(c.Type.Base, c.Unsigned)
	case Tinyint:
		data = generateRandomTinyint()
	case Timestamp, Datetime:
		data = generateRandomDate()
	case Json:
		data = generateRandomJson()
	}
	return data
}
