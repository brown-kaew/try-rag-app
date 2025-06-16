package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/ollama/ollama/api"
)

// Simple in-memory document store
var documents = []string{
	"Golang is a statically typed, compiled programming language designed at Google.",
	"Retrieval-Augmented Generation (RAG) combines retrieval and generation for better answers.",
	"Ollama provides an API for running large language models locally.",
}

// Simple keyword-based retrieval
func retrieveRelevantDocs(query string, docs []string, topK int) []string {
	var results []string
	query = strings.ToLower(query)
	for _, doc := range docs {
		if strings.Contains(strings.ToLower(doc), query) {
			results = append(results, doc)
		}
		if len(results) >= topK {
			break
		}
	}
	return results
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: go run main.go <prompt>")
	}

	prompt := os.Args[1]
	ollamaUrl := "http://localhost:11434"
	model := "llama3.2:latest"

	parsedUrl, err := url.Parse(ollamaUrl)
	if err != nil {
		log.Fatalf("failed to parse URL: %v", err)
	}

	apiClient := api.NewClient(parsedUrl, http.DefaultClient)

	// Retrieve relevant documents
	retrieved := retrieveRelevantDocs(prompt, documents, 2)
	contextText := strings.Join(retrieved, "\n")

	// Agment the prompt with retrieved context
	augmentedPrompt := prompt
	if contextText != "" {
		augmentedPrompt = fmt.Sprintf("Context:\n%s\n\nQuestion: %s", contextText, prompt)
		fmt.Println("Augmented Prompt:")
	}
	fmt.Println(augmentedPrompt)

	// Generate response using the Ollama API
	ctx := context.Background()
	request := &api.GenerateRequest{
		Model:  model,
		Prompt: augmentedPrompt,
	}
	responseFunc := func(resp api.GenerateResponse) error {
		fmt.Print(resp.Response)
		return nil
	}

	if err := apiClient.Generate(ctx, request, responseFunc); err != nil {
		log.Fatalf("failed to generate response: %v", err)
	}

}
