package handlers

import (
	"context"
	"fmt"
	"github.com/easyXpat/procedure-service/data"
	"github.com/stripe/stripe-go"
	c "github.com/stripe/stripe-go/charge"
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
	if procedure.StepsMapping != nil {
		err = ph.db.MapProcedureSteps(context.Background(), procedure.ID, procedure.StepsMapping)
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

// swagger:route POST /charge charge createCharge
// Create a new charge
//
// responses:
//	200: chargeResponse

// CreateCharge handles POST requests to create new charge
func (pt *Charge) CreateCharge(w http.ResponseWriter, r *http.Request) {
	// fetch the payment attributes from the context
	w.Header().Set("Content-Type", "application/json")

	pt.logger.Debug("charge JSON: %s", r.Body)
	payment := &data.Charge{}

	// deserialize procedure
	err := data.FromJSON(payment, r.Body)
	if err != nil {
		pt.logger.Error("deserialization of charge json failed", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, w)
		return
	}

	pt.logger.Debug(fmt.Sprintf("Creating payment: %v", payment))
	apiKey := "sk_test_51JPx98CQzOA1I1zmcKCa7y77j4bWvl7GZHobBqJp623iTha6WQ9IU4IbMmXDRKEC90mEmkS1RTXtrNA4VMMM8qmR00Q62zyLNd"
	stripe.Key = apiKey

	_, err = c.New(&stripe.ChargeParams{
		Amount:       stripe.Int64(payment.Amount),
		Currency:     stripe.String(string(stripe.CurrencyUSD)),
		Description:  stripe.String(payment.ProcedureID),
		Source:       &stripe.SourceParams{Token: stripe.String("tok_visa")},
		ReceiptEmail: stripe.String(payment.ReceiptEmail)})

	if err != nil {
		data.ToJSON(&GenericError{Message: err.Error()}, w)
		return
	}

	err = pt.SavePayment(payment)
	if err != nil {
		data.ToJSON(&GenericError{Message: err.Error()}, w)
	}
	data.ToJSON(&payment, w)
}