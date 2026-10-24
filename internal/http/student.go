package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/theonlysroy/attendai-be-go/internal/storage"
)

type StudentHandler struct {
	store storage.StudentStore
}

func NewStudentHandler(store storage.StudentStore) *StudentHandler {
	return &StudentHandler{store: store}
}

func (h *StudentHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", MakeHTTPHandler(h.GetAllStudents))
	return r

}

func (h *StudentHandler) GetAllStudents(w http.ResponseWriter, r *http.Request) error {
	// h.store.List()
	WriteJson(w, http.StatusOK, ApiResponse{Success: true, Message: "Students fetched", Data: nil})
	return nil
}
