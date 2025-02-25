package handlers

import (
	"fmt"
	"rag-pipeline/models"
	"rag-pipeline/services"
	"rag-pipeline/utils"

	"github.com/gin-gonic/gin"
)

type QueryHandler struct {
	embedder  *services.HuggingFaceEmbedder
	llm       *services.LLM
	documents map[string]*models.Document
	cache     *utils.Cache
}

func NewQueryHandler(embedder *services.HuggingFaceEmbedder, llm *services.LLM, documents map[string]*models.Document, cache *utils.Cache) *QueryHandler {
	return &QueryHandler{
		embedder:  embedder,
		llm:       llm,
		documents: documents,
		cache:     cache,
	}
}

func (h *QueryHandler) Handle(c *gin.Context) {
	var query struct {
		Text string `json:"text"`
	}

	if err := c.BindJSON(&query); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// Check cache
	if response, found := h.cache.Get(query.Text); found {
		c.JSON(200, gin.H{"response": response})
		return
	}

	// Generate embeddings for query
	queryEmbedding, err := h.embedder.GetEmbeddings(query.Text)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to generate query embeddings: %v", err)})
		return
	}

	// Find relevant document using cosine similarity
	var bestMatch string
	var maxSimilarity float32 = -1

	for _, doc := range h.documents {
		similarity := utils.CosineSimilarity(queryEmbedding, doc.Embedding)
		if similarity > maxSimilarity {
			maxSimilarity = similarity
			bestMatch = doc.Content
		}
	}

	if bestMatch == "" {
		c.JSON(404, gin.H{"error": "No relevant context found"})
		return
	}

	// Generate response using LLM
	prompt := fmt.Sprintf("Context: %s\n\nQuestion: %s", bestMatch, query.Text)
	response, err := h.llm.GenerateResponse(prompt)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to generate response: %v", err)})
		return
	}

	// Cache the response
	h.cache.Set(query.Text, response)

	c.JSON(200, gin.H{"response": response})
}
