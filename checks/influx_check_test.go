package checks

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func TestInfluxCheck_Name(t *testing.T) {
	check := NewInfluxV2Check(2, nil)
	if check.Name() != "influxdb" {
		t.Errorf("Expected InfluxV2Check.Name to return 'influxdb', got '%s'", check.Name())
	}
}

func TestInfluxCheck_Pass(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))

	influxClient := influxdb2.NewClient(server.URL, "")

	check := NewInfluxV2Check(1, influxClient)
	if !check.Pass() {
		t.Error("InfluxCheck.Pass() returned false, want true")
	}
}

func TestInfluxCheck_Timeout(t *testing.T) {
	http.HandleFunc("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
	}))

	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			t.Errorf("server.ListenAndServe() error = %v", err)
			return
		}
	}()

	influxClient := influxdb2.NewClient("localhost:8080", "")
	check := NewInfluxV2Check(1, influxClient)
	if check.Pass() {
		t.Error("InfluxCheck.Pass() returned true, want false")
	}
}

func TestInfluxCheck_Fail(t *testing.T) {
	check := NewInfluxV2Check(1, nil)
	if check.Pass() {
		t.Error("InfluxCheck.Pass() returned true, want false")
	}
}
