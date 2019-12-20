package Router

import (
	"file-upload-server/Controllers"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter() {
	router := gin.Default()
	router.Use(CorsMiddleware())
	v1 := router.Group("v1")
	{
		v1.POST("/testinsert", Controllers.TestInsert)
		v1.POST("/upload", Controllers.Upload)
		v1.POST("/merge", Controllers.MergeFile)
		v1.POST("/delete", Controllers.DeleteFile)
	}
	router.Run(":8622")
}

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

			// 核心处理方式
			c.Header("Access-Control-Allow-Origin", "http://47.102.149.201:8620")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, content-type, Accept")
			c.Header("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT, DELETE")
			c.Header("Access-Control-Max-Age", "1800")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}

		c.Next()
	}
}
