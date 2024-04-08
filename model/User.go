package model

type User struct {
	id         int
	Email      string `json:"email" gorm:"column:email"`
	Password   string `json:"password" gorm:"column:password"`
	LastName   string `json:"lastName" gorm:"column:lastName"`
	FirstName  string `json:"firstName" gorm:"column:firstName"`
	MiddleName string `json:"middleName" gorm:"column:middleName"`
	Specialty  string `json:"specialty" gorm:"column:specialty"`
}

func (user User) getEmail() string {
	return user.Email
}
