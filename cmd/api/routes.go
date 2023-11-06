package main

import (
	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()
	router.GET("/v1/tax-calculator", app.calculateIncomeTaxHandler)
	return router
}
