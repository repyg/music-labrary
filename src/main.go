package main

import (
	"fmt"
	"log"
	"os"

	"music-library/src/api"
	"music-library/src/config"
	"music-library/src/repository"
	"music-library/src/service"

	_ "music-library/docs" // Импорт сгенерированной документации

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title Music Library API
// @version 1.0
// @description API для управления библиотекой песен
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Загрузка .env файла
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Подключение к базе данных
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Инициализация компонентов
	repo := repository.NewSongRepository(db)
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	service := service.NewSongService(repo, config.ExternalAPIURL, config)
	handler := api.NewHandler(service)

	// Настройка роутера
	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		songs := v1.Group("/songs")
		{
			songs.GET("", handler.GetSongs)
			songs.POST("", handler.AddSong)
			songs.GET("/:id", handler.GetSongById)
			songs.PUT("/:id", handler.UpdateSong)
			songs.DELETE("/:id", handler.DeleteSong)
			songs.GET("/:id/lyrics", handler.GetSongLyrics)
		}
	}

	// Подключение Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Запуск сервера
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
