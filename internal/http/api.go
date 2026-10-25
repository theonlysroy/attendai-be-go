package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/theonlysroy/attendai-be-go/internal/storage"
)

type APIServer struct {
	listenAddr string
	student    *StudentHandler
}

func NewAPIServer(listenAddr string, studentStore storage.StudentStore) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		student:    NewStudentHandler(studentStore),
	}
}

// start the server
func (s *APIServer) Run() {
	r := chi.NewRouter()

	// student handler
	// sh := NewStuden

	r.Route("/api", func(r chi.Router) {
		r.Get("/healthz", HealthzHandler)
		r.Mount("/students", s.student.Routes())
	})

	log.Printf("Server running on %s", s.listenAddr)
	if err := http.ListenAndServe(s.listenAddr, r); err != nil {
		log.Fatal(err)
	}
}

func HealthzHandler(w http.ResponseWriter, r *http.Request) {
	WriteJson(w, http.StatusOK, ApiResponse{Success: true, Message: "Api running...", Data: nil})
}

//---------------- helpers

type ApiResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type ApiError struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Details any    `json:"details"`
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func WriteJson(w http.ResponseWriter, status int, data any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func MakeHTTPHandler(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJson(w, http.StatusBadRequest, ApiError{Success: false, Message: "Something went wrong", Details: err.Error()})
		}
	}
}
