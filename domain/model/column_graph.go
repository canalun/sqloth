package model

import (
	"github.com/pkg/errors"
)

// Graph structure of columns expressed by adjacency matrix and nodes
type ColumnGraph struct {
	AdjacencyMatrix AdjacencyMatrix
	ColumnNodes     []ColumnNode
}

type AdjacencyMatrix [][]int

func newAdjacencyMatrix(n int) AdjacencyMatrix {
	am := make(AdjacencyMatrix, 0, n)
	i := 0
	for i < n {
		am = append(am, make([]int, n))
		i++
	}
	return am
}

type ColumnNode struct {
	column Column
	isDone bool
	index  int
}

func (cn *ColumnNode) Done() {
	cn.isDone = true
}

func (cn ColumnNode) IsDone() bool {
	return cn.isDone
}

func (cn ColumnNode) GetColumn() Column {
	return cn.column
}

func GenerateColumnGraph(schema Schema) ColumnGraph {
	columnNodes := []ColumnNode{}
	columnToIndex := map[string]int{}
	i := 0
	for _, table := range schema.Tables {
		for _, column := range table.Columns {
			columnToIndex[string(table.Name)+"."+string(column.Name)] = i
			columnNodes = append(columnNodes, ColumnNode{
				column: column,
				isDone: false,
				index:  i,
			})
			i += 1
		}
	}

	am := newAdjacencyMatrix(len(columnToIndex))
	for _, table := range schema.Tables {
		for _, column := range table.Columns {
			if column.HasForeignKey() {
				if i, ok := columnToIndex[string(table.Name)+"."+string(column.Name)]; ok {
					for _, foreignKey := range column.ForeignKeys {
						if j, ok := columnToIndex[string(foreignKey.TableName)+"."+string(foreignKey.ColumnName)]; ok {
							am[i][j] = 1
						}
					}
				}
			}
		}
	}

	return ColumnGraph{
		AdjacencyMatrix: am,
		ColumnNodes:     columnNodes,
	}
}

func (cg ColumnGraph) isAllDone() bool {
	for _, cn := range cg.ColumnNodes {
		if !cn.isDone {
			return false
		}
	}
	return true
}

//TODO: adopt error handling such as Stacktrace
func (cg ColumnGraph) HasParentNodes(i int) (bool, error) {
	if i >= len(cg.AdjacencyMatrix) {
		return false, errors.New("invalid index")
	}
	for _, v := range cg.AdjacencyMatrix[i] {
		if v == 1 {
			return true, nil
		}
	}
	return false, nil
}

func (cg ColumnGraph) IsParentNodesAreAllDone(i int) (bool, error) {
	if i >= len(cg.AdjacencyMatrix) {
		return false, errors.New("invalid index")
	}
	for parentIndex, v := range cg.AdjacencyMatrix[i] {
		if v == 1 {
			if !cg.ColumnNodes[parentIndex].IsDone() {
				return false, nil
			}
		}
	}
	return true, nil
}

func (cg ColumnGraph) ParentNodeIndexes(i int) ([]int, error) {
	if i >= len(cg.AdjacencyMatrix) {
		return []int{}, errors.New("invalid index")
	}
	parentNodeIndexes := []int{}
	for i, v := range cg.AdjacencyMatrix[i] {
		if v == 1 {
			parentNodeIndexes = append(parentNodeIndexes, i)
		}
	}
	return parentNodeIndexes, nil
}

func (cg ColumnGraph) HasChildrenNodes(i int) (bool, error) {
	if i >= len(cg.AdjacencyMatrix) {
		return false, errors.New("invalid index")
	}
	for _, r := range cg.AdjacencyMatrix {
		if r[i] == 1 {
			return true, nil
		}
	}
	return false, nil
}

func (cg ColumnGraph) ChildrenNodeIndexes(i int) ([]int, error) {
	if i >= len(cg.AdjacencyMatrix) {
		return []int{}, errors.New("invalid index")
	}
	childrenNodeIndexes := []int{}
	for ri, r := range cg.AdjacencyMatrix {
		if r[i] == 1 {
			childrenNodeIndexes = append(childrenNodeIndexes, ri)
		}
	}
	return childrenNodeIndexes, nil
}
