package driver

import "github.com/canalun/sqloth/domain/model"

type Driver interface {
	GetSchema() model.Schema
}
