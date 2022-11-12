package controller

import "time"

type SamlResponse struct {
	Response Response
}

type Response struct {
	Destination string `xml:"Destination,attr"`
	Issuer      string
	Status      Status
	Assertion   Assertion
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
	NotBefore           string
	NotOnOrAfter        string
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
	Name           string `xml:"Name,attr"`
	NameFormat     string `xml:"NameFormat,attr"`
	AttributeValue string
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

func (r SamlResponse) SessionNotOnOrAfter() (time.Time, error) {
	return time.Parse(time.RFC3339, r.Response.Assertion.AuthnStatement.SessionNotOnOrAfter)
}
