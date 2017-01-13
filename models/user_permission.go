package models

type UserPermission int

const (
	Guest UserPermission = 1 << iota
	Admin
)

func (u User) IsAdmin() bool {
	return u.Permission&Admin == Admin
}
