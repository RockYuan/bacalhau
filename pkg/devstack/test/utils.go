package test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/filecoin-project/bacalhau/pkg/devstack"
	"github.com/filecoin-project/bacalhau/pkg/executor"
	_ "github.com/filecoin-project/bacalhau/pkg/logger"
	"github.com/filecoin-project/bacalhau/pkg/system"
	"github.com/stretchr/testify/assert"

	"github.com/rs/zerolog/log"
)

// run the job on 2 nodes
const TEST_CONCURRENCY = 1

// both nodes must agree on the result
const TEST_CONFIDENCE = 1

// the results must be within 10% of each other
const TEST_TOLERANCE = 0.1

func setupTest(
	t *testing.T,
	nodes int,
	badActors int,
	executors map[string]executor.Executor,
) (*devstack.DevStack, context.CancelFunc) {
	ctx, cancelFunction := system.GetCancelContext()
	stack, err := devstack.NewDevStack(
		ctx,
		nodes,
		badActors,
		executors,
	)
	assert.NoError(t, err)
	if err != nil {
		log.Fatal().Msg(fmt.Sprintf("Unable to create devstack: %s", err))
	}
	// we need a better method for this - i.e. waiting for all the ipfs nodes to be ready
	time.Sleep(time.Second * 2)
	return stack, cancelFunction
}

// this might be called multiple times if KEEP_STACK is active
// the first time - once the test has completed, this function will be called
// it will reset the KEEP_STACK variable so the user can ctrl+c the running stack
func teardownTest(stack *devstack.DevStack, cancelFunction context.CancelFunc) {
	if os.Getenv("KEEP_STACK") == "" {
		cancelFunction()
		// need some time to let ipfs processes shut down
		time.Sleep(time.Second * 1)
	} else {
		stack.PrintNodeInfo()
		os.Setenv("KEEP_STACK", "")
		select {}
	}
}
