package handlers

import (
	"context"
	"fmt"
	"github.com/easyXpat/procedure-service/data"
	"net/http"
)

// swagger:route POST /procedures procedures createProcedure
// Create a new procedure
//
// responses:
//	200: procedureResponse
// 	400: errorValidation
//  422: errorValidation
//  501: errorResponse

// CreateProcedure handles POST requests to add new procedures
func (ph *Procedure) CreateProcedure(w http.ResponseWriter, r *http.Request) {
	// fetch the product from the context
	procedure := r.Context().Value(ProcedureKey{}).(data.Procedure)
	ph.logger.Debug(fmt.Sprintf("Inserting procedure: %v", procedure))
	err := ph.db.AddProcedure(context.Background(), &procedure)
	if err != nil {
		data.ToJSON(&GenericError{Message: err.Error()}, w)
		return
	}
	// insert steps mapping is present
	if procedure.Steps != nil {
		err = ph.db.MapProcedureSteps(context.Background(), procedure.ID, procedure.Steps)
		if err != nil {
			data.ToJSON(&GenericError{Message: err.Error()}, w)
			return
		}
	}
	data.ToJSON(&procedure, w)
}

// swagger:route POST /steps steps createStep
// Create a new step
//
// responses:
//	200: stepResponse

// CreateStep handles POST requests to add new step
func (st *Step) CreateStep(w http.ResponseWriter, r *http.Request) {
	// fetch the product from the context
	step := r.Context().Value(StepKey{}).(data.Step)
	st.logger.Debug(fmt.Sprintf("Inserting step: %v", step))
	err := st.db.AddStep(context.Background(), &step)
	if err != nil {
		data.ToJSON(&GenericError{Message: err.Error()}, w)
		return
	}
	data.ToJSON(&step, w)
}