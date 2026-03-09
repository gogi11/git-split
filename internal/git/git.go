package git

import (
	"bytes"
	"fmt"
	"git-split/helpers"
	"os/exec"
	"strings"
)

func Run(args ...string) error {
	_, err := runGit(args...)
	return err
}

func runGit(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("git %v failed: %v\n%s",
			args,
			err,
			stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func GetCommitsBetween(base_branch, target_branch string) ([]string, error) {
	out, err := runGit("rev-list", "--reverse", fmt.Sprintf("%s..%s", base_branch, target_branch))
	if err != nil {
		return nil, err
	}
	if out == "" {
		return []string{}, nil
	}
	return strings.Split(out, "\n"), nil
}

func CreateBranch(base, newBranch string) error {
	_, err := runGit("checkout", base)
	if err != nil {
		return err
	}
	_, err = runGit("checkout", "-b", newBranch)
	return err
}

func CherryPickCommits(commits []string) error {
	for _, c := range commits {
		fmt.Printf("Cherry-picking %s\n", c)
		_, err := runGit("cherry-pick", c)
		if err != nil {
			return err
		}
	}

	return nil
}

func ChunkCommits(commits []string, size int) [][]string {
	var chunks [][]string
	for size < len(commits) {
		commits, chunks = commits[size:], append(chunks, commits[0:size:size])
	}
	chunks = append(chunks, commits)
	return chunks
}

func Push(remote, branch string) error {
	_, err := runGit("push", "-u", remote, branch)
	return err
}

func GetRemoteURL() (string, error) {
	remote, err := SelectRemote()
	if err != nil {
		return "", err
	}
	out, err := runGit("remote", "get-url", remote)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func GetRemotes() ([]string, error) {

	out, err := runGit("remote")
	if err != nil {
		return nil, err
	}

	if out == "" {
		return []string{}, nil
	}

	return strings.Split(strings.TrimSpace(out), "\n"), nil
}

func GetChangedFiles(base, target string) ([]string, error) {
	// DEPRECATED, WE USE GetChangedFilesWithStatus INSTEAD
	out, err := runGit(
		"diff",
		"--name-only",
		base+".."+target,
	)

	if err != nil {
		return nil, err
	}

	if out == "" {
		return []string{}, nil
	}

	return strings.Split(out, "\n"), nil
}
func GetCommitsForFiles(base, target string, files []string) ([]string, error) {
	args := []string{
		"log",
		"--pretty=format:%H",
		base + ".." + target,
	}
	args = append(args, "--")
	args = append(args, files...)
	out, err := runGit(args...)
	if err != nil {
		return nil, err
	}

	if out == "" {
		return []string{}, nil
	}

	lines := strings.Split(out, "\n")
	return helpers.Unique(lines), nil
}

func Commit(message string) error {
	_, err := runGit("add", ".")
	if err != nil {
		return err
	}
	_, err = runGit("commit", "-m", message)
	return err
}

func Checkout(branch string) error {
	_, err := runGit("checkout", branch)
	return err
}

func ApplyPathFromBranch(targetBranch, path string) error {
	_, err := runGit(
		"checkout",
		targetBranch,
		"--",
		path,
	)
	return err
}
