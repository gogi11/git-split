package git

func ApplyPathFromBranch(targetBranch, path string) error {

	_, err := runGit(
		"checkout",
		targetBranch,
		"--",
		path,
	)

	return err
}
