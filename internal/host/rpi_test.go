//go:build rpi

package host

import (
	"strings"
	"testing"
)

func TestParseMemoryUsage(t *testing.T) {
	usage, err := parseMemoryUsage(strings.NewReader("MemTotal:       1000 kB\nMemAvailable:    250 kB\n"))
	if err != nil {
		t.Fatalf("parseMemoryUsage returned an error: %v", err)
	}
	if usage != 75 {
		t.Fatalf("parseMemoryUsage = %v, want 75", usage)
	}
}

func TestParseMemoryUsageMissingTotal(t *testing.T) {
	_, err := parseMemoryUsage(strings.NewReader("MemAvailable: 250 kB\n"))
	if err == nil {
		t.Fatal("parseMemoryUsage returned nil error for missing total")
	}
}

func TestParseCPUTemperature(t *testing.T) {
	temperature, err := parseCPUTemperature([]byte("52500\n"))
	if err != nil {
		t.Fatalf("parseCPUTemperature returned an error: %v", err)
	}
	if temperature != 52.5 {
		t.Fatalf("parseCPUTemperature = %v, want 52.5", temperature)
	}
}

func TestParseCPUTemperatureInvalidInput(t *testing.T) {
	_, err := parseCPUTemperature([]byte("not a temperature"))
	if err == nil {
		t.Fatal("parseCPUTemperature returned nil error for invalid input")
	}
}
