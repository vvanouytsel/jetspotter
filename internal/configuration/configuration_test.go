package configuration

import (
	"os"
	"testing"
)

// TestScanRangeDefaultsToMaxRangeKilometers tests that MaxScanRangeKilometers defaults to MaxRangeKilometers
// when MAX_SCAN_RANGE_KILOMETERS is not set
func TestScanRangeDefaultsToMaxRangeKilometers(t *testing.T) {
	// Save current environment
	oldMaxRange := os.Getenv("MAX_RANGE_KILOMETERS")
	oldMaxScanRange := os.Getenv("MAX_SCAN_RANGE_KILOMETERS")
	defer func() {
		os.Setenv("MAX_RANGE_KILOMETERS", oldMaxRange)
		os.Setenv("MAX_SCAN_RANGE_KILOMETERS", oldMaxScanRange)
	}()

	// Set MAX_RANGE_KILOMETERS, but not MAX_SCAN_RANGE_KILOMETERS
	os.Setenv("MAX_RANGE_KILOMETERS", "50")
	os.Unsetenv("MAX_SCAN_RANGE_KILOMETERS")

	config, err := GetConfig()
	if err != nil {
		t.Fatalf("Failed to get config: %v", err)
	}

	if config.MaxRangeKilometers != 50 {
		t.Fatalf("expected MaxRangeKilometers to be 50, got %d", config.MaxRangeKilometers)
	}

	if config.MaxScanRangeKilometers != config.MaxRangeKilometers {
		t.Fatalf("expected MaxScanRangeKilometers to equal MaxRangeKilometers (%d), got %d",
			config.MaxRangeKilometers, config.MaxScanRangeKilometers)
	}
}

// TestScanRangeCanBeDifferentFromMaxRange tests that MaxScanRangeKilometers can be set to a different
// value than MaxRangeKilometers
func TestScanRangeCanBeDifferentFromMaxRange(t *testing.T) {
	// Save current environment
	oldMaxRange := os.Getenv("MAX_RANGE_KILOMETERS")
	oldMaxScanRange := os.Getenv("MAX_SCAN_RANGE_KILOMETERS")
	defer func() {
		os.Setenv("MAX_RANGE_KILOMETERS", oldMaxRange)
		os.Setenv("MAX_SCAN_RANGE_KILOMETERS", oldMaxScanRange)
	}()

	// Set both environment variables to different values
	os.Setenv("MAX_RANGE_KILOMETERS", "30")
	os.Setenv("MAX_SCAN_RANGE_KILOMETERS", "100")

	config, err := GetConfig()
	if err != nil {
		t.Fatalf("Failed to get config: %v", err)
	}

	if config.MaxRangeKilometers != 30 {
		t.Fatalf("expected MaxRangeKilometers to be 30, got %d", config.MaxRangeKilometers)
	}

	if config.MaxScanRangeKilometers != 100 {
		t.Fatalf("expected MaxScanRangeKilometers to be 100, got %d", config.MaxScanRangeKilometers)
	}
}
