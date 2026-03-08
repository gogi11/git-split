package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func CreateGitBucketPR(repo Repo, title, base, head string) error {

	token := os.Getenv("GITBUCKET_TOKEN")
	if token == "" {
		return fmt.Errorf("GITBUCKET_TOKEN not set")
	}

	url := fmt.Sprintf(
		"%s/repos/%s/%s/pulls",
		repo.APIBase,
		repo.Owner,
		repo.Name,
	)

	body := map[string]string{
		"title": title,
		"head":  head,
		"base":  base,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("gitbucket API returned %s", resp.Status)
	}

	fmt.Println("GitBucket PR created")

	return nil
}
