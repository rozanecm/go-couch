package couchdb

import "errors"

var (
	ErrNotFound = errors.New("not found")

	codeToError = map[int]error{
		404: ErrNotFound,
	}
)
