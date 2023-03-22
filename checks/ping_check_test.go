package checks

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewPingCheckWithDefaultMethod(t *testing.T) {
	url := "https://example.com"
	timeout := 1000
	body := bytes.NewBufferString("request body")
	headers := map[string]string{
		"Authorization": "Bearer xxx",
	}

	check := NewPingCheck(url, "", timeout, body, headers)

	if check.Method != "GET" {
		t.Errorf("NewPingCheck() Method = %s, want GET", check.Method)
	}
}

func TestNewPingCheckWithDefaultTimeout(t *testing.T) {
	url := "https://example.com"
	method := "POST"
	body := bytes.NewBufferString("request body")
	headers := map[string]string{
		"Authorization": "Bearer xxx",
	}

	check := NewPingCheck(url, method, 0, body, headers)

	if check.Timeout != 500 {
		t.Errorf("NewPingCheck() Timeout = %d, want 500", check.Timeout)
	}
}

func TestPingCheckPassWithSuccessfulRequest(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	check := NewPingCheck(server.URL, "GET", 1000, nil, nil)

	if !check.Pass() {
		t.Error("PingCheck.Pass() returned false, want true")
	}
}

func TestPingCheckPassWithFailedRequest(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	check := NewPingCheck(server.URL, "GET", 1000, nil, nil)

	if check.Pass() {
		t.Error("PingCheck.Pass() returned true, want false")
	}
}

func TestPingCheckName(t *testing.T) {
	url := "https://example.com"
	check := NewPingCheck(url, "GET", 1000, nil, nil)

	want := "ping-https://example.com"
	if name := check.Name(); name != want {
		t.Errorf("PingCheck.Name() = %s, want %s", name, want)
	}
}

func TestNewPingCheckDefaultValues(t *testing.T) {
	url := "https://example.com"
	check := NewPingCheck(url, "", 0, nil, nil)

	if check.Method != "GET" {
		t.Errorf("NewPingCheck() default Method = %s, want GET", check.Method)
	}

	if check.Timeout != 500 {
		t.Errorf("NewPingCheck() default Timeout = %d, want 500", check.Timeout)
	}
}

func TestPingCheck_Pass_InvalidRequest(t *testing.T) {
	url := "http://localhost:8080/test"
	method := "GET"
	check := NewPingCheck(url, method, 3000, nil, nil)

	if check.Pass() {
		t.Error("PingCheck.Pass() returned true, want false")
	}
}

func TestPingCheck_Pass_Timeout(t *testing.T) {
	http.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
	}))

	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			t.Errorf("server.ListenAndServe() error = %v", err)
			return
		}
	}()

	url := "http://localhost:8080/test"
	method := "GET"
	check := NewPingCheck(url, method, 1000, nil, nil)

	if check.Pass() {
		t.Error("PingCheck.Pass() returned true, want false")
	}
}

func TestPingCheck_Name_EmptyURL(t *testing.T) {
	url := ""
	method := "GET"
	check := NewPingCheck(url, method, 1000, nil, nil)

	want := "ping-"
	if name := check.Name(); name != want {
		t.Errorf("PingCheck.Name() = %s, want %s", name, want)
	}
}

func TestNewPingCheck_Headers(t *testing.T) {
	url := "https://example.com"
	method := "GET"
	timeout := 500
	body := bytes.NewBufferString("request body")
	headers := map[string]string{
		"Authorization": "Bearer xxx",
		"Content-Type":  "application/json",
	}

	check := NewPingCheck(url, method, timeout, body, headers)

	if len(check.Headers) != len(headers) {
		t.Errorf("NewPingCheck() Headers length = %d, want %d", len(check.Headers), len(headers))
	}

	for key, value := range headers {
		if check.Headers[key] != value {
			t.Errorf("NewPingCheck() Header %s = %s, want %s", key, check.Headers[key], value)
		}
	}
}

func TestNewPingCheck_NoHeaders(t *testing.T) {
	url := "https://example.com"
	method := "GET"
	timeout := 500
	body := bytes.NewBufferString("request body")

	check := NewPingCheck(url, method, timeout, body, nil)

	if check.Headers != nil {
		t.Errorf("NewPingCheck() Headers = %v, want nil", check.Headers)
	}
}
