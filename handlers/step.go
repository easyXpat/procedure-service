package handlers

import (
	"github.com/easyXpat/procedure-service/data"
	"github.com/hashicorp/go-hclog"
)

// StepKey is used as a key for storing the Step object in context at middleware
type StepKey struct{}

// Step wraps instances needed to perform operations on step object
type Step struct {
	logger          hclog.Logger
	db 	data.StepDB
	validator         *data.Validation
}

// NewProcedure creates a new procedure handler
func NewStep(l hclog.Logger, pdb data.StepDB, v *data.Validation) *Step {
	return &Step{
		l,
		pdb,
		v,
	}
}

