package main

import (
	"did/method"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tjfoc/gmsm/x509"
)

func main() {
	r := gin.Default() // 初始化Gin框架，默认中间件包括logger和recovery

	// 加载HTML模板，用于渲染页面
	r.LoadHTMLFiles("templates/first.html", "templates/register.html", "templates/did.html", "templates/login.html", "templates/main.html")

	// 设置静态文件目录，用于加载CSS、JS、图片等资源
	r.Static("/static", "./static")

	// 路由：主页，渲染"first.html"
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "first.html", nil)
	})

	// 路由：注册页面，渲染"register.html"
	r.GET("/register", func(c *gin.Context) {
		c.HTML(200, "register.html", nil)
	})

	// 路由：登录页面，渲染"login.html"
	r.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", nil)
	})

	// POST路由：处理注册请求，生成DID和私钥，返回"did.html"
	r.POST("/registersuccess", func(c *gin.Context) {
		// 获取表单数据
		name := c.PostForm("name")
		cardid := c.PostForm("idCard")
		school := c.PostForm("school")
		schoolrecord := c.PostForm("education")

		// 调用自定义方法包的Register函数
		did, _ := method.Register(name, cardid, school, schoolrecord)

		// 私钥文件路径
		privateKeyPath := "private_key.pem"

		// 检查私钥文件是否存在
		_, err := os.Stat(privateKeyPath)
		if os.IsNotExist(err) {
			// 如果私钥文件不存在，返回500错误
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Private key file does not exist"})
			return
		}

		// 将DID和私钥文件路径传递给前端
		c.HTML(http.StatusOK, "did.html", gin.H{
			"Did":            did,
			"PrivateKeyPath": privateKeyPath,
		})
	})

	// GET路由：下载私钥文件
	r.GET("/download/private_key", func(c *gin.Context) {
		// 私钥文件路径
		privateKeyPath := "private_key.pem"

		// 检查私钥文件是否存在
		_, err := os.Stat(privateKeyPath)
		if os.IsNotExist(err) {
			// 如果私钥文件不存在，返回404错误
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		// 设置下载响应头，开始文件下载
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", privateKeyPath))
		c.FileAttachment(privateKeyPath, "private_key.pem")
	})

	// POST路由：处理登录请求，验证DID和私钥
	r.POST("/home", func(c *gin.Context) {
		// 获取DID
		did := c.PostForm("did")

		// 获取上传的私钥文件
		privateKeyFile, _ := c.FormFile("privatekey")

		// 打开上传的私钥文件
		file, err := privateKeyFile.Open()
		if err != nil {
			// 文件打开失败，记录错误并退出
			log.Fatalf("打开上传的私钥文件时出错: %v", err)
		}
		defer file.Close() // 关闭文件

		// 读取私钥文件内容
		privateKeyPem, err := ioutil.ReadAll(file)
		if err != nil {
			// 读取文件内容失败，记录错误并退出
			log.Fatalf("读取上传的私钥文件内容时出错: %v", err)
		}

		// 解码PEM格式的私钥
		block, _ := pem.Decode(privateKeyPem)
		if block == nil || block.Type != "PRIVATE KEY" {
			// PEM解码失败，记录错误并退出
			log.Fatalf("解码私钥时出错")
		}

		// 从PEM格式解析私钥
		priv, err := x509.ReadPrivateKeyFromPem(privateKeyPem, nil)
		if err != nil {
			// 私钥解析失败，记录错误并退出
			log.Fatalf("解析私钥时出错: %v", err)
		}

		// 登录逻辑，调用自定义方法包的Login函数
		if method.Login(did, priv) {
			// 登录成功，渲染"main.html"
			c.HTML(200, "main.html", nil)
		} else {
			c.JSON(200, gin.H{
				"massage": "信息有误,请重新输入",
			})
		}
	})

	// 启动服务器，监听9090端口
	r.Run(":9090")
}
