package models

import (
	"github.com/buildbeaver/buildbeaver/common/gerror"
)

type JobSearch struct {
	Pagination
	// RunnerID is the Runner associated with Jobs being searched for, or nil to include jobs for any Runner.
	RunnerID      *RunnerID `json:"runner_id"`
	IncludeBuild  bool      `json:"include_build"`
	IncludeRepo   bool      `json:"include_repo"`
	IncludeCommit bool      `json:"include_commit"`
}

func NewJobSearch() *JobSearch {
	return &JobSearch{Pagination: Pagination{}}
}

// NewJobSearchForRunner returns search criteria to search for jobs for a particular runner.
// Other search criteria can be specified to narrow the search.
func NewJobSearchForRunner(runnerID RunnerID, limit int) *JobSearch {
	return &JobSearch{
		Pagination: Pagination{
			Limit:  limit,
			Cursor: nil,
		},
		RunnerID:      &runnerID,
		IncludeRepo:   true,
		IncludeBuild:  true,
		IncludeCommit: true,
	}
}

func (j *JobSearch) Validate() error {
	// TODO: Do we need to validate anything at the moment?
	if j.RunnerID == nil {
		return gerror.NewErrValidationFailed("RunnerID must be specified when a ref is specified")
	}

	return nil
}
