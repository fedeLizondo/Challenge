package items

import (
	"github.com/gin-gonic/gin"
)

// Configure for files
func ConfigureForFiles(r *gin.Engine) {
	//r.Use(AuthGoogle())
	r.GET("/search-in-doc/:id", AuthDriveApi(), SearchForFile)
	r.POST("/file", AuthDriveApi(), PostFile)
	r.GET("callback", callback)
}
