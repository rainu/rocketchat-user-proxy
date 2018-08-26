package messages

import (
	"crypto/sha256"
	"fmt"
)

type User struct {
	Username string `json:"username"`
}

type Password struct {
	Digest    string `json:"digest"`
	Algorithm string `json:"algorithm"`
}

type LoginParam struct {
	User     User     `json:"user"`
	Password Password `json:"password"`
}

func NewLoginPlain(username, password string) *MethodCall {
	return NewLoginHash(username, fmt.Sprintf("%x", sha256.Sum256([]byte(password))))
}

func NewLoginHash(username, passwordHash string) *MethodCall {
	params := []interface{}{
		LoginParam{
			User: User{
				Username: username,
			},
			Password: Password{
				Digest:    passwordHash,
				Algorithm: "sha-256",
			},
		},
	}

	return NewMethodCall("login", params)
}
