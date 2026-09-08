package queue_server_test

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/buildbeaver/buildbeaver/common/gerror"
	"github.com/buildbeaver/buildbeaver/common/models"
	"github.com/buildbeaver/buildbeaver/server/app/server_test"
	"github.com/buildbeaver/buildbeaver/server/dto/dto_test/referencedata"
)

// TestDequeueRace exercises many runners concurrently polling Dequeue for the same single eligible
// job, simulating the scenario described in the double-dequeue race: two (or more) runners polling
// at the same time must not both come away believing they've been given the same job.
//
// Before the fix (SELECT ... FOR UPDATE SKIP LOCKED in JobStore.FindQueuedJob), a runner that lost
// the race could still match the initial unlocked SELECT and do all the work of Dequeue, only to
// fail later with an optimistic lock error when its update to the job row's ETag conflicted with the
// winner's already-committed update. That's not silent double-execution (the ETag check does catch
// it eventually) but it is the wrong outcome: the loser should cleanly find no eligible job, not fail
// with an unrelated-looking error after doing all of Dequeue's other work. This test asserts exactly
// that: exactly one concurrent caller gets the job, and every other caller gets a plain "not found",
// not any other kind of error.
func TestDequeueRace(t *testing.T) {
	app, cleanup, err := server_test.New(server_test.TestConfig(t))
	require.NoError(t, err)
	defer cleanup()
	ctx := context.Background()

	legalEntity, _ := server_test.CreatePersonLegalEntity(t, ctx, app, "", "", "")
	repo := server_test.CreateRepo(t, ctx, app, legalEntity.ID)

	const numRunners = 10
	runnerIDs := make([]models.RunnerID, numRunners)
	for i := 0; i < numRunners; i++ {
		runner := server_test.CreateRunner(t, ctx, app, models.ResourceName(fmt.Sprintf("race-runner-%d", i)), legalEntity.ID, nil)
		runnerIDs[i] = runner.ID
	}

	// CreateAndQueueBuild queues the standard 4-job reference build. ReferenceJob2-4 all depend on
	// ReferenceJob1, so it's the only job immediately eligible for dequeue - giving exactly one
	// contended job for all the runners below to race over.
	build := server_test.CreateAndQueueBuild(t, ctx, app, repo.ID, legalEntity.ID, "")
	var eligibleJobID models.JobID
	for _, job := range build.Jobs {
		if job.Name == referencedata.ReferenceJob1.Name {
			eligibleJobID = job.ID
		}
	}
	require.NotEqual(t, models.JobID{}, eligibleJobID, "expected to find ReferenceJob1 in the queued build")

	var (
		wg          sync.WaitGroup
		mu          sync.Mutex
		dequeuedIDs []models.JobID
		otherErrs   []error
	)
	for _, runnerID := range runnerIDs {
		runnerID := runnerID
		wg.Add(1)
		go func() {
			defer wg.Done()
			job, err := app.QueueService.Dequeue(ctx, runnerID)
			mu.Lock()
			defer mu.Unlock()
			if err == nil {
				require.NotNil(t, job)
				dequeuedIDs = append(dequeuedIDs, job.ID)
			} else if gerror.ToNotFound(err) == nil {
				otherErrs = append(otherErrs, err)
			}
		}()
	}
	wg.Wait()

	require.Empty(t, otherErrs, "every runner that didn't win the race should see a plain not-found error, not: %v", otherErrs)
	require.Len(t, dequeuedIDs, 1, "exactly one runner should have successfully dequeued the contended job")
	require.Equal(t, eligibleJobID, dequeuedIDs[0])
}
