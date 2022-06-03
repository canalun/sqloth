package model

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTable_GenerateQuery(t *testing.T) {
	type fields struct {
		Name    TableName
		Columns []Column
	}
	type args struct {
		numOfData int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "generate insert queries of dummy data for tables ",
			fields: fields{
				Name: "table-1",
				Columns: []Column{
					{
						Type: ColumnType{
							Base:  Varchar,
							Param: 255,
						},
					},
					{
						Type: ColumnType{
							Base:  Int,
							Param: 11,
						},
					},
				},
			},
			args: args{
				numOfData: 3,
			},
			// TODO: mod want data
			want: `INSERT INTO table-1 VALUES ('varchar','int'),('varchar','int'),('varchar','int');`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := Table{
				Name:    tt.fields.Name,
				Columns: tt.fields.Columns,
			}
			got := tr.GenerateQuery(tt.args.numOfData)
			diff := cmp.Diff(got, tt.want)
			if diff != "" {
				t.Error("-:got, +:want", diff)
			}
		})
	}
}
