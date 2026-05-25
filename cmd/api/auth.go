package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/seminhnva/e-com/internal/mailer"
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
type CreateUserTokenPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
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

	userWithToken := UserWithToken{
		User:  user,
		Token: plainToken,
	}
	activationURL := fmt.Sprintf("%s/confirm/%s", app.config.frontendURL, plainToken)
	isProEnv := app.config.env == "production"
	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: activationURL,
	}
	//sned mail
	status, err := app.mailer.Send(mailer.UserWelcomeTemplate, user.Username, user.Email, vars, !isProEnv)

	if err != nil {
		app.logger.Errorw("error sending welcome email", "error", err)
		// rollback user reaction if email fails (saga pattern)
		if err := app.store.Users.Delete(context.Background(), user.ID); err != nil {
			app.logger.Errorw("error deleting user", "error", err)
		}
		app.handleError(w, r, err)
		return
	}
	app.logger.Info("Email sent", "status code", status)

	if err := app.jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		app.handleError(w, r, err)
	}

}

// @Summary		Create authentication token
// @Description	Authenticates a user with email and password, returns a JWT token
// @Tags			authentication
// @Accept			json
// @Produce		json
// @Param			payload	body		CreateUserTokenPayload	true	"User credentials"
// @Success		201		{string}	string					"JWT token"
// @Failure		400		{object}	error					"Invalid request body"
// @Failure		401		{object}	error					"Unauthorized - invalid credentials"
// @Failure		500		{object}	error					"Internal server error"
// @Router			/authentication/token [post]
func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {
	//parses payload credentials
	var payload CreateUserTokenPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.handleError(w, r, err)
		return
	}
	if err := Validator.Struct(payload); err != nil {
		app.handleError(w, r, err)
		return
	}
	//fetch the user(check if the user exist) from the payload
	user, err := app.store.Users.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.unauthorizedErrorResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	// verify password
	match, err := user.Password.Compare(payload.Password)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if !match {
		app.unauthorizedErrorResponse(w, r, fmt.Errorf("invalid credentials"))
		return
	}

	// check account is activated
	if !user.IsActive {
		app.unauthorizedErrorResponse(w, r, fmt.Errorf("account is not activated"))
		return
	}

	// generate the token -> add claimns
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(app.config.auth.token.exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": app.config.auth.token.iss,
		"aud": app.config.auth.token.iss,
	}
	token, err := app.authenticator.GenerateToken(claims)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	// send it to the client

	if err := app.jsonResponse(w, http.StatusCreated, token); err != nil {
		app.handleError(w, r, err)
	}

}
