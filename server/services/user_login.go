package services

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserLogin struct {
	Username       string
	HashedPassword string
	Role           string
}

func NewUser(username string, password string, role string) (*UserLogin, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("cannot hash password: %w", err)
	}

	userLogin := &UserLogin{
		Username:       username,
		HashedPassword: string(hashedPassword),
		Role:           role,
	}

	return userLogin, nil
}

func (userLogin *UserLogin) IsCorrectPassword(password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(userLogin.HashedPassword), []byte(password))
    return err == nil
}

func (userLogin *UserLogin) Clone() *UserLogin {
    return &UserLogin{
		HashedPassword: userLogin.HashedPassword,
        Username:       userLogin.Username,
        Role:           userLogin.Role,
    }
}

