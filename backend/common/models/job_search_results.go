package models

// JobSearchResult represents the content for a search of jobs.
type JobSearchResult struct {
	// Note that a JobSearchResult instance is passed to ResourceTable.ListIn() and this function does
	// not support composing/nesting the primary resource table, so we must embed Job here.
	// (See comment inside ResourceTable.ListIn() for details).
	*Job
	Build  *Build  `db:"builds"`
	Commit *Commit `db:"commits"`
	Repo   *Repo   `db:"repos"`
}
