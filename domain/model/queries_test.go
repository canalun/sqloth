package model

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGenerateQuery(t *testing.T) {
	type args struct {
		rft map[TableName][]Record
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
					"table2": []Record{{"table2-v1", "table2-v2", "table2-v3"}, {"table2-v4", "table2-v5", "table2-v6"}},
				},
			},
			want: []string{
				"SET FOREIGN KEY = 0;",
				"INSERT INTO table1 VALUES ('table1-v1','table1-v2','table1-v3'),('table1-v4','table1-v5','table1-v6');",
				"INSERT INTO table2 VALUES ('table2-v1','table2-v2','table2-v3'),('table2-v4','table2-v5','table2-v6');",
				"SET FOREIGN KEY = 1;",
			},
		},
		{
			name: "return empty string from an empty map",
			args: args{
				rft: map[TableName][]Record{},
			},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateQuery(tt.args.rft)
			diff := cmp.Diff(got, tt.want)
			if diff != "" {
				t.Errorf("GenerateQuery(); -got, +want\n%v", diff)
			}
		})
	}
}
