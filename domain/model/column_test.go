package model

import (
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestStrToColumnTypeBase(t *testing.T) {
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
			got, err := StrToColumnTypeBase(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("StrToColumnTypeBase() error = %v, wantErr %v", err, tt.wantErr)
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
		Constraints   []Constraint
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
			name: "return int slice {0...n} when AutoIncrement is true",
			fields: fields{
				Name:     "test",
				FullName: "test",
				Type: ColumnType{
					Base: Int,
				},
				AutoIncrement: true,
				Constraints:   []Constraint{},
			},
			args: args{n: 3},
			want: []Value{Value(strconv.Itoa(0)), Value(strconv.Itoa(1)), Value(strconv.Itoa(2))},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Column{
				Name:          tt.fields.Name,
				FullName:      tt.fields.FullName,
				Type:          tt.fields.Type,
				AutoIncrement: tt.fields.AutoIncrement,
				Constraints:   tt.fields.Constraints,
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
		cg ColumnGraph
		n  int
	}
	tests := []struct {
		name     string
		args     args
		assertFn func(map[ColumnFullName][]Value)
	}{
		{
			name: "return map of values for columns considering foreign key constraints",
			args: args{
				cg: ColumnGraph{
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
