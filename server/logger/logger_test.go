package logger

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	// Create a temporary log file for testing
	tempFile, err := ioutil.TempFile("", "test-log")
	if err != nil {
		t.Fatalf("Failed to create temporary log file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Create a new logger with the temporary log file
	log, err := NewLogger(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer log.Close()

	// Write a log message to the logger
	log.Log("test message")

	// Read the contents of the log file and verify that it contains the log message
	contents, err := ioutil.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}
	if string(contents) != "test message\n" {
		t.Errorf("Log file contents (%q) do not match expected contents (%q)", contents, "test message\n")
	}
}

func TestLoggerError(t *testing.T) {
	// Attempt to create a logger with an invalid filename
	log, err := NewLogger("/invalid/filename")
	if err == nil {
		log.Close()
		t.Errorf("Expected error when creating logger with invalid filename, but got nil")
	}
}
