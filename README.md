# try-rag-app
An simple RAG application

## Install Ollama
```bash
brew install ollama
```

### Serve Ollama
```bash
brew services start ollama
```

### Try to run model
```bash
ollama run llama3.2:latest
```
Prompt it something and use `/exit` to exit


### Try to server endpoint
Generate mode
```bash
curl http://localhost:11434/api/generate -d '{
  "model": "llama3.2",
  "stream": false,
  "prompt":"What is NMIXX in K-pop?"
}'
```
Chat mode
```bash
curl http://localhost:11434/api/chat -d '{
  "model": "llama3.2",
  "stream": false,
  "messages": [
    { "role": "user", "content": "What is NMIXX in K-pop?" }
  ]
}'
```
** Stream = false means return response in a single message

## Ran simple RAG application
```Bash
go run main.go <prompt>
```
