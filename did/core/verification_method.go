package core

// VerificationMethod 公钥
type VerificationMethod struct {
	Address      string `gorm:"type:varchar(256)" json:"address,omitempty"`
	Controller   string `gorm:"type:varchar(256)" json:"controller,omitempty"`
	Id           string `gorm:"type:varchar(256);primary_key" json:"id,omitempty"`
	PublicKeyPem string `gorm:"type:varchar(1024)" json:"publicKeyPem,omitempty"`
	Type         string `gorm:"type:varchar(256)" json:"type,omitempty"`
}
