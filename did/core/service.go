package core

// Service DID服务
type Service struct {
	Id              string `gorm:"type:varchar(256)" json:"id,omitempty"`
	ServiceEndpoint string `gorm:"type:varchar(256)" json:"serviceEndpoint,omitempty"`
	Type            string `gorm:"type:varchar(256)" json:"type,omitempty"`
}
