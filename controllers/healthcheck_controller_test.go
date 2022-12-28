package controllers

import (
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tavsec/gin-healthcheck/checks"
	"github.com/tavsec/gin-healthcheck/config"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

var (
	res  *httptest.ResponseRecorder
	conf config.Config
)

func init() {
	gin.SetMode(gin.TestMode)
	conf = config.DefaultConfig()
}

type FailingCheck struct{}

func (c FailingCheck) Pass() bool {
	return false
}
func (c FailingCheck) Name() string {
	return "Failing Check"
}

func TestHealthcheckController(t *testing.T) {
	router := gin.New()
	router.GET("/healthcheck", HealthcheckController([]checks.Check{}, conf))
	assertRequest(t, router, "GET", "/healthcheck", "", 200, "[]")

}

func TestHealthcheckControllerWithSqlCheck(t *testing.T) {
	router := gin.New()
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	router.GET("/healthcheck", HealthcheckController([]checks.Check{checks.SqlCheck{Sql: db}}, conf))

	response, err := json.Marshal([]CheckStatus{{
		Name: "mysql",
		Pass: true,
	}})
	assertRequest(t, router, "GET", "/healthcheck", "", 200, string(response))

}

func TestNoPass(t *testing.T) {
	router := gin.New()
	router.GET("/healthcheck", HealthcheckController([]checks.Check{FailingCheck{}}, conf))

	response, _ := json.Marshal([]CheckStatus{{
		Name: "Failing Check",
		Pass: false,
	}})
	assertRequest(t, router, "GET", "/healthcheck", "", 503, string(response))
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
