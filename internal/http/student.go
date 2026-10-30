package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/theonlysroy/attendai-be-go/internal/model"
	"github.com/theonlysroy/attendai-be-go/internal/storage"
)

type StudentHandler struct {
	store storage.StudentStore
}

var (
	ErrStudentIdMissing = errors.New("Student id is required")
	ErrInvalidReqBody   = errors.New("Invalid request body")
)

func NewStudentHandler(store storage.StudentStore) *StudentHandler {
	return &StudentHandler{store: store}
}

func (h *StudentHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", MakeHTTPHandler(h.GetAllStudents))
	r.Get("/{id}", MakeHTTPHandler(h.GetStudentById))
	r.Post("/", MakeHTTPHandler(h.CreateStudent))
	r.Post("/update/{id}", MakeHTTPHandler(h.UpdateStudentById))
	r.Delete("/{id}", MakeHTTPHandler(h.DeleteStudentById))

	return r

}

func (h *StudentHandler) DeleteStudentById(w http.ResponseWriter, r *http.Request) error {
	studentId, err := UrlVars(r, "id")
	if err != nil {
		return ErrStudentIdMissing
	}

	err = h.store.DeleteByID(r.Context(), studentId)
	if err != nil {
		return err
	}

	WriteJson(w, http.StatusOK, ApiResponse{Success: true, Message: "Student deleted.", Data: studentId})
	return nil
}

func (h *StudentHandler) UpdateStudentById(w http.ResponseWriter, r *http.Request) error {
	studentId, err := UrlVars(r, "id")
	if err != nil {
		return ErrStudentIdMissing
	}

	var req model.UpdateStudentRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		return ErrInvalidReqBody
	}

	updatedStudent, err := h.store.UpdateByID(r.Context(), studentId, req)
	if err != nil {
		return err
	}
	WriteJson(w, http.StatusAccepted, ApiResponse{Success: true, Message: "Student details updated.", Data: updatedStudent})
	return nil
}

func (h *StudentHandler) CreateStudent(w http.ResponseWriter, r *http.Request) error {
	var studentObj model.User

	err := json.NewDecoder(r.Body).Decode(&studentObj)
	if err != nil {
		return err
	}

	// fmt.Printf("req body -> %v\n", studentObj)

	newStudent, err := h.store.Create(r.Context(), studentObj)
	if err != nil {
		return err
	}

	WriteJson(w, http.StatusCreated, ApiResponse{Success: true, Message: "Student created", Data: newStudent})
	return nil
}

func (h *StudentHandler) GetAllStudents(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()

	pageStr := query.Get("page")
	limitStr := query.Get("limit")

	page := 1
	limit := 5

	if pageStr != "" {
		p, _ := strconv.Atoi(pageStr)
		page = p
	}

	if limitStr != "" {
		l, _ := strconv.Atoi(limitStr)
		limit = l
	}

	offset := (page - 1) * limit

	students, total, err := h.store.List(r.Context(), limit, offset)
	if err != nil {
		return err
	}

	meta := NewPaginationMeta(total, page, limit)
	WriteJson(w, http.StatusOK, PaginatedResponse[model.User]{Success: true, Message: "Students fetched.", Data: students, Meta: meta})
	return nil
}

func (h *StudentHandler) GetStudentById(w http.ResponseWriter, r *http.Request) error {
	studentId, err := UrlVars(r, "id")
	if err != nil {
		return ErrStudentIdMissing
	}
	student, err := h.store.GetByID(r.Context(), studentId)
	if err != nil {
		return err
	}
	WriteJson(w, http.StatusOK, ApiResponse{Success: true, Message: "Student fetched.", Data: student})
	return nil
}

// helpers

func UrlVars(r *http.Request, param string) (int64, error) {
	id := chi.URLParam(r, param)
	studentId, err := strconv.Atoi(id)
	if err != nil {
		return -1, err
	}
	return int64(studentId), nil
}
