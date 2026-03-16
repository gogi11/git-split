package filechanges

import (
	"fmt"
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

func ConvertFileWithStatusLinesToFileChange(actions []string, paths [][]string) ([]FileChange, error) {
	var result []FileChange
	if len(actions) != len(paths) {
		return result, fmt.Errorf("actions and paths length mismatch")
	}
	for i := range actions {
		var action FileChangeAction
		switch actions[i] {
		case "A":
			action = ADDED
		case "M":
			action = MODIFIED
		case "D":
			action = DELETED
		case "R":
			action = RENAMED
		default:
			return nil, fmt.Errorf("unknown action: %s", actions[i])
		}
		switch action {
		case MODIFIED, ADDED, DELETED:
			result = append(result, FileChange{
				Path:   paths[i][0],
				Action: action,
			})
		case RENAMED:
			result = append(result, FileChange{
				OldPath: paths[i][0],
				Path:    paths[i][1],
				Action:  action,
			})
		}
	}
	return result, nil
}
