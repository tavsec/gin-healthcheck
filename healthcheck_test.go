package gin_healthcheck

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/tavsec/gin-healthcheck/checks"
	config2 "github.com/tavsec/gin-healthcheck/config"
	"github.com/tavsec/gin-healthcheck/controllers"
)

var (
	res *httptest.ResponseRecorder
)

func init() {
	gin.SetMode(gin.TestMode)
}

type SucceedingCheck struct{}

func (c SucceedingCheck) Pass() bool {
	return true
}
func (c SucceedingCheck) Name() string {
	return "Succeeding Check"
}

var _ checks.Check = (*SucceedingCheck)(nil)

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
	assert.NoError(t, err)
	assertRequest(t, router, "GET", config.HealthPath, "", 200, string(response))
}

func TestHealthcheckSwitchStateBetweenCall(t *testing.T) {
	router := gin.Default()
	config := config2.DefaultConfig()
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPing().WillReturnError(nil)

	c := []checks.Check{checks.SqlCheck{Sql: db}}
	New(router, config, c)

	response, err := json.Marshal([]controllers.CheckStatus{{
		Name: "mysql",
		Pass: true,
	}})
	assert.NoError(t, err)
	assertRequest(t, router, "GET", config.HealthPath, "", 200, string(response))

	mock.ExpectPing().WillReturnError(driver.ErrBadConn)

	response, err = json.Marshal([]controllers.CheckStatus{{
		Name: "mysql",
		Pass: false,
	}})
	assert.NoError(t, err)
	assertRequest(t, router, "GET", config.HealthPath, "", 503, string(response))
}

func TestHealthcheckContext(t *testing.T) {
	router := gin.Default()
	config := config2.DefaultConfig()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := []checks.Check{checks.NewContextCheck(ctx), checks.NewContextCheck(context.Background(), "test"), SucceedingCheck{}}
	New(router, config, c)

	response, err := json.Marshal([]controllers.CheckStatus{
		{
			Name: "github.com/tavsec/gin-healthcheck.TestHealthcheckContext",
			Pass: true,
		},
		{
			Name: "test",
			Pass: true,
		},
		{
			Name: "Succeeding Check",
			Pass: true,
		},
	})
	assert.NoError(t, err)
	assertRequest(t, router, "GET", config.HealthPath, "", 200, string(response))

	cancel()
	<-ctx.Done()

	// We need to give time to the goroutine to get scheduled before checking the status
	time.Sleep(1 * time.Millisecond)

	response, err = json.Marshal([]controllers.CheckStatus{
		{
			Name: "github.com/tavsec/gin-healthcheck.TestHealthcheckContext",
			Pass: false,
		},
		{
			Name: "test",
			Pass: true,
		},
		{
			Name: "Succeeding Check",
			Pass: true,
		},
	})
	assert.NoError(t, err)
	assertRequest(t, router, "GET", config.HealthPath, "", 503, string(response))
}

func TestNotification(t *testing.T) {
	router := gin.Default()

	config := config2.DefaultConfig()
	config.FailureNotification.Chan = make(chan error, 1)
	defer close(config.FailureNotification.Chan)
	config.FailureNotification.Threshold = 3

	var errNotification error
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		errNotification = <-config.FailureNotification.Chan
	}()

	ctx, cancel := context.WithCancel(context.Background())

	c := []checks.Check{checks.NewContextCheck(ctx)}
	New(router, config, c)

	response, err := json.Marshal([]controllers.CheckStatus{
		{
			Name: "github.com/tavsec/gin-healthcheck.TestNotification",
			Pass: false,
		},
	})
	assert.NoError(t, err)

	cancel()
	<-ctx.Done()

	// We need to give time to the goroutine to get scheduled before checking the status
	time.Sleep(1 * time.Millisecond)

	assertRequest(t, router, "GET", config.HealthPath, "", 503, string(response))
	assertRequest(t, router, "GET", config.HealthPath, "", 503, string(response))
	assertRequest(t, router, "GET", config.HealthPath, "", 503, string(response))

	wg.Wait()
	assert.Equal(t, controllers.ErrHealthcheckFailed, errNotification)

	assertRequest(t, router, "GET", config.HealthPath, "", 503, string(response))
	errNotification = <-config.FailureNotification.Chan
	assert.Equal(t, controllers.ErrHealthcheckFailed, errNotification)
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
