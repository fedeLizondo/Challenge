package items

import (
	"github.com/gin-gonic/gin"
)

// Routes para File
func ConfigureForFiles(r *gin.Engine) {
	r.GET("/search-in-doc/:id", AuthDriveApi(), SearchForFile)
	r.POST("/file", AuthDriveApi(), PostFile)
	r.GET("callback", callback)
}
