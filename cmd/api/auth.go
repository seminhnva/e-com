package main

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/google/uuid"
	"github.com/seminhnva/e-com/internal/store"
)

type RegisterUserPayLoad struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type UserWithToken struct {
	*store.User
	Token string `json:"token"`
}

// @Summary		Register a new user
// @Description	Register a new user and send an invitation email
// @Tags			authentication
// @Accept			json
// @Produce		json
// @Param			payload	body		RegisterUserPayLoad	true	"User registration payload"
// @Success		201		{object}	store.User
// @Failure		400		{object}	error	"Invalid request body"
// @Failure		500		{object}	error	"Internal server error"
// @Router			/authentication/user [post]
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
	if err := user.Password.Set(payload.Password); err != nil {
		app.handleError(w, r, err)
		return
	}
	plainToken := uuid.New().String()
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])
	// store hash token
	if err := app.store.Users.CreateAndInvite(r.Context(), user, hashToken, app.config.mail.exp); err != nil {
		app.handleError(w, r, err)
		return
	}

	//mail
	userWithToken := UserWithToken{
		User:  user,
		Token: plainToken,
	}

	if err := app.jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		app.handleError(w, r, err)
	}

}
