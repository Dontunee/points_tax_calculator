package main

import (
	"fmt"
	"net/http"
)

func (app *application) logError(err error) {
	app.logger.Println(err)
}

func (app *application) errorResponse(w http.ResponseWriter, status int, err error) {
	app.logError(err)
	err = app.writeJSON(w, status, err.Error(), nil)
	if err != nil {
		app.logError(fmt.Errorf("an error occurred writing error json response"))
		w.WriteHeader(status)
	}
}
