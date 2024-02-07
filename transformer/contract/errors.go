package contract

import (
	"errors"
)

type ErrorCode int

type Error struct {
	Err     error
	Code    ErrorCode
	Payload interface{}
}

const (
	Unknown ErrorCode = iota
	UnreadableInputData
	InvalidInputDataStructure
	InvalidFiltersStructure
	NonWriteableOutputData
)

var ErrorCodes = map[ErrorCode]string{
	Unknown:                   "unknown error",
	UnreadableInputData:       "unreadable input data",
	InvalidInputDataStructure: "invalid input data structure",
	InvalidFiltersStructure:   "invalid filters structure",
	NonWriteableOutputData:    "can't write output data",
}

func NewError(code ErrorCode, payload interface{}) *Error {
	return &Error{
		Err:     errors.New(ErrorCodes[code]),
		Code:    code,
		Payload: payload,
	}
}
