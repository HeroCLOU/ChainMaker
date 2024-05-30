package services

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/tjfoc/gmsm/sm2"
	"log"
	"time"
)

type VP struct {
	ID         string `json:"id"`
	Credential string `json:"credential"`
	Signature  string `json:"signature"`
	IssuedAt   string `json:"issued_at"`
}

func CreateVP(vc string, privateKey *sm2.PrivateKey) (*VP, error) {
	log.Printf("正在为VC创建VP: %s", vc)
	signature, err := privateKey.Sign(rand.Reader, []byte(vc), nil)
	if err != nil {
		log.Printf("签署VC时出错: %v", err)
		return nil, err
	}
	vp := &VP{
		ID:         "vp-" + time.Now().Format("20060102150405"),
		Credential: vc,
		Signature:  hex.EncodeToString(signature),
		IssuedAt:   time.Now().Format(time.RFC3339),
	}
	log.Printf("VP创建成功: %+v", vp)
	return vp, nil
}

func VerifyVP(vp *VP, publicKey *sm2.PublicKey) bool {
	log.Printf("正在验证VP: %+v", vp)
	signature, err := hex.DecodeString(vp.Signature)
	if err != nil {
		log.Printf("解码签名时出错: %v", err)
		return false
	}

	// 验证签名
	result := publicKey.Verify([]byte(vp.Credential), signature)
	log.Printf("验证结果: %v", result)
	return result
}
