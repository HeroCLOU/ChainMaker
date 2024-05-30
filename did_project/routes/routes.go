package routes

import (
	"did_project/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// 将数据库连接传递给Gin上下文
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
		log.Println("数据库连接已设置")
	})

	// 设置路由
	r.POST("/did", controllers.CreateDID)
	r.POST("/vc", controllers.CreateVC)
	r.POST("/vp/generate", controllers.GenerateVP)
	r.POST("/vp/verify", controllers.VerifyVP)

	log.Println("路由设置完成")
	return r
}
