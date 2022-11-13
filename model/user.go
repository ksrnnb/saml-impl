package model

import "errors"

var userStore []*User

const defaultCompanyID = 1

type User struct {
	ID        string
	Password  string
	CompanyID int
	Email     string
}

func init() {
	demoUser := &User{
		ID:        "demo",
		Password:  "password",
		CompanyID: defaultCompanyID,
		Email:     "demo@test.com",
	}
	userStore = append(userStore, demoUser)
}

func FindUser(id string) *User {
	for _, u := range userStore {
		if u.ID == id {
			return u
		}
	}
	return nil
}

func FindUserByEmail(email string) *User {
	for _, u := range userStore {
		if u.Email == email {
			return u
		}
	}
	return nil
}

func (u *User) ValidatePassword(pwd string) error {
	if u.Password != pwd {
		return errors.New("password is invalid")
	}
	return nil
}
