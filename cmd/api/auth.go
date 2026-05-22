package main

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/google/uuid"
	"github.com/seminhnva/e-com/internal/password"
	"github.com/seminhnva/e-com/internal/store"
)

type RegisterUserPayLoad struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

//

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayLoad
	if err := readJSON(w, r, &payload); err != nil {
		app.badReqestResponse(w, r, err)
		return
	}
	if err := Validator.Struct(payload); err != nil {
		app.badReqestResponse(w, r, err)
	}

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}
	hassPw, err := password.Hash(payload.Password)
	if err != nil {
		app.handleError(w, r, err)
	}
	user.Password = hassPw
	plainToken := uuid.New().String()
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])
	// store hash token
	err = app.store.Users.CreateAndInvite(r.Context(), user, hashToken, app.config.mail.exp)
	if err != nil {
		app.handleError(w, r, err)
	}

	//mail
	if err := app.jsonResponse(w, http.StatusCreated, nil); err != nil {
		app.handleError(w, r, err)
	}

}
