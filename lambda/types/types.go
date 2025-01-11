package types

import (
	"golang.org/x/crypto/bcrypt"
)

type RegisterUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password"`
}

type Event struct {
	Username string `json:"username"`
}

func NewUser(u *RegisterUser) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	return &User{
		Username:     u.Username,
		PasswordHash: string(hashedPassword),
	}, nil
}
