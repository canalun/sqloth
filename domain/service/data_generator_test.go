package service

import (
	"reflect"
	"testing"

	"github.com/canalun/sqloth/domain/model"
)

// TODO: write test
func TestGenerateQuery(t *testing.T) {
	type args struct {
		schema model.Schema
		data   map[string][]model.ColumnData
		n      int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateQuery(tt.args.schema, tt.args.data, tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO: write test
func TestGenerateColumnData(t *testing.T) {
	type args struct {
		sortedColumnList []int
		schema           model.Schema
		n                int
	}
	tests := []struct {
		name string
		args args
		want map[string][]model.ColumnData
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateColumnData(tt.args.sortedColumnList, tt.args.schema, tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateColumnData() = %v, want %v", got, tt.want)
			}
		})
	}
}
