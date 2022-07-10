package model

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

////// Column //////////////////////////

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
