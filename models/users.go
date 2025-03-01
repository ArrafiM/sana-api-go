package models

import (
	"html"
	"sana-api/db"
	"sana-api/utils/token"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"size:255;not null;unique" json:"email"`
	Password string `gorm:"size:255;not null;" json:"password"`
	Phone    string `json:"phone"`
	Name     string `gorm:"type:varchar(255)" json:"name"`
	Role_id  int    `json:"role_id"`
	Picture  string `json:"picture"`
}

func (u *User) SaveUser() (*User, error) {

	var err = db.CON.Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) BeforeSave() error {

	//turn password into hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)

	//remove spaces in email
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))

	return nil

}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

type LoginResponse struct {
	Token  string `json:"token"`
	UserId uint   `json:"user_id"`
}

func LoginCheck(email string, password string) (LoginResponse, error) {

	var err error

	u := User{}

	err = db.CON.Model(User{}).Where("email = ?", email).Take(&u).Error

	if err != nil {
		return LoginResponse{}, err
	}

	err = VerifyPassword(password, u.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return LoginResponse{}, err
	}

	token, err := token.GenerateToken(u.ID)

	if err != nil {
		return LoginResponse{}, err
	}

	res := LoginResponse{
		Token:  token,
		UserId: u.ID,
	}
	return res, nil

}

type ChangePass struct {
	Oldpass        string `json:"oldpass" binding:"required"`
	Newpass        string `json:"newpass" binding:"required"`
	ComfirmNewpass string `json:"confirm_newpass" binding:"required"`
}
