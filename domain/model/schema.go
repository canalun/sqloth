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
