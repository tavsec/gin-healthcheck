package config

type Config struct {
	HealthPath  string
	Method      string
	StatusOK    int
	StatusNotOK int

	FailureNotification struct {
		Threshold uint32
		Chan      chan error
	}
}

func DefaultConfig() Config {
	return Config{
		HealthPath:  "/healthz",
		Method:      "GET",
		StatusOK:    200,
		StatusNotOK: 503,
		FailureNotification: struct {
			Threshold uint32
			Chan      chan error
		}{
			Threshold: 1,
		},
	}
}
