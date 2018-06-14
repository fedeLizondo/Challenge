package items

import (
	"api/app/models"
	"net/http"
	"strings"

	_ "github.com/gin-contrib/sessions"
	_ "github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/zalando/gin-oauth2/google"
	//"github.com/zalando/gin-oauth2/google"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	//_ "google.golang.org/api/drive/v3"
	drive "google.golang.org/api/drive/v3"
)

var (
	googleOauthConfig *oauth2.Config = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/callback",
		ClientID:     "265169521968-cdu9f25cbim1qa86c4lidnv44dffhti0.apps.googleusercontent.com", // os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: "ls6bFVb9Xqpj_f9vKz12vNpk",                                                 //os.getEnv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/drive"},
		Endpoint:     google.Endpoint,
	}
)

func AuthDriveApi() gin.HandlerFunc {
	return func(c *gin.Context) {

		//GENERAR UNA CLAVE RANDOM Y GUARDAR EN UNA session
		claveGeneradaRandom := "state"
		//SESSION guardo el context (?) o guardo el request
		//NECESITO GUARDAR EL POST Y EL JSON(O EL OBJETO) Y EL GET

		conf := googleOauthConfig
		url := conf.AuthCodeURL(claveGeneradaRandom, oauth2.AccessTypeOffline)
		c.Redirect(http.StatusTemporaryRedirect, url)

	}
}

func callback(c *gin.Context) {
	//id := "0BweY51pBQw37bzdzTFo3Mkw4dTg"

	code := c.Query("code")

	if code == "" {
		c.JSON(http.StatusBadRequest, "ERROR obteniendo el code")
		return
	}

	tok, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, "ERROR obteniendo el token")
		return
	}

	client := googleOauthConfig.Client(oauth2.NoContext, tok)
	if client == nil {
		c.JSON(http.StatusBadRequest, "ERROR obteniendo el cliente")
		return
	}

}

func SearchForFile(c *gin.Context) {

	fileID := strings.TrimSpace(c.Param("id"))
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id_error"})
		return
	}
	word := strings.TrimSpace(c.Param("word"))
	if word != "" {
		word = word
	}
	//CODIGO DEL HANDLER
	conf := googleOauthConfig
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
	return

	// client := c.get("client")
	// //REQUEST
	// srv, err := drive.New(client)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, "ERROR obteniendo el cliente de drive")
	// 	return
	// }
	//
	// file, err := srv.Files.Get(fileID).Fields("id,name,createdTime").Do()
	// if err != nil || file == nil {
	// 	c.Status(http.StatusBadRequest)
	// 	return
	// }
	//
	// if word != "" {
	// 	r, err := srv.Files.List().Q("fullText contains '" + word + "' and name = '" + file.Name + "' and createdTime = '" + file.CreatedTime + "'").Fields("files(id)").Do()
	// 	if err != nil {
	// 		c.Status(http.StatusBadRequest)
	// 		return
	// 	}
	// 	//c.Status(http.StatusOK)
	//
	// 	if len(r.Files) > 0 {
	// 		for _, i := range r.Files {
	// 			if i.Id == file.Id {
	// 				c.Status(http.StatusOK)
	// 				return
	// 			}
	// 		}
	// 	}
	// 	c.Status(http.StatusBadRequest)
	// 	return
	// } else {
	// 	c.Status(http.StatusOK) //El ID existe
	// 	return
	// }

}

// PostFile ...
func PostFile(c *gin.Context) {
	i := &models.File{}
	if erro := c.BindJSON(i); c.Request.ContentLength == 0 || erro != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error en los parametros", "description": erro.Error()})
		return
	}

	client, ok := c.Get("client")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ERROR obteniendo el cliente de drive", "description": ""})
		return
	}

	srv, err := drive.New(client.(*http.Client))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ERROR obteniendo el cliente de drive", "description": err.Error()})
		return
	}

	file := &drive.File{Name: i.Title, Description: i.Description, MimeType: "text/plain"}

	resp, err := srv.Files.Create(file).Do()
	if err != nil || resp == nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	i.ID = resp.Id
	c.JSON(http.StatusOK, i)
	return
}
