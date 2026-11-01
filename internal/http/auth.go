package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type AuthHandler struct {
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

func (h *AuthHandler) RegisterController(w http.ResponseWriter, r *http.Request) error {
	var payload RegisterReqPayload

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&payload); err != nil {
		return ErrInvalidRegisterPayload
	}

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
