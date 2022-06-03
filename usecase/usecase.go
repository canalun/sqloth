package usecase

import (
	"github.com/canalun/sqloth/domain/driver"
	"github.com/canalun/sqloth/domain/service"
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

	adjacencyMatrix := service.GenerateAdjacencyMatrix(schema)
	sortedColumnList := service.GenerateSortedColumnList(schema, adjacencyMatrix)
	data := service.GenerateColumnData(sortedColumnList, schema, num)
	query := service.GenerateQuery(schema, data, num)

	return query
}
