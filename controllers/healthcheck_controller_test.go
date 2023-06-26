package controllers

import (
	"database/sql/driver"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
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
	return "Slow Check"
}

type ControlledCheck struct{ willPass bool }

var _ checks.Check = (*ControlledCheck)(nil)

func (c *ControlledCheck) Pass() bool {
	return c.willPass
}

func (c *ControlledCheck) Name() string {
	return "Controlled Check"
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

func TestNotification(t *testing.T) {
	router := gin.New()

	controlled := &ControlledCheck{willPass: true}
	conf := config.DefaultConfig()
	conf.FailureNotification.Chan = make(chan error, 1)
	defer close(conf.FailureNotification.Chan)
	conf.FailureNotification.Threshold = 3

	router.GET("/healthcheck", HealthcheckController([]checks.Check{controlled}, conf))

	successResponse, err := json.Marshal([]CheckStatus{{
		Name: "Controlled Check",
		Pass: true,
	}})
	assert.NoError(t, err)
	failureResponse, err := json.Marshal([]CheckStatus{{
		Name: "Controlled Check",
		Pass: false,
	}})
	assert.NoError(t, err)

	var wg sync.WaitGroup
	var errNotification error
	goWaitOnChan := func() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errNotification = <-conf.FailureNotification.Chan
		}()
	}

	goWaitOnChan()

	// First, everything is good at start
	assertRequest(t, router, "GET", "/healthcheck", "", 200, string(successResponse))

	// Testing things when health goes bad
	controlled.willPass = false
	assertRequest(t, router, "GET", "/healthcheck", "", 503, string(failureResponse))
	assertRequest(t, router, "GET", "/healthcheck", "", 503, string(failureResponse))
	assertRequest(t, router, "GET", "/healthcheck", "", 503, string(failureResponse))

	wg.Wait()
	assert.Error(t, ErrHealthcheckFailed, errNotification)

	// Testing things are going back to normal
	goWaitOnChan()

	controlled.willPass = true
	assertRequest(t, router, "GET", "/healthcheck", "", 200, string(successResponse))

	wg.Wait()
	assert.NoError(t, errNotification)

	// Testing finally that health is going back to bad
	goWaitOnChan()

	controlled.willPass = false
	assertRequest(t, router, "GET", "/healthcheck", "", 503, string(failureResponse))
	assertRequest(t, router, "GET", "/healthcheck", "", 503, string(failureResponse))
	assertRequest(t, router, "GET", "/healthcheck", "", 503, string(failureResponse))

	wg.Wait()
	assert.Error(t, ErrHealthcheckFailed, errNotification)
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
