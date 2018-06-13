package items

import (
	"github.com/gin-gonic/gin"
)

// Configure for files
func ConfigureForFiles(r *gin.Engine) {
	//r.Use(AuthGoogle())
	r.GET("/search-in-doc/:id", SearchForFile)
	r.POST("/file", PostFile)
	r.GET("callback", callback)
}
