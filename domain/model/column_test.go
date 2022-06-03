package model

import (
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
