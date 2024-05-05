package model

type User struct {
	Id            string `gorm:"column:id"`
	Email         string `json:"email" gorm:"column:email"`
	Password      string `json:"password" gorm:"column:password"`
	LastName      string `json:"lastName" gorm:"column:lastName"`
	FirstName     string `json:"firstName" gorm:"column:firstName"`
	MiddleName    string `json:"middleName" gorm:"column:middleName"`
	Specialty     string `json:"specialty" gorm:"column:specialty"`
	AccessTokenID string `json:"accessTokenID" gorm:"column:accessTokenID"`
}

func NewUser(email, password, lastName, firstName, middleName, specialty string) *User {
	return &User{
		Email:      email,
		Password:   password,
		LastName:   lastName,
		FirstName:  firstName,
		MiddleName: middleName,
		Specialty:  specialty,
	}
}
