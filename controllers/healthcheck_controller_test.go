package controllers

import (
	"database/sql/driver"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/tavsec/gin-healthcheck/checks"
	"github.com/tavsec/gin-healthcheck/config"

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

var _ checks.Check = FailingCheck{}

func (c FailingCheck) Pass() bool {
	return false
}
func (c FailingCheck) Name() string {
	return "Failing Check"
}

type SlowCheck struct{}

var _ checks.Check = SlowCheck{}

func (c SlowCheck) Pass() bool {
	time.Sleep(2 * time.Second)
	return true
}

func (c SlowCheck) Name() string {
	return "Sloc Check"
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
	assert.NoError(t, err)
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

func TestSwitchResultBetweenCall(t *testing.T) {
	router := gin.New()
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	router.GET("/healthcheck", HealthcheckController([]checks.Check{checks.SqlCheck{Sql: db}}, conf))

	mock.ExpectPing().WillReturnError(nil)

	response, err := json.Marshal([]CheckStatus{{
		Name: "mysql",
		Pass: true,
	}})
	assert.NoError(t, err)
	assertRequest(t, router, "GET", "/healthcheck", "", 200, string(response))

	mock.ExpectPing().WillReturnError(driver.ErrBadConn)

	response, err = json.Marshal([]CheckStatus{{
		Name: "mysql",
		Pass: false,
	}})
	assert.NoError(t, err)
	assertRequest(t, router, "GET", "/healthcheck", "", 503, string(response))
}

func TestParallelCheck(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	router := gin.New()
	router.GET("/healthcheck", HealthcheckController([]checks.Check{checks.SqlCheck{Sql: db}, FailingCheck{}}, conf))

	response, _ := json.Marshal([]CheckStatus{{
		Name: "mysql",
		Pass: true,
	}, {
		Name: "Failing Check",
		Pass: false,
	}})
	assertRequest(t, router, "GET", "/healthcheck", "", 503, string(response))
}

func TestParallelCheckSpeed(t *testing.T) {
	f := HealthcheckController([]checks.Check{SlowCheck{}, SlowCheck{}}, conf)

	start := time.Now()
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	f(c)
	assert.Less(t, time.Since(start), 3*time.Second)
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

func BenchmarkCheckLatency(b *testing.B) {
	f := HealthcheckController([]checks.Check{SlowCheck{}, SlowCheck{}}, conf)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		f(c)
	}
}
