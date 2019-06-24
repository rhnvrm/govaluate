package govaluate

import (
	"errors"
)

// Parameters is a collection of named parameters that can be used by an Expression to retrieve parameters
// when an expression tries to use them.
type Parameters interface {

	// Get gets the parameter of the given name, or an error if the parameter is unavailable.
	// Failure to find the given parameter should be indicated by returning an error.
	Get(name string) (interface{}, error)
}

// MapParameters is an implementation of Parameters interface with map as store for parameters.
type MapParameters map[string]interface{}

// Get returns the parameter of the given name, or an error if the parameter is unavailable.
func (p MapParameters) Get(name string) (interface{}, error) {

	value, found := p[name]

	if !found {
		errorMessage := "No parameter '" + name + "' found."
		return nil, errors.New(errorMessage)
	}

	return value, nil
}
