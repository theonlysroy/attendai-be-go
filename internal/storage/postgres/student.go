package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/theonlysroy/attendai-be-go/internal/model"
)

type StudentStore struct {
	db *pgxpool.Pool
}

// constuctor
func NewStudentStore(db *pgxpool.Pool) *StudentStore {
	return &StudentStore{db: db}
}

func (s *StudentStore) GetByID(ctx context.Context, id int64) (model.User, error) {
	var student model.User

	query := `
		select id, full_name, email
		from users
		where id = $1
	`
	err := s.db.QueryRow(ctx, query, id).Scan(
		&student.ID,
		&student.FullName,
		&student.Email,
	)

	return student, err
}

func (s *StudentStore) List(ctx context.Context, limit int, offset int) ([]model.User, error) {
	query := `
        SELECT *
        FROM users
        WHERE role='student' 
        LIMIT $1
        OFFSET $2
    `

	rows, err := s.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.User])
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *StudentStore) Create(ctx context.Context, st model.User) (model.User, error) {

	fullName := st.FirstName
	if st.LastName != nil {
		fullName = fullName + " " + *st.LastName
	}

	query := `
		INSERT INTO users (first_name, last_name, full_name, email, role, department)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, first_name, last_name, full_name, email, role, department, created_at
	`
	rows, err := s.db.Query(ctx, query, st.FirstName, st.LastName, fullName, st.Email, "student", st.Department)
	if err != nil {
		return model.User{}, err
	}

	newStudent, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.User])
	if err != nil {
		return model.User{}, nil
	}

	return newStudent, nil
}

func (s *StudentStore) UpdateByID(ctx context.Context, id int64) (model.User, error) {
	return model.User{}, nil
}

func (s *StudentStore) DeleteByID(ctx context.Context, id int64) error {
	return nil
}
