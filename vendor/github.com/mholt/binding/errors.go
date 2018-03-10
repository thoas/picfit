package binding

import (
	"encoding/json"
	"fmt"
	"strings"
)

// This file shamelessly adapted from martini-contrib/binding

type (
	// Errors may be generated during deserialization, binding,
	// or validation.
	Errors []Error

	// An Error is an error that is associated with 0 or more fields of a
	// request.
	//
	// Fields should return the fields associated with the error. When the return
	// value's length is 0, something is wrong with the request as a whole.
	//
	// Kind should return a string that can be used like an error code to process
	// or categorize the Error.
	//
	// Message should return the error message.
	Error interface {
		error
		Fields() []string
		Kind() string
		Message() string
	}

	fieldsError struct {
		// A fieldError supports zero or more field names, because an error can
		// morph three ways:
		// 		* it can indicate something wrong with the request as a whole.
		//		* it can point to a specific problem with a particular input field
		//		* it can span multiple related input fields.
		fields []string

		// The classification is like an error code, convenient to
		// use when processing or categorizing an error programmatically.
		// It may also be called the "kind" of error.
		kind string

		// Message should be human-readable and detailed enough to
		// pinpoint and resolve the problem, but it should be brief. For
		// example, a payload of 100 objects in a JSON array might have
		// an error in the 41st object. The message should help the
		// end user find and fix the error with their request.
		message string
	}
)

// Add adds an Error associated with the fields indicated by fieldNames, with
// the given kind and message.
//
// Use a fieldNames value of length 0 to indicate that the error is about the
// request as a whole, and not necessarily any of the fields.
//
// kind should be a string that can be used like an error code to process or
// categorize the error being added.
//
// message should be human-readable and detailed enough to pinpoint and resolve
// the problem, but it should be brief. For example, a payload of 100 objects
// in a JSON array might have an error in the 41st object. The message should
// help the end user find and fix the error with their request.
func (e *Errors) Add(fieldNames []string, kind, message string) {
	*e = append(*e, NewError(fieldNames, kind, message))
}

// Len returns the number of errors.
func (e *Errors) Len() int {
	return len(*e)
}

// Has determines whether kind matches the return value of Kind() of any Error
// in e.
func (e *Errors) Has(kind string) bool {
	for _, err := range *e {
		if err.Kind() == kind {
			return true
		}
	}
	return false
}

// Error returns a concatenation of all its error messages.
func (e Errors) Error() string {
	messages := []string{}
	for _, err := range e {
		messages = append(messages, err.Message())
	}
	return strings.Join(messages, "\n")
}

func NewError(fieldNames []string, kind, message string) Error {
	return fieldsError{
		fields:  fieldNames,
		kind:    kind,
		message: message,
	}
}

// Fields returns the names of the fields associated with e.
func (e fieldsError) Fields() []string {
	return e.fields
}

// Kind returns a string that can be used to categorize e.
func (e fieldsError) Kind() string {
	return e.kind
}

// Message returns the message of e.
func (e fieldsError) Message() string {
	return e.message
}

func (e fieldsError) Error() string {
	if len(e.fields) == 0 {
		return e.message
	}

	fields := make([]string, len(e.fields))
	for i, f := range e.fields {
		fields[i] = fmt.Sprintf("* %s", f)
	}

	return fmt.Sprintf("%s\n\t%s", e.message, strings.Join(fields, "\n\t"))
}

func (e fieldsError) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		FieldNames     []string `json:"fieldNames,omitempty"`
		Classification string   `json:"classification,omitempty"`
		Message        string   `json:"message,omitempty"`
	}{
		FieldNames:     e.fields,
		Classification: e.kind,
		Message:        e.message,
	})
}

const (
	RequiredError        = "RequiredError"
	ContentTypeError     = "ContentTypeError"
	DeserializationError = "DeserializationError"
	TypeError            = "TypeError"
)
