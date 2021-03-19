package data

import (
	"context"
	"time"
)

// Procedure defines the structure for an API procedure
// swagger:model
type Procedure struct {
	// unique id for the procedure
	//
	// required: false
	// min: 1
	// max length: 255
	ID string `json:"id" sql:"id"`

	// name for the procedure
	//
	// required: true
	// max length: 255
	Name string `json:"name" validate:"required" sql:"name"`

	// description for the procedure
	//
	// required: false
	// max length: 10000
	Description string `json:"description" validate:"required" sql:"description"`

	// city for the procedure
	//
	// required: false
	// max length: 32
	City string `json:"city" validate:"required" sql:"city"`

	// creation time for the procedure
	//
	// required: false
	CreatedAt  time.Time `json:"created_at" sql:"created_at"`

	// last update time for the procedure
	//
	// required: false
	UpdatedAt  time.Time `json:"updated_at" sql:"updated_at"`
}

type Procedures []*Procedure

// ProcedureDB is an interface for the storage implementation of the procedure service
type ProcedureDB interface {
	AddProcedure(ctx context.Context, p *Procedure) error
	GetAllProcedures(ctx context.Context) (Procedures, error)
}
