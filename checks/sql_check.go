package checks

import (
	"database/sql"
	"reflect"
)

type SqlCheck struct {
	Sql *sql.DB
}

func (s SqlCheck) Pass() bool {
	if s.Sql == nil {
		return false
	}

	err := s.Sql.Ping()

	return err == nil
}

func (s SqlCheck) Name() string {
	if s.Sql == nil {
		return "no_driver"
	}

	driver := s.Sql.Driver()
	return reflect.TypeOf(driver).String()
}
