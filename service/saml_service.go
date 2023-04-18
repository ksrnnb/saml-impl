package service

import (
	"errors"
	"net/url"

	"github.com/crewjam/saml"
	"github.com/crewjam/saml/samlsp"
	"github.com/ksrnnb/saml-impl/model"
)

// SamlService is wrapper of saml
type SamlService struct {
	*samlsp.Middleware
	ss SamlSPService
	is SamlIdPService
	md *model.IdPMetadata
}

const (
	supportIdPInitiated = true

	supportedNameIDFormat = saml.EmailAddressNameIDFormat
)

func NewSamlService(companyID string) (*SamlService, error) {
	md, err := model.FindMetadtaByCompanyID(companyID)
	if err != nil {
		return nil, err
	}
	if md == nil {
		return nil, errors.New("metadata is not found")
	}

	ss := NewSamlSPService(md.CompanyID)
	is := NewSamlIdPService()
	ied, err := is.BuildIdPEntityDescriptor(md)
	if err != nil {
		return nil, err
	}

	samlsp, _ := samlsp.New(samlsp.Options{
		EntityID:          ss.SPEntityID().String(),
		AllowIDPInitiated: supportIdPInitiated,
		IDPMetadata:       ied,
	})
	samlsp.ServiceProvider.AcsURL = *ss.ACSURL()
	samlsp.ServiceProvider.SloURL = *ss.SLOURL()
	samlsp.ServiceProvider.AuthnNameIDFormat = supportedNameIDFormat

	return &SamlService{samlsp, ss, is, md}, nil
}

func (s *SamlService) MakeAuthnRequestURL(relayState string) (*url.URL, error) {
	return s.ServiceProvider.MakeRedirectAuthenticationRequest(relayState)
}
