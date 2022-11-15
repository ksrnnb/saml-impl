package controller

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"time"

	"github.com/ksrnnb/saml/model"
	"github.com/labstack/echo/v4"
)

type SamlResponse struct {
	Response Response
}

type Response struct {
	Destination string `xml:"Destination,attr"`
	Issuer      string
	Status      Status
	Assertion   Assertion
	Signature   Signature
}

type Status struct {
	StatusCode StatusCode
}

type StatusCode struct {
	Value string `xml:"Value,attr"`
}

type Assertion struct {
	Issuer             string
	Subject            Subject
	Conditions         Conditions
	AuthnStatement     AuthnStatement
	AttributeStatement AttributeStatement
	Signature          Signature
}

type Subject struct {
	NameID              NameID
	SubjectConfirmation SubjectConfirmation
}

type NameID struct {
	Format string `xml:"Format,attr"`
	Value  string `xml:",chardata"`
}

type SubjectConfirmation struct {
	Method                  string `xml:"Method,attr"`
	SubjectConfirmationData SubjectConfirmationData
}

type SubjectConfirmationData struct {
	NotOnOrAfter string `xml:"NotOnOrAfter,attr"`
	Recipient    string `xml:"Recipient,attr"`
}

type Conditions struct {
	NotBefore           string `xml:"NotBefore,attr"`
	NotOnOrAfter        string `xml:"NotOnOrAfter,attr"`
	AudienceRestriction AudienceRestriction
}

type AudienceRestriction struct {
	Audience string
}

type AuthnStatement struct {
	AuthnInstant        string `xml:"AuthnInstant,attr"`
	SessionIndex        string `xml:"SessionIndex,attr"`
	SessionNotOnOrAfter string `xml:"SessionNotOnOrAfter,attr"`
	AuthnContext        AuthnContext
}

type AuthnContext struct {
	AuthnContextClassRef string
}

type AttributeStatement struct {
	Attributes []Attribute `xml:"Attribute"`
}

type Attribute struct {
	Name           string         `xml:"Name,attr"`
	NameFormat     string         `xml:"NameFormat,attr"`
	AttributeValue AttributeValue `xml:"AttributeValue"`
}

type AttributeValue struct {
	XS    string `xml:"xs,attr"`
	XSI   string `xml:"xsi,attr"`
	Type  string `xml:"https://www.w3.org/2001/XMLSchema-instance type,attr"`
	Value string `xml:",chardata"`
}

func (r SamlResponse) NameID() string {
	return r.Response.Assertion.Subject.NameID.Value
}

func (r SamlResponse) Email() string {
	for _, attr := range r.Response.Assertion.AttributeStatement.Attributes {
		if attr.Name == "email" {
			return attr.AttributeValue.Value
		}
	}
	return ""
}

func (r SamlResponse) Destination() string {
	return r.Response.Destination
}

func (r SamlResponse) Issuer() string {
	return r.Response.Issuer
}

func (r SamlResponse) StatusCode() string {
	return r.Response.Status.StatusCode.Value
}

func (r SamlResponse) Recipient() string {
	return r.Response.Assertion.Subject.SubjectConfirmation.SubjectConfirmationData.Recipient
}

func (r SamlResponse) SubjectNotOnOrAfter() (time.Time, error) {
	return time.Parse(time.RFC3339, r.Response.Assertion.Subject.SubjectConfirmation.SubjectConfirmationData.NotOnOrAfter)
}

func (r SamlResponse) ConditionNotOnOrAfter() (time.Time, error) {
	return time.Parse(time.RFC3339, r.Response.Assertion.Conditions.NotOnOrAfter)
}

func (r SamlResponse) ConditionNotBefore() (time.Time, error) {
	return time.Parse(time.RFC3339, r.Response.Assertion.Conditions.NotBefore)
}

func (r SamlResponse) Audience() string {
	return r.Response.Assertion.Conditions.AudienceRestriction.Audience
}

func (r SamlResponse) SessionIndex() string {
	return r.Response.Assertion.AuthnStatement.SessionIndex
}

func (r SamlResponse) SessionNotOnOrAfter() (time.Time, error) {
	return time.Parse(time.RFC3339, r.Response.Assertion.AuthnStatement.SessionNotOnOrAfter)
}

func (r SamlResponse) ResponseSignature() Signature {
	return r.Response.Signature
}

func (r SamlResponse) AssertionSignature() Signature {
	return r.Response.Assertion.Signature
}

func (r SamlResponse) Validate(md *model.Metadata) error {
	if r.Destination() != md.ACSURL() {
		return errors.New("destination is invalid")
	}

	if r.Issuer() != md.EntityID {
		return errors.New("issuer is invalid")
	}

	if r.StatusCode() != StatusSeccess {
		return errors.New("status is not success")
	}

	if r.Recipient() != md.ACSURL() {
		return errors.New("recipient is invalid")
	}

	cnb, err := r.ConditionNotBefore()
	if err != nil {
		return fmt.Errorf("condition NotBefore: %v, now: %v", err, time.Now().Format(time.RFC3339))
	}
	cnooa, err := r.ConditionNotOnOrAfter()
	if err != nil {
		return fmt.Errorf("parse error: condition NotOnOrAfter: %v, now: %v", err, time.Now().Format(time.RFC3339))
	}
	snooa, err := r.SubjectNotOnOrAfter()
	if err != nil {
		return fmt.Errorf("parse error: subject NotOnOrAfter: %v, now: %v", err, time.Now().Format(time.RFC3339))
	}
	_, err = r.SessionNotOnOrAfter()
	if err != nil {
		return fmt.Errorf("parse error: session NotOnOrAfter: %v, now: %v", err, time.Now().Format(time.RFC3339))
	}
	now := time.Now()

	if now.Before(cnb) {
		return fmt.Errorf("condition NotBefore: %v, now: %v", cnb, time.Now().Format(time.RFC3339))
	}
	if !now.Before(cnooa) {
		return fmt.Errorf("condition NotOnOrAfter: %v, now: %v", cnb, time.Now().Format(time.RFC3339))
	}
	if !now.Before(snooa) {
		return fmt.Errorf("subject NotOnOrAfter: %v, now: %v", cnb, time.Now().Format(time.RFC3339))
	}
	return nil
}

func getSAMLResponse(c echo.Context) (SamlResponse, error) {
	encRes := c.FormValue("SAMLResponse")
	decRes, _ := base64.StdEncoding.DecodeString(encRes)
	res := SamlResponse{}
	if err := xml.Unmarshal(decRes, &res.Response); err != nil {
		return SamlResponse{}, err
	}
	return res, nil
}
