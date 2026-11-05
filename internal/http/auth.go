package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/theonlysroy/attendai-be-go/internal/model"
	"github.com/theonlysroy/attendai-be-go/internal/storage"
)

type AuthHandler struct {
	userStore storage.StudentStore
	secret    []byte
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
		Subject:   strconv.Itoa(int(userID)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
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
	user := model.User{
		FirstName:  payload.FirstName,
		LastName:   &payload.LastName,
		Role:       payload.Type,
		Email:      payload.Email,
		Department: payload.Department,
	}
	user, err := h.userStore.Create(r.Context(), user)
	if err != nil {
		return err
	}

	// generate tokens
	token, err := h.sign(user.ID, 24*time.Hour)
	if err != nil {
		return err
	}

	responseData := &struct {
		AccessToken string `json:"accessToken"`
	}{
		AccessToken: token,
	}

	WriteJson(w, http.StatusCreated, ApiResponse{Success: true, Message: "in register", Data: responseData})

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
