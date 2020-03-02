package main

import (
	"github.com/gin-gonic/gin"
	"github.com/heroku/changemomentum/db"
	"github.com/heroku/changemomentum/handlers"
	"github.com/heroku/changemomentum/logger"
	"log"
	"os"
)


func setupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "static")

	handler := handlers.NewHandler()

	router.GET("/", handler.List)
	router.GET("/contacts", handler.List)
	router.GET("/contacts/new", handler.AddForm)
	router.POST("/contacts/new", handler.Add)
	router.GET("/contacts/addphone/:id", handler.AddFormPhone)
	router.POST("/contacts/addphone/:id", handler.AddPhone)
	router.GET("/contacts/edit/:id", handler.Edit)
	router.POST("/contacts/update/:id", handler.Update)
	router.DELETE("/contacts/delete/:id", handler.Delete)


	router.GET("/search", handler.Search)


	router.GET("/participants", handler.List)

	router.GET("/participants/new", handler.AddForm)
	router.POST("/participants/new", handler.Add)

	router.GET("/participants/{id}/move", func(c *gin.Context) {
		//TODO move
	})

	router.GET("/participants/edit/:id", handler.Edit)
	router.POST("/participants/update/:id", handler.Update)
	router.DELETE("/participants/delete/:id", handler.Delete)

	router.POST("/participants/registration/:id", handler.Registration)

	router.POST("/participants/unregistration/:id", handler.UnRegistration)






	return router
}


func main() {
	logger.NewLogger()
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	DatabaseURL := os.Getenv("DATABASE_URL")

	postgresrepo, err := db.NewPostgresrepo(&DatabaseURL)
	if err != nil {
		logger.Error("Error DB. Please check your connect for DB", err, DatabaseURL)
		log.Fatal()
	}
	db.SetRepository(postgresrepo)
	err = postgresrepo.Db.Ping()
	if err != nil {
		logger.Error("Error DB. Please check your connect for DB", err, DatabaseURL)
		log.Fatal()
	}

	defer db.Close()

	logger.Info("Connect to DB ", DatabaseURL)

	router := setupRouter()

	router.Run(":" + port)
}