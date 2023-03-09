package main

import (
	"encoding/json"
	"fmt"
	"os"

	types "github.com/enriquenc/orderer-map-client-server-go/shared"
)

func parseRequestData(fileName, action, key, value string) ([]types.TestDataAction, error) {
	var testData []types.TestDataAction

	if fileName != "" {
		err := parseDataFromFile(&testData, fileName)
		if err != nil {
			return nil, fmt.Errorf("Error in file passing: %v", err)
		}
		println(testData[0].RequestData.Action)
	} else {
		err := validateCommandLineArguments(action, key, value)
		if err != nil {
			return nil, fmt.Errorf("Error in arguments parsing: %v", err)
		}
		testData = append(testData, types.TestDataAction{RequestData: types.Request{Key: key, Value: value, Action: action}})
	}

	return testData, nil
}

func parseDataFromFile(testData *[]types.TestDataAction, fileName string) error {
	file, err := os.Open(fileName)

	if err != nil {
		return fmt.Errorf("Failed to open %s file: %v", fileName, err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	for decoder.More() {
		var testDataAction types.TestDataAction
		err := decoder.Decode(&testDataAction)
		if err != nil {
			return fmt.Errorf("Failed to decode test data: %v", err)
		}
		// println(testDataAction.RequestData.Action)
		*testData = append(*testData, testDataAction)
	}

	return nil
}

func validateCommandLineArguments(action, key, value string) error {
	// Validate command line arguments
	if action == "" || !isValidAction(action) {
		return fmt.Errorf("Invalid action. Must be one of: %s, %s, %s, %s.", types.AddItem, types.GetItem, types.RemoveItem, types.GetAll)
	}
	if (action == types.AddItem || action == types.GetItem || action == types.RemoveItem) && key == "" {
		return fmt.Errorf("Key is required for %s, %s, and %s actions.", types.AddItem, types.GetItem, types.RemoveItem)
	}
	if (action == types.AddItem) && value == "" {
		return fmt.Errorf("Value is required for %s action.", types.AddItem)
	}
	return nil
}

func isValidAction(action string) bool {
	switch action {
	case types.AddItem, types.RemoveItem, types.GetItem, types.GetAll:
		return true
	default:
		return false
	}
}
