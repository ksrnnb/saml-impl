package service

import (
	"errors"

	"github.com/crewjam/saml"
	"github.com/crewjam/saml/samlsp"
)

type SamlIdPService struct{}

func NewSamlIdPService() SamlIdPService {
	return SamlIdPService{}
}

func (s SamlIdPService) Parse(data []byte) (map[string]string, error) {
	md, err := samlsp.ParseMetadata(data)
	if err != nil {
		return nil, errors.New("error on parse metadata")
	}
	if len(md.IDPSSODescriptors) == 0 {
		return nil, errors.New("no IDPSSODescriptor")
	}

	idp := md.IDPSSODescriptors[0]

	found := false
	for _, fm := range idp.NameIDFormats {
		if fm == saml.EmailAddressNameIDFormat {
			found = true
			break
		}
	}
	if !found {
		return nil, errors.New("name id format is not email address")
	}

	// SSO URL は HTTP-Redirect のみ対応
	ssourl := ""
	for _, s := range idp.SingleSignOnServices {
		if s.Binding == saml.HTTPRedirectBinding {
			ssourl = s.Location
			break
		}
	}
	if ssourl == "" {
		return nil, errors.New("sso url is not specified")
	}

	cert := ""
	for _, k := range idp.KeyDescriptors {
		if k.Use == "signing" {
			cert = k.KeyInfo.X509Data.X509Certificates[0].Data
			break
		}
	}
	if ssourl == "" {
		return nil, errors.New("idp certificate is not specified")
	}

	return map[string]string{
		"idpEntityId":    md.EntityID,
		"ssoUrl":         ssourl,
		"idpCertificate": cert,
	}, nil
}
