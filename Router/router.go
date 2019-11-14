package Router

import (
	"file-upload-server/Controllers"
	"github.com/gin-gonic/gin"
	"regexp"
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
		origin := c.Request.Header.Get("Origin")
		var filterHost = [...]string{"http://*.hfjy.com"}
		// filterHost 做过滤器，防止不合法的域名访问
		var isAccess = false
		for _, v := range(filterHost) {
			match, _ := regexp.MatchString(v, origin)
			if match {
				isAccess = true
			}
		}
		if isAccess {
			// 核心处理方式
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
			c.Header("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT, DELETE")
			c.Set("content-type", "application/json")
		}
		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}

		c.Next()
	}
}
