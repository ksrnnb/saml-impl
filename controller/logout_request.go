package controller

type LogoutRequest struct {
	Destination  string `xml:"Destination,attr"`
	Issuer       string
	NameID       string
	SessionIndex string
}
