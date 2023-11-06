package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type config struct {
	port             int
	env              string
	taxCalculatorUrl string
}

type application struct {
	config config
	logger *log.Logger
}

func main() {
	startServer()
}

func startServer() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 8000, "server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.taxCalculatorUrl, "taxCalculatorUrl", "http://localhost:5000/tax-calculator", "")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	app := &application{
		config: cfg,
		logger: logger,
	}
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// start the server
	logger.Printf("starting %s server on  %s", cfg.env, srv.Addr)
	err := srv.ListenAndServe()
	if err != nil {
		logger.Fatal(err, "Failure to serve")
	}
}
