package usecase

import (
	"testing"

	"github.com/canalun/sqloth/domain/driver"
	"github.com/canalun/sqloth/domain/driver/mock_driver"
	"github.com/canalun/sqloth/domain/model"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
)

func TestGenerateQueryOfDummyData(t *testing.T) {
	type fields struct {
		driver func(ctrl *gomock.Controller) driver.Driver
	}
	type args struct {
		num int
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		assertFn func([]string)
	}{
		{
			name: "can generate query of dummy data from sql schema file with foreign keys",
			fields: fields{
				driver: func(ctrl *gomock.Controller) driver.Driver {
					m := mock_driver.NewMockDriver(ctrl)
					m.EXPECT().GetSchema().Return(model.Schema{
						Tables: []model.Table{
							{
								Name: "customer",
								Columns: []model.Column{
									{
										Name:     "id",
										FullName: "customer.id",
										Type: model.ColumnType{
											Base:  model.Int,
											Param: model.ColumnTypeParam(3),
										},
										AutoIncrement: true,
									},
									{
										Name:     "name",
										FullName: "customer.name",
										Type: model.ColumnType{
											Base:  model.Varchar,
											Param: model.ColumnTypeParam(3),
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
											Param: model.ColumnTypeParam(3),
										},
									},
									{
										Name:     "customeridname",
										FullName: "product.customeridname",
										Type: model.ColumnType{
											Base: model.Text,
										},
										ForeignKeys: []model.ForeignKey{
											{
												TableName:  "customer",
												ColumnName: "id",
											},
											{
												TableName:  "customer",
												ColumnName: "name",
											},
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
								},
							},
						},
					})
					return m
				},
			},
			args: args{num: 3},
			//TODO: mod assertFn
			assertFn: func(s []string) {
				diff := cmp.Diff(s[0], "SET foreign_key_checks = 0;")
				if diff != "" {
					t.Errorf(diff)
				}
				diff = cmp.Diff(s[len(s)-1], "SET foreign_key_checks = 1;")
				if diff != "" {
					t.Errorf(diff)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u := Usecase{
				driver: tt.fields.driver(ctrl),
			}
			got := u.GenerateQueryOfDummyData(tt.args.num)
			tt.assertFn(got)
		})
	}
}
