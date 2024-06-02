package core

// VcTemplate VC模板
type VcTemplate struct {
	Id      string `gorm:"type:varchar(256);primary_key" json:"id,omitempty"`
	Name    string `gorm:"type:varchar(256)" json:"name,omitempty"`
	VcType  string `gorm:"type:varchar(256)" json:"vcType,omitempty"`
	Version string `gorm:"type:varchar(256)" json:"version,omitempty"`
}
