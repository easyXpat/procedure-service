package handlers

import (
	"context"
	"github.com/easyXpat/procedure-service/data"
	"github.com/gorilla/mux"
	"net/http"
)

// swagger:route GET /procedures procedures listProcedures
// return all procedures from the database
// responses:
// 	200: proceduresResponse

// ListAll handles GET requests and returns all current procedures
func (ph *Procedure) ListAll(w http.ResponseWriter, r *http.Request) {
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