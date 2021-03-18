package main

import (
	"context"
	"fmt"
	"github.com/easyXpat/procedure-service/config"
	"github.com/easyXpat/procedure-service/store"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	logger := config.NewLogger()
	logger.Info("Initializing Procedure API")

	configs := config.NewConfiguration(logger)
	fmt.Println(configs.DBHost)

	sm := mux.NewRouter()

	pg := store.NewPGClient(logger)
	pg.CreateProcedureDB()

	svr := http.Server{
		Addr:         configs.ServerAddress,
		Handler:      sm,
		ErrorLog:     logger.StandardLogger(&hclog.StandardLoggerOptions{}),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// start the server
	go func() {
		logger.Info("starting the server at port", configs.ServerAddress)

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
	svr.Shutdown(ctx)



}