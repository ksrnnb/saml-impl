package model

import "errors"

const defaultUserID = "demo"
const defaultPassword = "&!6Z9@K3f"
const defaultEmail = "demo@test.com"

const adminUserID = "admin"
const adminPassword = "k4s60#lkf"
const adminEmail = "admin@test.com"

type UserType uint8

const (
	UserTypeUnknown UserType = iota
	UserTypeNormal
	UserTypeAdmin
)

func (t UserType) String() string {
	switch t {
	case UserTypeNormal:
		return "Normal"
	case UserTypeAdmin:
		return "Admin"
	default:
		return "Unknown"
	}
}

// User belongs to only one Company
type User struct {
	ID           string `gorm:"primaryKey"`
	Password     string
	CompanyID    string
	Email        string
	PersistentID string
	UserType     UserType
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

func ListAllUsers() ([]*User, error) {
	var users []*User
	if err := db.Order("user_type DESC").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (u *User) ValidatePassword(pwd string) error {
	if u.Password != pwd {
		return errors.New("password is invalid")
	}
	return nil
}
