package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/iancoleman/orderedmap"
)

type Action struct {
	Command          string
	Key              string
	Value            string
	ExpectedResponse interface{}
}

func main() {
	dataAmount := flag.Int("amount", 1000, "Value to use for the item")
	fileName := flag.String("file", "testdata.json", "Value to use for the item")
	min := flag.Int("min", 0, "min number from range generated keys and values")
	max := flag.Int("max", 20, "max number from range generated keys and values")
	flag.Parse()
	// Define the actions that the client can perform
	actions := []string{"add", "remove", "get", "getAll"}

	if *min < 0 || *max < 0 || *max < *min {
		log.Fatalf("Error min/max values for the range")
	}
	// Write the test data to a file
	file, err := os.Create(*fileName)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	defer file.Close()

	// Generate 1000 random test data
	var testData []Action
	serverMap := orderedmap.New()

	for i := 0; i < *dataAmount; i++ {
		// Generate a random action
		action := actions[rand.Intn(len(actions))]

		// Generate a random key and value
		randKeyValue := rand.Intn(*max-*min) + *min
		key := fmt.Sprintf("key%d", randKeyValue)

		// Add the action and data to the test data slice
		var expectedResponse interface{}
		var value string
		switch action {
		case "add":
			value = fmt.Sprintf("value%d", randKeyValue)
			serverMap.Set(key, value)
		case "remove":
			serverMap.Delete(key)
		case "get":
			val, ok := serverMap.Get(key)
			if ok {
				value = val.(string)
				expectedResponse = value
			} else {
				value = ""
				expectedResponse = nil
			}
		case "getAll":
			// create a new ordered map to hold the expected response
			expectedResponse = orderedmap.New()
			keys := serverMap.Keys()
			for _, k := range keys {
				v, _ := serverMap.Get(k)
				expectedResponse.(*orderedmap.OrderedMap).Set(k, v)
			}
		}

		testData = append(testData, Action{Command: action, Key: key, Value: value, ExpectedResponse: expectedResponse})
	}

	encoder := json.NewEncoder(file)
	for _, data := range testData {
		err = encoder.Encode(data)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Test data generated successfully.")
}
