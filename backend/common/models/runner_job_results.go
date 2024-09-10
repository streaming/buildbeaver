package models

// RunnerJobResult represents the content for a search of jobs for a given Runner.
type RunnerJobResult struct {
	// Note that a RunnerJobResult instance is passed to ResourceTable.ListIn() and this function does
	// not support composing/nesting the primary resource table, so we must embed Job here.
	// (See comment inside ResourceTable.ListIn() for details).
	*Job
	Build  *Build  `db:"builds"`
	Commit *Commit `db:"commits"`
	Repo   *Repo   `db:"repos"`
}
