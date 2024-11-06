package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"git.sr.ht/~rehandaphedar/tilawah-hub/internal/db"
	"git.sr.ht/~rehandaphedar/tilawah-hub/internal/sqlc"
	"git.sr.ht/~rehandaphedar/tilawah-hub/internal/validators"
	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
)

type registerDTO struct {
	Username string `json:"username" validate:"required,min=3,max=64"`
	Password string `json:"password" validate:"required,min=3,max=64"`
}

type loginDTO struct {
	Username string `json:"username" validate:"required,min=3,max=64"`
	Password string `json:"password" validate:"required,min=3,max=64"`
}

// Register godoc
//
//	@Tags		Auth
//	@Accept		json
//	@Produce	json
//	@Param		request	body		registerDTO	true	"Register"
//	@Success	200		{object}	sqlc.AuthInsertUserRow
//	@Failure	400		{object}	models.Error
//	@Failure	500		{object}	models.Error
//	@Router		/register [post]
func Register(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Error reading request body",
			"error":   err.Error(),
		})
		return
	}

	var request registerDTO
	err = json.Unmarshal(body, &request)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Error parsing request JSON",
			"error":   err.Error(),
		})
		return
	}

	err = validators.ValidateStruct(request)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Incorrect username or password",
			"error":   err.Error(),
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, render.M{
			"message": "Error hashing password",
			"error":   err.Error(),
		})
		return
	}

	user, err := db.Queries.AuthInsertUser(context.Background(),
		sqlc.AuthInsertUserParams{
			Username:    request.Username,
			Password:    string(hash),
			Displayname: request.Username})
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "User already exists",
			"error":   err.Error(),
		})
		return
	}

	render.JSON(w, r, user)
}

// Login godoc
//
//	@Tags		Auth
//	@Accept		json
//	@Produce	json
//	@Param		request	body		loginDTO	true	"Login"
//	@Success	200		{object}	sqlc.AuthSelectUserRow
//	@Failure	400		{object}	models.Error
//	@Failure	500		{object}	models.Error
//	@Header		200		{string}	Session	Token	""
//	@Header		200		{string}	CSRF	Token	""
//	@Router		/login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Error reading request body",
			"error":   err.Error(),
		})
		return
	}

	var request loginDTO
	err = json.Unmarshal(body, &request)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Error parsing request JSON",
			"error":   err.Error(),
		})
		return
	}

	err = validators.ValidateStruct(request)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Incorrect username or password",
			"error":   err.Error(),
		})
		return
	}

	user, err := db.Queries.AuthSelectUser(context.Background(), request.Username)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "User does not exist",
			"error":   err.Error(),
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Wrong username or password",
			"error":   err.Error(),
		})
		return
	}

	sessionToken, err := generateToken(32)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, render.M{
			"message": "Error generating session token",
			"error":   err.Error(),
		})
		return
	}

	csrfToken, err := generateToken(32)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, render.M{
			"message": "Error generating csrf token",
			"error":   err.Error(),
		})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(time.Hour * 24),
		Secure:   true,
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Expires:  time.Now().Add(time.Hour * 24),
		Secure:   true,
		HttpOnly: false,
	})

	_, err = db.Queries.AuthInsertSession(context.Background(), sqlc.AuthInsertSessionParams{
		SessionToken: sessionToken,
		CsrfToken:    csrfToken,
		Username:     request.Username,
	})
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, render.M{
			"message": "Error creating session",
			"error":   err.Error(),
		})
		return
	}

	session, err := db.Queries.AuthSelectSession(context.Background(), sessionToken)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "User does not exist",
			"error":   err.Error(),
		})
		return
	}

	if session.SessionToken != sessionToken {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Invalid session token",
			"error":   "Discrepancy between provided and actual session token",
		})
		return
	}

	if session.CsrfToken != csrfToken {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Invalid csrf token",
			"error":   "Discrepancy between session token and csrf token",
		})
		return
	}

	render.JSON(w, r, user)
}

// Logout godoc
//
//	@Tags		Auth
//	@Produce	json
//
//	@Param		X-CSRF-TOKEN	header		string	true	"CSRF Token"
//
//	@Success	200				{object}	sqlc.Session
//	@Failure	400				{object}	models.Error
//	@Failure	401				{object}	models.Error
//	@Failure	500				{object}	models.Error
//	@Router		/logout [post]
func Logout(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Invalid session token",
			"error":   err.Error(),
		})
		return
	}

	session, err := db.Queries.AuthDeleteSession(context.Background(), sessionCookie.Value)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Error while deleting session",
			"error":   err.Error(),
		})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		Secure:   true,
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		Secure:   true,
		HttpOnly: false,
	})

	render.JSON(w, r, session)
}

func generateToken(length int) (string, error) {
	bytes := make([]byte, length)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}
