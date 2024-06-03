package method

import (
	"did/core"
	"did/models"
	"fmt"

	"github.com/tjfoc/gmsm/sm2"
)

func CreateVP(did string, priv *sm2.PrivateKey) models.VP { //创建VP
	var vc models.VC
	DB.Table("vc").Where("holder = ?", did).First(&vc)
	randomPart := fmt.Sprintf("%010d", randomInt(1000000000))
	id := "vc" + randomPart
	var context core.JSONList
	context = append(context, "https://www.w3.org/2018/credentials/v1")
	var vptype core.JSONList
	vptype = append(vptype, "vp")
	purpose := "签发VP"
	proof := CreateProof(vc.Holder, purpose, priv)
	vp := models.VP{
		ID:                   id,
		Holder:               vc.Holder,
		VerifiableCredential: []models.VC{vc},
		Context:              context,
		Type:                 vptype,
		Proof:                proof,
	}
	return vp
}

func verifyVP(vp models.VP, pub *sm2.PublicKey) bool { //验证VP
	csign := vp.Proof.ProofValue
	sign := vp.Proof.Created + vp.Proof.ProofPurpose
	result := pub.Verify([]byte(sign), csign)

	csign2 := vp.VerifiableCredential[0].Proof.ProofValue
	sign2 := vp.VerifiableCredential[0].Proof.Created + vp.VerifiableCredential[0].Proof.ProofPurpose
	result2 := pub.Verify([]byte(sign2), csign2)
	credential := vp.VerifiableCredential[0].Credential
	if result && result2 {
		fmt.Printf("身份信息:\n")
		for _, item := range credential {
			fmt.Printf("%s \n", item)
		}
		fmt.Println()
		return true
	}
	return false
}

func Login(did string, privatekey *sm2.PrivateKey) bool { //登录
	vp := CreateVP(did, privatekey)
	fmt.Println("VP registered successfully")
	fmt.Println("-*******************************************-")

	result := verifyVP(vp, sm2PublicKey)
	if result {
		fmt.Println("VP verify successfully")
		return true
	} else {
		fmt.Println("VP verify failed")
		return false
	}
}
