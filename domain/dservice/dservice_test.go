package dservice

import (
	"strconv"
	"testing"

	"github.com/canalun/sqloth/domain/model"
	"github.com/google/go-cmp/cmp"
)

func TestGenerateSchemaGraph(t *testing.T) {
	type args struct {
		schema model.Schema
	}
	tests := []struct {
		name string
		args args
		want model.SchemaGraph
	}{
		{
			name: "generate correct column graph",
			args: args{
				schema: model.Schema{
					Tables: []model.Table{
						{
							Name: "customer",
							Columns: []model.Column{
								{
									Name:     "id",
									FullName: "customer.id",
									Type: model.ColumnType{
										Base:  model.Int,
										Param: model.ColumnTypeParam(10),
									},
									ForeignKeys: []model.ForeignKey{
										{
											TableName:  model.TableName("product"),
											ColumnName: model.ColumnName("id"),
										},
									},
								},
								{
									Name:     "name",
									FullName: "customer.name",
									Type: model.ColumnType{
										Base:  model.Varchar,
										Param: model.ColumnTypeParam(255),
									},
								},
							},
						},
						{
							Name: "product",
							Columns: []model.Column{
								{
									Name:     "id",
									FullName: "product.id",
									Type: model.ColumnType{
										Base:  model.Int,
										Param: model.ColumnTypeParam(14),
									},
								},
								{
									Name:     "owner",
									FullName: "product.owner",
									Type: model.ColumnType{
										Base:  model.Varchar,
										Param: model.ColumnTypeParam(255),
									},
									ForeignKeys: []model.ForeignKey{
										{
											TableName:  model.TableName("customer"),
											ColumnName: model.ColumnName("name"),
										},
									},
								},
							},
						},
					},
				},
			},
			want: model.SchemaGraph{
				AdjacencyMatrix: model.AdjacencyMatrix{
					{0, 0, 1, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 1, 0, 0},
				},
				ColumnNodes: []model.ColumnNode{
					model.NewColumnNode(model.Column{
						Name:     "id",
						FullName: "customer.id",
						Type: model.ColumnType{
							Base:  model.Int,
							Param: model.ColumnTypeParam(10),
						},
						ForeignKeys: []model.ForeignKey{
							{
								TableName:  "product",
								ColumnName: "id",
							},
						},
					},
						false,
						0,
					),
					model.NewColumnNode(model.Column{
						Name:     "name",
						FullName: "customer.name",
						Type: model.ColumnType{
							Base:  model.Varchar,
							Param: model.ColumnTypeParam(255),
						},
					},
						false,
						1,
					),
					model.NewColumnNode(model.Column{
						Name:     "id",
						FullName: "product.id",
						Type: model.ColumnType{
							Base:  model.Int,
							Param: model.ColumnTypeParam(14),
						},
					},
						false,
						2,
					),
					model.NewColumnNode(model.Column{
						Name:     "owner",
						FullName: "product.owner",
						Type: model.ColumnType{
							Base:  model.Varchar,
							Param: model.ColumnTypeParam(255),
						},
						ForeignKeys: []model.ForeignKey{
							{
								TableName:  model.TableName("customer"),
								ColumnName: model.ColumnName("name"),
							},
						},
					},
						false,
						3,
					),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateSchemaGraph(tt.args.schema)
			diff := cmp.Diff(got, tt.want, cmp.AllowUnexported(model.ColumnNode{}))
			if diff != "" {
				t.Errorf("GenerateSchemaGraph(); -got, +want\n%v", diff)
			}
		})
	}
}

func TestGenerateValuesForColumns(t *testing.T) {
	type args struct {
		cg model.SchemaGraph
		n  int
	}
	tests := []struct {
		name     string
		args     args
		assertFn func(map[model.ColumnFullName][]model.Value)
	}{
		{
			name: "return map of model.values for columns considering foreign keys",
			args: args{
				cg: model.SchemaGraph{
					AdjacencyMatrix: model.AdjacencyMatrix{
						{0, 0, 0, 1},
						{0, 0, 0, 0},
						{0, 0, 0, 0},
						{0, 1, 1, 0},
					},
					ColumnNodes: []model.ColumnNode{
						model.NewColumnNode(model.Column{
							FullName: "test0",
							Type: model.ColumnType{
								Base:  model.Int,
								Param: model.ColumnTypeParam(1), //TODO: if the given schema is wrong(e.g. the param for column test0 was Int(3)), should sqloth validate it and notice users the inconsistency?
							},
						},
							false,
							0,
						),
						model.NewColumnNode(model.Column{
							FullName: "test1",
							Type: model.ColumnType{
								Base:  model.Int,
								Param: model.ColumnTypeParam(1),
							},
						},
							false,
							1,
						),
						model.NewColumnNode(model.Column{
							FullName: "test2",
							Type: model.ColumnType{
								Base:  model.Int,
								Param: model.ColumnTypeParam(1),
							},
						},
							false,
							2,
						),
						model.NewColumnNode(model.Column{
							FullName: "test3",
							Type: model.ColumnType{
								Base:  model.Int,
								Param: model.ColumnTypeParam(2),
							},
						},
							false,
							3,
						),
					},
				},
				n: 3,
			},
			assertFn: func(m map[model.ColumnFullName][]model.Value) {
				for key, vs := range m {
					switch key {
					case "test0":
						for idx, v := range vs {
							if v != m["test3"][idx] {
								t.Errorf("model.values of test1 is not valid; idx: %v, model.value: %v", idx, v)
							}
						}
					case "test1", "test2":
						for idx, v := range vs {
							n, err := strconv.Atoi(string(v))
							if err != nil {
								t.Errorf("cannot convert model.value to int; model.value: %v", v)
							}
							// validation using intRangeMap[Int]
							if !(-2147483648 <= n || n <= 2147483647) {
								t.Errorf("model.values of %v is out of range; idx: %v, model.value: %v", key, idx, v)
							}
						}
					case "test3":
						for idx, v := range vs {
							if v != m["test1"][idx]+m["test2"][idx] {
								t.Errorf("model.values of test3 is not valid; idx: %v, model.value: %v", idx, v)
							}
						}
					}
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateValuesForColumns(tt.args.cg, tt.args.n)
			tt.assertFn(got)
		})
	}
}

func TestGenerateRecordsForTables(t *testing.T) {
	type args struct {
		vfc    map[model.ColumnFullName][]model.Value
		schema model.Schema
		n      int
	}
	tests := []struct {
		name string
		args args
		want map[model.TableName][]model.Record
	}{
		{
			name: "generate model.records for tables from map of column full-names and values",
			args: args{
				vfc: map[model.ColumnFullName][]model.Value{
					"table1.test1": []model.Value{"1", "2", "3"},
					"table1.test2": []model.Value{"aaa", "bbb", "ccc"},
					"table2.test1": []model.Value{"v1", "v2", "v3"},
				},
				schema: model.Schema{
					Tables: []model.Table{
						{Name: "table1", Columns: []model.Column{{FullName: "table1.test1"}, {FullName: "table1.test2"}}},
						{Name: "table2", Columns: []model.Column{{FullName: "table2.test1"}}},
					},
				},
				n: 3,
			},
			want: map[model.TableName][]model.Record{
				"table1": []model.Record{{"1", "aaa"}, {"2", "bbb"}, {"3", "ccc"}},
				"table2": []model.Record{{"v1"}, {"v2"}, {"v3"}},
			},
		},
		{
			name: "skip auto increment columns",
			args: args{
				vfc: map[model.ColumnFullName][]model.Value{
					"table1.test1": []model.Value{"1", "2", "3"},
					"table1.test2": []model.Value{"aaa", "bbb", "ccc"},
					"table2.test1": []model.Value{"v1", "v2", "v3"},
				},
				schema: model.Schema{
					Tables: []model.Table{
						{Name: "table1", Columns: []model.Column{{FullName: "table1.test1", AutoIncrement: true}, {FullName: "table1.test2"}}},
						{Name: "table2", Columns: []model.Column{{FullName: "table2.test1"}}},
					},
				},
				n: 3,
			},
			want: map[model.TableName][]model.Record{
				"table1": []model.Record{{"aaa"}, {"bbb"}, {"ccc"}},
				"table2": []model.Record{{"v1"}, {"v2"}, {"v3"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateRecordsForTables(tt.args.vfc, tt.args.schema, tt.args.n)
			diff := cmp.Diff(got, tt.want)
			if diff != "" {
				t.Errorf("Generatemodel.RecordsForTables(); -got, +want\n%v", diff)
			}
		})
	}
}
