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

// TestOverheadPredictionDefaults verifies that overhead config defaults to
// disabled with the documented default values.
func TestOverheadPredictionDefaults(t *testing.T) {
	allKeys := []string{
		"OVERHEAD_PREDICTION_ENABLED", "OVERHEAD_RADIUS_KILOMETERS",
		"OVERHEAD_LOOK_AHEAD_MINUTES", "OVERHEAD_CONFIRM_SECONDS",
		"OVERHEAD_SUBLOOP_POLL_SECONDS", "OVERHEAD_INBOUND_MARGIN_DEGREES",
		"MAX_RANGE_KILOMETERS", "MAX_SCAN_RANGE_KILOMETERS",
		"MAX_ALTITUDE_FEET", "FETCH_INTERVAL",
		"LOCATION_LATITUDE", "LOCATION_LONGITUDE",
	}
	for _, key := range allKeys {
		old := os.Getenv(key)
		defer os.Setenv(key, old)
		os.Unsetenv(key)
	}

	config, err := GetConfig()
	if err != nil {
		t.Fatalf("Failed to get config: %v", err)
	}
	if config.OverheadPredictionEnabled {
		t.Fatal("expected OverheadPredictionEnabled=false by default")
	}
	if config.OverheadRadiusKilometers != 3 {
		t.Fatalf("expected OverheadRadiusKilometers=3, got %d", config.OverheadRadiusKilometers)
	}
	if config.OverheadLookAheadMinutes != 10 {
		t.Fatalf("expected OverheadLookAheadMinutes=10, got %d", config.OverheadLookAheadMinutes)
	}
	if config.OverheadConfirmSeconds != 5 {
		t.Fatalf("expected OverheadConfirmSeconds=5, got %d", config.OverheadConfirmSeconds)
	}
	if config.OverheadSubloopPollSeconds != 5 {
		t.Fatalf("expected OverheadSubloopPollSeconds=5, got %d", config.OverheadSubloopPollSeconds)
	}
	if config.OverheadInboundMarginDegrees != 10 {
		t.Fatalf("expected OverheadInboundMarginDegrees=10, got %d", config.OverheadInboundMarginDegrees)
	}
}

// TestOverheadPredictionEnabled verifies that overhead config can be set via
// environment variables.
func TestOverheadPredictionEnabled(t *testing.T) {
	allKeys := []string{
		"OVERHEAD_PREDICTION_ENABLED", "OVERHEAD_RADIUS_KILOMETERS",
		"OVERHEAD_LOOK_AHEAD_MINUTES",
		"MAX_RANGE_KILOMETERS", "MAX_SCAN_RANGE_KILOMETERS",
		"MAX_ALTITUDE_FEET", "FETCH_INTERVAL",
		"LOCATION_LATITUDE", "LOCATION_LONGITUDE",
	}
	for _, key := range allKeys {
		old := os.Getenv(key)
		defer os.Setenv(key, old)
		os.Unsetenv(key)
	}
	os.Setenv("OVERHEAD_PREDICTION_ENABLED", "true")
	os.Setenv("OVERHEAD_RADIUS_KILOMETERS", "8")
	os.Setenv("OVERHEAD_LOOK_AHEAD_MINUTES", "20")

	config, err := GetConfig()
	if err != nil {
		t.Fatalf("Failed to get config: %v", err)
	}
	if !config.OverheadPredictionEnabled {
		t.Fatal("expected OverheadPredictionEnabled=true")
	}
	if config.OverheadRadiusKilometers != 8 {
		t.Fatalf("expected OverheadRadiusKilometers=8, got %d", config.OverheadRadiusKilometers)
	}
	if config.OverheadLookAheadMinutes != 20 {
		t.Fatalf("expected OverheadLookAheadMinutes=20, got %d", config.OverheadLookAheadMinutes)
	}
}
