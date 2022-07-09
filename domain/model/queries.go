package model

import "strings"

// Value is literally a value of a column.
type Value string

// Record is a set of Value, that is "(XX, XX, ....)" in SQL query.
type Record []Value

// Query is string as an actually executable SQL query.
type Query string

func GenerateQuery(rft map[TableName][]Record, schema Schema) []string {
	if len(rft) == 0 {
		return []string{}
	}
	re := []string{"SET foreign_key_checks = 0;"}
	for _, table := range schema.Tables {
		q := "INSERT INTO " + string(table.Name) + "(" + strings.Join(listColumnsForQuery(table), ", ") + ")" + " VALUES "
		for _, record := range rft[table.Name] {
			q += querizeRecord(record)
		}
		q = q[:len(q)-1] + ";"
		re = append(re, q)
	}
	re = append(re, "SET foreign_key_checks = 1;")
	return re
}

func querizeRecord(record Record) string {
	re := "("
	for _, v := range record {
		re += "'" + string(v) + "',"
	}
	re = re[:len(re)-1] + "),"
	return re
}

func listColumnsForQuery(table Table) []string {
	re := []string{}
	for _, c := range table.Columns {
		if !c.AutoIncrement {
			re = append(re, "`"+string(c.Name)+"`")
		}
	}
	return re
}
