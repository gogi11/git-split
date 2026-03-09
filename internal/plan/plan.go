package plan

import "git-split/internal/git"

type OperationType string

const (
	OpCherryPick OperationType = "cherry-pick"
	OpApplyPath  OperationType = "apply-path"
)

type Operation struct {
	Type        OperationType
	Commits     []string
	FileChanges []git.FileChange
	FromRef     string
}

type BranchPlan struct {
	Branch        string
	Base          string
	Operations    []Operation
	Push          bool
	CreateMR      bool
	MRTitle       string
	MRDescription string
}

type Plan struct {
	Remote   string
	Provider string
	Repo     string
	Branches []BranchPlan
}
