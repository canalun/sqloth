package file_driver

import (
	"testing"

	"github.com/canalun/sqloth/domain/model"
	"github.com/google/go-cmp/cmp"
)

func TestGetSchema(t *testing.T) {
	type fields struct {
		FilePath string
	}
	tests := []struct {
		name   string
		fields fields
		want   model.Schema
	}{
		{
			name: "can get schema from sql schema file with correct auto_increment settings",
			fields: fields{
				FilePath: "./testSchema.sql",
			},
			want: model.Schema{
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
								AutoIncrement: true,
								Unsigned:      true,
							},
							{
								Name:     "created_at",
								FullName: "customer.created_at",
								Type: model.ColumnType{
									Base: model.Timestamp,
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
							{
								Name:     "material",
								FullName: "customer.material",
								Type: model.ColumnType{
									Base: model.Json,
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
								AutoIncrement: true,
							},
							{
								Name:     "name",
								FullName: "product.name",
								Type: model.ColumnType{
									Base:  model.Varchar,
									Param: model.ColumnTypeParam(255),
								},
							},
							{
								Name:     "owner",
								FullName: "product.owner",
								Type: model.ColumnType{
									Base:  model.Varchar,
									Param: model.ColumnTypeParam(255),
								},
								Constraints: []model.Constraint{
									{
										TableName:  model.TableName("customer"),
										ColumnName: model.ColumnName("name"),
									},
								},
							},
							{
								Name:     "description",
								FullName: "product.description",
								Type: model.ColumnType{
									Base:  model.Text,
									Param: model.ColumnTypeParam(100),
								},
							},
							{
								Name:     "stock",
								FullName: "product.stock",
								Type: model.ColumnType{
									Base:  model.Tinyint,
									Param: model.ColumnTypeParam(1),
								},
							},
							{
								Name:     "sale_day",
								FullName: "product.sale_day",
								Type: model.ColumnType{
									Base: model.Datetime,
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fd := FileDriver{
				FilePath: tt.fields.FilePath,
			}
			got := fd.GetSchema()
			diff := cmp.Diff(got, tt.want)
			if diff != "" {
				t.Error("-:got, +:want", diff)
			}
		})
	}
}
