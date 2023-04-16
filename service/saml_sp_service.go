package service

import (
	"fmt"
	"net/url"

	"github.com/ksrnnb/saml-impl/model"
)

const (
	baseURL = "http://localhost:3000"
)

type SamlSPService struct {
	CompanyID string
}

func NewSamlSPService(companyID string) SamlSPService {
	return SamlSPService{
		CompanyID: companyID,
	}
}

func (s SamlSPService) SPMetadata() *model.SPMetadata {
	return model.NewSPMetadata(s.SPEntityID().String(), s.ACSURL().String(), s.SLOURL().String())
}

func (ss SamlSPService) ACSURL() *url.URL {
	u, _ := url.Parse(fmt.Sprintf("%s/acs/%s", baseURL, ss.CompanyID))
	return u
}

func (ss SamlSPService) SLOURL() *url.URL {
	u, _ := url.Parse(fmt.Sprintf("%s/slo/%s", baseURL, ss.CompanyID))
	return u
}

func (ss SamlSPService) SPEntityID() *url.URL {
	u, _ := url.Parse(fmt.Sprintf("%s/%s", baseURL, ss.CompanyID))
	return u
}
