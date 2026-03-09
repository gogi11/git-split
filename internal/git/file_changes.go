package git

import (
	"fmt"
	"strings"
)

type FileChangeAction string

const (
	ADDED    FileChangeAction = "A"
	MODIFIED FileChangeAction = "M"
	DELETED  FileChangeAction = "D"
	RENAMED  FileChangeAction = "R"
)

type FileChange struct {
	Path    string // destination path
	OldPath string // optional, for rename
	Action  FileChangeAction
}

func GetChangedFilesWithStatus(base, target string) ([]FileChange, error) {
	out, err := runGit("diff", "--name-status", base+".."+target)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(out, "\n")
	changes, err := convertFileWithStatusLinesToFileChange(lines)
	if err != nil {
		return nil, err
	}
	return changes, nil
}

func convertFileWithStatusLinesToFileChange(lines []string) ([]FileChange, error) {
	var result []FileChange
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		detected_action := parts[0]
		var action FileChangeAction
		switch detected_action {
		case "A":
			action = ADDED
		case "M":
			action = MODIFIED
		case "D":
			action = DELETED
		case "R":
			action = RENAMED
		default:
			return nil, fmt.Errorf("unknown action: %s", detected_action)
		}
		switch action {
		case MODIFIED, ADDED, DELETED:
			result = append(result, FileChange{
				Path:   parts[1],
				Action: action,
			})
		case RENAMED:
			result = append(result, FileChange{
				Path:    parts[2],
				OldPath: parts[1],
				Action:  action,
			})
		}
	}
	return result, nil
}
