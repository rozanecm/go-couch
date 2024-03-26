package couchdb

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
)

var ErrMissingID = errors.New("missing _id field")
var ErrMissingRev = errors.New("missing _rev field")
var ErrMissingDocumentFields = errors.New("missing document fields")

// checkParameter checks if the parameter is a struct or a map[string]interface{}
// and if it contains the fields "_id" and "_rev". It returns custom errors for each missing field.
func checkParameter(param interface{}) error {
	value := reflect.ValueOf(param)
	kind := value.Kind()

	switch kind {
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
		// Iterate over fields of the struct and check if Document is embedded
		docType := reflect.TypeOf(Document{})
		baseType := value.Type()
		for i := 0; i < baseType.NumField(); i++ {
			field := baseType.Field(i)
			if field.Type == docType && field.Anonymous {
				return nil // Document is embedded
			}
		}

		return ErrMissingDocumentFields // Document is not embedded
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

// isValidURLScheme checks if a given string represents a valid URL scheme.
func isValidURLScheme(s string) bool {
	parsedURL, err := url.Parse(s)
	if err != nil || parsedURL.Scheme == "" {
		return false
	}
	return true
}

// formAuthenticatedURL forms a URL with the provided base URL, username, and password.
// It returns the formatted URL string.
func formAuthenticatedURL(baseURL, username, password string) (string, error) {
	// Parse the base URL to ensure it's valid
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("error parsing base URL: %w", err)
	}

	// Add username and password to the URL
	if username != "" && password != "" {
		parsedURL.User = url.UserPassword(username, password)
	}

	authURL := parsedURL.String()
	return authURL, nil
}

// testURLWithHEAD sends a HEAD request to the specified URL and checks the response status code.
// It returns true if the response status code is within the 200-299 range, indicating a successful request.
func testURLWithHEAD(url string) error {
	// Send a HEAD request to the URL
	resp, err := http.Head(url)
	if err != nil {
		return fmt.Errorf("error sending HEAD request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	return errors.New("invalid response status code")
}

// addSlashIfNeeded checks if the last character of a string is a slash '/' and adds it if it's not present.
// If the input string is empty or already ends with a slash, it returns the input string as is.
func addSlashIfNeeded(s string) string {
	if len(s) == 0 || s[len(s)-1] != '/' {
		return s + "/"
	}
	return s
}
