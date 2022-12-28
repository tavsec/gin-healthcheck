package config

type Config struct {
	HealthPath  string
	Method      string
	StatusOK    int
	StatusNotOK int
}

func DefaultConfig() Config {
	return Config{
		HealthPath:  "/healthz",
		Method:      "GET",
		StatusOK:    200,
		StatusNotOK: 503,
	}
}
