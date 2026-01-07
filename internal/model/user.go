package model

type User struct {
<<<<<<< Updated upstream
	NickName string `gorm:"column:nick_name;type:varchar(255);not null"`
	Email    string `gorm:"column:email;type:varchar(255);not null;unique"`
	Phone    *string `gorm:"column:phone;type:varchar(20);unique"`
	Password string `gorm:"column:password;type:varchar(255);not null"`
	CommonPartUnique
=======
	UUID     string  `gorm:"primaryKey;type:varchar(36)"`
	NickName string  `gorm:"column:nick_name;type:varchar(255);not null"`
	Email    string  `gorm:"column:email;type:varchar(255);not null;unique"`
	Phone    *string `gorm:"column:phone;type:varchar(20);unique"`
	Password string  `gorm:"column:password;type:varchar(255);not null"`
	CommonPartNoUnique
>>>>>>> Stashed changes
}

func (User) TableName() string { return "user" }
