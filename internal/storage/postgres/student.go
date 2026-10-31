package postgres

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/theonlysroy/attendai-be-go/internal/model"
)

var (
	ErrInvalidStudentID         = errors.New("invalid id")
	ErrStudentNotFound          = errors.New("no student found with this student id")
	ErrNoStudentDetailsToUpdate = errors.New("no student details found to update")
	ErrEmailAlreadyExists       = errors.New("email address already exists")
)

type StudentStore struct {
	db *pgxpool.Pool
}

// constuctor
func NewStudentStore(db *pgxpool.Pool) *StudentStore {
	return &StudentStore{db: db}
}

func (s *StudentStore) GetByID(ctx context.Context, id int64) (model.User, error) {
	if id <= 0 {
		return model.User{}, ErrInvalidStudentID
	}

	query := `
		select *
		from users
		where id=$1 AND role='student'
	`
	rows, err := s.db.Query(ctx, query, id)
	if err != nil {
		return model.User{}, err
	}
	defer rows.Close()

	student, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, ErrStudentNotFound
		}
		return model.User{}, err
	}

	return student, nil
}

func (s *StudentStore) List(ctx context.Context, limit int, offset int) ([]model.User, int64, error) {
	countQuery := `
		SELECT COUNT(*)
		FROM users
		WHERE role='student'
	`
	query := `
        SELECT *
        FROM users
        WHERE role='student' 
        LIMIT $1
        OFFSET $2
    `
	var total int64
	err := s.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.User])
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (s *StudentStore) Create(ctx context.Context, st model.User) (model.User, error) {

	var existingStudentId int64

	findEmailQuery := `
		SELECT id
		FROM users
		WHERE email = @email
	`
	s.db.QueryRow(ctx, findEmailQuery, pgx.NamedArgs{"email": st.Email}).Scan(&existingStudentId)

	log.Printf("Existing student id: %d", existingStudentId)
	if existingStudentId != 0 {
		return model.User{}, ErrEmailAlreadyExists
	}

	fullName := st.FirstName
	if st.LastName != nil {
		fullName = fullName + " " + *st.LastName
	}

	query := `
		INSERT INTO users (first_name, last_name, full_name, email, role, department)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, first_name, last_name, full_name, email, role, department, created_at, deleted_at
	`
	rows, err := s.db.Query(ctx, query, st.FirstName, st.LastName, fullName, st.Email, "student", st.Department)
	if err != nil {
		return model.User{}, err
	}

	newStudent, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.User])
	if err != nil {
		return model.User{}, err
	}

	return newStudent, nil
}

func (s *StudentStore) UpdateByID(ctx context.Context, id int64, st model.UpdateStudentRequest) (model.User, error) {
	var student model.User

	findUserQuery := `
		SELECT first_name, last_name, email, deleted_at
		FROM users
		WHERE id = $1
	`
	updateUserquery := `
		UPDATE users SET
		first_name = COALESCE(@first_name, first_name),
		last_name = COALESCE(@last_name, last_name),
		email = COALESCE(@email, email)
		WHERE id = @id
		RETURNING id, first_name, last_name, email
	`

	if st.FirstName == nil && st.LastName == nil && st.Email == nil {
		return model.User{}, ErrNoStudentDetailsToUpdate
	}

	err := s.db.QueryRow(ctx, findUserQuery, id).Scan(&student.FirstName, &student.LastName, &student.Email, &student.DeletedAt)

	if err != nil {
		return model.User{}, ErrStudentNotFound
	}

	if student.DeletedAt != nil {
		return model.User{}, ErrStudentNotFound
	}

	// update query
	err = s.db.QueryRow(ctx, updateUserquery, pgx.NamedArgs{
		"id": id, "first_name": st.FirstName, "last_name": st.LastName, "email": st.Email}).Scan(&student.ID, &student.FirstName, &student.LastName, &student.Email)

	if err != nil {
		return model.User{}, err
	}

	return student, nil
}

func (s *StudentStore) DeleteByID(ctx context.Context, id int64) error {
	query := `
		UPDATE users
		SET deleted_at = NOW()
		WHERE id = $1
	`
	_, err := s.db.Query(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
