package items

import (
	//"api/app/models"
	"github.com/gin-gonic/gin"
)

// Configure for items
func ConfigureForFiles(r *gin.Engine) {
	r.GET("/search-in-doc/:id", SearchForFile)
	r.POST("/file", PostFile)
}
