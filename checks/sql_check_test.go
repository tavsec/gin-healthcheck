package checks

import (
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
)

func TestSqlCheck_Pass(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	check := SqlCheck{Sql: db}
	if !check.Pass() {
		t.Errorf("Expected SqlCheck.Pass to return true, got false")
	}

	check = SqlCheck{Sql: nil}
	if check.Pass() {
		t.Errorf("Expected SqlCheck.Pass to return false, got true")
	}
}

func TestSqlCheck_Name(t *testing.T) {
	check := SqlCheck{Sql: nil}
	if check.Name() != "mysql" {
		t.Errorf("Expected SqlCheck.Name to return 'mysql', got '%s'", check.Name())
	}
}

func TestSqlCheck_Fail(t *testing.T) {
	db, _, err := sqlmock.New()
	db.Close()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	check := SqlCheck{Sql: db}
	if check.Pass() {
		t.Errorf("Expected SqlCheck.Pass to return false, got true")
	}

}
