package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func CreateGitHubPR(repo Repo, title, description, base, head string) error {

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return fmt.Errorf("GITHUB_TOKEN not set")
	}

	body := map[string]string{
		"title": title,
		"head":  head,
		"base":  base,
		"body":  description,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	url := fmt.Sprintf(
		"%s/repos/%s/%s/pulls",
		repo.APIBase,
		repo.Owner,
		repo.Name,
	)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("github API returned %s", resp.Status)
	}

	fmt.Println("GitHub PR created")

	return nil
}
