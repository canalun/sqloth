package dservice

import "github.com/canalun/sqloth/domain/model"

func GenerateSchemaGraph(schema model.Schema) model.SchemaGraph {
	columnNodes := []model.ColumnNode{}
	columnToIndex := map[string]int{}
	i := 0
	for _, table := range schema.Tables {
		for _, column := range table.Columns {
			columnToIndex[string(table.Name)+"."+string(column.Name)] = i
			columnNodes = append(columnNodes, model.NewColumnNode(column, false, i))
			i += 1
		}
	}

	am := model.NewAdjacencyMatrix(len(columnToIndex))
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

	return model.SchemaGraph{
		AdjacencyMatrix: am,
		ColumnNodes:     columnNodes,
	}
}

//TODO: better to be defined as a method of map[ColumnFullName][]Value?
//TODO: shuffle values. currently, values with foreignKeys are just simple sum of strings in the order.
func GenerateValuesForColumns(cg model.SchemaGraph, n int) map[model.ColumnFullName][]model.Value {
	dict := map[model.ColumnFullName][]model.Value{}
	for i := range cg.ColumnNodes {
		if !cg.IsAllDone() {
			generateValuesForColumnsByRecursion(&cg, i, n, dict)
		}
	}
	return dict
}

//TODO: better to be defined as a method with side-effect of map[ColumnFullName][]Value?
func generateValuesForColumnsByRecursion(cg *model.SchemaGraph, i, n int, dict map[model.ColumnFullName][]model.Value) {
	if cg.ColumnNodes[i].IsDone() {
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
			values := []model.Value{}
			for j := 0; j < n; j++ {
				d := ""
				parentNodeIndexes, _ := cg.ParentNodeIndexes(i)
				for _, parentNodeIndex := range parentNodeIndexes {
					fn := cg.ColumnNodes[parentNodeIndex].GetColumn().FullName
					d += string(dict[fn][j])
				}
				values = append(values, model.Value(d))
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

func GenerateRecordsForTables(vfc map[model.ColumnFullName][]model.Value, schema model.Schema, n int) map[model.TableName][]model.Record {
	rft := map[model.TableName][]model.Record{}
	for _, table := range schema.Tables {
		records := []model.Record{}
		for i := 0; i < n; i++ {
			var record model.Record
			for _, column := range table.Columns {
				//skip auto increment column
				if !column.AutoIncrement {
					record = append(record, vfc[column.FullName][i])
				}
			}
			records = append(records, record)
		}
		rft[table.Name] = records
	}
	return rft
}
