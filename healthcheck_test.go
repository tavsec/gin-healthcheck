package gin_healthcheck

import (
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/tavsec/gin-healthcheck/checks"
	config2 "github.com/tavsec/gin-healthcheck/config"
	"github.com/tavsec/gin-healthcheck/controllers"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	res *httptest.ResponseRecorder
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestInitHealthcheckWithDefaultConfig(t *testing.T) {
	router := gin.Default()
	config := config2.DefaultConfig()
	err := New(router, config, make([]checks.Check, 0))
	if err != nil {
		t.Fatal(err)
	}

	healthRoute := router.Routes()[0]

	if healthRoute.Path != config.HealthPath {
		t.Errorf("Healthcheck path route is not equal to config path")
	}
}

func TestHealthcheckResponseNoChecks(t *testing.T) {
	router := gin.Default()
	config := config2.DefaultConfig()
	New(router, config, []checks.Check{})

	assertRequest(t, router, "GET", config.HealthPath, "", 200, "[]")
}

func TestHealthcheckResponseMySqlCheck(t *testing.T) {
	router := gin.Default()
	config := config2.DefaultConfig()
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	c := []checks.Check{checks.SqlCheck{Sql: db}}
	New(router, config, c)

	response, err := json.Marshal([]controllers.CheckStatus{{
		Name: "mysql",
		Pass: true,
	}})
	assertRequest(t, router, "GET", config.HealthPath, "", 200, string(response))
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
