package entity

type Apichats struct {
	Name      string `gorm:"column:name"`
	Token     string `gorm:"column:token"`
	ExpiresAt int    `gorm:"column:expiresAt"`
}
