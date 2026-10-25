package http

import (
	"net/http"
	"strconv"

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
	r.Get("/{id}", MakeHTTPHandler(h.GetStudentById))
	return r

}

func (h *StudentHandler) GetAllStudents(w http.ResponseWriter, r *http.Request) error {
	students, err := h.store.List(r.Context())
	if err != nil {
		return err
	}
	WriteJson(w, http.StatusOK, ApiResponse{Success: true, Message: "Students fetched.", Data: students})
	return nil
}

func (h *StudentHandler) GetStudentById(w http.ResponseWriter, r *http.Request) error {
	studentId := UrlVars(r, "id")
	student, err := h.store.GetByID(r.Context(), studentId)
	if err != nil {
		return err
	}
	WriteJson(w, http.StatusOK, ApiResponse{Success: true, Message: "Student fetched.", Data: student})
	return nil
}

func UrlVars(r *http.Request, param string) int64 {
	id := chi.URLParam(r, param)
	studentId, _ := strconv.Atoi(id)
	return int64(studentId)
}
