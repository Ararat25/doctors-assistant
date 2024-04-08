package model

type User struct {
	id       int
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (user User) getEmail() string {
	return user.Email
}
