package workflow_entities

import (
	"github.com/stretchr/testify/mock"
	"go.temporal.io/sdk/workflow"
	"sync/atomic"
)

func (s *UnitTestSuite) Test_BatchProducer_Normal() {
	env := s.NewTestWorkflowEnvironment()

	env.RegisterWorkflow(BatchProducerWorkflow)
	env.RegisterWorkflow(MessageProducerWorkflow)
	value := int32(0)
	env.OnWorkflow(MessageProducerWorkflow, mock.Anything, mock.Anything).Return(
		func(ctx workflow.Context, param MessageProducerParam) (*MessageProducerResult, error) {
			atomic.AddInt32(&value, 1)
			return nil, nil
		})

	env.ExecuteWorkflow(BatchProducerWorkflow, SetupParams{
		ProducerScheduleID: "local_test_schedule",
		ProducerNumber:     100,
	})

	s.True(env.IsWorkflowCompleted())
	s.NoError(env.GetWorkflowError())
	s.Equal(100, int(value))
}
