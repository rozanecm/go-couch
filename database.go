package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

type Database struct {
	httpClient *CustomHTTPClient
	dbName     string
}

// CreateDoc creates a new document in the database.
//
// This function sends an HTTP POST request to create a new document in the database with the provided context and document data.
// It returns an error if there was a problem sending the request or if the response status code is not 200 (OK).
// If an error occurs during the HTTP request, it is wrapped and returned with additional context.
// If the response status code indicates an error, an error message is constructed using the status code and response body.
//
// Parameters:
//   - ctx: The context.Context for the HTTP request.
//   - doc: The document data to be created in the database. It can be of any type.
//
// Returns:
//   - An error, if any, encountered during the creation of the document.
//     If the creation is successful, it returns nil.
//
// Example:
//
//	err := db.CreateDoc(ctx, map[string]interface{}{
//	    "name": "John Doe",
//	    "age":  30,
//	})
//	if err != nil {
//	    log.Fatalf("Error creating document: %v", err)
//	}
//
// Note: This function assumes that db.httpClient is a CustomHTTPClient instance with methods for sending HTTP requests.
// The response body is expected to contain additional information in case of errors.
func (db *Database) CreateDoc(ctx context.Context, doc any) (*CreateDocResponseType, error) {
	respCode, respBody, err := db.httpClient.Post(ctx, db.dbName, doc)
	if err != nil {
		return nil, fmt.Errorf("error creating doc: %w", err)
	}

	if respCode != 200 && respCode != 201 {
		return nil, fmt.Errorf("error creating doc: %d - %s", respCode, string(respBody))
	}

	var createDocResponse CreateDocResponseType

	err = json.Unmarshal(respBody, &createDocResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling create doc response: %w", err)
	}

	return &createDocResponse, nil
}

type CreateDocResponseType struct {
	ID  string `json:"id"`
	Ok  bool   `json:"ok"`
	Rev string `json:"rev"`
}

// GetDoc retrieves a document from the database by its ID and populates the provided struct with its data.
//
// This function sends an HTTP GET request to retrieve a document from the database based on the provided ID.
// It populates the provided struct pointer with the retrieved document data.
// If the provided document parameter is not a pointer to a struct, an error is returned.
// It returns an error if there was a problem sending the request, if the response status code is not 200 (OK),
// or if there was an error unmarshalling the response body into the provided struct.
//
// Parameters:
//   - ctx: The context.Context for the HTTP request.
//   - id: The ID of the document to retrieve from the database.
//   - doc: A pointer to a struct where the retrieved document data will be populated.
//
// Returns:
//   - An error, if any, encountered during the retrieval and unmarshalling of the document.
//     If the retrieval and unmarshalling are successful, it returns nil.
//
// Example:
//
//	type Person struct {
//	    Name string `json:"name"`
//	    Age  int    `json:"age"`
//	}
//
//	var person Person
//	err := db.GetDoc(ctx, "document_id", &person)
//	if err != nil {
//	    log.Fatalf("Error getting document: %v", err)
//	}
func (db *Database) GetDoc(ctx context.Context, id string, doc any) error {
	if !isValidParam(doc) {
		return fmt.Errorf("doc parameter must be a pointer to a struct")
	}

	respCode, respBody, err := db.httpClient.Get(ctx, fmt.Sprintf("%s/%s", db.dbName, id))
	if err != nil {
		return fmt.Errorf("error getting doc: %w", err)
	}

	if respCode != 200 {
		if errFromMap, ok := codeToError[respCode]; ok {
			return errFromMap
		}
		return fmt.Errorf("error getting doc: %d - %s", respCode, string(respBody))
	}

	err = json.Unmarshal(respBody, doc)
	if err != nil {
		return fmt.Errorf("error unmarshalling doc: %w", err)
	}

	return nil
}

// UpdateDoc creates or updates a document in the database.
//
// This function either creates a new document with the specified ID or updates an existing document with a new revision.
// To update an existing document, the current revision must be provided in the document body, as a query parameter ("rev"),
// or in the "If-Match" request header.
//
// Parameters:
//   - ctx: The context.Context for the HTTP request.
//   - doc: The document data to be created or updated. It can be of any type, but it must contain the current revision information for updates.
//   - id: The ID of the document to be created or updated in the database.
//
// Returns:
//   - An error, if any, encountered during the creation or update of the document.
//     If the operation is successful, it returns nil.
//
// Example:
//
//	// Update an existing document
//	err := db.UpdateDoc(ctx, map[string]interface{}{
//	    "_id":  "existing_doc_id",
//	    "_rev": "current_revision",
//	    "key":  "new_value",
//	}, "existing_doc_id")
//	if err != nil {
//	    log.Fatalf("Error updating document: %v", err)
//	}
//
//	// Create a new document
//	err = db.UpdateDoc(ctx, map[string]interface{}{
//	    "_id":  "new_doc_id",
//	    "key":  "value",
//	}, "new_doc_id")
//	if err != nil {
//	    log.Fatalf("Error creating document: %v", err)
//	}
func (db *Database) UpdateDoc(ctx context.Context, id string, doc any) error {
	if err := checkParameter(doc); err != nil {
		return fmt.Errorf("doc check failed: %w", err)
	}

	respCode, respBody, err := db.httpClient.Put(ctx, fmt.Sprintf("%s/%s", db.dbName, id), doc)
	if err != nil {
		return fmt.Errorf("error updating doc: %w", err)
	}
	if respCode != 200 && respCode != 201 {
		return fmt.Errorf("error updating doc: %d - %s", respCode, string(respBody))
	}

	return nil
}

// DeleteDoc deletes a document from the database using its ID.
//
// It takes a context object (ctx) for cancellation and deadline propagation.
// The function first retrieves the document with the given ID to obtain its revision ID (_rev).
// Then it sends a DELETE request to the database to delete the document using its ID and revision ID.
// If the response status code is not 200 (OK) or 202 (Accepted), an error is returned.
//
// Parameters:
//   - ctx: The context.Context for the HTTP request.
//   - id: The ID of the document to be deleted from the database.
//
// Returns:
//   - An error, if any, encountered during the deletion of the document.
//     If the deletion is successful, it returns nil.
//
// Example:
//
//	err := db.DeleteDoc(ctx, "document_id")
//	if err != nil {
//	    log.Fatalf("Error deleting document: %v", err)
//	}
func (db *Database) DeleteDoc(ctx context.Context, id string) error {
	var doc map[string]interface{}
	err := db.GetDoc(ctx, id, &doc)
	if err != nil {
		return fmt.Errorf("error getting doc to delete: %w", err)
	}

	rev, _ := doc["_rev"].(string)

	respCode, respBody, err := db.httpClient.Delete(ctx, fmt.Sprintf("%s/%s?rev=%s", db.dbName, id, rev))
	if err != nil {
		return fmt.Errorf("error deleting doc: %w", err)
	}

	if respCode != 200 && respCode != 202 {
		return fmt.Errorf("error deleting doc: %d - %s", respCode, string(respBody))
	}

	return nil
}

func (db *Database) CreateDesignDoc(ctx context.Context, designDoc string, views map[string]ViewDefinition) error {
	docID := fmt.Sprintf("_design/%s", designDoc)
	body := designDocument{
		ID:         docID,
		Language:   "javascript",
		Autoupdate: true,
		Views:      views,
	}

	var prevDoc designDocument
	err := db.GetDoc(ctx, docID, &prevDoc)
	if !errors.Is(err, errNotFound) {
		body.Rev = prevDoc.Rev
	}

	code, responseBytes, err := db.httpClient.Put(ctx, fmt.Sprintf("%s/_design/%s", db.dbName, designDoc), body)
	if err != nil {
		return fmt.Errorf("error creating design doc: %w", err)
	}

	if code != 200 && code != 201 {
		return fmt.Errorf("error creating design doc: %d - %s", code, string(responseBytes))
	}
	return nil
}

type designDocument struct {
	ID                string                    `json:"_id"`
	Rev               string                    `json:"_rev,omitempty"`
	Language          string                    `json:"language"`
	Options           map[string]any            `json:"options,omitempty"`
	Filters           map[string]string         `json:"filters,omitempty"`
	Lists             map[string]string         `json:"lists,omitempty"`    // Deprecated
	Rewrites          any                       `json:"rewrites,omitempty"` // Deprecated. Array or string
	Shows             map[string]string         `json:"shows,omitempty"`    // Deprecated
	Updates           map[string]string         `json:"updates,omitempty"`
	ValidateDocUpdate string                    `json:"validate_doc_update,omitempty"`
	Views             map[string]ViewDefinition `json:"views,omitempty"`
	Autoupdate        bool                      `json:"autoupdate,omitempty"`
}

type ViewDefinition struct {
	Map    string `json:"map"`
	Reduce string `json:"reduce,omitempty"`
}
