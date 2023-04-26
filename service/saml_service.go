package service

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/beevik/etree"
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

func (s *SamlService) MakeAuthnRequestURL(relayState string) (string, *url.URL, error) {
	req, err := s.ServiceProvider.MakeAuthenticationRequest(s.ServiceProvider.GetSSOBindingLocation(saml.HTTPRedirectBinding), saml.HTTPRedirectBinding, saml.HTTPPostBinding)
	if err != nil {
		return "", nil, err
	}
	u, err := req.Redirect(relayState, &s.ServiceProvider)
	return req.ID, u, err
}

func (s *SamlService) ValidateInResponseTo(samlResponse string, requestID string) error {
	// AllowIdPInitiated == true の場合は InResponseTo を検証しないようになっているので自前で検証する
	samlRes, err := s.BuildSamlResponse(samlResponse)
	if err != nil {
		return err
	}
	if samlRes.InResponseTo == "" {
		return nil
	}

	if samlRes.InResponseTo != requestID {
		return errors.New("invalid Request ID")
	}
	return nil
}

func (s *SamlService) BuildSamlResponse(samlResponse string) (saml.Response, error) {
	decodedResponseXML, err := base64.StdEncoding.DecodeString(samlResponse)
	if err != nil {
		return saml.Response{}, err
	}

	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(decodedResponseXML); err != nil {
		return saml.Response{}, err
	}

	var response saml.Response
	if err := unmarshalElement(doc.Root(), &response); err != nil {
		return saml.Response{}, fmt.Errorf("cannot unmarshal response: %v", err)
	}
	return response, nil
}

func (s *SamlService) ParseResponse(r *http.Request, possibleRequestIDs []string) (*Assertion, error) {
	sa, err := s.ServiceProvider.ParseResponse(r, possibleRequestIDs)
	if err != nil {
		return nil, err
	}
	return &Assertion{sa: sa}, nil
}

func (s *SamlService) ParseLogoutRequest(samlrequest string) (*LogoutRequest, error) {
	rawRequestBuf, err := base64.StdEncoding.DecodeString(samlrequest)
	if err != nil {
		return nil, err
	}
	var logoutRequest saml.LogoutRequest
	if err := xml.Unmarshal(rawRequestBuf, &logoutRequest); err != nil {
		return nil, err
	}
	if logoutRequest.Destination != s.ss.SLOURL().String() {
		return nil, errors.New("invalid Destination")
	}
	if logoutRequest.Issuer.Value != s.md.EntityID {
		return nil, errors.New("invalid Issuer")
	}
	// TODO: 署名の検証
	return &LogoutRequest{lr: logoutRequest}, nil
}

type Assertion struct {
	sa *saml.Assertion
}

func (a *Assertion) NameID() string {
	return a.sa.Subject.NameID.Value
}

func (a *Assertion) SessionTTL() int {
	for _, stmt := range a.sa.AuthnStatements {
		if !stmt.SessionNotOnOrAfter.IsZero() {
			ttl := stmt.SessionNotOnOrAfter.Sub(time.Now())
			return int(ttl.Seconds())
		}
	}
	return 0
}

type LogoutRequest struct {
	lr saml.LogoutRequest
}

func (r *LogoutRequest) NameID() string {
	return r.lr.NameID.Value
}

func (r *LogoutRequest) ID() string {
	return r.lr.ID
}

func elementToBytes(el *etree.Element) ([]byte, error) {
	namespaces := map[string]string{}
	for _, childEl := range el.FindElements("//*") {
		ns := childEl.NamespaceURI()
		if ns != "" {
			namespaces[childEl.Space] = ns
		}
	}

	doc := etree.NewDocument()
	doc.SetRoot(el.Copy())
	for space, uri := range namespaces {
		doc.Root().CreateAttr("xmlns:"+space, uri)
	}

	return doc.WriteToBytes()
}

// unmarshalElement serializes el into v by serializing el and then parsing it with xml.Unmarshal.
func unmarshalElement(el *etree.Element, v interface{}) error {
	buf, err := elementToBytes(el)
	if err != nil {
		return err
	}
	return xml.Unmarshal(buf, v)
}
