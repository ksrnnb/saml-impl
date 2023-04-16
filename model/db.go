package model

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const dbFileName = "sqlite.db"

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open(sqlite.Open(dbFileName), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(
		&Company{},
		&IdPMetadata{},
		&User{},
	)
	if err != nil {
		panic("failed to migrate database")
	}
	createDemoDataIfNeeded()
}

func createDemoDataIfNeeded() {
	c := &Company{}
	res := db.Limit(1).Find(c, "id = ?", defaultCompanyID)
	if res.Error != nil {
		panic(res.Error)
	}
	if res.RowsAffected == 0 {
		c = &Company{ID: defaultCompanyID}
		if err := db.Create(c).Error; err != nil {
			panic(err)
		}
	}

	u := &User{}
	res = db.Limit(1).Find(u, "id = ?", defaultUserID)
	if res.Error != nil {
		panic(res.Error)
	}
	if res.RowsAffected == 0 {
		u := User{
			ID:        defaultUserID,
			Password:  defaultPassword,
			CompanyID: c.ID,
			Email:     defaultEmail,
		}
		if err := db.Create(&u).Error; err != nil {
			panic(err)
		}
	}

	u = &User{}
	res = db.Limit(1).Find(u, "id = ?", adminUserID)
	if res.Error != nil {
		panic(res.Error)
	}
	if res.RowsAffected == 0 {
		u := User{
			ID:        adminUserID,
			Password:  adminPassword,
			CompanyID: c.ID,
			Email:     adminEmail,
		}
		if err := db.Create(&u).Error; err != nil {
			panic(err)
		}
	}
}
