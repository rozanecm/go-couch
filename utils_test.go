package couchdb

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

// Define structs for testing
type structWithTags struct {
	ID  string `json:"_id"`
	Rev string `json:"_rev"`
}

type structWithoutTags struct {
	ID  string
	Rev string
}

type structWithEmptyTags struct {
	ID  string `json:""`
	Rev string `json:""`
}

type structWithMixedTags struct {
	ID  string `json:"_id"`
	Rev string
}

// TestCheckParameter tests the checkParameter function
func TestCheckParameter(t *testing.T) {
	tests := []struct {
		name     string
		param    interface{}
		expected error
	}{
		{
			name:     "Test map[string]interface{} with _id and _rev",
			param:    map[string]interface{}{"_id": "123", "_rev": "456"},
			expected: nil,
		},
		{
			name:     "Test map[string]interface{} without _id",
			param:    map[string]interface{}{"_rev": "456"},
			expected: ErrMissingID,
		},
		{
			name:     "Test map[string]interface{} without _rev",
			param:    map[string]interface{}{"_id": "123"},
			expected: ErrMissingRev,
		},
		{
			name:     "Test structWithTags with _id and _rev",
			param:    structWithTags{ID: "123", Rev: "456"},
			expected: nil,
		},
		{
			name:     "Test structWithoutTags without _id and _rev",
			param:    structWithoutTags{},
			expected: ErrMissingID,
		},
		{
			name:     "Test structWithEmptyTags with empty _id and _rev",
			param:    structWithEmptyTags{},
			expected: ErrMissingID,
		},
		{
			name:     "Test structWithMixedTags with _id and without _rev",
			param:    structWithMixedTags{ID: "123"},
			expected: ErrMissingRev,
		},
		{
			name:     "Test unsupported type",
			param:    123,
			expected: errors.New("unsupported type"),
		},
		{
			name:     "Test nil parameter",
			param:    nil,
			expected: errors.New("unsupported type"), // Adjust this expected error message if needed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checkParameter(tt.param)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("checkParameter() got = %v, want %v", got, tt.expected)
			}
		})
	}
}
func TestIsValidDBName(t *testing.T) {
	testCases := []struct {
		name     string
		expected bool
	}{
		{"my_database_123", true},     // Valid database name
		{"Database_123", false},       // Database name doesn't start with a lowercase letter
		{"my-database", true},         // Valid database name with hyphen
		{"my(database)", true},        // Valid database name with parentheses
		{"my+database", true},         // Valid database name with plus sign
		{"my_database$", true},        // Valid database name with dollar sign
		{"my_database&", false},       // Invalid character (&)
		{"MyDatabase", false},         // Database name doesn't start with a lowercase letter
		{"123database", false},        // Database name starts with a digit
		{"_database", false},          // Database name doesn't start with a lowercase letter
		{"my database", false},        // Database name contains whitespace
		{"my/database", true},         // Valid database name with slash
		{"my+database-123_$()", true}, // Valid database name with various special characters
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := isValidDBName(tc.name)
			if result != tc.expected {
				t.Errorf("Expected isValidDBName(%s) to be %v, got %v", tc.name, tc.expected, result)
			}
		})
	}
}

func TestIsValidParam(t *testing.T) {
	tests := []struct {
		name     string
		param    interface{}
		expected bool
	}{
		{name: "nil value", param: nil, expected: false},
		{name: "non-pointer value", param: "string", expected: false},
		{name: "pointer to empty struct", param: struct{}{}, expected: false},
		{name: "pointer to struct", param: &struct{ Name string }{}, expected: true},
		{name: "pointer to map[string]interface{}", param: &map[string]interface{}{}, expected: true},
		{name: "pointer to map[string]string", param: &map[string]string{}, expected: false},
		{name: "pointer to map[int]interface{}", param: &map[int]interface{}{}, expected: false},
		{name: "pointer to slice", param: &[]string{"a", "b", "c"}, expected: false},
		{name: "pointer to slice of map[string]interface{}", param: &[]map[string]interface{}{{}}, expected: false},
		{name: "pointer to slice of struct", param: &[]struct{ Name string }{{}}, expected: false},
		{name: "pointer to slice of interface{}", param: &[]interface{}{"a", 1, true}, expected: false},
		{name: "pointer to map[string]int", param: &map[string]int{"a": 1}, expected: false},
		{name: "pointer to map[int]struct{}", param: &map[int]struct{}{}, expected: false},
		{name: "pointer to map[int]string", param: &map[int]string{}, expected: false},
		{name: "pointer to map[string]struct{}", param: &map[string]struct{ Name string }{"a": {Name: "b"}}, expected: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := isValidParam(test.param)
			if result != test.expected {
				t.Errorf("Expected %t, got %t", test.expected, result)
			}
		})
	}
}

func TestAddSlashIfNeeded(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"https://example.com/api/", "https://example.com/api/"},
		{"https://example.com/api", "https://example.com/api/"},
		{"", "/"},
		{"/", "/"},
		{"test/", "test/"},
		{"test", "test/"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Input: %s", tc.input), func(t *testing.T) {
			output := addSlashIfNeeded(tc.input)
			if output != tc.expected {
				t.Errorf("Expected: %s, Got: %s", tc.expected, output)
			}
		})
	}
}
