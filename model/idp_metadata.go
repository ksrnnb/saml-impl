package model

import (
	"errors"

	"gorm.io/gorm"
)

// Metadata is metadata of IdP
type IdPMetadata struct {
	ID          int `gorm:"primaryKey"`
	CompanyID   string
	EntityID    string
	Certificate string
	SSOURL      string
	SLOURL      string
}

func (m IdPMetadata) TableName() string {
	return "idp_metadatas"
}

func NewIdPMetadata(cid, eid, cert, ssourl, slourl string) *IdPMetadata {
	return &IdPMetadata{
		CompanyID:   cid,
		EntityID:    eid,
		Certificate: cert,
		SSOURL:      ssourl,
		SLOURL:      slourl,
	}
}

func FindMetadtaByCompanyID(cid string) (*IdPMetadata, error) {
	var m IdPMetadata
	err := db.First(&m).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (m *IdPMetadata) Save() error {
	mm, err := FindMetadtaByCompanyID(m.CompanyID)
	if err != nil {
		return err
	}
	if mm != nil {
		m.ID = mm.ID
		return db.Save(m).Error
	}

	return db.Create(m).Error
}

func (m *IdPMetadata) Valid() bool {
	return m.ID != 0 && m.CompanyID != "" && m.EntityID != "" && m.Certificate != "" && m.SSOURL != "" && m.SLOURL != ""
}
