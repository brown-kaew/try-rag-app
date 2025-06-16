# try-rag-app
A simple Retrieval-Augmented Generation (RAG) application  
Basic RAG knowledge: [RAG](https://www.geeksforgeeks.org/nlp/what-is-retrieval-augmented-generation-rag/)

## Install and Run Ollama
```bash
brew install ollama
brew services start ollama
ollama run llama3.2:latest
```
Prompt it something and use `/exit` to exit.

## Try Ollama API
Generate mode:
```bash
curl http://localhost:11434/api/generate -d '{
  "model": "llama3.2",
  "stream": false,
  "prompt":"What is NMIXX in K-pop?"
}'
```
Chat mode:
```bash
curl http://localhost:11434/api/chat -d '{
  "model": "llama3.2",
  "stream": false,
  "messages": [
    { "role": "user", "content": "What is NMIXX in K-pop?" }
  ]
}'
```
*Note: `stream = false` returns the response in a single message.*

## Run the RAG application
```bash
go run main.go <prompt>
```

## Thanks
Special thanks to [KProfiles](https://kprofiles.com/k-pop-girl-groups/) for K-pop group data used in this project.