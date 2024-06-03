package method

import (
	"crypto/rand"
	"did/core"
	"did/models"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
)

var (
	DB            *gorm.DB
	sm2PrivateKey *sm2.PrivateKey
	sm2PublicKey  *sm2.PublicKey
)

func initMySQL() { //连接数据库
	var err error
	des := "root:xcr20020304@(127.0.0.1:3306)/did?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open("mysql", des)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	err = DB.DB().Ping()
	if err != nil {
		log.Fatalf("数据库Ping失败: %v", err)
	}
}

func randomInt(max int) int { //生成随机数
	nBig, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(nBig.Int64())
}

func CreatePublicKey(did string, publickey string) { //创建publickey
	randomPart := fmt.Sprintf("%010d", randomInt(1000000000))
	id := "publickey" + randomPart
	public := models.DIDPublicKey{
		ID:         id,
		Type:       "ificationKey2019",
		Controller: "did:example:123456",
		DID:        did,
		Address:    "0x742d35Cc6634C0532925a3",
		PublicKey:  publickey,
		Created:    time.Now(),
	}
	if err := DB.Table("did_public_key").Create(&public).Error; err != nil {
		log.Fatalf("Failed to create PublicKey: %v", err)
	}
	fmt.Println("PublicKey registered successfully")
	fmt.Println("-*******************************************-")
}

func CreateProof(did string, purpose string, priv *sm2.PrivateKey) core.Proof { //创建proof
	now := time.Now().UTC().Format(time.RFC3339)
	msg := now + purpose
	bymsg := []byte(msg)

	sign, err := priv.Sign(rand.Reader, bymsg, nil)
	if err != nil {
		log.Fatalf("签名时出错: %v", err)
	}

	proof := core.Proof{
		Challenge:          "aUniqueChallengeString",
		Created:            now,
		ExcludedFields:     core.JSONList{"field1", "field2"},
		ProofPurpose:       purpose,
		ProofValue:         sign,
		SignedFields:       core.JSONList{"created", "proofPurpose"},
		Type:               "JsonWebSignature2020",
		VerificationMethod: "did:example:123456#key-1",
	}
	return proof
}

func CreateDIDDocument(did string, priv *sm2.PrivateKey) core.Document { //创建DIDDocument
	var context core.JSONList
	context = append(context, "https://www.w3.org/ns/did/v1")

	var authentication core.JSONList
	authentication = append(authentication, "/authentication/key-1")

	var controller core.JSONList
	controller = append(controller, "did:example:123456")

	purpose := "生成DIDDocument"
	proof := CreateProof(did, purpose, priv)

	var verificationMethod core.VerificationMethod
	verificationMethod.Id = "did:example:123456#key-1"
	verificationMethod.Type = "Ed25519VerificationKey2018"
	verificationMethod.Controller = "did:example:123456"

	document := core.Document{
		Context:            context,
		Authentication:     authentication,
		Controller:         controller,
		Created:            time.Now().Format(time.RFC3339),
		Id:                 did,
		Proof:              []core.Proof{proof},
		Service:            []core.Service{},
		UnionId:            core.JSONMap{"example": "value"},
		Updated:            time.Now().Format(time.RFC3339),
		VerificationMethod: []core.VerificationMethod{verificationMethod},
	}
	return document
}

func CreateDID(priv *sm2.PrivateKey) models.DID { //创建DID
	timestamp := time.Now().Unix()
	id := fmt.Sprintf("did:example:%d", timestamp)
	diddocument := CreateDIDDocument(id, priv)
	did := models.DID{
		ID:          id,
		Document:    diddocument,
		Txid:        "abcdef1234567890",
		DocumentURL: "https://example.com/dids/did:example:123456/document",
		Created:     time.Now(),
		Updated:     time.Now(),
		Version:     1,
	}
	return did
}

func CreateVC(did string, priv *sm2.PrivateKey, credential core.JSONList) models.VC { //创建VC
	randomPart := fmt.Sprintf("%010d", randomInt(1000000000))
	id := "vc" + randomPart

	now := time.Now()
	proof := CreateProof(did, "生成vc", priv)
	vc := models.VC{
		ID:         id,
		Holder:     did,
		Issuer:     did,
		TemplateID: "tpl-001",
		Issuance:   now,
		Expiration: now.AddDate(1, 0, 0),
		Status:     0,
		Created:    now,
		Credential: credential,
		Proof:      proof,
	}
	return vc
}

func generateKeys() { //生成密钥对
	privateKey, err := sm2.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatalf("生成私钥时出错: %v", err)
	}

	privBytes, err := x509.WritePrivateKeyToPem(privateKey, nil)
	if err != nil {
		log.Fatalf("编码私钥为PEM时出错: %v", err)
	}

	pubBytes, err := x509.WritePublicKeyToPem(&privateKey.PublicKey)
	if err != nil {
		log.Fatalf("编码公钥为PEM时出错: %v", err)
	}

	err = os.WriteFile("private_key.pem", privBytes, 0600)
	if err != nil {
		log.Fatalf("写入私钥文件时出错: %v", err)
	}
	err = os.WriteFile("public_key.pem", pubBytes, 0644)
	if err != nil {
		log.Fatalf("写入公钥文件时出错: %v", err)
	}

	fmt.Println("私钥和公钥已生成并存储在 PEM 文件中")
}

func loadKeys() { //加载密钥
	privateKeyPem, err := os.ReadFile("private_key.pem")
	if err != nil {
		log.Fatalf("读取私钥文件时出错: %v", err)
	}

	block, _ := pem.Decode(privateKeyPem)
	if block == nil || block.Type != "PRIVATE KEY" {
		log.Fatalf("解码私钥时出错")
	}

	sm2PrivateKey, err = x509.ReadPrivateKeyFromPem(privateKeyPem, nil)
	if err != nil {
		log.Fatalf("解析私钥时出错: %v", err)
	}

	publicKeyPem, err := os.ReadFile("public_key.pem")
	if err != nil {
		log.Fatalf("读取公钥文件时出错: %v", err)
	}

	block, _ = pem.Decode(publicKeyPem)
	if block == nil || block.Type != "PUBLIC KEY" {
		log.Fatalf("解码公钥时出错")
	}

	sm2PublicKey, err = x509.ReadPublicKeyFromPem(publicKeyPem)
	if err != nil {
		log.Fatalf("解析公钥时出错: %v", err)
	}

	fmt.Println("私钥和公钥已成功加载")
}

func Input(name string, cardid string, school string, schoolrecord string) core.JSONList { //修改输入
	var prename = "姓名:"
	var precardid = "身份证号:"
	var preschool = "学校:"
	var preschoolrecord = "学历:"

	name = prename + name
	cardid = precardid + cardid
	school = preschool + school
	schoolrecord = preschoolrecord + schoolrecord
	var credential core.JSONList
	credential = append(credential, name)
	credential = append(credential, cardid)
	credential = append(credential, school)
	credential = append(credential, schoolrecord)
	return credential
}

func Register(name string, cardid string, school string, schoolrecord string) (string, *sm2.PrivateKey) { //注册
	initMySQL()
	credential := Input(name, cardid, school, schoolrecord)
	generateKeys()
	loadKeys()

	fmt.Println("-*******************************************-")

	did := CreateDID(sm2PrivateKey)
	if err := DB.Table("did").Create(&did).Error; err != nil {
		log.Fatalf("Failed to create DID: %v", err)
	}

	id := did.ID
	fmt.Println("DID registered successfully")
	fmt.Println("your did:", id)
	fmt.Println("-*******************************************-")

	public, err := os.ReadFile("public_key.pem")
	if err != nil {
		log.Fatalf("读取公钥文件时出错: %v", err)
	}

	CreatePublicKey(id, string(public))

	vc := CreateVC(id, sm2PrivateKey, credential)
	if err := DB.Table("vc").Create(&vc).Error; err != nil {
		log.Fatalf("Failed to create VC: %v", err)
	}
	fmt.Println("VC registered successfully")
	fmt.Println("-*******************************************-")
	return id, sm2PrivateKey
}
