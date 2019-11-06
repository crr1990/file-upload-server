package Router

import (
	"file-upload-server/Controllers"
	"github.com/gin-gonic/gin"
)

func InitRouter() {
	router := gin.Default()
	v1 := router.Group("v1")
	{
		v1.POST("/testinsert", Controllers.TestInsert)
		v1.POST("/upload", Controllers.Upload)
		v1.POST("/merge", Controllers.MergeFile)
		v1.POST("/delete", Controllers.DeleteFile)
	}
	router.Run(":8622")
}
