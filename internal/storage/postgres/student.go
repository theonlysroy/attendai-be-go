package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/theonlysroy/attendai-be-go/internal/model"
)

var (
	ErrInvalidStudentID = errors.New("invalid id")
	ErrStudentNotFound  = errors.New("no student found with this student id")
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
		return model.User{}, err
	}

	return newStudent, nil
}

func (s *StudentStore) UpdateByID(ctx context.Context, id int64) (model.User, error) {
	return model.User{}, nil
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
