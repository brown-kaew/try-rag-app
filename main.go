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

	"github.com/bbalet/stopwords"
	"github.com/ollama/ollama/api"
)

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

// Improved retrieval: rank by word overlap, using stopwords lib
func retrieveRelevantDocs(query string, docs []string, topK int) []string {
	cleanQuery := stopwords.CleanString(query, "en", false)
	queryTokens := strings.Fields(strings.ToLower(cleanQuery))
	type docScore struct {
		doc   string
		score int
	}
	var scored []docScore
	for _, doc := range docs {
		cleanDoc := stopwords.CleanString(doc, "en", false)
		docTokens := strings.Fields(strings.ToLower(cleanDoc))
		score := 0
		for _, qt := range queryTokens {
			for _, dt := range docTokens {
				if qt == dt {
					score++
				}
			}
		}
		if score > 0 {
			scored = append(scored, docScore{doc, score})
		}
	}
	// Sort by score descending
	for i := 0; i < len(scored)-1; i++ {
		for j := i + 1; j < len(scored); j++ {
			if scored[j].score > scored[i].score {
				scored[i], scored[j] = scored[j], scored[i]
			}
		}
	}
	var results []string
	for i := 0; i < len(scored) && i < topK; i++ {
		results = append(results, scored[i].doc)
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
