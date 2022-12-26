package gin_healthcheck

import (
	"github.com/gin-gonic/gin"
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
	config := DefaultConfig()
	err := New(router, config)
	if err != nil {
		t.Fatal(err)
	}

	healthRoute := router.Routes()[0]

	if healthRoute.Path != config.HealthPath {
		t.Errorf("Healthcheck path route is not equal to config path")
	}
}

func TestHealthcheckResponse(t *testing.T) {
	router := gin.Default()
	config := DefaultConfig()
	New(router, config)

	assertRequest(t, router, "GET", config.HealthPath, "", 200, "")
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
