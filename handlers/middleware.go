package handlers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/easyXpat/procedure-service/data"
)

func (ph *Procedure) MiddlewareValidateProcedure(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		ph.logger.Debug("procedure JSON: %s", r.Body)
		procedure := &data.Procedure{}

		// deserialize procedure
		err := data.FromJSON(procedure, r.Body)
		if err != nil {
			buf := new(bytes.Buffer)
			buf.ReadFrom(r.Body)
			newStr := buf.String()

			ph.logger.Debug(fmt.Sprintf("String to deserialize %s", newStr))
			ph.logger.Debug(fmt.Sprintf("JSON to deserialize: %s", r.Body))
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

func (st *Step) MiddlewareValidateStep(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		st.logger.Debug("step ", "json", r.Body)
		step := &data.Step{}

		// deserialize step
		err := data.FromJSON(step, r.Body)
		if err != nil {
			st.logger.Error("deserialization of step json failed", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: err.Error()}, w)
			return
		}

		// validate the step
		errs := st.validator.Validate(step)
		if len(errs) != 0 {
			st.logger.Error("validation of step JSON failed", "error", errs)
			w.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&ValidationError{Messages: errs.Errors()}, w)
			return
		}

		// add procedure to the context
		ctx := context.WithValue(r.Context(), StepKey{}, *step)
		r = r.WithContext(ctx)

		// call next handler
		next.ServeHTTP(w, r)
	})
}
