package utils

import (
	"errors"
	"fmt"
	"log"
	"reflect"
)

func GetStructFields(s interface{}) (string, error) {
	v := reflect.ValueOf(s)

	// Check if the input is a struct
	if v.Kind() != reflect.Struct {
		log.Println("Provided value is not a struct.")
		return "", errors.New("Provided value is not a struct.")
	}

	// Iterate over the struct fields
	var msg string
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		msg += fmt.Sprintf("%s: %v, ", field.Name, value)
	}

	// Delete the last comma and space
	msg = msg[:len(msg)-2]

	return msg, nil
}
