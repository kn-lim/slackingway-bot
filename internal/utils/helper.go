package utils

import (
	"encoding/json"
	"errors"
	"log"
	"reflect"
)

// PrintStructFields returns a struct as a readable JSON string
func PrintStructFields(s interface{}) (string, error) {
	v := reflect.ValueOf(s)

	// Check if the input is a struct
	if v.Kind() != reflect.Struct {
		log.Println("Provided value is not a struct.")
		return "", errors.New("provided value is not a struct")
	}

	// Marshal the struct into a compact JSON string
	jsonData, err := json.Marshal(s)
	if err != nil {
		log.Printf("Error marshalling struct to JSON: %v", err)
		return "", err
	}

	return string(jsonData), nil
}
