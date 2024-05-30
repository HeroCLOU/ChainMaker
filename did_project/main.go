package main

import (
	"crypto/rand"
	"did_project/models"
	"did_project/routes" // 替换为您的项目路径
	"encoding/hex"
	"fmt"
	"github.com/tjfoc/gmsm/sm2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"math/big"
	"os"
)

var sm2PrivateKey *sm2.PrivateKey
var sm2PublicKey *sm2.PublicKey

// 生成并存储密钥
func generateKeys() {
	privateKey, err := sm2.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatalf("生成私钥时出错: %v", err)
	}

	publicKey := &privateKey.PublicKey

	// 将私钥和公钥分别编码为十六进制字符串
	privateKeyHex := hex.EncodeToString(privateKey.D.Bytes())
	publicKeyHex := hex.EncodeToString(append(publicKey.X.Bytes(), publicKey.Y.Bytes()...))

	// 将私钥和公钥存储到文件
	err = ioutil.WriteFile("private_key.hex", []byte(privateKeyHex), 0600)
	if err != nil {
		log.Fatalf("写入私钥文件时出错: %v", err)
	}
	err = ioutil.WriteFile("public_key.hex", []byte(publicKeyHex), 0644)
	if err != nil {
		log.Fatalf("写入公钥文件时出错: %v", err)
	}

	fmt.Println("私钥和公钥已生成并存储在文件中")
	fmt.Printf("私钥: %s\n", privateKeyHex)
	fmt.Printf("公钥: %s\n", publicKeyHex)
}

// 加载私钥和公钥
func loadKeys() {
	// 从文件加载私钥
	privateKeyHex, err := ioutil.ReadFile("private_key.hex")
	if err != nil {
		log.Fatalf("读取私钥文件时出错: %v", err)
	}

	privateKeyBytes, err := hex.DecodeString(string(privateKeyHex))
	if err != nil {
		log.Fatalf("解码私钥时出错: %v", err)
	}

	sm2PrivateKey = new(sm2.PrivateKey)
	sm2PrivateKey.D = new(big.Int).SetBytes(privateKeyBytes)
	sm2PrivateKey.PublicKey.Curve = sm2.P256Sm2()
	sm2PrivateKey.PublicKey.X, sm2PrivateKey.PublicKey.Y = sm2PrivateKey.PublicKey.Curve.ScalarBaseMult(privateKeyBytes)

	// 从文件加载公钥
	publicKeyHex, err := ioutil.ReadFile("public_key.hex")
	if err != nil {
		log.Fatalf("读取公钥文件时出错: %v", err)
	}

	publicKeyBytes, err := hex.DecodeString(string(publicKeyHex))
	if err != nil {
		log.Fatalf("解码公钥时出错: %v", err)
	}

	sm2PublicKey = new(sm2.PublicKey)
	sm2PublicKey.Curve = sm2.P256Sm2()
	sm2PublicKey.X = new(big.Int).SetBytes(publicKeyBytes[:len(publicKeyBytes)/2])
	sm2PublicKey.Y = new(big.Int).SetBytes(publicKeyBytes[len(publicKeyBytes)/2:])

	fmt.Println("私钥和公钥已成功加载")
}

func init() {
	// 检查密钥文件是否存在
	if _, err := os.Stat("private_key.hex"); os.IsNotExist(err) {
		// 如果密钥文件不存在，则生成密钥
		generateKeys()
	}

	// 加载密钥
	loadKeys()
}

func main() {
	dsn := "root:g28781212828@tcp(127.0.0.1:3306)/did?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("连接数据库失败: ", err)
	} else {
		log.Println("成功连接到数据库")
	}

	err = db.AutoMigrate(&models.DID{}, &models.DIDPublicKey{}, &models.VC{})
	if err != nil {
		log.Fatal("数据库迁移失败: ", err)
	} else {
		log.Println("数据库迁移成功")
	}

	r := routes.SetupRouter(db)
	log.Println("服务器正在8080端口启动")
	r.Run(":8080")
}
