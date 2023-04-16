package service

import (
	"fmt"

	"github.com/ksrnnb/saml-impl/model"
)

const (
	baseURL = "http://localhost:3000"
)

type SamlService struct {
	CompanyID string
}

func NewSamlService(companyID string) SamlService {
	return SamlService{
		CompanyID: companyID,
	}
}

func (s SamlService) SPMetadata() *model.SPMetadata {
	return model.NewSPMetadata(s.SPEntityID(), s.ACSURL(), s.SLOURL())
}

func (ss SamlService) ACSURL() string {
	return fmt.Sprintf("%s/acs/%s", baseURL, ss.CompanyID)
}

func (ss SamlService) SLOURL() string {
	return fmt.Sprintf("%s/slo/%s", baseURL, ss.CompanyID)
}

func (ss SamlService) SPEntityID() string {
	return fmt.Sprintf("%s/%s", baseURL, ss.CompanyID)
}
