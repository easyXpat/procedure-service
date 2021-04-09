package data

import (
	"time"
)

// Step defines the structure for an API step
// swagger:model
type Step struct {
	// unique id for the step
	//
	// required: false
	// min: 1
	// max length: 255
	ID string `json:"id" sql:"id"`

	// unique procedure id for the step
	//
	// required: true
	// min: 1
	// max length: 255
	ProcedureID string `json:"procedure_id" sql:"procedure_id"`

	// name for the step
	//
	// required: true
	// min: 1
	// max length: 255
	Name string `json:"name" validate:"required" sql:"name"`

	// description for the step
	//
	// required: false
	// min: 1
	// max length: 1000
	Description string `name:"description" validate:"required" sql:"description"`

	// creation time for the step
	//
	// required: false
	CreatedAt time.Time `name:"created_at" sql:"created_at"`

	// last update time for the step
	//
	// required: false
	UpdatedAt time.Time `name:"updated_at" sql:"updated_at"`
}

type Steps []*Step