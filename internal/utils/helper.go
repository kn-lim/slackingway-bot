package utils

import (
	"encoding/json"
	"errors"
	"log"
	"net/url"
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

// PrintPayloadFields prints the fields of a payload as a readable JSON string
func PrintPayloadFields(payload string) (string, error) {
	// Decode the payload
	decodedPayload, err := url.QueryUnescape(payload)
	if err != nil {
		log.Printf("Error decoding URL: %v", err)
		return "", err
	}

	// Unmarshal the payload into a map
	var payloadMap map[string]interface{}
	if err := json.Unmarshal([]byte(decodedPayload), &payloadMap); err != nil {
		log.Printf("Error parsing payload: %v", err)
		return "", err
	}

	// Marshal the map into a compact JSON string
	jsonData, err := json.Marshal(payloadMap)
	if err != nil {
		log.Printf("Error marshalling payload to JSON: %v", err)
		return "", err
	}

	return string(jsonData), nil
}
