package model

import "errors"

const defaultUserID = "demo"
const defaultPassword = "&!6Z9@K3f"
const defaultEmail = "demo@test.com"

// User belongs to only one Company
type User struct {
	ID           string `gorm:"primaryKey"`
	Password     string
	CompanyID    string
	Email        string
	PersistentID string
}

func FindUser(id string) (*User, error) {
	u := &User{}
	if err := db.Limit(1).Find(u, "id=?", id).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func FindUserByEmail(email string) (*User, error) {
	u := &User{}
	if err := db.Limit(1).Find(u, "email=?", email).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (u *User) ValidatePassword(pwd string) error {
	if u.Password != pwd {
		return errors.New("password is invalid")
	}
	return nil
}
