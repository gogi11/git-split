package planner

import (
	"git-split/internal/plan"
	"git-split/internal/provider"
)

type Planner interface {
	// Build returns a plan.Plan based on the implementation's logic. The remote parameter is provided for planners that need to know the remote name for branch conflict checks or MR creation.
	Build(remote string) (plan.Plan, error)
}

func InitializePlan(remote string) (plan.Plan, error) {
	repo, err := provider.BuildRepo(remote)
	if err != nil {
		return plan.Plan{}, err
	}
	return plan.Plan{Remote: remote, Repo: repo}, nil
}
