package main

import "testing"

func TestParseSeedFlags_Defaults(t *testing.T) {
	force, err := parseSeedFlags([]string{})
	if err != nil {
		t.Fatalf("parseSeedFlags() error = %v", err)
	}
	if force {
		t.Fatalf("force = true, want false")
	}
}

func TestParseSeedFlags_ForceTrue(t *testing.T) {
	force, err := parseSeedFlags([]string{"--force"})
	if err != nil {
		t.Fatalf("parseSeedFlags() error = %v", err)
	}
	if !force {
		t.Fatalf("force = false, want true")
	}
}

func TestParseSeedFlags_InvalidFlag(t *testing.T) {
	_, err := parseSeedFlags([]string{"--unknown"})
	if err == nil {
		t.Fatalf("parseSeedFlags() error = nil, want non-nil")
	}
}

func TestDefaultSeedOptions(t *testing.T) {
	opts := defaultSeedOptions(true)

	if opts.Users != 24 || opts.Categories != 8 || opts.Products != 64 || opts.Orders != 50 {
		t.Fatalf("defaultSeedOptions() values are unexpected: %+v", opts)
	}
	if !opts.Force {
		t.Fatalf("defaultSeedOptions(true).Force = false, want true")
	}
}
