package main

import (
	"log"
	"os"
	"rag-pipeline/api/handlers"
	"rag-pipeline/services"
	"rag-pipeline/utils"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, checking system environment variables")
	}

	apiKey := os.Getenv("HUGGINGFACE_API_KEY")
	if apiKey == "" {
		log.Fatal("HUGGINGFACE_API_KEY not found in environment")
	}

	// Initialize services with dependency injection
	embedder := services.NewHuggingFaceEmbedder(apiKey)
	llm := services.NewLLM(apiKey)
	cache := utils.NewCache()

	// Initialize handlers with dependencies
	uploadHandler := handlers.NewUploadHandler(embedder, cache)
	queryHandler := handlers.NewQueryHandler(embedder, llm, uploadHandler.GetDocuments(), cache)

	// Setup Gin router
	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1", "::1"})

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Routes
	r.POST("/upload", uploadHandler.Handle)
	r.POST("/query", queryHandler.Handle)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
