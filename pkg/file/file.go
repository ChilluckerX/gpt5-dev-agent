package file

import (
	"encoding/json"
	"os"
)

// WriteJSONFile writes data to a JSON file with proper formatting
func WriteJSONFile(filename string, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(filename, jsonData, 0644)
}

// ReadJSONFile reads data from a JSON file
func ReadJSONFile(filename string, data interface{}) error {
	fileData, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	
	return json.Unmarshal(fileData, data)
}