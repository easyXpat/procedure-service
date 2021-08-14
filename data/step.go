package data

import (
	"context"
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

	// procedure name of the step. In case the step is procedure specific.
	//
	// required: false
	// min: 1
	// max length: 255
	ProcedureName string `json:"procedure_name" sql:"procedure_name"`

	// city for the step. In case the step is city specific
	//
	// required: false
	// min: 1
	// max length: 255
	City string `json:"city" sql:"city"`

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

type StepDB interface {
	AddStep(ctx context.Context, p *Step) error
	GetAllSteps(ctx context.Context) (Steps, error)
	GetStep(ctx context.Context, id string) (Steps, error)
	UpdateStep(ctx context.Context, p *Step) (Steps, error)
	DeleteStep(ctx context.Context, id string) error
}

type Steps []*Step