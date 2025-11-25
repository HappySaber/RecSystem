package models

type User struct {
	ID       string
	Email    string
	Name     string
	Surname  string
	PassHash []byte
	Role     string
}
