package couchdb

import (
	"errors"
	"reflect"
	"regexp"
)

var ErrMissingID = errors.New("missing _id field")
var ErrMissingRev = errors.New("missing _rev field")

// checkParameter checks if the parameter is a struct or a map[string]interface{}
// and if it contains the fields "_id" and "_rev". It returns custom errors for each missing field.
func checkParameter(param interface{}) error {
	switch reflect.TypeOf(param).Kind() {
	case reflect.Map:
		paramMap := param.(map[string]interface{})
		if _, ok := paramMap["_id"]; !ok {
			return ErrMissingID
		}
		if _, ok := paramMap["_rev"]; !ok {
			return ErrMissingRev
		}
		return nil
	case reflect.Struct:
		fields := reflect.TypeOf(param)
		_, idExist := fields.FieldByName("_id")
		_, revExist := fields.FieldByName("_rev")
		if !idExist {
			return ErrMissingID
		}
		if !revExist {
			return ErrMissingRev
		}
		return nil
	default:
		return errors.New("unsupported type")
	}
}

// isValidDBName checks if the provided name is a valid database name according to the specified rules.
//
// The database name must adhere to the following rules:
// - Name must begin with a lowercase letter (a-z)
// - Allowed characters: lowercase letters (a-z), digits (0-9), and any of the characters _, $, (, ), +, -, and /
// - Regular expression representation: ^[a-z][a-z0-9_$()+/-]*$
//
// Parameters:
//   - name: The name to be validated as a database name.
//
// Returns:
//   - A boolean value indicating whether the provided name is a valid database name or not.
//
// Example:
//
//	isValid := isValidDBName("my_database_123")
//	fmt.Println(isValid) // Output: true
func isValidDBName(name string) bool {
	// Regular expression pattern for validating database names
	pattern := `^[a-z][a-z0-9_$()+/-]*$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(name)
}

// isValidParam checks if the provided parameter is a pointer to a struct or a map[string]interface{}.
//
// Parameters:
//   - param: The parameter to be checked.
//
// Returns:
//   - A boolean value indicating whether the parameter is a valid type.
func isValidParam(param interface{}) bool {
	if param == nil {
		return false
	}
	t := reflect.TypeOf(param)
	if t.Kind() != reflect.Ptr {
		return false
	}
	elemType := t.Elem()
	switch elemType.Kind() {
	case reflect.Struct:
		return true
	case reflect.Map:
		return elemType.Key().Kind() == reflect.String && elemType.Elem().Kind() == reflect.Interface
	default:
		return false
	}
}
