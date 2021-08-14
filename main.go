package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/easyXpat/procedure-service/config"
	"github.com/easyXpat/procedure-service/data"
	"github.com/easyXpat/procedure-service/handlers"
	"github.com/easyXpat/procedure-service/store"
	"github.com/go-openapi/runtime/middleware"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	//"github.com/rs/cors"
)

const (
	UUIDv4Format = "^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$"
)

func main() {
	logger := config.NewLogger()
	logger.Info("Initializing Procedure API")

	configs := config.NewConfiguration(logger)

	// validator contains all the methods that are needed to validate procedure requests
	validator := data.NewValidation()

	// create a new postgres connection
	db, err := store.NewConnection(logger, configs)
	if err != nil {
		logger.Error("unable to connect to db", "error", err)
		panic(err)
	}

	// creation of procedure table
	db.Exec(context.Background(), store.ProcedureTableDDL)
	db.Exec(context.Background(), store.StepTableDDL)

	// procedure service contains all methods that interact with DB to perform CRUD operations for procedure
	procedureDB := data.NewProcedurePostgres(logger, db)
	stepDB := data.NewStepPostgres(logger, db)

	// handlers encapsulates procedure related services.
	ph := handlers.NewProcedure(logger, procedureDB, validator)
	st := handlers.NewStep(logger, stepDB, validator)

	// handler for documentation
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	//jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options {
	//	ValidationKeyGetter: handlers.ValidationKeyGetter,
	//	SigningMethod: jwt.SigningMethodRS256,
	//})

	// create mux server
	sm := mux.NewRouter()

	//	=== GET methods ===
	getR := sm.Methods(http.MethodGet).Subrouter()
	// docs
	getR.Handle("/docs", sh)
	getR.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))
	// procedures
	getR.HandleFunc("/procedure/{id}", ph.GetProcedure)
	getR.HandleFunc("/procedures", ph.GetProcedures)
	// steps
	getR.HandleFunc("/steps/{id}", st.GetStep)
	getR.HandleFunc("/steps", st.GetSteps)

	//getR.HandleFunc("/procedures/{city}", ph.GetProceduresFromCity)
	//getR.Handle("/procedures", jwtMiddleware.Handler(http.HandlerFunc(ph.GetProcedures)))
	//getR.HandleFunc("/steps/{procedure}", st.GetProcedureSteps)
	//getR.HandleFunc(fmt.Sprintf("/procedures/{id:%s}", UUIDv4Format), ph.GetProcedure)
	//getR.HandleFunc("/procedures", ph.ListAll)

	//	=== POST methods ===
	postR := sm.Methods(http.MethodPost).Subrouter()
	postStepR := sm.Methods(http.MethodPost).Subrouter()
	postR.Use(ph.MiddlewareValidateProcedure)
	postStepR.Use(st.MiddlewareValidateStep)
	// procedures
	postR.HandleFunc("/procedures", ph.CreateProcedure)
	// steps
	postStepR.HandleFunc("/steps", st.CreateStep)

	//	=== PUT methods ===
	putR := sm.Methods(http.MethodPut).Subrouter()
	putR.Use(ph.MiddlewareValidateProcedure)
	putStepR := sm.Methods(http.MethodPut).Subrouter()
	putStepR.Use(st.MiddlewareValidateStep)
	// procedure
	putR.HandleFunc("/procedures", ph.UpdateProcedure)
	// step
	putStepR.HandleFunc("/steps", st.UpdateStep)

	//	=== DELETE methods ===
	deleteR := sm.Methods(http.MethodDelete).Subrouter()
	// procedure
	deleteR.HandleFunc("/procedures/{id}", ph.DeleteProcedure)
	deleteR.HandleFunc("/steps/{id}", st.DeleteStep)


	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))
	//ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}), gohandlers.AllowedMethods([]string{"*"}), gohandlers.AllowedHeaders([]string{"*"}))
	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}

	//corsWrapper := cors.New(cors.Options{
	//	AllowedMethods: []string{"GET", "POST"},
	//	AllowedHeaders: []string{"Content-Type", "Origin", "Accept", "*"},
	//})

	logger.Info("Starting web server", "port", port)
	logger.Info("Test Heroku", "port", port)
	svr := http.Server{
		Addr:         ":"+port,
		//Handler:      corsWrapper.Handler(sm),
		Handler:      ch(sm),
		ErrorLog:     logger.StandardLogger(&hclog.StandardLoggerOptions{}),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// start the server
	go func() {
		err := svr.ListenAndServe()
		if err != nil {
			logger.Error("could not start the server", "error", err)
			os.Exit(1)
		}
	}()

	// look for interrupts for graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c
	logger.Info("shutting down the server", "received signal", sig)

	//gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	defer db.Close(ctx)
	svr.Shutdown(ctx)

}
