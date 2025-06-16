package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/blevesearch/bleve/v2"
	"github.com/ollama/ollama/api"
)

var bleveIndex bleve.Index

// Read all .txt files from the kpop-data directory as documents
func loadDocumentsFromDir(dir string) ([]string, error) {
	var docs []string
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".txt") {
			continue
		}
		content, err := os.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}
		docs = append(docs, string(content))
	}
	return docs, nil
}

// Build Bleve index once after loading documents
func buildBleveIndex(docs []string) (bleve.Index, error) {
	indexMapping := bleve.NewIndexMapping()
	index, err := bleve.NewMemOnly(indexMapping)
	if err != nil {
		return nil, err
	}
	for i, doc := range docs {
		err := index.Index(fmt.Sprintf("%d", i), doc)
		if err != nil {
			return nil, err
		}
	}
	return index, nil
}

// Use Bleve for document retrieval (reuse global index)
func retrieveRelevantDocs(query string, docs []string, topK int) []string {
	if bleveIndex == nil {
		log.Fatalf("bleve index not initialized")
	}
	searchRequest := bleve.NewSearchRequestOptions(bleve.NewQueryStringQuery(query), topK, 0, false)
	searchResult, err := bleveIndex.Search(searchRequest)
	if err != nil {
		log.Printf("bleve search error: %v", err)
		return nil
	}
	var results []string
	for _, hit := range searchResult.Hits {
		idx := hit.ID
		i := 0
		fmt.Sscanf(idx, "%d", &i)
		results = append(results, docs[i])
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

	// Load documents from kpop-data directory
	documents, err := loadDocumentsFromDir("kpop-data")
	if err != nil {
		log.Fatalf("failed to load documents: %v", err)
	}

	// Build Bleve index once
	bleveIndex, err = buildBleveIndex(documents)
	if err != nil {
		log.Fatalf("failed to build bleve index: %v", err)
	}

	parsedUrl, err := url.Parse(ollamaUrl)
	if err != nil {
		log.Fatalf("failed to parse URL: %v", err)
	}

	apiClient := api.NewClient(parsedUrl, http.DefaultClient)

	// Retrieve relevant documents
	retrieved := retrieveRelevantDocs(prompt, documents, 1)
	contextText := strings.Join(retrieved, "\n")

	// Agment the prompt with retrieved context
	augmentedPrompt := prompt
	if contextText != "" {
		augmentedPrompt = fmt.Sprintf("Context:\n%s\n\nQuestion: %s", contextText, prompt)
		fmt.Println("Augmented Prompt:")
	}
	fmt.Println(augmentedPrompt)
	fmt.Println("---")
	fmt.Println("Answer:")

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
