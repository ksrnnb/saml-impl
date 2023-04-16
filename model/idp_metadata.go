package model

// Metadata is metadata of IdP
type IdPMetadata struct {
	ID          int `gorm:"primaryKey"`
	CompanyID   string
	EntityID    string
	Certificate string
	SSOURL      string
}

func (m IdPMetadata) TableName() string {
	return "idp_metadatas"
}

func NewIdPMetadata(cid string, eid string, cert string, ssourl string) *IdPMetadata {
	return &IdPMetadata{
		CompanyID:   cid,
		EntityID:    eid,
		Certificate: cert,
		SSOURL:      ssourl,
	}
}

func FindMetadtaByCompanyID(cid string) (*IdPMetadata, error) {
	m := &IdPMetadata{}
	if err := db.Limit(1).Find(m).Error; err != nil {
		return nil, err
	}
	return m, nil
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
	return m.ID != 0 && m.CompanyID != "" && m.EntityID != "" && m.Certificate != "" && m.SSOURL != ""
}
