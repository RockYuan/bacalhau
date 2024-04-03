package requester

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/bacalhau-project/bacalhau/pkg/jobstore"
	"github.com/bacalhau-project/bacalhau/pkg/models"
)

type HousekeepingParams struct {
	Endpoint Endpoint
	JobStore jobstore.Store
	NodeID   string
	Interval time.Duration
}

type Housekeeping struct {
	endpoint Endpoint
	jobStore jobstore.Store
	nodeID   string
	interval time.Duration

	stopChannel chan struct{}
	stopOnce    sync.Once
}

func NewHousekeeping(params HousekeepingParams) *Housekeeping {
	h := &Housekeeping{
		endpoint:    params.Endpoint,
		jobStore:    params.JobStore,
		nodeID:      params.NodeID,
		interval:    params.Interval,
		stopChannel: make(chan struct{}),
	}

	return h
}

func (h *Housekeeping) Start(ctx context.Context) {
	go h.housekeepingBackgroundTask(ctx)
}

func (h *Housekeeping) housekeepingBackgroundTask(ctx context.Context) {
	ticker := time.NewTicker(h.interval)
	for {
		select {
		case <-ctx.Done():
			// TODO(forrest): [correctness] this needs better start/stop logic
			h.Stop()
		case <-ticker.C:
			jobs, err := h.jobStore.GetInProgressJobs(ctx)
			if err != nil {
				log.Ctx(ctx).Err(err).Msg("failed to get in progress jobs")
				continue
			}
			now := time.Now()
			for _, job := range jobs {
				// in case the job store is shared between multiple nodes, we only want to clean up jobs that are owned by this node
				requesterID, ok := job.Meta[models.MetaRequesterID]
				if !ok {
					log.Ctx(ctx).Warn().Msgf("job %s has no requester ID. Skipping", job.ID)
					continue
				}
				if requesterID != h.nodeID {
					continue
				}

				// timeout is only applicable to batch and ops jobs
				if job.Type != models.JobTypeBatch && job.Type != models.JobTypeOps {
					continue
				}

				// cancel jobs that have been in progress beyond the timeout period
				if now.Sub(job.GetCreateTime()) > job.Task().Timeouts.GetExecutionTimeout() {
					log.Ctx(ctx).Info().Msgf("job %s timed out. Canceling", job.ID)
					go func(jobID string) {
						_, innerErr := h.endpoint.CancelJob(ctx, CancelJobRequest{
							JobID:  jobID,
							Reason: "timed out",
						})
						if innerErr != nil {
							log.Ctx(ctx).Err(innerErr).Msgf("failed to cancel job %s", jobID)
						}
					}(job.ID)
				}
			}
		case <-h.stopChannel:
			log.Ctx(ctx).Debug().Msg("stopped housekeeping task")
			ticker.Stop()
			return
		}
	}
}

func (h *Housekeeping) Stop() {
	h.stopOnce.Do(func() {
		h.stopChannel <- struct{}{}
	})
}
