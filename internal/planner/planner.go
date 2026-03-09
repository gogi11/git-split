package planner

import "git-split/internal/plan"

type Planner interface {
	Build() (plan.Plan, error)
}
