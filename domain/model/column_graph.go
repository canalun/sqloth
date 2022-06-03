package model

type ColumnGraph struct {
	AdjacencyMatrix AdjacencyMatrix
	ColumnNodes     []ColumnNode
}

type ColumnNode struct {
	column Column
	isDone bool
	index  int
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
			if column.HasConstraint() {
				if i, ok := columnToIndex[string(table.Name)+"."+string(column.Name)]; ok {
					if j, ok := columnToIndex[string(column.Constraint.TableName)+"."+string(column.Constraint.ColumnName)]; ok {
						am[i][j] = 1
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
