package handlers

import (
	"context"
	"github.com/easyXpat/procedure-service/data"
	"net/http"
)

// swagger:route GET /procedures procedures listProcedures
// return all procedures from the database

// responses
// 	200: proceduresResponse

// ListAll handles GET requests and returns all current procedures
func (p *Procedure) ListAll(w http.ResponseWriter, r *http.Request) {
	p.logger.Debug("Fetch all procedures")
	w.Header().Add("Content-Type", "application/json")

	procedures, err := p.db.GetAllProcedures(context.Background())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, w)
		return
	}

	err = data.ToJSON(procedures, w)
	if err != nil {
		p.logger.Error("Unable to serializing product", "error", err)
	}
}