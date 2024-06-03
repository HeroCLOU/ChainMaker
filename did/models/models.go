package models

import (
	"did/core"
	"time"

	"gorm.io/gorm"
)

// DID 表示一个去中心化标识符（Decentralized Identifier）模型。
type DID struct {
	// ID 是 DID 的主键和唯一标识符。
	ID string `gorm:"primaryKey;type:varchar(128)" json:"id"` // DID 的唯一标识符
	// Document 存储 DID 文档的 JSON 格式数据。
	Document core.Document `gorm:"type:json" json:"document"` // DID 文档的 JSON 格式数据
	// Txid 存储与 DID 创建或更新相关的交易 ID。
	Txid string `gorm:"type:varchar(64)" json:"txid"` // 与 DID 创建或更新相关的交易 ID
	// DocumentURL 提供可以访问 DID 文档的 URL。
	DocumentURL string `gorm:"type:varchar(255)" json:"document_url"` // 可以访问 DID 文档的 URL
	// Created 记录 DID 创建的时间戳。
	Created time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created"` // DID 创建的时间戳
	// Updated 记录 DID 最后一次更新的时间戳。
	Updated time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP" json:"updated"` // DID 最后一次更新的时间戳
	// Version 记录 DID 的版本号，默认为 0。
	Version int64 `gorm:"default:0" json:"version"` // DID 的版本号
}

// DIDPublicKey 代表与 DID 关联的公钥信息。
func (*DID) TableName() string {
	return "did"
}

// DIDPublicKey 表示与 DID 相关联的公钥模型。
type DIDPublicKey struct {
	// ID 是公钥的唯一标识符。
	ID string `gorm:"primaryKey;type:varchar(64)" json:"id"` // 公钥的唯一标识符
	// Type 指定公钥的类型。
	Type string `gorm:"type:varchar(128)" json:"type"` // 公钥的类型
	// Controller 指定控制该公钥的 DID。
	Controller string `gorm:"type:varchar(64)" json:"controller"` // 控制该公钥的 DID
	// DID 关联的去中心化标识符。
	DID string `gorm:"type:varchar(128)" json:"did"` // 关联的去中心化标识符
	// Address 存储公钥对应的区块链地址。
	Address string `gorm:"type:varchar(128)" json:"address"` // 公钥对应的区块链地址
	// PublicKey 存储公钥的字符串表示形式。
	PublicKey string `gorm:"type:varchar(256)" json:"public_key"` // 公钥的字符串表示形式
	// Created 记录公钥创建的时间戳。
	Created time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created"` // 公钥创建的时间戳
	// DeletedAt 记录公钥被删除的时间戳，用于软删除。
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"` // 公钥被删除的时间戳
}

// TableName 方法返回表名。
func (DIDPublicKey) TableName() string {
	return "did_public_key"
}

// VC 代表可验证凭证（Verifiable Credential）模型。
type VC struct {
	// ID 是凭证的唯一标识符。
	ID string `gorm:"primaryKey;type:varchar(32)" json:"id"` // 凭证的唯一标识符
	// Holder 存储凭证持有者的 DID。
	Holder string `gorm:"type:varchar(64)" json:"holder"` // 凭证持有者的 DID
	// Issuer 存储凭证发行者的 DID。
	Issuer string `gorm:"type:varchar(64)" json:"issuer"` // 凭证发行者的 DID
	// TemplateID 存储凭证模板的标识符。
	TemplateID string `gorm:"type:varchar(32)" json:"template_id"` // 凭证模板的标识符
	// Credential 存储凭证的详细信息。
	Credential core.JSONList `gorm:"type:json" json:"credential"` // 凭证的详细信息
	// Issuance 存储凭证的发行时间。
	Issuance time.Time `json:"issuance"` // 凭证的发行时间
	// Expiration 存储凭证的过期时间。
	Expiration time.Time `json:"expiration"` // 凭证的过期时间
	// Status 表示凭证的状态，默认为 0。
	Status int `gorm:"default:0" json:"status"` // 凭证的状态
	// Created 记录凭证创建的时间戳。
	Created time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created"` // 凭证创建的时间戳
	// Proof 存储凭证的证明信息。
	Proof core.Proof `gorm:"type:json" json:"proof"` // 凭证的证明信息
}

// TableName 方法返回表名。
func (*VC) TableName() string {
	return "vc"
}

// VP 代表可验证展示（Verifiable Presentation）模型。
type VP struct {
	// ID 是展示的唯一标识符。
	ID string `gorm:"primaryKey;type:varchar(32)" json:"id"` // 展示的唯一标识符
	// Context 存储展示的上下文信息。
	Context core.JSONList `gorm:"type:json" json:"@context"` // 展示的上下文信息
	// Type 存储展示的类型信息。
	Type core.JSONList `gorm:"type:json" json:"type"` // 展示的类型信息
	// Holder 存储展示持有者的 DID。
	Holder string `gorm:"type:varchar(64)" json:"holder"` // 展示持有者的 DID
	// VerifiableCredential 存储展示中包含的可验证凭证列表。
	VerifiableCredential []VC `gorm:"type:json" json:"verifiableCredential"` // 展示中包含的可验证凭证列表
	// Proof 存储展示的证明信息。
	Proof core.Proof `gorm:"type:json" json:"proof"` // 展示的证明信息
}

// TableName 方法返回表名。
func (*VP) TableName() string {
	return "vp"
}
