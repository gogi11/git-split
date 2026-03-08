package plan

type BranchPlan struct {
	Branch   string
	Base     string
	Commits  []string
	Push     bool
	CreateMR bool
	MRTitle  string
}

type Plan struct {
	Remote   string
	Provider string
	Repo     string
	Branches []BranchPlan
}
