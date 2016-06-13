package users

import (
	"github.com/crob1140/codewiz/datastore"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	datastore.BaseRecord
	Username string `db:"Username"`
	Email string 	`db:"Email"`
	HashedPassword string `db:"Password"`
}

func NewUser(username string, password string, email string) *User {
	user := &User{Username : username, Email : email}
	user.SetPassword(password)
	return user
}

func (user *User) SetPassword(password string) error {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.HashedPassword = string(pass)
	return nil
}

func (user *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	return (err == nil)
}