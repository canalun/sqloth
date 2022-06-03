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
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			name: "can generate query of dummy data from sql schema file",
			fields: fields{
				driver: func(ctrl *gomock.Controller) driver.Driver {
					m := mock_driver.NewMockDriver(ctrl)
					m.EXPECT().GetSchema().Return(model.Schema{

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
									},
									{
										Name: "created_at",
										Type: model.ColumnType{
											Base: model.Timestamp,
										},
									},
									{
										Name: "name",
										Type: model.ColumnType{
											Base:  model.Varchar,
											Param: model.ColumnTypeParam(255),
										},
									},
									{
										Name: "material",
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
									{
										Name: "description",
										Type: model.ColumnType{
											Base: model.Text,
										},
									},
									{
										Name: "stock",
										Type: model.ColumnType{
											Base:  model.Tinyint,
											Param: model.ColumnTypeParam(1),
										},
									},
									{
										Name: "sale_day",
										Type: model.ColumnType{
											Base: model.Datetime,
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
			// TODO: mod want data
			want: []string{""},
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
			diff := cmp.Diff(got, tt.want)
			if diff != "" {
				t.Error("-:got, +:want", diff)
			}
		})
	}
}
