package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

type LLMRequest struct {
	Prompt string `json:"prompt"`
	Model  string `json:"model"`
}

type LLMResponse struct {
	Response string `json:"response"`
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
	jsonData, _ := json.Marshal(request)
	fmt.Println("Request body:")
	fmt.Println(string(jsonData))

	resp, err := http.Post(llmURL, "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		return "", fmt.Errorf("error calling LLM: %s", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	// Parse the JSON response into LLMResponse struct
	var llmResp LLMResponse
	if err := json.Unmarshal(body, &llmResp); err != nil {
		return "", fmt.Errorf("error parsing LLM response: %w", err)
	}
	return llmResp.Response, nil
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
			panic(err)
		}
		fmt.Println(suggestedCommitMessage)
	}

}
