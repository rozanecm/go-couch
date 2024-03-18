package main

import "errors"

var (
	errorDBNotFound = errors.New("db not found")
	errNotFound     = errors.New("not found")

	codeToError = map[int]error{
		400: errNotFound,
	}
)
