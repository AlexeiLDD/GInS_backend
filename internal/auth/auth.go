package auth

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserCredentials struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Username string    `json:"name"`
	Password string    `json:"password"`
	Phone    string    `json:"phone"`
	Telegram string    `json:"telegram"`
	AvatarId uuid.UUID `json:"avatarId"`
}

type LoginResponce struct {
	ID       string `json:"Id"`
	Email    string `json:"Email"`
	Username string `json:"Name"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
