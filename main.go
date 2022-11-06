package main

import (
	"encoding/xml"
	"fmt"
	"os"
)

const httpPostBinding = "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST"

func main() {
	f, err := os.ReadFile("./metadata.xml")
	if err != nil {
		panic(err)
	}

	md := Metadata{}
	xml.Unmarshal(f, &md.EntityDescriptor)

	fmt.Printf("%+v\n", md)
}
