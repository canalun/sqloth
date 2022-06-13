package model

import (
	"strings"
)

type TableName string

type Table struct {
	Name    TableName
	Columns []Column
}

func NewTable(name TableName, columns []Column) Table {
	return Table{
		Name:    TableName(name),
		Columns: columns,
	}
}

func (t *Table) AddColumns(column Column) {
	t.Columns = append(t.Columns, column)
}

// deprecated
func (t Table) GenerateQuery(numOfData int) string {
	records := []string{}
	for i := 0; i < numOfData; i++ {

		// data look like {"'aaa'", "'aaaaaa'"...}
		data := []string{}
		for _, column := range t.Columns {
			randomData := column.GenerateRandomData()
			data = append(data, "'"+string(randomData)+"'")
		}

		// records look like {"('aaa','aaaaaa', ...)", ...}
		records = append(records, "("+strings.Join(data, ",")+")")
	}

	query := "INSERT INTO " + string(t.Name) + " VALUES " + strings.Join(records, ",") + ";"

	return query
}

func GenerateRecordsForTables(vfc map[ColumnFullName][]Value, schema Schema, n int) map[TableName][]Record {
	rft := map[TableName][]Record{}
	for _, table := range schema.Tables {
		records := []Record{}
		for i := 0; i < n; i++ {
			var record Record
			for _, column := range table.Columns {
				record = append(record, vfc[column.FullName][i])
			}
			records = append(records, record)
		}
		rft[table.Name] = records
	}
	return rft
}
