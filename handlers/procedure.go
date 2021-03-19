package handlers

import (
	"github.com/easyXpat/procedure-service/data"
	"github.com/hashicorp/go-hclog"
)

// ProcedureKey is used as a key for storing the Procedure object in context at middleware
type ProcedureKey struct{}

// Procedure wraps instances needed to perform operations on procedure object
type Procedure struct {
	logger          hclog.Logger
	db 	data.ProcedureDB
	validator         *data.Validation
}

type GenericError struct {
	Message string `json:"message"`
}

// ValidationError is a collection of validation error messages
type ValidationError struct {
	Messages []string `json:"messages"`
}

// NewProcedure creates a new procedure handler
func NewProcedure(l hclog.Logger, pdb data.ProcedureDB, v *data.Validation) *Procedure {
	return &Procedure{
		l,
		pdb,
		v,
	}
}

