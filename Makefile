generate-mock:
	go install github.com/golang/mock/mockgen@latest
	mockgen -destination=./domain/driver/mock_driver/mock_driver.go -package=mock_driver github.com/canalun/sqloth/domain/driver Driver

test:
	go test -v ./...

run-ci:
	circleci config process .circleci/config.yml > process.yml
	circleci local execute -c process.yml --job test_and_build
	rm process.yml