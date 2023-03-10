package main

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/enriquenc/orderer-map-client-server-go/shared"
	types "github.com/enriquenc/orderer-map-client-server-go/shared"
	"github.com/stretchr/testify/require"
)

func TestParseDataFromFile(t *testing.T) {
	// create a test file with some test data
	testFile, err := os.Create("test.json")
	if err != nil {
		t.Errorf("Failed to create test file: %v", err)
	}
	defer testFile.Close()
	defer os.Remove("test.json")

	// create some test data
	testData := []shared.TestDataAction{
		{RequestData: shared.Request{Key: "key1", Value: "value1", Action: shared.AddItem}},
		{RequestData: shared.Request{Key: "key2", Value: "value2", Action: shared.AddItem}},
		{RequestData: shared.Request{Key: "key3", Value: "value3", Action: shared.AddItem}},
	}

	// encode test data and write to test file
	encoder := json.NewEncoder(testFile)
	for _, data := range testData {
		err := encoder.Encode(&data)
		if err != nil {
			t.Errorf("Failed to encode test data: %v", err)
		}
	}

	// reset file offset
	_, err = testFile.Seek(0, 0)
	if err != nil {
		t.Errorf("Failed to reset test file offset: %v", err)
	}

	// call parseDataFromFile function to parse test data from file
	var parsedTestData []shared.TestDataAction
	err = parseDataFromFile(&parsedTestData, "test.json")
	if err != nil {
		t.Errorf("Failed to parse test data from file: %v", err)
	}

	// compare parsed test data with original test data
	for i, data := range testData {
		if data.RequestData.Key != parsedTestData[i].RequestData.Key {
			t.Errorf("Key mismatch for test data item %d: expected %s, got %s", i, data.RequestData.Key, parsedTestData[i].RequestData.Key)
		}
		if data.RequestData.Value != parsedTestData[i].RequestData.Value {
			t.Errorf("Value mismatch for test data item %d: expected %s, got %s", i, data.RequestData.Value, parsedTestData[i].RequestData.Value)
		}
		if data.RequestData.Action != parsedTestData[i].RequestData.Action {
			t.Errorf("Action mismatch for test data item %d: expected %s, got %s", i, data.RequestData.Action, parsedTestData[i].RequestData.Action)
		}
	}
}

func TestMain_InvalidJSONFile(t *testing.T) {
	// Create temporary file with invalid JSON
	invalidJSON := []byte(`[{"requestData": {"key": "foo", "value": "bar", "action": "add"}}, {"requestData": {"key": "baz", "value": "qux", "action": "remove"}}, {"requestData": {"key": "hello", "value": "world", "action": "get"}}}`)
	file, err := os.Create("test.json")
	if err != nil {
		t.Errorf("Failed to create test file: %v", err)
	}
	defer file.Close()
	defer os.Remove("test.json")

	require.NoError(t, err)
	_, err = file.Write(invalidJSON)
	require.NoError(t, err)

	// Call function under test
	_, err = parseRequestData(file.Name(), "", "", "")

	// Assert error is returned with expected message
	require.Error(t, err)
	require.Contains(t, err.Error(), "Failed to decode test data")
}

func TestParseRequestData_ValidArgs(t *testing.T) {
	// Set up test case with valid command-line arguments
	fileName := ""
	action := "add"
	key := "foo"
	value := "bar"

	// Call parseRequestData
	testData, err := parseRequestData(fileName, action, key, value)

	// Assert that no error occurred
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Assert that testData contains the expected data
	expectedTestData := []types.TestDataAction{{
		RequestData: types.Request{
			Action: action,
			Key:    key,
			Value:  value,
		},
	}}
	if !reflect.DeepEqual(testData, expectedTestData) {
		t.Errorf("Unexpected testData: got %v, want %v", testData, expectedTestData)
	}
}

func TestParseRequestDataWithInvalidArgs(t *testing.T) {
	// Set up invalid command-line arguments
	fileName := "testFile.json"
	action := "invalid_action"
	key := ""
	value := ""

	// Call parseRequestData with invalid arguments
	_, err := parseRequestData(fileName, action, key, value)

	// Check if an error was returned
	if err == nil {
		t.Errorf("Expected an error to be returned for invalid arguments, but got nil")
	}
}

func TestIsValidActionWithValidActions(t *testing.T) {
	validActions := []string{types.AddItem, types.GetItem, types.RemoveItem, types.GetAll}

	for _, action := range validActions {
		if !isValidAction(action) {
			t.Errorf("isValidAction returned false for a valid action: %s", action)
		}
	}
}

func TestIsValidActionWithInvalidActions(t *testing.T) {
	invalidActions := []string{"invalid", "ADD", "get_item", "", " ", "1234"}

	for _, action := range invalidActions {
		if isValidAction(action) {
			t.Errorf("isValidAction returned true for an invalid action: %s", action)
		}
	}
}
