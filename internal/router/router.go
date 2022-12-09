package router

import (
	"{{cookiecutter.project_slug}}/internal/api"
	"{{cookiecutter.project_slug}}/internal/storage"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Router ...
type Router struct {
	router *gin.Engine
	store  storage.Store
	creds  credentials
}

type credentials struct {
	username string
	password string
}

// New ...
func New(s storage.Store) *Router {
	r := &Router{
		router: gin.Default(),
		store:  s,
		creds: credentials{
			username: viper.GetString("web.username"),
			password: viper.GetString("web.password"),
		},
	}

	r.router.Use(gin.BasicAuth(gin.Accounts{
		r.creds.username: r.creds.password,
	}))
	r.initRoutes()
	r.setCors()
	return r
}

// ServeHTTP ...
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}

func (r *Router) setCors() {
	r.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:80", "http://127.0.0.1:80"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}

func (r *Router) initRoutes() {
	// API Router Group
	apiRouter := r.router.Group("/api")

	// Pet controllers
	{
		apiRouter.GET("/pet/:id", api.GetPet)
		apiRouter.POST("/pet", api.CreatePet)
		apiRouter.PUT("/pet/:id", api.UpdatePet)
		apiRouter.DELETE("/pet/:id", api.DeletePet)
		apiRouter.GET("/pets", api.GetPets)
	}
}
