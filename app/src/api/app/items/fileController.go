package items

import (
	"api/app/models"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	sessions "github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	drive "google.golang.org/api/drive/v3"
)

func SearchForFile(c *gin.Context) {

	fileID := strings.TrimSpace(c.Param("id"))
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id_error"})
		return
	}

	word := strings.TrimSpace(c.Query("word"))

	session := sessions.Default(c)
	//Limpio si quedo el error cuando se logueo la primera vez
	session.Delete(errorFile)
	session.Delete(archivo)

	srv, err := getDriveApiService(session, tokenAuth)
	if err != nil || srv == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "obteniendo el cliente de drive ", "Descripcion": err.Error()})
		return
	}

	file, err := srv.Files.Get(fileID).Fields("id,name,createdTime").Do()
	if err != nil || file == nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if word != "" {
		r, err := srv.Files.List().Q("fullText contains '" + word + "' and name = '" + file.Name + "' and createdTime = '" + file.CreatedTime + "'").Fields("files(id)").Do()
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		if len(r.Files) > 0 {
			for _, i := range r.Files {
				if i.Id == file.Id {
					c.Status(http.StatusOK)
					return
				}
			}
		}
		c.Status(http.StatusBadRequest)
		return
	} else {
		c.Status(http.StatusOK) //El ID existe pero no buscaban un string vacio
		return
	}
}

func ObtenerFile(c *gin.Context) (*models.File, error) {
	session := sessions.Default(c)
	errorFiles := session.Get(errorFile)
	marshalFile := session.Get(archivo)
	i := &models.File{}
	session.Delete(errorFile)
	session.Delete(archivo)
	session.Save()

	if errorFiles != nil || marshalFile != nil {
		if errorFiles != nil {
			return nil, errors.New(errorFiles.(string))
		} else {
			stringFileJson := marshalFile.(string)
			if err := json.Unmarshal([]byte(stringFileJson), i); err != nil {
				return nil, err
			}
		}
	} else {
		if err := c.BindJSON(i); c.Request.ContentLength == 0 || err != nil {
			return nil, err
		}
	}
	return i, nil
}

// PostFile ...
func PostFile(c *gin.Context) {
	i, err := ObtenerFile(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error en los parametros", "description": err.Error()})
	}

	session := sessions.Default(c)
	srv, err := getDriveApiService(session, tokenAuth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "obteniendo el cliente de drive" + err.Error()})
		return
	}

	file := &drive.File{Name: i.Title, Description: i.Description, MimeType: "text/plain"}

	resp, err := srv.Files.Create(file).Do()
	if err != nil || resp == nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	i.ID = resp.Id //El ID se genera despues de que se guarda el archivo en el drive
	c.JSON(http.StatusOK, i)
	return
}
