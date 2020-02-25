package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
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
	router.GET("/contacts/{id}", handler.AddFormPhone)
	router.POST("/contacts/{id}", handler.AddPhone)
	router.GET("/search", handler.Search)

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