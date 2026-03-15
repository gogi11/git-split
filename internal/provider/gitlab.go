package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

func CreateGitLabMR(repo Repo, title, description, base, head string) error {
	token := os.Getenv("GITLAB_TOKEN")
	if token == "" {
		return fmt.Errorf("GITLAB_TOKEN not set")
	}
	project := url.PathEscape(repo.Owner + "/" + repo.Name)
	apiURL := fmt.Sprintf(
		"https://gitlab.com/api/v4/projects/%s/merge_requests",
		project,
	)
	body := map[string]interface{}{
		"title":                title,
		"description":          description,
		"source_branch":        head,
		"target_branch":        base,
		"remove_source_branch": true,
	}
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("PRIVATE-TOKEN", token)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("gitlab API returned %s", resp.Status)
	}
	fmt.Println("GitLab MR created")
	return nil
}
