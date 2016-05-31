package models

import (
	"github.com/crob1140/codewiz/datastore"
)

type UserDao struct {
	DataStore *datastore.SQLDataStore
}

type User struct {
	Username string `db:"USERNAME"`
	HashedPassword string `db:"PASSWORD"`
}

func NewUserDao(dataStore *datastore.SQLDataStore) *UserDao {
	return &UserDao{DataStore : dataStore}
}

func NewUser(username string, password string, email string) *User {
	return &User{Username : username}
}


func (dao *UserDao) GetUser(username string) *User {
	return nil
}

func (dao *UserDao) AddUser(user *User) error {
	return nil
}

func (user *User) Password() {

}

func (user *User) SetPassword(password string) {

}