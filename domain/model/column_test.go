package model

import (
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_NewColumnType(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		wantCt  ColumnType
		wantErr bool
	}{
		{
			name: "set base of type without param(e.g. text, json)",
			args: args{str: "JSON"},
			wantCt: ColumnType{
				Base: Json,
			},
			wantErr: false,
		},
		{
			name: "separate and set base and param of type with param(e.g. varchar, int)",
			args: args{str: "VARCHAR(255)"},
			wantCt: ColumnType{
				Base:  Varchar,
				Param: ColumnTypeParam(255),
			},
			wantErr: false,
		},
		{
			name: "set base and additional param for type TEXT",
			args: args{str: "TEXT"},
			wantCt: ColumnType{
				Base:  Text,
				Param: ColumnTypeParam(100),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCt, err := NewColumnType(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewColumnType() error = %v, wantErr %v", err, tt.wantErr)
			}
			diff := cmp.Diff(gotCt, tt.wantCt)
			if diff != "" {
				t.Error("-:got, +:want", diff)
			}
		})
	}
}

func Test_newColumnTypeBase(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    ColumnTypeBase
		wantErr bool
	}{
		{
			name:    "return error when handling unregistered type",
			args:    args{str: "aaaa"},
			wantErr: true,
		},
		{
			name:    "return ColumnTypeBase when handling registered type",
			args:    args{str: "varchar"},
			want:    Varchar,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newColumnTypeBase(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("newColumnTypeBase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			diff := cmp.Diff(got, tt.want)
			if diff != "" {
				t.Error("-:got, +:want", diff)
			}
		})
	}
}

func TestColumn_GenerateData(t *testing.T) {
	type fields struct {
		Name          ColumnName
		FullName      ColumnFullName
		Type          ColumnType
		AutoIncrement bool
		ForeignKeys   []ForeignKey
	}
	type args struct {
		n int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []Value
	}{
		{
			name: "return slice of 'NULL' when AutoIncrement is true",
			fields: fields{
				Name:     "test",
				FullName: "test",
				Type: ColumnType{
					Base: Int,
				},
				AutoIncrement: true,
				ForeignKeys:   []ForeignKey{},
			},
			args: args{n: 3},
			want: []Value{Value("NULL"), Value("NULL"), Value("NULL")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Column{
				Name:          tt.fields.Name,
				FullName:      tt.fields.FullName,
				Type:          tt.fields.Type,
				AutoIncrement: tt.fields.AutoIncrement,
				ForeignKeys:   tt.fields.ForeignKeys,
			}
			got := c.GenerateData(tt.args.n)
			diff := cmp.Diff(got, tt.want)
			if diff != "" {
				t.Error("Column.GenerateData(); -:got, +:want", diff)
			}
		})
	}
}

func TestGenerateValuesForColumns(t *testing.T) {
	type args struct {
		cg SchemaGraph
		n  int
	}
	tests := []struct {
		name     string
		args     args
		assertFn func(map[ColumnFullName][]Value)
	}{
		{
			name: "return map of values for columns considering foreign keys",
			args: args{
				cg: SchemaGraph{
					AdjacencyMatrix: AdjacencyMatrix{
						{0, 0, 0, 1},
						{0, 0, 0, 0},
						{0, 0, 0, 0},
						{0, 1, 1, 0},
					},
					ColumnNodes: []ColumnNode{
						ColumnNode{
							column: Column{
								FullName: "test0",
								Type: ColumnType{
									Base:  Int,
									Param: ColumnTypeParam(1), //TODO: if the given schema is wrong(e.g. the param for column test0 was Int(3)), should sqloth validate it and notice users the inconsistency?
								},
							},
							isDone: false,
							index:  0,
						},
						ColumnNode{
							column: Column{
								FullName: "test1",
								Type: ColumnType{
									Base:  Int,
									Param: ColumnTypeParam(1),
								},
							},
							isDone: false,
							index:  1,
						},
						ColumnNode{
							column: Column{
								FullName: "test2",
								Type: ColumnType{
									Base:  Int,
									Param: ColumnTypeParam(1),
								},
							},
							isDone: false,
							index:  2,
						},
						ColumnNode{
							column: Column{
								FullName: "test3",
								Type: ColumnType{
									Base:  Int,
									Param: ColumnTypeParam(2),
								},
							},
							isDone: false,
							index:  3,
						},
					},
				},
				n: 3,
			},
			assertFn: func(m map[ColumnFullName][]Value) {
				for key, vs := range m {
					switch key {
					case "test0":
						for idx, v := range vs {
							if v != m["test3"][idx] {
								t.Errorf("values of test1 is not valid; idx: %v, value: %v", idx, v)
							}
						}
					case "test1", "test2":
						for idx, v := range vs {
							n, err := strconv.Atoi(string(v))
							if err != nil {
								t.Errorf("cannot convert value to int; value: %v", v)
							}
							if !(intRangeMap[Int][0] <= n || n <= intRangeMap[Int][1]) {
								t.Errorf("values of %v is out of range; idx: %v, value: %v", key, idx, v)
							}
						}
					case "test3":
						for idx, v := range vs {
							if v != m["test1"][idx]+m["test2"][idx] {
								t.Errorf("values of test3 is not valid; idx: %v, value: %v", idx, v)
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
