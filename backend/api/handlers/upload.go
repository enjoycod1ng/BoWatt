package handlers

import (
	"io"
	"log"
	"net/http"
	"rag-pipeline/models"
	"rag-pipeline/services"
	"rag-pipeline/utils"
	"sync"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	embedder  *services.HuggingFaceEmbedder
	documents map[string]*models.Document
	cache     *utils.Cache
	mu        sync.RWMutex
}

func NewUploadHandler(embedder *services.HuggingFaceEmbedder, cache *utils.Cache) *UploadHandler {
	return &UploadHandler{
		embedder:  embedder,
		documents: make(map[string]*models.Document),
		cache:     cache,
	}
}

func (h *UploadHandler) Handle(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("Error getting form file: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Validate file size (10MB limit)
	if file.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large"})
		return
	}

	// Open and read file
	uploadedFile, err := file.Open()
	if err != nil {
		log.Printf("Error opening file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}
	defer uploadedFile.Close()

	// Read file content
	content, err := io.ReadAll(uploadedFile)
	if err != nil {
		log.Printf("Error reading file content: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file content"})
		return
	}

	textContent := string(content)
	if textContent == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is empty"})
		return
	}

	log.Printf("Generating embeddings for file: %s (size: %d bytes)", file.Filename, len(textContent))
	embedding, err := h.embedder.GetEmbeddings(textContent)
	if err != nil {
		log.Printf("Error generating embeddings: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate embeddings"})
		return
	}

	// Store document and embedding
	document := models.Document{
		Content:   textContent,
		Embedding: embedding,
		Filename:  file.Filename,
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	h.documents[file.Filename] = &document

	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded successfully",
		"filename": file.Filename,
	})
}

func (h *UploadHandler) GetDocuments() map[string]*models.Document {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.documents
}
