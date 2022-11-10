package controller

type SamlResponse struct {
	Response Response
}

type Response struct {
	Issuer    string
	Status    Status
	Assertion Assertion
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
