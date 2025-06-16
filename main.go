package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/ollama/ollama/api"
)

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

	ctx := context.Background()
	request := &api.GenerateRequest{
		Model:  model,
		Prompt: prompt,
	}
	responseFunc := func(resp api.GenerateResponse) error {
		fmt.Print(resp.Response)
		return nil
	}

	if err := apiClient.Generate(ctx, request, responseFunc); err != nil {
		log.Fatalf("failed to generate response: %v", err)
	}

}
