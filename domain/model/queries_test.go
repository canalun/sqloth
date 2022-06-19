package model

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestGenerateQuery(t *testing.T) {
	type args struct {
		rft    map[TableName][]Record
		schema Schema
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "generate correct query from a map from table names to records",
			args: args{
				rft: map[TableName][]Record{
					"table1": []Record{{"table1-v1", "table1-v2", "table1-v3"}, {"table1-v4", "table1-v5", "table1-v6"}},
					"table2": []Record{{"table2-v1", "table2-v2"}, {"table2-v4", "table2-v5"}},
				},
				schema: Schema{
					Tables: []Table{
						{
							Name: "table1",
							Columns: []Column{
								{Name: "table1-column1"},
								{Name: "table1-column2"},
								{Name: "table1-column3"},
							},
						},
						{
							Name: "table2",
							Columns: []Column{
								{Name: "table2-column1"},
								{Name: "table2-column2", AutoIncrement: true},
								{Name: "table2-column3"},
							},
						},
					},
				},
			},
			want: []string{
				"SET foreign_key_checks = 0;",
				"INSERT INTO table1(`table1-column1`, `table1-column2`, `table1-column3`) VALUES ('table1-v1','table1-v2','table1-v3'),('table1-v4','table1-v5','table1-v6');",
				"INSERT INTO table2(`table2-column1`, `table2-column3`) VALUES ('table2-v1','table2-v2'),('table2-v4','table2-v5');",
				"SET foreign_key_checks = 1;",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateQuery(tt.args.rft, tt.args.schema)
			diff1 := cmp.Diff(got[0], tt.want[0])
			diff2 := cmp.Diff(got[len(got)-1], tt.want[len(tt.want)-1])
			diff3 := cmp.Diff(got[1:len(got)-1], tt.want[1:len(tt.want)-1],
				cmpopts.SortSlices(func(i, j string) bool {
					return i < j
				}),
			)
			if diff1+diff2+diff3 != "" {
				t.Errorf("GenerateQuery(); -got, +want\n%v\n\n%v\n\n%v", diff1, diff2, diff3)
			}
		})
	}
}

func Test_listColumnsForQuery(t *testing.T) {
	type args struct {
		table Table
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "return column list excluding auto increment column",
			args: args{
				table: Table{
					Columns: []Column{
						{Name: "test1"},
						{Name: "test2", AutoIncrement: true},
						{Name: "test3"},
					},
				},
			},
			want: []string{"`test1`", "`test3`"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := listColumnsForQuery(tt.args.table)
			diff := cmp.Diff(got, tt.want)
			if diff != "" {
				t.Errorf("listColumnsForQuery(); -got, +want %v", diff)
			}
		})
	}
}
