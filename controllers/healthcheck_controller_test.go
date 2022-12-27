package controllers

import (
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
	router.GET("/healthcheck", HealthcheckController)
	assertRequest(t, router, "GET", "/healthcheck", "", 200, "")

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
