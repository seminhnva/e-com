package main

import (
	"net/http"
)

// healthCheckHandler godoc
//
//	@Summary		Health check
//	@Description	Returns the health status of the API
//	@Tags			ops
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": app.config.version,
	}
	if err := writeJSON(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
