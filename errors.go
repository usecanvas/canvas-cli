package main

import "strings"

type ErrorPayload struct {
	Errors []ApiError
}

func (p ErrorPayload) Error() string {
	var errors []string
	for _, err := range p.Errors {
		errors = append(errors, err.Error())
	}
	return strings.Join(errors, "\n")
}

type ApiError struct {
	Detail string
}

func (e ApiError) Error() string {
	return e.Detail
}
