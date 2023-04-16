package model

// Metadata is metadata of IdP
type IdPMetadata struct {
	ID          int
	CompanyID   int
	EntityID    string // idp entityID
	Certificate string
	SSOURL      string
}

var IdPMetadataRepo []*IdPMetadata

func NewIdPMetadata(cid int, eid string, cert string, ssourl string) *IdPMetadata {
	return &IdPMetadata{
		CompanyID:   cid,
		EntityID:    eid,
		Certificate: cert,
		SSOURL:      ssourl,
	}
}

func FindMetadtaByCompanyID(cid int) *IdPMetadata {
	for _, m := range IdPMetadataRepo {
		if m.CompanyID == cid {
			return m
		}
	}
	return nil
}

func (m *IdPMetadata) Save() {
	if m.ID == 0 {
		m.ID = len(IdPMetadataRepo) + 1
		IdPMetadataRepo = append(IdPMetadataRepo, m)
		return
	}

	for i, mInRepo := range IdPMetadataRepo {
		if m.ID == mInRepo.ID {
			IdPMetadataRepo[i] = m
			return
		}
	}
}

func (m *IdPMetadata) Valid() bool {
	return m.ID != 0 && m.CompanyID != 0 && m.EntityID != "" && m.Certificate != "" && m.SSOURL != ""
}
