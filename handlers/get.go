package handlers

import (
	"context"
	"github.com/easyXpat/procedure-service/data"
	"github.com/gorilla/mux"
	"net/http"
)

// swagger:route GET /procedures procedures getSteps
// return all procedures from the database
// responses:
// 	200: proceduresResponse

// ListAll handles GET requests and returns all current procedures
func (ph *Procedure) GetProcedures(w http.ResponseWriter, r *http.Request) {
	ph.logger.Debug("Fetch all procedures")
	w.Header().Add("Content-Type", "application/json")

	procedures, err := ph.db.GetAllProcedures(context.Background())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, w)
		return
	}

	err = data.ToJSON(procedures, w)
	if err != nil {
		ph.logger.Error("Unable to serializing product", "error", err)
	}
}

// swagger:route GET /procedures/{id} procedures getProcedure
// list single procedure from db
// responses:
// 	200: procedureResponse

// GetProcedure handles GET request for getProcedure
func (ph *Procedure) GetProcedure(w http.ResponseWriter, r *http.Request) {
	ph.logger.Debug("handle getProcedure")
	w.Header().Add("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	procedure, err := ph.db.GetProcedure(context.Background(), id)
	if err != nil {
		ph.logger.Error("Unable to fetch procedure")
		w.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, w)
		return
	}
	ph.logger.Info("Procedure created successfully")
	err = data.ToJSON(procedure, w)
}

// swagger:route GET /steps steps getSteps
// return all steps from the database
// responses:
// 	200: stepsResponse

// ListAll handles GET requests and returns all current procedures
func (st *Step) GetSteps(w http.ResponseWriter, r *http.Request) {
	st.logger.Debug("Fetch all steps")
	w.Header().Add("Content-Type", "application/json")

	steps, err := st.db.GetAllSteps(context.Background())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, w)
		return
	}

	err = data.ToJSON(steps, w)
	if err != nil {
		st.logger.Error("Unable to serializing product", "error", err)
	}
}

// getProcedureSteps returns all steps for a procedure
func (st *Step) GetProcedureSteps(w http.ResponseWriter, r *http.Request) {
	st.logger.Debug("handle getProcedureSteps")
	w.Header().Add("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["procedure"]

	steps, err := st.db.GetProcedureSteps(context.Background(), id)
	if err != nil {
		st.logger.Error("Unable to fetch steps for procedure")
		w.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, w)
		return
	}
	st.logger.Info("Procedure created successfully")
	err = data.ToJSON(steps, w)
}

func (st *Step) GetStep(w http.ResponseWriter, r *http.Request) {
	st.logger.Debug("handle GetStep")
	w.Header().Add("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	step, err := st.db.GetStep(context.Background(), id)
	if err != nil {
		st.logger.Error("Unable to fetch steps for procedure")
		w.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, w)
		return
	}
	st.logger.Info("Procedure created successfully")
	err = data.ToJSON(step, w)
}