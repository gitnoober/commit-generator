package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

type LLMRequest struct {
	Prompt string `json:"prompt"`
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
	llmURL := "http://localhost:11400/completions"

	prompt := fmt.Sprintf("Here are the changes from `git diff`:\n%s\nGenerate a concise Git commit message.", diff)
	requestBody, _ := json.Marshal(LLMRequest{Prompt: prompt})
	fmt.Println("Request body:")
	fmt.Println(string(requestBody))

	resp, err := exec.Command("curl", "-XPOST", llmURL, "--data", string(requestBody)).Output()

	if err != nil {
		return "", fmt.Errorf("error calling LLM: %s", err)
	}

	var llmResp LLMResponse
	if err := json.Unmarshal(resp, &llmResp); err != nil {
		return "", fmt.Errorf("error parsing LLM response: %s", err)
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
