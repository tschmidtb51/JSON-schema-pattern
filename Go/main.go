//go:build go1.18
// +build go1.18

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
)

// JSONSchema represents the structure of the schema file
type JSONSchema struct {
	Properties struct {
		TestValues struct {
			Items struct {
				Properties struct {
					Value struct {
						Pattern string `json:"pattern"`
					} `json:"value"`
				} `json:"properties"`
			} `json:"items"`
		} `json:"test_values"`
	} `json:"properties"`
}

// TestValue represents the structure of each test value in the JSON file
type TestValue struct {
	Value     string `json:"value"`
	Assertion bool   `json:"assertion"`
}

// TestFile represents the structure of the JSON file to be tested
type TestFile struct {
	TestValues []TestValue `json:"test_values"`
}

func main() {
	// Get file paths from command line arguments
	if len(os.Args) < 3 {
		fmt.Println("Usage: <executable> json/file/path.json schema/file/path.json")
		os.Exit(1)
	}

	filePath := os.Args[1]
	schemaPath := os.Args[2]

	// Read and parse JSON file
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading JSON file: %v\n", err)
		os.Exit(1)
	}

	var testFile TestFile
	err = json.Unmarshal(fileData, &testFile)
	if err != nil {
		fmt.Printf("Error parsing JSON file: %v\n", err)
		os.Exit(1)
	}

	// Read and parse schema file
	schemaData, err := os.ReadFile(schemaPath)
	if err != nil {
		fmt.Printf("Error reading schema file: %v\n", err)
		os.Exit(1)
	}

	var schema JSONSchema
	err = json.Unmarshal(schemaData, &schema)
	if err != nil {
		fmt.Printf("Error parsing schema file: %v\n", err)
		os.Exit(1)
	}

	// Extract pattern from schema
	pattern := schema.Properties.TestValues.Items.Properties.Value.Pattern
	r, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Printf("Error compiling regex pattern: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Current regex to test:\n", pattern)

	failed := false

	// Test each value in the JSON file
	for _, element := range testFile.TestValues {
		result := r.MatchString(element.Value)
		failed = failed || (result != element.Assertion)
		fmt.Println(result, "\t", element)
	}

	fmt.Println("\nOverall result:\t", func() string {
		if failed {
			return "Error occurred"
		}
		return "Test successful"
	}())

	if failed {
		os.Exit(1)
	}
}
