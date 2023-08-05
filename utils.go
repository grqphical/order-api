package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
)

// Removes an item from a slice based on an index
func remove(slice []Order, s int) []Order {
	return append(slice[:s], slice[s+1:]...)
}

// Saves the curent JSON data to the "database" which is just a JSON file
func saveDatabase(orders []Order) {
	bytes, err := json.Marshal(orders)

	if err != nil {
		panic(err)
	}

	os.WriteFile("orders.json", bytes, 0644)
}

func ValidateStruct(s interface{}) (err error) {
	// first make sure that the input is a struct
	// having any other type, especially a pointer to a struct,
	// might result in panic
	structType := reflect.TypeOf(s)
	if structType.Kind() != reflect.Struct {
		return errors.New("input param should be a struct")
	}

	// now go one by one through the fields and validate their value
	structVal := reflect.ValueOf(s)
	fieldNum := structVal.NumField()

	for i := 0; i < fieldNum; i++ {
		// Field(i) returns i'th value of the struct
		field := structVal.Field(i)
		fieldName := structType.Field(i).Name

		// CAREFUL! IsZero interprets empty strings and int equal 0 as a zero value.
		// To check only if the pointers have been initialized,
		// you can check the kind of the field:
		// if field.Kind() == reflect.Pointer { // check }

		// IsZero panics if the value is invalid.
		// Most functions and methods never return an invalid Value.
		isSet := field.IsValid() && !field.IsZero()

		if !isSet {
			err = errors.New(fmt.Sprintf("%v%s in not set; ", err, fieldName))
		}

	}

	return err
}
