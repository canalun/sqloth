package model

type Schema struct {
	Tables []Table
}

func (s *Schema) AddTable(table Table) {
	s.Tables = append(s.Tables, table)
}

func (s *Schema) LastTable() *Table {
	return &s.Tables[len(s.Tables)-1]
}

func (s Schema) GetMapFromColumnToIndex() map[string]int {
	columnToIndex := map[string]int{}
	i := 0
	for _, table := range s.Tables {
		for _, column := range table.Columns {
			columnToIndex[string(table.Name)+"."+string(column.Name)] = i
			i += 1
		}
	}
	return columnToIndex
}

func (s Schema) GetMapFromIndexToColumnName() map[int]string {
	indexToColumn := map[int]string{}
	i := 0
	for _, table := range s.Tables {
		for _, column := range table.Columns {
			indexToColumn[i] = string(table.Name) + "." + string(column.Name)
			i += 1
		}
	}
	return indexToColumn
}

func (s *Schema) GetMapFromIndexToColumn() map[int]Column {
	indexToColumn := map[int]Column{}
	i := 0
	for _, table := range s.Tables {
		for _, column := range table.Columns {
			indexToColumn[i] = column
			i += 1
		}
	}
	return indexToColumn
}
