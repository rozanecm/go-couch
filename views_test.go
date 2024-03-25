package couchdb

import (
	"testing"
)

// Define a test struct for holding test cases
type testCase struct {
	Name      string
	Input     interface{}
	ShouldErr bool
}

func TestCheckStructForJSONFields(t *testing.T) {
	testCases := []testCase{
		{
			Name:      "Valid struct with required fields and JSON tags",
			Input:     &validStruct{},
			ShouldErr: false,
		},
		{
			Name:      "Struct missing 'Rows' field",
			Input:     &missingRowsStruct{},
			ShouldErr: true,
		},
		{
			Name:      "Struct with 'Rows' field of wrong type",
			Input:     &wrongTypeRowsStruct{},
			ShouldErr: true,
		},
		{
			Name:      "Struct with 'Rows' field missing JSON tag",
			Input:     &missingRowsTagStruct{},
			ShouldErr: true,
		},
		{
			Name:      "Struct missing 'ID' field",
			Input:     &missingIDStruct{},
			ShouldErr: true,
		},
		{
			Name:      "Struct missing 'Key' field",
			Input:     &missingKeyStruct{},
			ShouldErr: true,
		},
		{
			Name:      "Struct with 'ID' field missing JSON tag",
			Input:     &missingIDTagStruct{},
			ShouldErr: true,
		},
		{
			Name:      "Struct with 'Key' field missing JSON tag",
			Input:     &missingKeyTagStruct{},
			ShouldErr: true,
		},
		{
			Name:      "Valid struct with 'Doc' field and JSON tag",
			Input:     &validDocStruct{},
			ShouldErr: false,
		},
		{
			Name:      "Struct with 'Doc' field missing JSON tag",
			Input:     &missingDocTagStruct{},
			ShouldErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := checkStructForJSONFields(tc.Input)
			if (err != nil) != tc.ShouldErr {
				t.Errorf("Expected error: %v, Got error: %v", tc.ShouldErr, err)
			}
		})
	}
}

// Define sample structs for testing

type validStruct struct {
	Rows []struct {
		ID  string   `json:"id"`
		Key string   `json:"key"`
		Doc struct{} `json:"doc"`
	} `json:"rows"`
}

type missingRowsStruct struct {
	ID  string   `json:"id"`
	Key string   `json:"key"`
	Doc struct{} `json:"doc"`
}

type wrongTypeRowsStruct struct {
	Rows string   `json:"rows"`
	ID   string   `json:"id"`
	Key  string   `json:"key"`
	Doc  struct{} `json:"doc"`
}

type missingRowsTagStruct struct {
	Rows []struct {
		ID  string   `json:"id"`
		Key string   `json:"key"`
		Doc struct{} `json:"doc"`
	}
}

type missingIDStruct struct {
	Rows []struct {
		Key string   `json:"key"`
		Doc struct{} `json:"doc"`
	} `json:"rows"`
}

type missingKeyStruct struct {
	Rows []struct {
		ID  string   `json:"id"`
		Doc struct{} `json:"doc"`
	} `json:"rows"`
}

type missingIDTagStruct struct {
	Rows []struct {
		ID  string
		Key string   `json:"key"`
		Doc struct{} `json:"doc"`
	} `json:"rows"`
}

type missingKeyTagStruct struct {
	Rows []struct {
		ID  string `json:"id"`
		Key string
		Doc struct{} `json:"doc"`
	} `json:"rows"`
}

type validDocStruct struct {
	Rows []struct {
		ID  string   `json:"id"`
		Key string   `json:"key"`
		Doc struct{} `json:"doc"`
	} `json:"rows"`
}

type missingDocTagStruct struct {
	Rows []struct {
		ID  string   `json:"id"`
		Key string   `json:"key"`
		Doc struct{} `json:"dock"`
	} `json:"rows"`
}
