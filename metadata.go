package main

type Metadata struct {
	EntityDescriptor EntityDescriptor
}

type EntityDescriptor struct {
	EntityID         string           `xml:"entityID,attr"`
	IDPSSODescriptor IDPSSODescriptor `xml:"IDPSSODescriptor"`
}

type IDPSSODescriptor struct {
	KeyDescriptor        KeyDescriptor         `xml:"KeyDescriptor"`
	SingleSignOnServices []SingleSignOnService `xml:"SingleSignOnService"`
}

type KeyDescriptor struct {
	KeyInfo KeyInfo `xml:"KeyInfo"`
}

type KeyInfo struct {
	X509Data X509Data `xml:"X509Data"`
}

type X509Data struct {
	X509Certificate string `xml:"X509Certificate"`
}

type SingleSignOnService struct {
	Binding  string `xml:"Binding,attr"`
	Location string `xml:"Location,attr"`
}
