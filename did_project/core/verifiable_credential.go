package core

// VerifiableCredential 简称“凭证”，遵循W3C Verifiable Credential规范的电子凭证
type VerifiableCredential struct {
	Context           JSONList   `gorm:"type:json" json:"@context,omitempty"`
	CredentialSubject []byte     `gorm:"type:json" json:"credentialSubject,omitempty"`
	ExpirationDate    string     `gorm:"type:varchar(256)" json:"expirationDate,omitempty"`
	Holder            string     `gorm:"type:varchar(256)" json:"holder,omitempty"`
	Id                string     `gorm:"type:varchar(256);primary_key" json:"id,omitempty"`
	IssuanceDate      string     `gorm:"type:varchar(256)" json:"issuanceDate,omitempty"`
	Issuer            string     `gorm:"type:varchar(256)" json:"issuer,omitempty"`
	Proof             []Proof    `gorm:"type:json" json:"proof,omitempty"`
	Template          VcTemplate `gorm:"type:json" json:"template,omitempty"`
	Type              JSONList   `gorm:"type:json" json:"type,omitempty"`
}
