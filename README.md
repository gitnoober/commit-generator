# Commit Message Generator

A Go client that uses Ollama's local LLM to generate Git commit messages based on diff content.

## Features

- Connects to local Ollama instance
- Streams responses for efficient processing
- Handles NDJSON response format
- Generates concise commit messages from git diff output

## Prerequisites

- Go 1.16 or higher
- Ollama running locally on port 11434
- llama2 model installed in Ollama

## Usage

1. Start your local Ollama instance
2. Run the application:

```bash
go run main.go
```

## API Structure

### Request Format

```go
type LLMRequest struct {
    Model  string `json:"model"`
    Prompt string `json:"prompt"`
}
```

### Response Format

```go
type LLMResponse struct {
    Model            string `json:"model"`
    CreatedAt        string `json:"created_at"`
    Response         string `json:"response"`
    Done             bool   `json:"done"`
    Context          []int  `json:"context,omitempty"`
    TotalDuration    int64  `json:"total_duration,omitempty"`
    LoadDuration     int64  `json:"load_duration,omitempty"`
    PromptEvalCount  int    `json:"prompt_eval_count,omitempty"`
    EvalCount        int    `json:"eval_count,omitempty"`
    EvalDuration     int64  `json:"eval_duration,omitempty"`
}
```

## Error Handling

The client includes comprehensive error handling for:
- Request marshaling
- HTTP communication
- Response parsing
- Stream processing

## Example

```go
diff := `diff --git a/file.txt b/file.txt
+ Added new feature
- Removed old code`

message, err := getSuggestedCommitMessage(diff)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Suggested commit message: %s\n", message)
```
