package controller

type Signature struct {
	SignedInfo     SignedInfo
	SignatureValue string
}

type SignedInfo struct {
	Reference Reference
}

type Reference struct {
	DigestValue string
}

func (s Signature) IsZero() bool {
	return s == Signature{}
}
