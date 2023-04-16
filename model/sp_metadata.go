package model

type SPMetadata struct {
	EntityID string
	ACSURL   string
	SLOURL   string
}

func NewSPMetadata(eid string, acsurl string, slourl string) *SPMetadata {
	return &SPMetadata{
		EntityID: eid,
		ACSURL:   acsurl,
		SLOURL:   slourl,
	}
}
