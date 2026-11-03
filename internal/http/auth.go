package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/theonlysroy/attendai-be-go/internal/storage"
)

type AuthHandler struct {
	userStore storage.StudentStore
	secret []byte
}

var (
	ErrInvalidRegisterPayload = errors.New("Invalid register payload")
)

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

/*
POST /register
POST /login
POST /refresh-token
GET /profile/me
*/

func (h *AuthHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/register", MakeHTTPHandler(h.RegisterController))

	return r
}

func (h *AuthHandler) sign(userID int64, ttl time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject: userID,
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl))
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(h.secret)
}

func (h *AuthHandler) RegisterController(w http.ResponseWriter, r *http.Request) error {
	var payload RegisterReqPayload

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&payload); err != nil {
		return ErrInvalidRegisterPayload
	}

	// create the user

	// generate tokens
	// save token
	// return response

	WriteJson(w, http.StatusCreated, ApiResponse{Success: true, Message: "in register", Data: payload})
	return nil
}

// helpers
type RegisterReqPayload struct {
	Type       string `json:"type"` // student, teacher
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Email      string `json:"email"`
	Department string `json:"department"`
}
