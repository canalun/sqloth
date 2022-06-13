package usecase

import (
	"fmt"

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
	fmt.Printf("%#v\n", schema)

	columnGraph := model.GenerateColumnGraph(schema)
	fmt.Printf("%+v\n", columnGraph)
	valuesForColumns := model.GenerateValuesForColumns(columnGraph, num)
	fmt.Printf("%#v\n", valuesForColumns)
	recordsForTables := model.GenerateRecordsForTables(valuesForColumns, schema, num)
	queries := model.GenerateQuery(recordsForTables)

	return queries
}
