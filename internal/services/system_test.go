package services

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

func TestParseDiskUsage(t *testing.T) {
	output := []byte("Filesystem 1024-blocks Used Available Capacity Mounted on\n/dev/root 1000 250 750 25% /\n")
	usage, err := parseDiskUsage(output)
	if err != nil {
		t.Fatalf("parseDiskUsage returned an error: %v", err)
	}
	if usage != 25 {
		t.Fatalf("parseDiskUsage = %v, want 25", usage)
	}
}

func TestParseDiskUsageInvalidOutput(t *testing.T) {
	_, err := parseDiskUsage([]byte("invalid output\n"))
	if err == nil {
		t.Fatal("parseDiskUsage returned nil error for invalid output")
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

func TestParseAvailableUpdates(t *testing.T) {
	output := []byte("Listing...\npackage-one/stable 1.0 amd64 [upgradable from: 0.9]\n\npackage-two/stable 2.0 amd64 [upgradable from: 1.9]\n")
	if updates := parseAvailableUpdates(output); updates != 2 {
		t.Fatalf("parseAvailableUpdates = %d, want 2", updates)
	}
}

func TestParseAvailableUpdatesEmpty(t *testing.T) {
	if updates := parseAvailableUpdates([]byte("Listing...\n\n")); updates != 0 {
		t.Fatalf("parseAvailableUpdates = %d, want 0", updates)
	}
}
