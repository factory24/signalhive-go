package connectivity

import "testing"

func TestEnvReadsDottedThenUnderscored(t *testing.T) {
	t.Setenv("CONNECTIVITY_BASE_URL", "https://underscored.example/v1")
	if got := env("BASE_URL"); got != "https://underscored.example/v1" {
		t.Fatalf("underscored name not read: %q", got)
	}
	t.Setenv("CONNECTIVITY.BASE_URL", "https://dotted.example/v1")
	if got := env("BASE_URL"); got != "https://dotted.example/v1" {
		t.Fatalf("dotted name does not win: %q", got)
	}
	t.Setenv("CONNECTIVITY_CELLULAR_API_KEY", "k")
	if got := env("CELLULAR.API_KEY"); got != "k" {
		t.Fatalf("nested key not underscored: %q", got)
	}
}

