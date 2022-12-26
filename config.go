package gin_healthcheck

type Config struct {
	HealthPath string
}

func DefaultConfig() Config {
	return Config{
		HealthPath: "/healthz",
	}
}
