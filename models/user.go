package models

type User struct {
	ID                   string `json:"id" db:"id"`
	Email                string `json:"email" db:"email"`
	PasswordHash         string `json:"password_hash" db:"password_hash"`
	Password             string `json:"-" db:"-"`
	PasswordConfirmation string `json:"-" db:"-"`
}

// Users is not required by pop and may be deleted
type Users []User
