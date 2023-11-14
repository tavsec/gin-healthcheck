package checks

import "database/sql"

type SqlCheck struct {
	Title string
	Sql   *sql.DB
}

func (s SqlCheck) Pass() bool {
	if s.Sql == nil {
		return false
	}
	return s.Sql.Ping() == nil
}

func (s SqlCheck) Name() string {
	if s.Title != "" {
		return s.Title
	}
	return "mysql"
}
