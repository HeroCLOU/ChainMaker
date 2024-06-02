package core

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// Document DID 文档
type Document struct {
	Context            JSONList             `gorm:"type:json" json:"@context,omitempty"`
	Authentication     JSONList             `gorm:"type:json" json:"authentication,omitempty"`
	Controller         JSONList             `gorm:"type:json" json:"controller,omitempty"`
	Created            string               `gorm:"type:varchar(256)" json:"created,omitempty"`
	Id                 string               `gorm:"type:varchar(256);primary_key" json:"id,omitempty"`
	Proof              []Proof              `gorm:"type:json" json:"proof,omitempty"`
	Service            []Service            `gorm:"type:json" json:"service,omitempty"`
	UnionId            JSONMap              `gorm:"type:json" json:"unionId,omitempty"`
	Updated            string               `gorm:"type:varchar(256)" json:"updated,omitempty"`
	VerificationMethod []VerificationMethod `gorm:"type:json" json:"verificationMethod,omitempty"`
}

// Scan implements the sql.Scanner interface for Document
func (d *Document) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot convert %v to Document", value)
	}
	return json.Unmarshal(bytes, d)
}

// Value implements the driver.Valuer interface for Document
func (d Document) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// JSONMap is a custom type for handling JSON maps in GORM
type JSONMap map[string]string

// Scan implements the sql.Scanner interface for JSONMap
func (j *JSONMap) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot convert %v to JSONMap", value)
	}
	return json.Unmarshal(bytes, j)
}

// Value implements the driver.Valuer interface for JSONMap
func (j JSONMap) Value() (driver.Value, error) {
	return json.Marshal(j)
}
