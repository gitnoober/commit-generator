package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

type LLMRequest struct {
	Prompt string `json:"prompt"`
	Model  string `json:"model"`
}

// LLMResponse represents a single response from Ollama API
type LLMResponse struct {
	Model           string `json:"model"`
	CreatedAt       string `json:"created_at"`
	Response        string `json:"response"`
	Done            bool   `json:"done"`
	Context         []int  `json:"context,omitempty"`
	TotalDuration   int64  `json:"total_duration,omitempty"`
	LoadDuration    int64  `json:"load_duration,omitempty"`
	PromptEvalCount int    `json:"prompt_eval_count,omitempty"`
	EvalCount       int    `json:"eval_count,omitempty"`
	EvalDuration    int64  `json:"eval_duration,omitempty"`
}

func runGitDiff() (string, error) {
	cmd := exec.Command("git", "diff")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.String(), err

}

func getSuggestedCommitMessage(diff string) (string, error) {
	llmURL := "http://localhost:11434/api/generate"
	model := "llama3.2"

	prompt := fmt.Sprintf("Here are the changes from `git diff`:\n%s\nGenerate a concise Git commit message.", diff)

	request := LLMRequest{
		Prompt: prompt,
		Model:  model,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %w", err)
	}

	resp, err := http.Post(llmURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error calling LLM: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse streaming NDJSON response
	scanner := bufio.NewScanner(resp.Body)
	var fullResponse string

	for scanner.Scan() {
		var llmResp LLMResponse
		if err := json.Unmarshal(scanner.Bytes(), &llmResp); err != nil {
			return "", fmt.Errorf("error parsing response: %w", err)
		}

		// Accumulate the response text
		fullResponse += llmResp.Response

		// Break if this is the final message
		if llmResp.Done {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	return fullResponse, nil
}

func main() {
	fmt.Println("Running git-helper")
	fmt.Println(os.Args)
	if len(os.Args) < 2 {
		panic("Usage: git-helper diff")
	}

	command := os.Args[1]

	switch command {
	case "diff":
		fmt.Println("Running git diff")
		diff, err := runGitDiff()
		if err != nil {
			fmt.Println("Error running git diff")
			panic(err)
		}
		fmt.Println(diff)
		suggestedCommitMessage, err := getSuggestedCommitMessage(diff)
		if err != nil {
			fmt.Println("Error getting suggested commit message", err)
			return
		}
		fmt.Println(suggestedCommitMessage)
	}

}
