package model

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGenerateColumnGraph(t *testing.T) {
	type args struct {
		schema Schema
	}
	tests := []struct {
		name string
		args args
		want ColumnGraph
	}{
		{
			name: "generate correct column graph",
			args: args{
				schema: Schema{
					Tables: []Table{
						{
							Name: "customer",
							Columns: []Column{
								{
									Name: "id",
									Type: ColumnType{
										Base:  Int,
										Param: ColumnTypeParam(10),
									},
									Constraint: Constraint{
										TableName:  "product",
										ColumnName: "id",
									},
								},
								{
									Name: "name",
									Type: ColumnType{
										Base:  Varchar,
										Param: ColumnTypeParam(255),
									},
								},
							},
						},
						{
							Name: "product",
							Columns: []Column{
								{
									Name: "id",
									Type: ColumnType{
										Base:  Int,
										Param: ColumnTypeParam(14),
									},
								},
								{
									Name: "owner",
									Type: ColumnType{
										Base:  Varchar,
										Param: ColumnTypeParam(255),
									},
									Constraint: Constraint{
										TableName:  TableName("customer"),
										ColumnName: ColumnName("name"),
									},
								},
							},
						},
					},
				},
			},
			want: ColumnGraph{
				AdjacencyMatrix: AdjacencyMatrix{
					{0, 0, 1, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 1, 0, 0},
				},
				ColumnNodes: []ColumnNode{
					ColumnNode{
						column: Column{
							Name: "id",
							Type: ColumnType{
								Base:  Int,
								Param: ColumnTypeParam(10),
							},
							Constraint: Constraint{
								TableName:  "product",
								ColumnName: "id",
							},
						},
						isDone: false,
						index:  0,
					},
					ColumnNode{
						column: Column{
							Name: "name",
							Type: ColumnType{
								Base:  Varchar,
								Param: ColumnTypeParam(255),
							},
						},
						isDone: false,
						index:  1,
					},
					ColumnNode{
						column: Column{
							Name: "id",
							Type: ColumnType{
								Base:  Int,
								Param: ColumnTypeParam(14),
							},
						},
						isDone: false,
						index:  2,
					},
					ColumnNode{
						column: Column{
							Name: "owner",
							Type: ColumnType{
								Base:  Varchar,
								Param: ColumnTypeParam(255),
							},
							Constraint: Constraint{
								TableName:  TableName("customer"),
								ColumnName: ColumnName("name"),
							},
						},
						isDone: false,
						index:  3,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateColumnGraph(tt.args.schema)
			diff := cmp.Diff(got, tt.want, cmp.AllowUnexported(ColumnNode{}))
			if diff != "" {
				t.Errorf("GenerateColumnGraph(); -got, +want\n%v", diff)
			}
		})
	}
}
