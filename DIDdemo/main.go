package main

import (
	"DIDdemo/core"

	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/tjfoc/gmsm/sm2"
)

var DB *gorm.DB // 全局变量

func Initdb() {
	// 设置数据库连接信息
	datainfo := "root:123456@tcp(127.0.0.1:3306)/did?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open("mysql", datainfo)

	fmt.Println("cuowu:%v\n", err)
	// 连接到数据库
	if err != nil {
		// 处理错误
		panic("连接失败！")
	}

	fmt.Println("连接数据库成功")
	DB = db // 将数据库连接保存到全局变量
}

// DID 结构体，与您的数据库表结构相对应
type DID struct {
	ID          string    `gorm:"primaryKey;type:varchar(128)" json:"id"`
	Document    []byte    `gorm:"type:json" json:"document"`
	Txid        string    `gorm:"type:varchar(64)" json:"txid"`
	DocumentURL string    `gorm:"type:varchar(255)" json:"document_url"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
	Version     int64     `json:"version"`
	Black       bool      `json:"black"`
}

// GenerateDID 使用当前的Unix时间戳生成一个简单的DID
func GenerateDID() string {
	timestamp := time.Now().Unix()
	did := fmt.Sprintf("did:zhw:%d", timestamp)
	return did
}

func RegisterDocument(publickey *sm2.PublicKey) (string, []byte, error) {
	//s := publickey
	did := GenerateDID()
	fmt.Println("Generated DID:", did)
	currentTimeStr := time.Now().Format(time.RFC3339)
	// 创建Proof实例
	proof := &core.Proof{
		Created:            currentTimeStr,
		ProofPurpose:       "study",
		ProofValue:         []byte("proof-value-binary"),
		SignedFields:       []string{"field1", "field2"},
		Type:               "SM2Signature2024",
		VerificationMethod: did + "#key1",
	}

	// 创建Service实例
	// service := &core.Service{
	// 	Id:              "did:example:123#service-endpoint",
	// 	ServiceEndpoint: "https://example.com/service-endpoint",
	// 	Type:            "ExampleService",
	// }

	pubBytes := sm2.Compress(publickey)
	pubBase64 := base64.StdEncoding.EncodeToString(pubBytes)

	// 创建VerificationMethod实例
	verificationMethod := &core.VerificationMethod{
		Address:      "did:example:123",
		Controller:   did,
		Id:           did + "#key1",
		PublicKeyPem: "-----BEGIN PUBLIC KEY-----\n" + pubBase64 + "\n-----END PUBLIC KEY-----",
		Type:         "SM2Signature2024",
	}

	// 创建Document实例
	document := &core.Document{
		Context:            []string{"https://www.w3.org/ns/did/v1"},
		Authentication:     []string{did + "#Key1"},
		Controller:         []string{"controller:bistu", did},
		Created:            currentTimeStr,
		Id:                 did,
		Proof:              []*core.Proof{proof},
		Updated:            currentTimeStr,
		VerificationMethod: []*core.VerificationMethod{verificationMethod},
	}

	fmt.Println(document)

	jsondata, err := json.MarshalIndent(document, "", " ")
	if err != nil {
		fmt.Printf("转化成json格式时出现错误")
	}
	return did, jsondata, err
}

// did 注册
func Register() {
	/*  SM2编码实现--生成密钥对 */
	privatekey, err := sm2.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatal(err)
	}
	publickey := &privatekey.PublicKey
	//fmt.Printf("\n---公钥---\n%v\n-----\n", publickey)
	//fmt.Printf("\n---私钥---\n%v\n-----\n", privatekey)

	/* 生成DID */
	//生成一个DID Document
	// 返回的did
	userdid, jsondata, err := RegisterDocument(publickey)
	if err != nil {
		fmt.Println("生成Document失败！")
	}

	str := "aacHJvb2YtdmFsdWUtYmluYXJ5"
	strby := []byte(str)
	fmt.Println(strby)

	// 构造DID
	did := DID{
		ID:          "1005",
		Document:    jsondata,
		Txid:        "",
		DocumentURL: "",
		Created:     time.Now(),
		Updated:     time.Now(),
		Version:     1,
		Black:       false,
	}
	//fmt.Println(did)

	// 使用智能合约发起上链请求

	// 将数据加入到数据库中
	DB.Table("did").Create(&did)

	/* 下面是将DID和公钥插入数据库 */
	// DIDPublicKey 表示did_public_key表的结构体
	type DIDPublicKey struct {
		ID             string     `json:"id"`                    // 唯一主键
		Type           string     `json:"type"`                  // 类型
		Controller     string     `json:"controller"`            // 控制方
		DID            string     `json:"did" gorm:"column:did"` // DID，这里添加gorm对数据进行强制定义
		Address        string     `json:"address"`               // 地址
		PublicKey      string     `json:"public_key"`            // PEM
		Created        time.Time  `json:"created"`               // 创建时间
		Authentication bool       `json:"authentication"`        // 是否是认证公钥
		DeletedAt      *time.Time `json:"deleted_at"`            // 删除时间，如果未删除则为nil
	}

	pubBytes := sm2.Compress(publickey)
	pubBase64 := base64.StdEncoding.EncodeToString(pubBytes)

	// 将did和公钥进行关联
	// 示例数据
	example := DIDPublicKey{
		ID:             "100087",
		Type:           "SM2Signature2024",
		Controller:     userdid,
		DID:            userdid,
		Address:        "bistu.edu.com",
		PublicKey:      "-----BEGIN PUBLIC KEY-----\n" + pubBase64 + "\n-----END PUBLIC KEY-----",
		Created:        time.Now(),
		Authentication: true,
		DeletedAt:      nil, // 如果未删除
	}

	// 将数据加入到数据库中
	DB.Table("did_public_key").Create(&example)
	fmt.Println(example)

	/* 向区块链发起上链请求 */

	/* 向数据库中插入数据 */

	// 插入记录
	// 创建VerifiableCredential实例
	vc := &core.VerifiableCredential{
		Context:           []string{"https://www.w3.org/2018/credentials/v1"},
		CredentialSubject: []byte(`{"id":"did:zhw:2021","name":"ZhangHouwen"}`),
		ExpirationDate:    "2024-05-29T12:00:00Z",
		Holder:            "did:example:studentID",
		Id:                "urn:uuid:3c0a79e7-55b2-4f4e-9f81-30a4e2ea1c2d",
		IssuanceDate:      "2021-05-29T12:00:00Z",
		Issuer:            "did:example:issuerID",
		//Proof:             []*core.Proof{proof},
		Template: &core.VcTemplate{ /* 需要定义VcTemplate结构 */ },
		Type:     []string{"VerifiableCredential"},
	}
	fmt.Println(vc)

	// 创建VerifiablePresentation实例
	vp := &core.VerifiablePresentation{
		Context:           []string{"https://www.w3.org/2018/credentials/v1"},
		ExpirationDate:    "2024-05-29T12:00:00Z",
		Extend:            []byte(`{"additionalData":"some additional data"}`),
		Id:                "urn:uuid:9b1da31a-7e69-46f9-9f4c-7a3f3a4d2b31",
		PresentationUsage: "求职",
		//Proof:                []*core.Proof{proof},
		Timestamp:            "2021-05-29T12:00:00Z",
		Type:                 "VerifiablePresentation",
		VerifiableCredential: []*core.VerifiableCredential{vc},
		Verifier:             "did:example:verifierID",
	}
	fmt.Println(vp)

}

// 用户生成公私钥对
func Test() {
	/*  SM2编码实现 */
	priv, err := sm2.GenerateKey(rand.Reader) // 生成密钥对
	if err != nil {
		log.Fatal(err)
	}
	msg := []byte("Tongji Fintech Research Institute")
	fmt.Println("原文：%x\n", msg)

	pub := &priv.PublicKey

	fmt.Printf("\n---私钥---\n%v\n-----\n", priv.D)
	fmt.Printf("\n---公钥---\n%v\n-----\n", &priv.PublicKey)

	/* 加密 */
	// 加密得到密文
	ciphertxt, err := pub.EncryptAsn1(msg, rand.Reader) //sm2加密
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("加密结果:%x\n", ciphertxt)

	// 解密得到明文
	plaintxt, err := priv.DecryptAsn1(ciphertxt) //sm2解密
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("解密结果：%x\n", plaintxt)

	if !bytes.Equal(msg, plaintxt) {
		log.Fatal("原文不匹配")
	}

	/* 数字签名 */
	// 验签
	sign, err := priv.Sign(rand.Reader, msg, nil) //sm2签名
	if err != nil {
		log.Fatal(err)
	}
	isok := pub.Verify(msg, sign) //sm2验签
	fmt.Println("Verified: %v\n", isok)
}

func main() {
	Initdb()
	Register()

	//Test()
}
