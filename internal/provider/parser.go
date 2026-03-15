package provider

import (
	"fmt"
	"git-split/internal/git"
	"net/url"
	"strings"
)

type Provider string

const (
	Unknown   Provider = "unknown"
	GitHub    Provider = "github"
	GitLab    Provider = "gitlab"
	GitBucket Provider = "gitbucket"
)

type Repo struct {
	Host      string
	Owner     string
	Name      string
	Namespace string
	APIBase   string
	Provider  Provider
}

func BuildRepo(remote string) (Repo, error) {
	remoteUrl, err := git.GetRemoteURL(remote)
	if err != nil {
		return Repo{}, err
	}
	fmt.Printf("Remote URL: %s\n", remoteUrl)
	repoInfo, err := parseRemote(remoteUrl)
	if err != nil {
		return Repo{}, err
	}
	return repoInfo, nil
}

func parseRemote(remoteUrl string) (Repo, error) {
	if strings.HasPrefix(remoteUrl, "git@") {
		return parseSSH(remoteUrl)
	}
	if strings.HasPrefix(remoteUrl, "http") {
		return parseHTTP(remoteUrl)
	}
	return Repo{}, fmt.Errorf("unsupported remote format: %s", remoteUrl)
}

func parseHTTP(remoteUrl string) (Repo, error) {
	u, err := url.Parse(remoteUrl)
	if err != nil {
		return Repo{}, err
	}
	p := strings.TrimSuffix(u.Path, ".git")
	p = strings.TrimPrefix(p, "/")
	parts := strings.Split(p, "/")
	if len(parts) < 2 {
		return Repo{}, fmt.Errorf("invalid remote path")
	}
	repo := parts[len(parts)-1]
	namespace := strings.Join(parts[:len(parts)-1], "/")
	owner := parts[0]
	host := u.Host
	return Repo{
		Host:      host,
		Owner:     owner,
		Name:      repo,
		Namespace: namespace,
		APIBase:   detectAPI(host),
		Provider:  detectProvider(host),
	}, nil
}

func parseSSH(remoteUrl string) (Repo, error) {
	// git@github.com:org/repo.git
	parts := strings.Split(remoteUrl, ":")
	if len(parts) != 2 {
		return Repo{}, fmt.Errorf("invalid ssh remote")
	}
	hostPart := strings.Split(parts[0], "@")
	host := hostPart[1]
	p := strings.TrimSuffix(parts[1], ".git")
	pathParts := strings.Split(p, "/")
	if len(pathParts) < 2 {
		return Repo{}, fmt.Errorf("invalid repo path")
	}
	repo := pathParts[len(pathParts)-1]
	namespace := strings.Join(pathParts[:len(pathParts)-1], "/")
	owner := pathParts[0]
	return Repo{
		Host:      host,
		Owner:     owner,
		Name:      repo,
		Namespace: namespace,
		APIBase:   detectAPI(host),
		Provider:  detectProvider(host),
	}, nil
}

func detectAPI(host string) string {
	switch {
	case strings.Contains(host, "github"):
		return "https://api.github.com"

	case strings.Contains(host, "gitlab"):
		return "https://" + host + "/api/v4"

	case strings.Contains(host, "gitbucket"):
		return "https://" + host + "/api/v3"

	default:
		return "https://" + host
	}
}

func detectProvider(host string) Provider {
	switch {
	case strings.Contains(host, "github"):
		return GitHub
	case strings.Contains(host, "gitlab"):
		return GitLab
	case strings.Contains(host, "gitbucket"):
		return GitBucket
	default:
		return Unknown
	}
}
