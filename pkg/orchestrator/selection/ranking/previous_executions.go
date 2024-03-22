package ranking

import (
	"context"
	"time"

	"github.com/bacalhau-project/bacalhau/pkg/jobstore"
	"github.com/bacalhau-project/bacalhau/pkg/models"
	"github.com/bacalhau-project/bacalhau/pkg/orchestrator"
	"github.com/rs/zerolog/log"
)

type PreviousExecutionsNodeRankerParams struct {
	JobStore jobstore.Store
}
type PreviousExecutionsNodeRanker struct {
	jobStore jobstore.Store
}

func NewPreviousExecutionsNodeRanker(params PreviousExecutionsNodeRankerParams) *PreviousExecutionsNodeRanker {
	return &PreviousExecutionsNodeRanker{
		jobStore: params.JobStore,
	}
}

// RankNodes ranks nodes based on whether the node has already executed the job, which is useful when ranking
// nodes when handling retries:
// - Rank 30: Node has never executed the job.
// - Rank 0: Node has already executed the job, or rejected a bid.
// - Rank -1: Node is still executing the job, or has rejected it within the retryDelay interval
func (s *PreviousExecutionsNodeRanker) RankNodes(ctx context.Context,
	job models.Job, retryDelay time.Duration, nodes []models.NodeInfo) ([]orchestrator.NodeRank, error) {
	ranks := make([]orchestrator.NodeRank, len(nodes)) // Rank of each node, indexes corresponding to those in the nodes array
	previousExecutors := make(map[string]int)          // Map from node ID to number of previous active or completed executions
	toFilterOut := make(map[string]bool)
	executions, err := s.jobStore.GetExecutions(ctx, jobstore.GetExecutionsOptions{
		JobID: job.ID,
	})
	now := time.Now()
	if err == nil {
		for _, execution := range executions {
			if _, ok := previousExecutors[execution.NodeID]; !ok {
				previousExecutors[execution.NodeID] = 0
			}
			previousExecutors[execution.NodeID]++
			if !execution.IsTerminalComputeState() {
				toFilterOut[execution.NodeID] = true
			}
			if execution.IsRejected() &&
				now.Sub(execution.GetModifyTime()) < retryDelay {
				toFilterOut[execution.NodeID] = true
			}
		}
	}
	for i, node := range nodes {
		rank := 3 * orchestrator.RankPreferred
		reason := "job not executed yet"
		if _, ok := previousExecutors[node.ID()]; ok {
			if _, filterOut := toFilterOut[node.ID()]; filterOut {
				rank = orchestrator.RankUnsuitable
				reason = "job still running, or rejected within retry interval, on this node"
			} else {
				// This will include cases where the execution was
				// ExecutionStateAskForBidRejected; this might be a transient error
				// (the node lacked capacity at the time) so we can try again. When
				// we have a way to distinguish transient from permanent errors,
				// this logic should be revised to rejected nodes that *permanently*
				// refused the job and retry nodes that *temporarily* refused it
				// (while still preferring nodes that haven't tried the job before,
				// of course)
				rank = orchestrator.RankPossible
				reason = "job previously attempted on this node"
			}
		}
		ranks[i] = orchestrator.NodeRank{
			NodeInfo: node,
			Rank:     rank,
			Reason:   reason,
		}
		log.Ctx(ctx).Trace().Object("Rank", ranks[i]).Msg("Ranked node")
	}
	return ranks, nil
}
