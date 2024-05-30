package controllers

import (
	"crypto/rand"
	"did_project/models"
	"did_project/pkg"
	"did_project/services"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/tjfoc/gmsm/sm2"
	"gorm.io/gorm"
	"log"
	"net/http"
)

var sm2PrivateKey, _ = sm2.GenerateKey(rand.Reader)
var sm2PublicKey = &sm2PrivateKey.PublicKey

func CreateDID(c *gin.Context) {
	var did models.DID
	if err := c.ShouldBindJSON(&did); err != nil {
		log.Printf("Error binding JSON: %v", err)
		pkg.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	db := c.MustGet("db").(*gorm.DB)
	if err := db.Create(&did).Error; err != nil {
		log.Printf("Error creating DID: %v", err)
		pkg.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	pkg.RespondWithJSON(c, http.StatusOK, did)
}

func CreateVC(c *gin.Context) {
	var vc models.VC
	if err := c.ShouldBindJSON(&vc); err != nil {
		log.Printf("绑定JSON数据时出错: %v", err)
		pkg.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	db := c.MustGet("db").(*gorm.DB)
	if err := db.Create(&vc).Error; err != nil {
		log.Printf("创建VC时出错: %v", err)
		pkg.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("VC创建成功: %+v", vc)
	pkg.RespondWithJSON(c, http.StatusOK, vc)
}

/*
func GenerateVP(c *gin.Context) {
	var req struct {
		VCID string `json:"vc_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		pkg.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	var vc models.VC
	if err := db.First(&vc, "id = ?", req.VCID).Error; err != nil {
		log.Printf("VC not found: %v", err)
		pkg.RespondWithError(c, http.StatusNotFound, "VC not found")
		return
	}

	vp, err := services.CreateVP(vc.Credential, sm2PrivateKey)
	if err != nil {
		log.Printf("Error creating VP: %v", err)
		pkg.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	pkg.RespondWithJSON(c, http.StatusOK, vp)
}*/

func GenerateVP(c *gin.Context) {
	var req struct {
		VCID string `json:"vc_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("绑定JSON数据时出错: %v", err)
		pkg.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	if req.VCID == "" {
		log.Println("VC ID 为空")
		pkg.RespondWithError(c, http.StatusBadRequest, "VC ID 不能为空")
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	var vc models.VC
	if err := db.First(&vc, "id = ?", req.VCID).Error; err != nil {
		log.Printf("未找到VC: %v", err)
		pkg.RespondWithError(c, http.StatusNotFound, "未找到VC")
		return
	}

	vp, err := services.CreateVP(vc.Credential, sm2PrivateKey)
	if err != nil {
		log.Printf("创建VP时出错: %v", err)
		pkg.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("VP创建成功: %+v", vp)
	pkg.RespondWithJSON(c, http.StatusOK, vp)
}

func VerifyVP(c *gin.Context) {
	var vp services.VP
	if err := c.ShouldBindJSON(&vp); err != nil {
		log.Printf("绑定JSON数据时出错: %v", err)
		pkg.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("收到的VP数据: %+v", vp)

	var credential map[string]interface{}
	if err := json.Unmarshal([]byte(vp.Credential), &credential); err != nil {
		log.Printf("解码凭证时出错: %v", err)
		pkg.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	valid := services.VerifyVP(&vp, sm2PublicKey)
	if !valid {
		log.Println("VP验证失败")
		pkg.RespondWithError(c, http.StatusUnauthorized, "无效的VP")
		return
	}

	log.Println("VP验证成功")
	pkg.RespondWithJSON(c, http.StatusOK, gin.H{"valid": valid})
}
