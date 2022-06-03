package service

import (
	"github.com/canalun/sqloth/domain/model"
)

// to handle constraints among columns, sqloth has to sort the columns.
// sqloth abstracts columns and their relationship as directed graph, and it uses topological sort.

// GenerateSortedColumnList returns one of possible topological-sorted column lists by calculating Adjacency Matrix.
// The returned list begins with root node and ends with the outest leaf node. Please check the unit test for details.

// TODO: fix to handle multi-index-key
// TODO: consider how to handle cyclic constraint
func GenerateSortedColumnList(schema model.Schema, am AdjacencyMatrix) []int {
	n := len(am)

	numsTodo := make([]int, len(am))
	for i := 0; i < len(numsTodo); i++ {
		numsTodo[i] = 0
	}
	sortedColumnList := []int{}

	for sumOfIntList(numsTodo) < n {
		for j, Done := range numsTodo {
			if Done != 1 {
				isRoot := 1
				for k := 0; k < len(am); k++ {
					if am[k][j] != 0 {
						isRoot = 0
					}
				}
				if isRoot == 1 {
					numsTodo[j] = 1
					sortedColumnList = append(sortedColumnList, j)
					//update adjacency matrix
					am.zeroizeColumn(j)
				}
			}
		}
	}

	return sortedColumnList
}

func sumOfIntList(l []int) int {
	re := 0
	for _, i := range l {
		re += i
	}
	return re
}

// Adjacency Matrix is crucial for the algorithm, but it's just intermediate product (so it's not included in domain).
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

func GenerateAdjacencyMatrix(schema model.Schema) AdjacencyMatrix {
	columnToIndex := schema.GetMapFromColumnToIndex()

	am := newAdjacencyMatrix(len(columnToIndex))
	for _, table := range schema.Tables {
		for _, column := range table.Columns {
			if column.HasConstraint() {
				if i, ok := columnToIndex[string(table.Name)+"."+string(column.Name)]; ok {
					if j, ok := columnToIndex[string(column.Constraint.TableName)+"."+string(column.Constraint.ColumnName)]; ok {
						am[i][j] = 1
					}
				}
			}
		}
	}
	return am
}

func (am *AdjacencyMatrix) zeroizeColumn(i int) {
	n := len(*am)
	for k := 0; k < n; k++ {
		(*am)[i][k] = 0
	}
}
