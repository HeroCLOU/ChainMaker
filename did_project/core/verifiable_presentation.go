package core

// VerifiablePresentation 可验证凭证/可验证声明 遵循W3C Verifiable Presentation规范的凭证打包集。一个Presentation可以包括一个或多个Credential。
type VerifiablePresentation struct {
	Context              JSONList               `gorm:"type:json" json:"@context,omitempty"`
	ExpirationDate       string                 `gorm:"type:varchar(256)" json:"expirationDate,omitempty"`
	Extend               []byte                 `gorm:"type:json" json:"extend,omitempty"`
	Id                   string                 `gorm:"type:varchar(256);primary_key" json:"id,omitempty"`
	PresentationUsage    string                 `gorm:"type:varchar(256)" json:"presentationUsage,omitempty"`
	Proof                []Proof                `gorm:"type:json" json:"proof,omitempty"`
	Timestamp            string                 `gorm:"type:varchar(256)" json:"timestamp,omitempty"`
	Type                 string                 `gorm:"type:varchar(256)" json:"type,omitempty"`
	VerifiableCredential []VerifiableCredential `gorm:"type:json" json:"verifiableCredential,omitempty"`
	Verifier             string                 `gorm:"type:varchar(256)" json:"verifier,omitempty"`
}
