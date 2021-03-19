package handlers

import (
	"context"
	"github.com/easyXpat/procedure-service/data"
	"net/http"
)

func (ph *Procedure) MiddlewareValidateProcedure(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "application/json")

		ph.logger.Debug("procedure JSON: %s", r.Body)
		procedure := &data.Procedure{}

		// deserialize procedure
		err := data.FromJSON(procedure, r.Body)
		if err != nil {
			ph.logger.Error("deserialization of procedure json failed", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: err.Error()}, w)
			return
		}

		// validate the procedure
		errs := ph.validator.Validate(procedure)
		if len(errs) != 0 {
			ph.logger.Error("validation of procedure JSON failed", "error", errs)
			w.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&ValidationError{Messages: errs.Errors()}, w)
			return
		}

		// add procedure to the context
		ctx := context.WithValue(r.Context(), ProcedureKey{}, *procedure)
		r = r.WithContext(ctx)

		// call next handler
		next.ServeHTTP(w, r)
	})
}
