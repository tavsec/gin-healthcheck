package controllers

import (
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tavsec/gin-healthcheck/checks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

var (
	res *httptest.ResponseRecorder
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestHealthcheckController(t *testing.T) {
	router := gin.New()
	router.GET("/healthcheck", HealthcheckController([]checks.Check{}))
	assertRequest(t, router, "GET", "/healthcheck", "", 200, "[]")

}

func TestHealthcheckControllerWithSqlCheck(t *testing.T) {
	router := gin.New()
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	router.GET("/healthcheck", HealthcheckController([]checks.Check{checks.SqlCheck{Sql: db}}))

	response, err := json.Marshal([]CheckStatus{{
		Name: "mysql",
		Pass: true,
	}})
	assertRequest(t, router, "GET", "/healthcheck", "", 200, string(response))

}

func assertRequest(t *testing.T, router *gin.Engine, method string, path string, body string, assertStatus int, assertBody string) {
	res = httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, strings.NewReader(body))

	router.ServeHTTP(res, req)

	if res.Code != assertStatus {
		t.Errorf("expected %d, got %d", assertStatus, res.Code)
	}
	if b := res.Body.String(); b != assertBody {
		t.Errorf("expected %q, got %q", assertBody, b)
	}
}
