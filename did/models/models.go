package models

import (
	"did/core"
	"time"

	"gorm.io/gorm"
)

type DID struct {
	ID          string        `gorm:"primaryKey;type:varchar(128)" json:"id"`
	Document    core.Document `gorm:"type:json" json:"document"`
	Txid        string        `gorm:"type:varchar(64)" json:"txid"`
	DocumentURL string        `gorm:"type:varchar(255)" json:"document_url"`
	Created     time.Time     `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created"`
	Updated     time.Time     `gorm:"type:timestamp;default:CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP" json:"updated"`
	Version     int64         `gorm:"default:0" json:"version"`
}

func (*DID) TableName() string {
	return "did"
}

type DIDPublicKey struct {
	ID         string         `gorm:"primaryKey;type:varchar(64)" json:"id"`
	Type       string         `gorm:"type:varchar(128)" json:"type"`
	Controller string         `gorm:"type:varchar(64)" json:"controller"`
	DID        string         `gorm:"type:varchar(128)" json:"did"`
	Address    string         `gorm:"type:varchar(128)" json:"address"`
	PublicKey  string         `gorm:"type:varchar(256)" json:"public_key"`
	Created    time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (DIDPublicKey) TableName() string {
	return "did_public_key"
}

type VC struct {
	ID         string        `gorm:"primaryKey;type:varchar(32)" json:"id"`
	Holder     string        `gorm:"type:varchar(64)" json:"holder"`
	Issuer     string        `gorm:"type:varchar(64)" json:"issuer"`
	TemplateID string        `gorm:"type:varchar(32)" json:"template_id"`
	Credential core.JSONList `gorm:"type:json" json:"credential"`
	Issuance   time.Time     `json:"issuance"`
	Expiration time.Time     `json:"expiration"`
	Status     int           `gorm:"default:0" json:"status"`
	Created    time.Time     `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created"`
	Proof      core.Proof    `gorm:"type:json" json:"proof"`
}

func (*VC) TableName() string {
	return "vc"
}

type VP struct {
	ID                   string        `gorm:"primaryKey;type:varchar(32)" json:"id"`
	Context              core.JSONList `gorm:"type:json" json:"@context"`
	Type                 core.JSONList `gorm:"type:json" json:"type"`
	Holder               string        `gorm:"type:varchar(64)" json:"holder"`
	VerifiableCredential []VC          `gorm:"type:json" json:"verifiableCredential"`
	Proof                core.Proof    `gorm:"type:json" json:"proof"`
}

func (*VP) TableName() string {
	return "vp"
}
