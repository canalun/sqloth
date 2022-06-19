package usecase

import (
	"github.com/canalun/sqloth/domain/driver"
	"github.com/canalun/sqloth/domain/model"
)

type Usecase struct {
	driver driver.Driver
}

func NewUsecase(driver driver.Driver) Usecase {
	return Usecase{
		driver: driver,
	}
}

//TODO: refactoring the entire
func (u Usecase) GenerateQueryOfDummyData(num int) []string {
	schema := u.driver.GetSchema()

	columnGraph := model.GenerateColumnGraph(schema)
	valuesForColumns := model.GenerateValuesForColumns(columnGraph, num)
	recordsForTables := model.GenerateRecordsForTables(valuesForColumns, schema, num)
	queries := model.GenerateQuery(recordsForTables, schema)

	return queries
}
