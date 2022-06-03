package service

import (
	"testing"

	"github.com/canalun/sqloth/domain/model"
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
		schema model.Schema
	}
	tests := []struct {
		name string
		args args
		want AdjacencyMatrix
	}{
		{
			name: "generate correct adjacency matrix",
			args: args{
				schema: model.Schema{
					Tables: []model.Table{
						{
							Name: "customer",
							Columns: []model.Column{
								{
									Name: "id",
									Type: model.ColumnType{
										Base:  model.Int,
										Param: model.ColumnTypeParam(10),
									},
									Constraint: model.Constraint{
										TableName:  "product",
										ColumnName: "id",
									},
								},
								{
									Name: "name",
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
									Name: "id",
									Type: model.ColumnType{
										Base:  model.Int,
										Param: model.ColumnTypeParam(14),
									},
								},
								{
									Name: "owner",
									Type: model.ColumnType{
										Base:  model.Varchar,
										Param: model.ColumnTypeParam(255),
									},
									Constraint: model.Constraint{
										TableName:  model.TableName("customer"),
										ColumnName: model.ColumnName("name"),
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
		schema model.Schema
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
				schema: model.Schema{
					Tables: []model.Table{
						{
							Name: "staff",
							Columns: []model.Column{
								{
									Name: "id",
									Type: model.ColumnType{
										Base:  model.Int,
										Param: model.ColumnTypeParam(14),
									},
								},
								{
									Name: "name",
									Type: model.ColumnType{
										Base:  model.Varchar,
										Param: model.ColumnTypeParam(255),
									},
								},
							},
						},
						{
							Name: "customer",
							Columns: []model.Column{
								{
									Name: "id",
									Type: model.ColumnType{
										Base:  model.Int,
										Param: model.ColumnTypeParam(10),
									},
									Constraint: model.Constraint{
										TableName:  "product",
										ColumnName: "id",
									},
								},
							},
						},
						{
							Name: "product",
							Columns: []model.Column{
								{
									Name: "id",
									Type: model.ColumnType{
										Base:  model.Int,
										Param: model.ColumnTypeParam(14),
									},
								},
								{
									Name: "owner",
									Type: model.ColumnType{
										Base:  model.Varchar,
										Param: model.ColumnTypeParam(255),
									},
									Constraint: model.Constraint{
										TableName:  model.TableName("staff"),
										ColumnName: model.ColumnName("name"),
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
