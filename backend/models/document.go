package models

type Document struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Embedding []float32 `json:"embedding"`
	Filename  string
}

type Query struct {
	Text string `json:"text"`
}

type Response struct {
	Text string `json:"response"`
}
