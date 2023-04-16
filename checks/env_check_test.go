package checks

import (
	"os"
	"testing"
)

func TestEnvCheck_Pass(t *testing.T) {
	os.Setenv("TEST_VAR", "12345")
	defer os.Unsetenv("TEST_VAR")

	tests := []struct {
		name        string
		envVariable string
		regex       string
		want        bool
	}{
		{
			name:        "variable exists, no regex",
			envVariable: "TEST_VAR",
			regex:       "",
			want:        true,
		},
		{
			name:        "variable does not exist, no regex",
			envVariable: "NOT_EXIST",
			regex:       "",
			want:        false,
		},
		{
			name:        "variable exists, regex match",
			envVariable: "TEST_VAR",
			regex:       "^[0-9]+$",
			want:        true,
		},
		{
			name:        "variable exists, regex does not match",
			envVariable: "TEST_VAR",
			regex:       "^[a-z]+$",
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := EnvCheck{
				EnvVariable: tt.envVariable,
				Regex:       tt.regex,
			}

			if got := e.Pass(); got != tt.want {
				t.Errorf("EnvCheck.Pass() = %v, want %v", got, tt.want)
			}
		})
	}
}
