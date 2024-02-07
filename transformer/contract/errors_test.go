package contract

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestError_Error(t *testing.T) {
	assertion := assert.New(t)

	err := NewError(InvalidInputDataStructure, nil)
	assertion.NotNil(err)
	assertion.Equal("invalid input data structure", err.Err.Error())
	assertion.Nil(err.Payload)

	err = NewError(InvalidInputDataStructure, "string")
	assertion.Equal("string", err.Payload)

	err = NewError(InvalidInputDataStructure, 123)
	assertion.Equal(123, err.Payload)

	err = NewError(InvalidInputDataStructure, map[string]float32{"test": 123})
	assertion.Equal(map[string]float32{"test": 123}, err.Payload)
}
