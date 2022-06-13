package model

import (
	"errors"
	"strconv"
	"strings"
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
type ColumnFullName string
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
	Name          ColumnName
	FullName      ColumnFullName
	Type          ColumnType
	AutoIncrement bool
	Constraints   []Constraint
}

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

func (c Column) HasConstraint() bool {
	return len(c.Constraints) > 0
}

func (c *Column) SetConstraint(constraint Constraint) {
	c.Constraints = append(c.Constraints, constraint)
}

func (c Column) GenerateData(n int) []Value {
	d := []Value{}
	switch c.AutoIncrement {
	case true:
		for i := 0; i < n; i++ {
			d = append(d, Value(strconv.Itoa(i)))
		}
	default:
		for i := 0; i < n; i++ {
			d = append(d, Value(c.GenerateRandomData()))
		}
	}
	return d
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

//TODO: better to be defined as a method of map[ColumnFullName][]Value?
//TODO: shuffle values. currently, values with constraints are just simple sum of strings in the order.
func GenerateValuesForColumns(cg ColumnGraph, n int) map[ColumnFullName][]Value {
	dict := map[ColumnFullName][]Value{}
	for i := range cg.ColumnNodes {
		if !cg.isAllDone() {
			generateValuesForColumnsByRecursion(&cg, i, n, dict)
		}
	}
	return dict
}

//TODO: better to be defined as a method with side-effect of map[ColumnFullName][]Value?
func generateValuesForColumnsByRecursion(cg *ColumnGraph, i, n int, dict map[ColumnFullName][]Value) {
	if cg.ColumnNodes[i].isDone {
		return
	}

	//TODO: error handling
	hasParentNodes, _ := cg.HasParentNodes(i)
	switch hasParentNodes {
	case false:
		c := cg.ColumnNodes[i].GetColumn()
		d := c.GenerateData(n)
		dict[c.FullName] = d
		cg.ColumnNodes[i].Done()
		if hasChildrenNodes, _ := cg.HasChildrenNodes(i); hasChildrenNodes {
			childrenNodesIndexes, _ := cg.ChildrenNodeIndexes(i)
			for _, childrenNodeIndex := range childrenNodesIndexes {
				if allDone, _ := cg.IsParentNodesAreAllDone(childrenNodeIndex); allDone {
					generateValuesForColumnsByRecursion(cg, childrenNodeIndex, n, dict)
				}
			}
		}
	default:
		allDone, _ := cg.IsParentNodesAreAllDone(i)
		switch allDone {
		case true:
			//TODO: distinct values and data
			values := []Value{}
			for j := 0; j < n; j++ {
				d := ""
				parentNodeIndexes, _ := cg.ParentNodeIndexes(i)
				for _, parentNodeIndex := range parentNodeIndexes {
					fn := cg.ColumnNodes[parentNodeIndex].GetColumn().FullName
					d += string(dict[fn][j])
				}
				values = append(values, Value(d))
			}
			dict[cg.ColumnNodes[i].GetColumn().FullName] = values
			cg.ColumnNodes[i].Done()
			if hasChildrenNodes, _ := cg.HasChildrenNodes(i); hasChildrenNodes {
				childrenNodesIndexes, _ := cg.ChildrenNodeIndexes(i)
				for _, childrenNodeIndex := range childrenNodesIndexes {
					if allDone, _ := cg.IsParentNodesAreAllDone(childrenNodeIndex); allDone {
						generateValuesForColumnsByRecursion(cg, childrenNodeIndex, n, dict)
					}
				}
			}
		default:
			parentNodeIndexes, _ := cg.ParentNodeIndexes(i)
			for _, parentIndex := range parentNodeIndexes {
				if !cg.ColumnNodes[parentIndex].IsDone() {
					generateValuesForColumnsByRecursion(cg, parentIndex, n, dict)
				}
			}
		}
	}
}
