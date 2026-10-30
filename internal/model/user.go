package model

import "time"

type User struct {
	ID         int64      `json:"id"`
	FirstName  string     `json:"firstName"`
	LastName   *string    `json:"lastName"`
	FullName   string     `json:"fullName"`
	Role       string     `json:"role"` // roles - student, teacher, admin
	Email      string     `json:"email"`
	Department string     `json:"department"`
	CreatedAt  time.Time  `json:"createdAt"`
	DeletedAt  *time.Time `json:"deletedAt"`
}

type UpdateStudentRequest struct {
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Email     *string `json:"email"`
}
