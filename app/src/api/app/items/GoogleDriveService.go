package items

import (
	"api/app/models"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	sessions "github.com/gin-contrib/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	drive "google.golang.org/api/drive/v3"
)

var (
	tokenAuth = "token"
	errorFile = "error"
	cliente = "client"
	archivo = "file"

	googleOauthConfig *oauth2.Config = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/drive"},
		Endpoint:     google.Endpoint,
	}

)

func AuthDriveApi() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if client := session.Get(cliente);client == nil {
			//Salvo el estado inicial del request
			session.Set("URL", c.Request.URL.RequestURI())
			session.Save()
			//Guardo un file , en caso de que sea un post para crear un archivo
			file := &models.File{}
			if erro := c.ShouldBindJSON(file); erro != nil {
				session.Set(errorFile, erro.Error())
				session.Save()
				} else {
						if fileJson, err := json.Marshal(file); err == nil {
							session.Set(archivo, string(fileJson))
							session.Save()
						}
				}
				//Redirecciono a la URL para authenticar el Usuario
				claveGeneradaRandom := "State"
				url := googleOauthConfig.AuthCodeURL(claveGeneradaRandom, oauth2.AccessTypeOffline)
				c.Redirect(http.StatusTemporaryRedirect, url)
			}
			c.Next()
		}
	}

func callback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error":"Obteniendo el code"})
		return
	}

	state := c.Query("state")
	if state == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error":"Obteniendo el state"})
		return
	}

	session := sessions.Default(c)
	session.Set(cliente, code)

	tok, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err == nil && tok != nil {
		 if tokenJson, err := json.Marshal(tok); err == nil  {
			session.Set(tokenAuth, string(tokenJson))
 			session.Save()
		}
	}

	URL := session.Get("URL")
	if URL == nil {
		c.JSON(400, gin.H{"ERROR":"No se puede conseguir la url original"})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, URL.(string))
}

func getDriveApiService(session sessions.Session,tokenKey string) (*drive.Service, error) {

	dataToken := session.Get(tokenKey);
	if dataToken == nil{
		return nil,errors.New("No se puede obtener el token")
	}

	token := dataToken.(string)
	if token == ""{
		return nil,errors.New("El token se encuentra vacio")
	}

	tok := &oauth2.Token{}
	 if err := json.Unmarshal([]byte(token),tok); err != nil{
		 return nil,err
	 }

	client := googleOauthConfig.Client(oauth2.NoContext, tok)
	if client == nil {
		return nil, errors.New("No se puede obtener el cliente")
	}

	driveService, err := drive.New(client)
	if err != nil {
		return nil, errors.New("No se puede obtener el Drive.Service* de drive "+err.Error())
	}

	return driveService, nil
}

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

	srv, err := getDriveApiService( session , tokenAuth )
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":"obteniendo el cliente de drive","Descripcion":err.Error()})
		return
	}

	file, err := srv.Files.Get(fileID).Fields("id,name,createdTime").Do()
	if err != nil || file == nil {
		c.Status(http.StatusBadRequest)
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
		c.Status(http.StatusOK) //El ID existe
		return
	}
}

// PostFile ...
func PostFile(c *gin.Context) {
	session := sessions.Default(c)
	errorFiles := session.Get(errorFile)
	marshalFile := session.Get(archivo)

	i := &models.File{}

	if errorFiles != nil || marshalFile != nil {
		if errorFiles != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error en los parametros", "description": errorFiles.(string)})
			session.Delete(errorFile)
			return
		}	else {
			stringFileJson := marshalFile.(string)
			if err := json.Unmarshal([]byte(stringFileJson), i); err != nil{
				c.JSON(http.StatusBadRequest, gin.H{"error": "Error en los parametros al convertir a json", "description": err.Error()})
				return
			}
			session.Delete(archivo)
		}
	} else {
		if erro := c.BindJSON(i); c.Request.ContentLength == 0 || erro != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error en los parametros", "description": erro.Error()})
			return
		}
	}
	session.Save()
	srv, err := getDriveApiService(session,tokenAuth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":"obteniendo el cliente de drive"+err.Error()})
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
