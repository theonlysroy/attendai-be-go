package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/theonlysroy/attendai-be-go/internal/model"
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
	r.Post("/", MakeHTTPHandler(h.CreateStudent))

	return r

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

// func GetPaginationVars(r *http.Request) (int64, int64, error) {
// 	// return page, limit
// 	query := r.URL.Query()

// 	pageStr := query.Get("page")
// 	limitStr := query.Get("limit")

// 	page := 1
// 	limit := 5

// 	if pageStr != "" {
// 		p, err := strconv.Atoi(pageStr)
// 		if err != nil || p < 1 {
// 			return nil, nil, fmt.Sprintf("Page must be a positive number. Got %v", p)
// 		}
// 		page = p
// 	}

// 	if limitStr != "" {
// 		l, err := strconv.Atoi(limitStr)
// 		if err != nil || l < 1 || l > 100 {
// 			retur
// 		}
// 	}
// }
