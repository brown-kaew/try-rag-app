package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/ollama/ollama/api"
)

func main() {
	ollamaUrl := "http://localhost:11434"
	model := "llama3.2:latest"

	url, err := url.Parse(ollamaUrl)
	if err != nil {
		fmt.Errorf("failed to parse URL: %v", err)
	}

	apiClient := api.NewClient(url, http.DefaultClient)

	context := context.Background()
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <prompt>")
		return
	}
	prompt := os.Args[1]
	request := &api.GenerateRequest{
		Model:  model,
		Prompt: prompt,
	}
	responseFunc := func(resp api.GenerateResponse) error {
		fmt.Print(resp.Response)
		return nil
	}

	if err := apiClient.Generate(context, request, responseFunc); err != nil {
		fmt.Errorf("failed to generate response: %v", err)
	}
}
