package service

import (
	"errors"

	"github.com/crewjam/saml"
	"github.com/crewjam/saml/samlsp"
	"github.com/ksrnnb/saml-impl/model"
)

const nameIDFormat = "urn:oasis:names:tc:SAML:2.0:nameid-format:emailAddress"

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

	// SLO URL は HTTP POST Binding のみ対応
	slourl := ""
	for _, s := range idp.SingleLogoutServices {
		if s.Binding == saml.HTTPPostBinding {
			slourl = s.Location
			break
		}
	}
	if slourl == "" {
		return nil, errors.New("slo url is not specified")
	}

	cert := ""
	for _, k := range idp.KeyDescriptors {
		if k.Use == "signing" {
			cert = k.KeyInfo.X509Data.X509Certificates[0].Data
			break
		}
	}
	if cert == "" {
		return nil, errors.New("idp certificate is not specified")
	}

	return map[string]string{
		"idpEntityId":    md.EntityID,
		"ssoUrl":         ssourl,
		"sloUrl":         slourl,
		"idpCertificate": cert,
	}, nil
}

func (s SamlIdPService) BuildIdPEntityDescriptor(md *model.IdPMetadata) (*saml.EntityDescriptor, error) {
	idpMD, err := model.FindMetadtaByCompanyID(md.CompanyID)
	if err != nil {
		return nil, err
	}

	return &saml.EntityDescriptor{
		EntityID: idpMD.EntityID,
		IDPSSODescriptors: []saml.IDPSSODescriptor{
			{
				SSODescriptor: saml.SSODescriptor{
					RoleDescriptor: saml.RoleDescriptor{
						ProtocolSupportEnumeration: "urn:oasis:names:tc:SAML:2.0:protocol",
						KeyDescriptors: []saml.KeyDescriptor{
							{
								Use: "signing",
								KeyInfo: saml.KeyInfo{
									X509Data: saml.X509Data{
										X509Certificates: []saml.X509Certificate{
											{Data: idpMD.Certificate},
										},
									},
								},
							},
						},
					},
					NameIDFormats: []saml.NameIDFormat{saml.NameIDFormat(saml.EmailAddressNameIDFormat)},
					SingleLogoutServices: []saml.Endpoint{
						{
							Binding:  saml.HTTPPostBinding,
							Location: md.SLOURL,
						},
					},
				},
				SingleSignOnServices: []saml.Endpoint{
					{
						Binding:  saml.HTTPRedirectBinding,
						Location: idpMD.SSOURL,
					},
				},
			},
		},
	}, nil
}
