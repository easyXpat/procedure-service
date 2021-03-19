// Package classification of Procedure API
//
// Documentation for Procedure API
//
//	Schemes: http
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package handlers

import "github.com/easyXpat/procedure-service/data"

//
// NOTE: Types defined here are purely for documentation purposes
// these types are not used by any of the handlers

// Generic error message returned as a string
// swagger:response errorResponse
type errorResponseWrapper struct {
	// Description of the error
	// in: body
	Body GenericError
}

// Validation errors defined as an array of strings
// swagger:response errorValidation
type errorValidationWrapper struct {
	// Collection of the errors
	// in: body
	Body ValidationError
}

// A list of procedures
// swagger:response proceduresResponse
type proceduresResponseWrapper struct {
	// All current procedures
	// in: body
	Body []data.Procedure
}

// Data structure representing a single procedure
// swagger:response procedureResponse
type procedureResponseWrapper struct {
	// Newly created procedure
	// in: body
	Body data.Procedure
}


