package core

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// Proof 证明
type Proof struct {
	Challenge          string   `gorm:"type:varchar(256)" json:"challenge,omitempty"`
	Created            string   `gorm:"type:varchar(256)" json:"created,omitempty"`
	ExcludedFields     JSONList `gorm:"type:json" json:"excludedFields,omitempty"`
	ProofPurpose       string   `gorm:"type:varchar(256)" json:"proofPurpose,omitempty"`
	ProofValue         []byte   `gorm:"type:blob" json:"proofValue,omitempty"`
	SignedFields       JSONList `gorm:"type:json" json:"signedFields,omitempty"`
	Type               string   `gorm:"type:varchar(256)" json:"type,omitempty"`
	VerificationMethod string   `gorm:"type:varchar(256)" json:"verificationMethod,omitempty"`
}

// Scan implements the sql.Scanner interface for Proof
func (p *Proof) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot convert %v to Proof", value)
	}
	return json.Unmarshal(bytes, p)
}

// Value implements the driver.Valuer interface for Proof
func (p Proof) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// JSONList is a custom type for handling JSON arrays in GORM
type JSONList []string

// Scan implements the sql.Scanner interface for JSONList
func (j *JSONList) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot convert %v to JSONList", value)
	}
	return json.Unmarshal(bytes, j)
}

// Value implements the driver.Valuer interface for JSONList
func (j JSONList) Value() (driver.Value, error) {
	return json.Marshal(j)
}
