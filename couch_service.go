package couchdb

import (
	"context"
	"fmt"
	"time"
)

type CouchServiceI interface {
	GetDB(ctx context.Context, name string, createIfItDoesntExist bool) (*Database, error)
}

type CouchService struct {
	baseURL string
}

func GetInstance(baseURL string) CouchServiceI {
	cs := &CouchService{
		baseURL: baseURL,
	}
	return cs
}

// GetDB retrieves a database with the specified name, optionally creating it if it doesn't exist.
//
// This function sends an HTTP HEAD request to check the existence of the database with the given name.
// If the database exists, it returns a *Database instance initialized with the provided name and a CustomHTTPClient.
// If the database doesn't exist and createIfItDoesntExist is true, it attempts to create the database using createDB function,
// then recursively calls itself with createIfItDoesntExist set to false to retrieve the created database.
// If createIfItDoesntExist is false and the database doesn't exist, it returns errorDBNotFound.
// It returns an error if there was a problem sending the request or if the response status code is not 200 (OK) or 400 (Bad Request).
//
// Parameters:
//   - ctx: The context.Context for the HTTP request.
//   - name: The name of the database to retrieve or create.
//   - createIfItDoesntExist: A boolean indicating whether to create the database if it doesn't exist.
//
// Returns:
//   - A *Database instance representing the retrieved or created database.
//   - An error, if any, encountered during the retrieval or creation of the database.
//     If the operation is successful, it returns nil.
func (c *CouchService) GetDB(ctx context.Context, name string, createIfItDoesntExist bool) (*Database, error) {
	httpClient := NewCustomHTTPClient(c.baseURL, 5, 2*time.Second, 30*time.Second)
	respCode, respBody, err := httpClient.Head(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("error getting database: %w", err)
	}
	if respCode != 200 {
		if respCode == 404 {
			if createIfItDoesntExist {
				err := createDB(ctx, httpClient, name)
				if err != nil {
					return nil, fmt.Errorf("error creating database: %w", err)
				}
				return c.GetDB(ctx, name, false)
			}
			return nil, errorDBNotFound
		}
		return nil, fmt.Errorf("error getting database: %d - %s", respCode, string(respBody))
	}
	return &Database{
		httpClient: httpClient,
		dbName:     name,
	}, nil
}

// createDB creates a new database with the specified name.
//
// This function sends an HTTP PUT request to create a new database with the given name.
// It returns an error if there was a problem sending the request or if the response status code is not 201 (Created) or 202 (Accepted).
//
// Parameters:
//   - ctx: The context.Context for the HTTP request.
//   - c: A *CustomHTTPClient instance used to send the HTTP request.
//   - dbName: The name of the database to create.
//
// Returns:
//   - An error, if any, encountered during the creation of the database.
//     If the operation is successful, it returns nil.
func createDB(ctx context.Context, c *CustomHTTPClient, dbName string) error {
	if !isValidDBName(dbName) {
		return fmt.Errorf("invalid database name: %s", dbName)
	}
	respCode, respBody, err := c.Put(ctx, dbName, nil)
	if err != nil {
		return fmt.Errorf("error creating db: %w", err)
	}
	if respCode != 201 && respCode != 202 {
		return fmt.Errorf("error creating db: %d - %s", respCode, string(respBody))
	}
	return nil
}
