package model

type TableName string

type Table struct {
	Name    TableName
	Columns []Column
}

func NewTable(name TableName) Table {
	return Table{
		Name:    TableName(name),
		Columns: []Column{},
	}
}

func (t *Table) AddColumns(column Column) {
	t.Columns = append(t.Columns, column)
}

func GenerateRecordsForTables(vfc map[ColumnFullName][]Value, schema Schema, n int) map[TableName][]Record {
	rft := map[TableName][]Record{}
	for _, table := range schema.Tables {
		records := []Record{}
		for i := 0; i < n; i++ {
			var record Record
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
