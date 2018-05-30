package items

import (
	"net/http"
	"strings"
	"os"
	"api/app/models"

	"github.com/gin-gonic/gin"
)

// GetItem ...
func GetItem(c *gin.Context) {
	itemID := strings.TrimSpace(c.Param("id"))
	if itemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id_error"})
		return
	}

	item, err := Is.Item(itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "find_error", "description": err.Error()})
		return
	}
	c.JSON(200, item)
	return
}

// PostItem ...
func PostItem(c *gin.Context) {
	i := &models.Item{}
	if err := c.BindJSON(i); c.Request.ContentLength == 0 || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bind_error", "description": err.Error()})
		return
	}
	err := Is.CreateItem(i)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "save_error", "description": err.Error()})
		return
	}
	c.JSON(201, i)
}

//Search For File then for word
func SearchForFile(c *gin.Context)  {
	fileId := strings.TrimSpace(c.Param("id"))
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id_error"})
		return
	}
	word := strings.TrimSpace(c.DefaultQuery("word",""))
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "word_error"})
		return
	}

	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read server secret file", "description": err.Error()})
		return
	}
	// If modifying these scopes, delete your previously saved client_secret.json.
	config, err := google.ConfigFromJSON(b, drive.DriveFileScope)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to parse server secret file to config", "description": err.Error()})
		return
	}

	srv, err := drive.New(getClient(config))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve Drive client", "description": err.Error()})
		return
	}

	r, err := srv.Files.get(fileId).List().
		Fields( fmt.Sprintf("fullText contains %s , nextPageToken, files(id, name)",word)).Do()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found files", "description": err.Error()})
		return
	}
	fmt.Println("Files:")
	if len(r.Files) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found files", "description": err.Error()})
		return
		}
		// else {
		// 	for _, i := range r.Files {
		// 		fmt.Printf("%s (%s)\n", i.Name, i.Id)
		// 	}
		// }

	c.JSON(200, r.Files )
	return
}
//Post Create a File
func PostFile(c *gin.Context) {
	i := &models.File{}
	if err := c.BindJSON(i); c.Request.ContentLength == 0 || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bind_error", "description": err.Error()})
		return
	}

	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read server secret file", "description": err.Error()})
		return
	}

	// If modifying these scopes, delete your previously saved client_secret.json.
	config, err := google.ConfigFromJSON(b, drive.DriveFileScope)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to parse server secret file to config", "description": err.Error()})
		return
	}

	f:=&drive.File{
		Name: i.Title,
		Description: i.Description
	}

	srv, err := drive.New(getClient(config))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve Drive client", "description": err.Error()})
		return
	}

	r, err := srv.Files.Create(f).PageSize(10).Do()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve files", "description": err.Error()})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "save_error", "description": err.Error()})
		return
	}

	c.JSON(http.StatusOK, i)
}
