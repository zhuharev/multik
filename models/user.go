package models

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordIsNil = fmt.Errorf("Password is nil")
)

type User struct {
	Id int64

	Permission UserPermission

	Name           string `xorm:"unique index"`
	Password       string `xorm:"-"`
	HashedPassword []byte
}

func (x *Engine) UserCreate(username, password string) (*User, error) {
	u := &User{
		Name:     username,
		Password: password,
	}
	e := x.userSave(u)
	return u, e
}

func (x *Engine) userSave(u *User) error {
	if u.Password == "" && u.HashedPassword == nil {
		return ErrPasswordIsNil
	}
	hp, e := HashPassword(u.Password)
	if e != nil {
		return e
	}
	u.HashedPassword = hp
	e = x.Save(u)
	return e
}

func CheckHashedPassword(pass string, hashed []byte) bool {
	e := bcrypt.CompareHashAndPassword(hashed, []byte(pass))
	if e != nil {
		return false
	}
	return true
}

func HashPassword(pass string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return []byte{}, err
	}
	return hashedPassword, nil
}
