package model

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGenerateRecordsForTables(t *testing.T) {
	type args struct {
		vfc    map[ColumnFullName][]Value
		schema Schema
		n      int
	}
	tests := []struct {
		name string
		args args
		want map[TableName][]Record
	}{
		{
			name: "generate records for tables from map of column full-names and values",
			args: args{
				vfc: map[ColumnFullName][]Value{
					"table1.test1": []Value{"1", "2", "3"},
					"table1.test2": []Value{"aaa", "bbb", "ccc"},
					"table2.test1": []Value{"v1", "v2", "v3"},
				},
				schema: Schema{
					Tables: []Table{
						{Name: "table1", Columns: []Column{{FullName: "table1.test1"}, {FullName: "table1.test2"}}},
						{Name: "table2", Columns: []Column{{FullName: "table2.test1"}}},
					},
				},
				n: 3,
			},
			want: map[TableName][]Record{
				"table1": []Record{{"1", "aaa"}, {"2", "bbb"}, {"3", "ccc"}},
				"table2": []Record{{"v1"}, {"v2"}, {"v3"}},
			},
		},
		{
			name: "skip auto increment columns",
			args: args{
				vfc: map[ColumnFullName][]Value{
					"table1.test1": []Value{"1", "2", "3"},
					"table1.test2": []Value{"aaa", "bbb", "ccc"},
					"table2.test1": []Value{"v1", "v2", "v3"},
				},
				schema: Schema{
					Tables: []Table{
						{Name: "table1", Columns: []Column{{FullName: "table1.test1", AutoIncrement: true}, {FullName: "table1.test2"}}},
						{Name: "table2", Columns: []Column{{FullName: "table2.test1"}}},
					},
				},
				n: 3,
			},
			want: map[TableName][]Record{
				"table1": []Record{{"aaa"}, {"bbb"}, {"ccc"}},
				"table2": []Record{{"v1"}, {"v2"}, {"v3"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateRecordsForTables(tt.args.vfc, tt.args.schema, tt.args.n)
			diff := cmp.Diff(got, tt.want)
			if diff != "" {
				t.Errorf("GenerateRecordsForTables(); -got, +want\n%v", diff)
			}
		})
	}
}
