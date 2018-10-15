package net

import (
	"fmt"
	"net/http"
	"runtime"
	"strconv"
)

// netError represent the error struct that contains the error message and the error code.
type netError struct {
	msg           string
	brokenElement string
	brokenIndex   int
	code          int
	stackTrace    string
}

// Error returns the error message of the error (also to be conform with the 'error' interface)
func (e netError) Error() string {
	return e.msg
}

// BrokenElement returns the key for a broken element (front-end can handle different elements better)
func (e netError) BrokenElement() string {
	return e.brokenElement
}

// BrokenIndex returns the index of the broken element (only available if a broken element was set)
func (e netError) BrokenIndex() int {
	return e.brokenIndex
}

// Code returns the http error code
func (e netError) Code() int {
	return e.code
}

// Stacktrace returns the Error stacktrace if available
func (e netError) Stacktrace() string {
	return e.stackTrace
}

// NewInternalServerError creates a new internal server error (500)
func NewInternalServerError(e error) error {
	// don't override old errors
	if netErr, ok := e.(netError); ok {
		return netErr
	}

	return netError{
		msg:        e.Error(),
		code:       http.StatusInternalServerError,
		stackTrace: createStackTrace(),
	}
}

// NewBadRequestError creates a new bad request error (400)
func NewBadRequestError(msg string, brokenElement string) error {
	return NewBadRequestErrorWithIndex(msg, brokenElement, 0)
}

// NewBadRequestError creates a new bad request error (400)
// this error should be used for all content the user has created
func NewBadRequestErrorWithIndex(msg string, brokenElement string, brokenIndex int) error {
	return netError{
		msg:           msg,
		brokenElement: brokenElement,
		brokenIndex:   brokenIndex,
		code:          http.StatusBadRequest,
		stackTrace:    "", // fixable by the user
	}
}

// NewNotFoundError creates a new not found error (404)
// this error should be used if there are missing element (the front end has a bug!)
func NewNotFoundError(msgArgs ...interface{}) error {
	return netError{
		msg:        convertMsgArgs(msgArgs),
		code:       http.StatusNotFound,
		stackTrace: createStackTrace(), // some code is broken
	}
}

// NewConflictError creates a new conflict error (409)
//
func NewConflictError(msgArgs ...interface{}) error {
	return netError{
		msg:        convertMsgArgs(msgArgs),
		code:       http.StatusConflict,
		stackTrace: createStackTrace(),
	}
}

// NewUnauthorizedError creates a new unauthorized error (401)
// this error should be used if the session is empty or expired
func NewUnauthorizedError(msg string) error {
	return netError{
		msg:        msg,
		code:       http.StatusUnauthorized,
		stackTrace: "", // missing or expired session error
	}
}

func convertMsgArgs(msgArgs []interface{}) string {
	if len(msgArgs) == 0 {
		return ""
	}
	msg, ok := msgArgs[0].(string)
	if !ok {
		return "Wrong 'msgArgs' format"
	}
	if len(msgArgs) == 1 {
		return msg
	}
	return fmt.Sprintf(msg, msgArgs[1:]...)
}

func createStackTrace() string {
	var result string
	for i := 2; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		result += file + ":" + strconv.Itoa(line) + "\n"
	}

	return result
}
