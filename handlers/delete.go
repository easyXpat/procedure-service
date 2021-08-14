package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/easyXpat/procedure-service/data"
	"github.com/gorilla/mux"
)

// swagger:route DELETE /procedures/{id} procedures deleteProcedures
// Delete an existing procedure
// responses:
// 	200: procedureResponse
// 	404: errorValidation

// DeleteProcedure handles DELETE requests to existing procedures
func (ph *Procedure) DeleteProcedure(w http.ResponseWriter, r *http.Request) {
	// Get the id from the URL
	vars := mux.Vars(r)
	id := vars["id"]
	// Get the procedure from the DB based on its id
	procedure, err := ph.db.GetProcedure(context.Background(), id)
	if err != nil {
		ph.logger.Error("Unable to fetch procedure")
		w.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, w)
		return
	}
	// Perform actual procedure deletion
	ph.logger.Debug(fmt.Sprintf("Deleting procedure: %v", procedure))
	err = ph.db.DeleteProcedure(context.Background(), procedure.ID)
	if err != nil {
		ph.logger.Error("Error ocurred while deleting procedure")
		data.ToJSON(&GenericError{Message: err.Error()}, w)
		return
	}
	data.ToJSON(&procedure, w)
}

// DeleteStep handles DELETE requests to existing steps
func (st *Step) DeleteStep(w http.ResponseWriter, r *http.Request) {
	// Get the id from the URL
	vars := mux.Vars(r)
	id := vars["id"]
	// Get the procedure from the DB based on its id
	step, err := st.db.GetStep(context.Background(), id)
	if err != nil {
		st.logger.Error("Unable to fetch step")
		w.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, w)
		return
	}
	// Perform actual procedure deletion
	st.logger.Debug(fmt.Sprintf("Deleting step: %v", step))
	if len(step) > 0 {
		err = st.db.DeleteStep(context.Background(), step[0].ID)
		if err != nil {
			st.logger.Error("Error ocurred while deleting procedure")
			data.ToJSON(&GenericError{Message: err.Error()}, w)
			return
		}
		data.ToJSON(&step[0], w)
	} else {
		st.logger.Info(fmt.Sprintf("No steps found using id: %s", id))
		return
	}

}
