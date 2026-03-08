package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

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
