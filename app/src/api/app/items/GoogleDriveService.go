package items

import (
	"api/app/models"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	sessions "github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	drive "google.golang.org/api/drive/v3"
)

var (
	tokenAuth                        = "token"
	errorFile                        = "error"
	cliente                          = "client"
	archivo                          = "file"
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
		if client := session.Get(cliente); client == nil {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Obteniendo el code"})
		return
	}

	state := c.Query("state")
	if state == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Obteniendo el state"})
		return
	}

	session := sessions.Default(c)
	session.Set(cliente, code)

	tok, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err == nil && tok != nil {
		if tokenJson, err := json.Marshal(tok); err == nil {
			session.Set(tokenAuth, string(tokenJson))
			session.Save()
		}
	}

	URL := session.Get("URL")
	if URL == nil {
		c.JSON(400, gin.H{"ERROR": "No se puede conseguir la url original"})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, URL.(string))
}

func getDriveApiService(session sessions.Session, tokenKey string) (*drive.Service, error) {

	dataToken := session.Get(tokenKey)
	if dataToken == nil {
		session.Delete(cliente)
		session.Save()
		return nil, errors.New("No se puede obtener el token")
	}

	token := dataToken.(string)
	if token == "" {
		session.Delete(cliente)
		session.Save()
		return nil, errors.New("El token se encuentra vacio")
	}

	tok := &oauth2.Token{}
	if err := json.Unmarshal([]byte(token), tok); err != nil {
		session.Delete(cliente)
		return nil, err
	}

	client := googleOauthConfig.Client(oauth2.NoContext, tok)
	if client == nil {
		session.Delete(cliente)
		session.Save()
		return nil, errors.New("No se puede obtener el cliente")
	}

	driveService, err := drive.New(client)
	if err != nil {
		session.Delete(cliente)
		session.Save()
		return nil, errors.New("No se puede obtener el Drive.Service* de drive " + err.Error())
	}

	return driveService, nil
}
