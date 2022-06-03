generate-mock:
	go install github.com/golang/mock/mockgen@latest
	mockgen -destination=./domain/driver/mock_driver/mock_driver.go -package=mock_driver github.com/canalun/sqloth/domain/driver Driver