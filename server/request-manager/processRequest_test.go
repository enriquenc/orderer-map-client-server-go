package requestmanager

import (
	// other imports
	"io/ioutil"
	"os"
	"server/logger"
	"testing"
	"time"

	types "github.com/enriquenc/orderer-map-client-server-go/shared"
	// other imports
)

func TestProcessRequests_GetItem(t *testing.T) {
	// Create a temporary file for the logger
	file, err := ioutil.TempFile("", "logger_test")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer os.Remove(file.Name())
	defer file.Close()

	// Create a logger
	myLogger, err := logger.NewLogger(file.Name())
	if err != nil {
		t.Fatalf("Error creating logger: %v", err)
	}
	defer myLogger.Close()

	// Create a channel of requests
	reqs := make(chan types.Request)

	// Call ProcessRequests in a separate goroutine
	go ProcessRequests(reqs, myLogger)

	// Add an item
	reqs <- types.Request{
		Action: types.AddItem,
		Key:    "foo",
		Value:  "bar",
	}

	// Get the item
	reqs <- types.Request{
		Action: types.GetItem,
		Key:    "foo",
	}

	// Close the channel
	close(reqs)

	// Wait for ProcessRequests to finish
	time.Sleep(time.Millisecond * 100)

	// Read the logger output from the file
	fileContent, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}

	// Verify that the logger output contains the expected message
	expected := "[add] Added key foo with value bar\n[get] Got key foo with value bar\n"
	if string(fileContent) != expected {
		t.Errorf("Expected logger output %q, but got %q", expected, string(fileContent))
	}
}

func TestProcessRequests_RemoveItem_NotExist(t *testing.T) {
	// Create a temporary file for the logger
	file, err := ioutil.TempFile("", "logger_test")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer os.Remove(file.Name())
	defer file.Close()

	// Create a logger
	myLogger, err := logger.NewLogger(file.Name())
	if err != nil {
		t.Fatalf("Error creating logger: %v", err)
	}
	defer myLogger.Close()

	// Create a channel of requests
	reqs := make(chan types.Request)

	// Call ProcessRequests in a separate goroutine
	go ProcessRequests(reqs, myLogger)

	// Remove a non-existent item
	reqs <- types.Request{
		Action: types.RemoveItem,
		Key:    "nonexistent",
	}

	// Close the channel
	close(reqs)

	// Wait for ProcessRequests to finish
	time.Sleep(time.Millisecond * 100)

	// Read the logger output from the file
	fileContent, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}

	// Verify that the logger output contains the expected message
	expected := "[remove] key nonexistent doesn't exist\n"
	if string(fileContent) != expected {
		t.Errorf("Expected logger output %q, but got %q", expected, string(fileContent))
	}
}

func TestProcessRequests_AddItem_Duplicate(t *testing.T) {
	// Create a temporary file for the logger
	file, err := ioutil.TempFile("", "logger_test")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer os.Remove(file.Name())
	defer file.Close()

	// Create a logger
	myLogger, err := logger.NewLogger(file.Name())
	if err != nil {
		t.Fatalf("Error creating logger: %v", err)
	}
	defer myLogger.Close()

	// Create a channel of requests
	reqs := make(chan types.Request)

	// Call ProcessRequests in a separate goroutine
	go ProcessRequests(reqs, myLogger)

	// Add an item
	reqs <- types.Request{
		Action: types.AddItem,
		Key:    "foo",
		Value:  "bar",
	}

	// Add the same item again
	reqs <- types.Request{
		Action: types.AddItem,
		Key:    "foo",
		Value:  "baz",
	}

	// Close the channel
	close(reqs)

	// Wait for ProcessRequests to finish
	time.Sleep(time.Millisecond * 100)

	// Read the logger output from the file
	fileContent, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}

	// Verify that the logger output contains the expected message
	expected := "[add] Added key foo with value bar\n[add] Added key foo with value baz\n"
	if string(fileContent) != expected {
		t.Errorf("Expected logger output %q, but got %q", expected, string(fileContent))
	}
}

func TestProcessRequests_GetAll(t *testing.T) {
	// Create a temporary file for the logger
	file, err := ioutil.TempFile("", "logger_test")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer os.Remove(file.Name())
	defer file.Close()

	// Create a logger
	myLogger, err := logger.NewLogger(file.Name())
	if err != nil {
		t.Fatalf("Error creating logger: %v", err)
	}
	defer myLogger.Close()

	// Create a channel of requests
	reqs := make(chan types.Request)

	// Call ProcessRequests in a separate goroutine
	go ProcessRequests(reqs, myLogger)

	// Add some items
	reqs <- types.Request{
		Action: types.AddItem,
		Key:    "foo",
		Value:  "bar",
	}
	reqs <- types.Request{
		Action: types.AddItem,
		Key:    "baz",
		Value:  "qux",
	}

	// Get all the items
	reqs <- types.Request{
		Action: types.GetAll,
	}

	// Close the channel
	close(reqs)

	// Wait for ProcessRequests to finish
	time.Sleep(time.Millisecond * 100)

	// Read the logger output from the file
	fileContent, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}

	// Verify that the logger output contains the expected message
	expected := "[add] Added key foo with value bar\n[add] Added key baz with value qux\n[getAll] All values [\"foo=bar\",\"baz=qux\"]\n"
	if string(fileContent) != expected {
		t.Errorf("Expected logger output %q, but got %q", expected, string(fileContent))
	}
}
