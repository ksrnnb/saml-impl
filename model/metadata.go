package model

type Metadata struct {
	ID          int
	CompanyID   int
	EntityID    string
	Certificate string
	SSOURL      string
}

var metadataRepo []*Metadata

func NewMetadata(cid int, eid string, cert string, ssourl string) *Metadata {
	return &Metadata{
		CompanyID:   cid,
		EntityID:    eid,
		Certificate: cert,
		SSOURL:      ssourl,
	}
}

func FindMetadtaByCompanyID(cid int) *Metadata {
	for _, m := range metadataRepo {
		if m.CompanyID == cid {
			return m
		}
	}
	return nil
}

func (m *Metadata) Save() {
	if m.ID == 0 {
		m.ID = len(metadataRepo) + 1
		metadataRepo = append(metadataRepo, m)
		return
	}

	for i, mInRepo := range metadataRepo {
		if m.ID == mInRepo.ID {
			metadataRepo[i] = m
			return
		}
	}
}

func (m *Metadata) Valid() bool {
	return m.ID != 0 && m.CompanyID != 0 && m.EntityID != "" && m.Certificate != "" && m.SSOURL != ""
}
