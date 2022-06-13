package model

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewAdjacencyMatrix(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want AdjacencyMatrix
	}{
		{
			name: "return n*n zero matrix",
			args: args{n: 5},
			want: AdjacencyMatrix{
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newAdjacencyMatrix(tt.args.n)
			diff := cmp.Diff(got, tt.want)
			if diff != "" {
				t.Error("-:got, +:want", diff)
			}
		})
	}
}

func TestGenerateAdjacencyMatrix(t *testing.T) {
	type args struct {
		schema Schema
	}
	tests := []struct {
		name string
		args args
		want AdjacencyMatrix
	}{
		{
			name: "generate correct adjacency matrix",
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
									Constraints: []Constraint{
										{
											TableName:  "product",
											ColumnName: "id",
										},
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
									Constraints: []Constraint{
										{
											TableName:  TableName("customer"),
											ColumnName: ColumnName("name"),
										},
									},
								},
							},
						},
					},
				},
			},
			want: AdjacencyMatrix{
				{0, 0, 1, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 1, 0, 0},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateAdjacencyMatrix(tt.args.schema)
			diff := cmp.Diff(got, tt.want)
			if diff != "" {
				t.Error("-:got, +:want", diff)
			}
		})
	}
}

func TestGenerateSortedColumnList(t *testing.T) {
	type args struct {
		schema Schema
		am     AdjacencyMatrix
	}
	tests := []struct {
		name     string
		args     args
		assertFn func([]int) bool
	}{
		{
			name: "return correct list",
			// the columns below is related like;
			// customer.id(2) --> product.id(3), product.owner(4) --> staff.name(1), staff.id(0)
			// so the expected sorted list should meet the following;
			// 2 is before 3, 4 is before 1
			args: args{
				schema: Schema{
					Tables: []Table{
						{
							Name: "staff",
							Columns: []Column{
								{
									Name: "id",
									Type: ColumnType{
										Base:  Int,
										Param: ColumnTypeParam(14),
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
							Name: "customer",
							Columns: []Column{
								{
									Name: "id",
									Type: ColumnType{
										Base:  Int,
										Param: ColumnTypeParam(10),
									},
									Constraints: []Constraint{
										{
											TableName:  "product",
											ColumnName: "id",
										},
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
									Constraints: []Constraint{
										{
											TableName:  TableName("staff"),
											ColumnName: ColumnName("name"),
										},
									},
								},
							},
						},
					},
				},
				am: AdjacencyMatrix{
					{0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0},
					{0, 0, 0, 1, 0},
					{0, 0, 0, 0, 0},
					{0, 1, 0, 0, 0},
				},
			},
			assertFn: func(i []int) bool {
				// the expected sorted list should meet the following, as said above;
				// 2 is before 3, 4 is before 1
				m := map[int]int{}
				for index, value := range i {
					m[value] = index
				}
				if m[2] < m[3] && m[4] < m[1] {
					return true
				}
				return false
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateSortedColumnList(tt.args.schema, tt.args.am)
			if !tt.assertFn(got) {
				t.Errorf("returned value is NOT as expected...\nGenerateSortedColumnList() = %v", got)
			}
		})
	}
}
