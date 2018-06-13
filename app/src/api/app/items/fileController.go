package items

//
// import (
// 	"net/http"
// 	"strings"
// 	"api/app/models"
// 	// "fmt"
// 	// "io/ioutil"
// 	// "os"
// 	"golang.org/x/oauth2"
// 	"github.com/gin-gonic/gin"
// )
//
// // GetFile ...
// func SearchForFile(c *gin.Context) {
// 	// fileID := strings.TrimSpace(c.Param("id"))
// 	// if fileID == "" {
//   //   c.JSON(http.StatusBadRequest, gin.H{"error": "id_error"})
// 	// 	return
// 	// }
//   // word := strings.TrimSpace(c.Param("word"))
//   //
//   //
//   // conf := googleOauthConfig
//   // url := conf.AuthCodeURL("state",oauth2.AccessTypeOffline)
//   //
//   //  // c.Use(google.Auth())
//   //  // file, err := Is.file(fileID)
// 	//  // if err != nil {
// 	//  // 	c.JSON(http.StatusInternalServerError, gin.H{"error": "find_error", "description": err.Error()})
// 	//  // 	return
//   //  //
//   //  // }
// 	// c.JSON(200, "file:"fileID+"word:"+word)
//
//
//   fileID := strings.TrimSpace(c.Param("id"))
// 	if fileID == "" {
//     c.JSON(http.StatusBadRequest, gin.H{"error": "id_error"})
// 		return
// 	}
//   word := strings.TrimSpace(c.Param("word"))
//
//   session := session.Default(c)
//   client := session.get("client")
//   if client == nil {
//     c.JSON(http.StatusInternalServerError, gin.H{"error": "Client not found"})
// 		return
//   }
//
//   resp, err := client.Get("https://www.googleapis.com/drive/v3/files/" + fileID )
// 	if err != nil {
// 		c.JSON(404, "ERROR")
// 	}
// 	defer resp.Body.Close()
//   c.JSON(http.StatusOK, "TODO:MOSTRAR RESPUESTA")
//
//   return
// }
//
// // PostFile ...
// func PostFile(c *gin.Context) {
// 	i := &models.File{}
// // 	if err := c.BindJSON(i); c.Request.ContentLength == 0 || err != nil {
// // 		c.JSON(http.StatusBadRequest, gin.H{"error": "bind_error", "description": err.Error()})
// // 		return
// // 	}
// // 	//err := Is.CreateFile(i)
// // 	if err != nil {
// // 		c.JSON(http.StatusBadRequest, gin.H{"error": "save_error", "description": err.Error()})
// // 		return
// // 	}
// 	c.JSON(201, i)
// }
