package service

import (
	"strings"

	"github.com/canalun/sqloth/domain/model"
)

// TODO: change way to use only adjacency matrix...!
func GenerateColumnData(sortedColumnList []int, schema model.Schema, n int) map[string][]model.ColumnData {
	indexToColumnName := schema.GetMapFromIndexToColumnName()
	indexToColumn := schema.GetMapFromIndexToColumn()

	// have to generate adjacency matrix one more time here.
	// this is because adjacency matrix is all-zeroized in SortedColumnList() that is called in usecase before this function.
	adjacencyMatrix := GenerateAdjacencyMatrix(schema)

	data := map[string][]model.ColumnData{}
	l := reverseSlice(sortedColumnList)
	for _, columnIndex := range l {
		hasConstraint := 0
		var constraintColumnIndex int
		for j := 0; j < len(adjacencyMatrix); j++ {
			if adjacencyMatrix[columnIndex][j] != 0 {
				hasConstraint = 1
				constraintColumnIndex = j
			}
		}
		switch hasConstraint {
		case 0:
			cd := []model.ColumnData{}
			for i := 0; i < n; i++ {
				cd = append(cd, indexToColumn[columnIndex].GenerateRandomData())
			}
			data[indexToColumnName[columnIndex]] = cd
		case 1:
			data[indexToColumnName[columnIndex]] = data[indexToColumnName[constraintColumnIndex]]
		}
	}
	return data
}

func reverseSlice(l []int) []int {
	re := make([]int, len(l))
	for i, v := range l {
		re[len(l)-1-i] = v
	}
	return re
}

func GenerateQuery(schema model.Schema, data map[string][]model.ColumnData, n int) []string {
	queries := []string{}
	for _, table := range schema.Tables {
		//records look like [["'column1data'", "'column2data'",...],["'column1data'", "'column2data'",...]]
		records := [][]string{}
		for i := 0; i < n; i++ {
			records = append(records, []string{})
		}
		for _, column := range table.Columns {
			for i := 0; i < n; i++ {
				records[i] = append(records[i], string("'"+data[string(table.Name)+"."+string(column.Name)][i]+"'"))
			}
		}
		//data looks like ["('column1data', 'column2data'...)", "('column1data', 'column2data'...)"...]
		data := []string{}
		for _, record := range records {
			re := "(" + strings.Join(record, ",") + ")"
			data = append(data, re)
		}
		query := "INSERT INTO " + string(table.Name) + " VALUES " + strings.Join(data, ",") + ";"
		queries = append(queries, query)
	}
	return queries
}
