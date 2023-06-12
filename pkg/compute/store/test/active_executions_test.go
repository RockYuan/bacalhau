//go:build unit || !integration

package test

import (
	"context"
	"testing"

	"github.com/bacalhau-project/bacalhau/pkg/compute/store"
	"github.com/bacalhau-project/bacalhau/pkg/compute/store/persistent"
	"github.com/bacalhau-project/bacalhau/pkg/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	executionStore store.ExecutionStore
	execution      store.Execution
	ctx            context.Context
}

func (s *Suite) SetupTest() {
	s.ctx = context.Background()
	s.executionStore, _ = persistent.NewStore(s.ctx, "")
	s.execution = newExecution()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestGetActiveExecution_Single() {
	ctx := context.Background()
	err := s.executionStore.CreateExecution(ctx, s.execution)
	s.NoError(err)

	active, err := store.GetActiveExecution(ctx, s.executionStore, s.execution.Job.ID())
	s.NoError(err)
	s.Equal(s.execution, active)
}

func (s *Suite) TestGetActiveExecution_Multiple() {
	ctx := context.Background()

	// create a newer execution with same job as the previous one
	newerExecution := s.execution
	newerExecution.ID = uuid.NewString()
	newerExecution.Job = s.execution.Job
	newerExecution.UpdateTime = s.execution.UpdateTime.Add(1)

	err := s.executionStore.CreateExecution(ctx, s.execution)
	s.NoError(err)

	err = s.executionStore.CreateExecution(ctx, newerExecution)
	s.NoError(err)

	active, err := store.GetActiveExecution(ctx, s.executionStore, s.execution.Job.ID())
	s.NoError(err)
	s.Equal(newerExecution, active)
}

func (s *Suite) TestGetActiveExecution_DoestExist() {
	_, err := store.GetActiveExecution(context.Background(), s.executionStore, s.execution.Job.ID())
	s.ErrorAs(err, &store.ErrExecutionsNotFoundForJob{})
}

func newExecution() store.Execution {
	return *store.NewExecution(
		uuid.NewString(),
		model.Job{
			Metadata: model.Metadata{
				ID: uuid.NewString(),
			},
		},
		"nodeID-1",
		model.ResourceUsageData{
			CPU:    1,
			Memory: 2,
		})
}
