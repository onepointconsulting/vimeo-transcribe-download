package vtt_parser

import (
	"os"
	"testing"
)

func TestExtractText(t *testing.T) {
	// Create a temporary test file
	content := `WEBVTT
X-TIMESTAMP-MAP=LOCAL:00:00:00.000,MPEGTS:90000

1
00:00:01.575 --> 00:00:03.635
Why are we not innovating at the moment

2
00:00:05.195 --> 00:00:07.745
we're in the middle of a crisis like none we've ever seen.

3
00:00:08.275 --> 00:00:10.065
Ordinarily, we've shown ourselves

4
00:00:10.065 --> 00:00:11.625
to be very innovative during a crisis.`

	tmpfile, err := os.CreateTemp("", "test-*.vtt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Test the function
	got, err := ExtractText(tmpfile.Name())
	if err != nil {
		t.Fatalf("ExtractText failed: %v", err)
	}

	want := "Why are we not innovating at the moment we're in the middle of a crisis like none we've ever seen. Ordinarily, we've shown ourselves to be very innovative during a crisis."
	if got != want {
		t.Errorf("ExtractText() = %q, want %q", got, want)
	}
}
