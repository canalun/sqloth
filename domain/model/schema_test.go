package model

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSchema_GetMapFromColumnToIndex(t *testing.T) {
	type fields struct {
		Tables []Table
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]int
	}{
		{
			name: "return correct map",
			fields: fields{
				Tables: []Table{
					{Name: "customer", Columns: []Column{{Name: "id"}, {Name: "name"}}},
					{Name: "product", Columns: []Column{{Name: "id"}, {Name: "owner"}, {Name: "name"}}},
				},
			},
			want: map[string]int{
				"customer.id":   0,
				"customer.name": 1,
				"product.id":    2,
				"product.owner": 3,
				"product.name":  4,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Schema{
				Tables: tt.fields.Tables,
			}
			got := s.GetMapFromColumnToIndex()
			diff := cmp.Diff(got, tt.want)
			if diff != "" {
				t.Error("-:got, +:want", diff)
			}
		})
	}
}

func TestSchema_GetMapFromIndexToColumnName(t *testing.T) {
	type fields struct {
		Tables []Table
	}
	tests := []struct {
		name   string
		fields fields
		want   map[int]string
	}{
		{
			name: "return correct map",
			fields: fields{
				Tables: []Table{
					{Name: "customer", Columns: []Column{{Name: "id"}, {Name: "name"}}},
					{Name: "product", Columns: []Column{{Name: "id"}, {Name: "owner"}, {Name: "name"}}},
				},
			},
			want: map[int]string{
				0: "customer.id",
				1: "customer.name",
				2: "product.id",
				3: "product.owner",
				4: "product.name",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Schema{
				Tables: tt.fields.Tables,
			}
			got := s.GetMapFromIndexToColumnName()
			diff := cmp.Diff(got, tt.want)
			if diff != "" {
				t.Error("-:got, +:want", diff)
			}
		})
	}
}

func TestSchema_GetMapFromIndexToColumn(t *testing.T) {
	type fields struct {
		Tables []Table
	}
	tests := []struct {
		name   string
		fields fields
		want   map[int]Column
	}{
		{
			name: "return correct map",
			fields: fields{
				Tables: []Table{
					{Name: "customer", Columns: []Column{{Name: "id"}, {Name: "name"}}},
					{Name: "product", Columns: []Column{{Name: "id"}, {Name: "owner"}, {Name: "name"}}},
				},
			},
			want: map[int]Column{
				0: Column{Name: "id"},
				1: Column{Name: "name"},
				2: Column{Name: "id"},
				3: Column{Name: "owner"},
				4: Column{Name: "name"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Schema{
				Tables: tt.fields.Tables,
			}
			got := s.GetMapFromIndexToColumn()
			diff := cmp.Diff(got, tt.want)
			if diff != "" {
				t.Error("-:got, +:want", diff)
			}
		})
	}
}
