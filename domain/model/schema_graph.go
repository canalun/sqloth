package model

import (
	"github.com/pkg/errors"
)

// Sqloth grasps schema structure as a graph of columns.
// The graph is expressed as nodes and the adjacency matrix.
type SchemaGraph struct {
	AdjacencyMatrix AdjacencyMatrix
	ColumnNodes     []ColumnNode
}

type AdjacencyMatrix [][]int

func NewAdjacencyMatrix(n int) AdjacencyMatrix {
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

func NewColumnNode(c Column, isDone bool, i int) ColumnNode {
	return ColumnNode{
		column: c,
		isDone: isDone,
		index:  i,
	}
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

func (cg SchemaGraph) IsAllDone() bool {
	for _, cn := range cg.ColumnNodes {
		if !cn.isDone {
			return false
		}
	}
	return true
}

//TODO: adopt error handling such as Stacktrace
func (cg SchemaGraph) HasParentNodes(i int) (bool, error) {
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

func (cg SchemaGraph) IsParentNodesAreAllDone(i int) (bool, error) {
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

func (cg SchemaGraph) ParentNodeIndexes(i int) ([]int, error) {
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

func (cg SchemaGraph) HasChildrenNodes(i int) (bool, error) {
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

func (cg SchemaGraph) ChildrenNodeIndexes(i int) ([]int, error) {
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
