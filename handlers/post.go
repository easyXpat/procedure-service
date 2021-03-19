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

// Create handles POST requests to add new procedures
func (p *Procedure) Create(w http.ResponseWriter, r *http.Request) {
	// fetch the product from the context
	procedure := r.Context().Value(ProcedureKey{}).(data.Procedure)
	p.logger.Debug(fmt.Sprintf("Inserting procedure: %v", procedure))
	err := p.db.AddProcedure(context.Background(), &procedure)
	if err != nil {
		data.ToJSON(&GenericError{Message: err.Error()}, w)
		return
	}
	data.ToJSON(&procedure, w)
}