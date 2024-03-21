package couchdb

import (
	"testing"
)

func TestGetInstance(t *testing.T) {
	testCases := []struct {
		Name          string
		BaseURL       string
		Username      string
		Password      string
		ExpectedError bool
	}{
		{
			Name:          "Valid URL and authentication",
			BaseURL:       "https://example.com",
			Username:      "testuser",
			Password:      "testpassword",
			ExpectedError: false,
		},
		{
			Name:          "Invalid URL scheme",
			BaseURL:       "example.com",
			Username:      "testuser",
			Password:      "testpassword",
			ExpectedError: true,
		},
		{
			Name:          "Invalid base URL",
			BaseURL:       "invalidurl",
			Username:      "testuser",
			Password:      "testpassword",
			ExpectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil && !tc.ExpectedError {
					t.Errorf("Unexpected panic: %v", r)
				}
			}()

			_ = GetInstance(tc.BaseURL, tc.Username, tc.Password)
		})
	}
}
