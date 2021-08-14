package handlers

import (
	"context"
	"github.com/easyXpat/procedure-service/data"
	"net/http"
)

// swagger:route PUT /procedures procedures updateProcedure
// Update existing procedure
//
// responses:
//	200: procedureResponse

// UpdateProcedure handles the update of a procedure
func (ph *Procedure) UpdateProcedure(w http.ResponseWriter, r *http.Request) {
	// fetch the procedure from the context
	ph.logger.Debug("handler for updateProcedure")
	procedure := r.Context().Value(ProcedureKey{}).(data.Procedure)
	ph.logger.Debug("Updating procedure", "id", procedure.ID)

	updatedProcedure, err := ph.db.UpdateProcedure(context.Background(), &procedure)
	if err != nil {
		ph.logger.Error("Error updating procedure", "error", err)
		data.ToJSON(&GenericError{Message: err.Error()}, w)
		return
	}
	data.ToJSON(&updatedProcedure, w)
}

// UpdateStep handles the update of a step
func (st *Step) UpdateStep(w http.ResponseWriter, r *http.Request) {
	// fetch the step from the context
	st.logger.Debug("handler for updateStep")
	step := r.Context().Value(StepKey{}).(data.Step)
	st.logger.Debug("Updating step", "id", step.ID)

	updatedStep, err := st.db.UpdateStep(context.Background(), &step)
	if err != nil {
		st.logger.Error("Error updating procedure", "error", err)
		data.ToJSON(&GenericError{Message: err.Error()}, w)
		return
	}
	data.ToJSON(&updatedStep, w)
}