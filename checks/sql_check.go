package checks

import "database/sql"

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
	return "mysql"
}
