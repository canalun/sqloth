package model

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_NewAdjacencyMatrix(t *testing.T) {
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
			got := NewAdjacencyMatrix(tt.args.n)
			diff := cmp.Diff(got, tt.want)
			if diff != "" {
				t.Error("-:got, +:want", diff)
			}
		})
	}
}

func TestSchemaGraph_IsAllDone(t *testing.T) {
	type fields struct {
		ColumnNodes []ColumnNode
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "true if isDone values of all column nodes are true",
			fields: fields{
				ColumnNodes: []ColumnNode{
					{
						isDone: true,
					},
					{
						isDone: true,
					},
				},
			},
			want: true,
		},
		{
			name: "false if isDone of one of column nodes is false",
			fields: fields{
				ColumnNodes: []ColumnNode{
					{
						isDone: true,
					},
					{
						isDone: false,
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cg := SchemaGraph{
				ColumnNodes: tt.fields.ColumnNodes,
			}
			if got := cg.IsAllDone(); got != tt.want {
				t.Errorf("SchemaGraph.IsAllDone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSchemaGraph_HasParentNodes(t *testing.T) {
	type fields struct {
		AdjacencyMatrix AdjacencyMatrix
	}
	type args struct {
		i int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "return true if the node with the given index has parent nodes",
			fields: fields{
				AdjacencyMatrix: AdjacencyMatrix{
					{0, 0, 1, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 1, 0, 0},
				},
			},
			args:    args{i: 3},
			want:    true,
			wantErr: false,
		},
		{
			name: "return false if the node with the given index does not have parent nodes",
			fields: fields{
				AdjacencyMatrix: AdjacencyMatrix{
					{0, 0, 1, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 1, 0, 0},
				},
			},
			args:    args{i: 2},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cg := SchemaGraph{
				AdjacencyMatrix: tt.fields.AdjacencyMatrix,
			}
			got, err := cg.HasParentNodes(tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("SchemaGraph.HasParentNodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SchemaGraph.HasParentNodes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSchemaGraph_IsParentNodesAreAllDone(t *testing.T) {
	type fields struct {
		AdjacencyMatrix AdjacencyMatrix
		ColumnNodes     []ColumnNode
	}
	type args struct {
		i int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "return true if isDone values of parent nodes of the node with the given index are all true",
			fields: fields{
				AdjacencyMatrix: AdjacencyMatrix{
					{0, 0, 1, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 1, 0, 0},
				},
				ColumnNodes: []ColumnNode{
					{isDone: true},
					{isDone: false},
					{isDone: true},
					{isDone: false},
				},
			},
			args:    args{i: 0},
			want:    true,
			wantErr: false,
		},
		{
			name: "return false if one of isDone values of parent nodes of the node with the given index is false",
			fields: fields{
				AdjacencyMatrix: AdjacencyMatrix{
					{0, 0, 1, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 1, 0, 0},
				},
				ColumnNodes: []ColumnNode{
					{isDone: true},
					{isDone: false},
					{isDone: true},
					{isDone: false},
				},
			},
			args:    args{i: 3},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cg := SchemaGraph{
				AdjacencyMatrix: tt.fields.AdjacencyMatrix,
				ColumnNodes:     tt.fields.ColumnNodes,
			}
			got, err := cg.IsParentNodesAreAllDone(tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("SchemaGraph.IsParentNodesAreAllDone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SchemaGraph.IsParentNodesAreAllDone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSchemaGraph_ParentNodeIndexes(t *testing.T) {
	type fields struct {
		AdjacencyMatrix AdjacencyMatrix
	}
	type args struct {
		i int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []int
		wantErr bool
	}{
		{
			name: "return indexes of parent nodes of the node with the given index",
			fields: fields{
				AdjacencyMatrix: AdjacencyMatrix{
					{0, 0, 1, 1},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 1, 1, 0},
				},
			},
			args:    args{i: 0},
			want:    []int{2, 3},
			wantErr: false,
		},
		{
			name: "return empty slice if the node with the given index has no parent nodes",
			fields: fields{
				AdjacencyMatrix: AdjacencyMatrix{
					{0, 0, 1, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 1, 0, 0},
				},
			},
			args:    args{i: 2},
			want:    []int{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cg := SchemaGraph{
				AdjacencyMatrix: tt.fields.AdjacencyMatrix,
			}
			got, err := cg.ParentNodeIndexes(tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("SchemaGraph.ParentNodeIndexes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SchemaGraph.ParentNodeIndexes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSchemaGraph_HasChildrenNodes(t *testing.T) {
	type fields struct {
		AdjacencyMatrix AdjacencyMatrix
	}
	type args struct {
		i int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "return true if the node with the given index has children nodes",
			fields: fields{
				AdjacencyMatrix: AdjacencyMatrix{
					{0, 0, 1, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 1, 0, 0},
				},
			},
			args:    args{i: 2},
			want:    true,
			wantErr: false,
		},
		{
			name: "return false if the node with the given index does not have children nodes",
			fields: fields{
				AdjacencyMatrix: AdjacencyMatrix{
					{0, 0, 1, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 1, 0, 0},
				},
			},
			args:    args{i: 0},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cg := SchemaGraph{
				AdjacencyMatrix: tt.fields.AdjacencyMatrix,
			}
			got, err := cg.HasChildrenNodes(tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("SchemaGraph.HasChildrenNodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SchemaGraph.HasChildrenNodes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSchemaGraph_ChildrenNodeIndexes(t *testing.T) {
	type fields struct {
		AdjacencyMatrix AdjacencyMatrix
	}
	type args struct {
		i int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []int
		wantErr bool
	}{
		{
			name: "return indexes of children nodes of the node with the given index",
			fields: fields{
				AdjacencyMatrix: AdjacencyMatrix{
					{0, 0, 1, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 1, 1, 0},
				},
			},
			args:    args{i: 2},
			want:    []int{0, 3},
			wantErr: false,
		},
		{
			name: "return empty slice if the node with the given index has no children nodes",
			fields: fields{
				AdjacencyMatrix: AdjacencyMatrix{
					{0, 0, 1, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 1, 0, 0},
				},
			},
			args:    args{i: 0},
			want:    []int{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cg := SchemaGraph{
				AdjacencyMatrix: tt.fields.AdjacencyMatrix,
			}
			got, err := cg.ChildrenNodeIndexes(tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("SchemaGraph.ChildrenNodeIndexes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SchemaGraph.ChildrenNodeIndexes() = %v, want %v", got, tt.want)
			}
		})
	}
}
