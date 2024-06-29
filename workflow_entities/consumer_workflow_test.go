package workflow_entities

import (
	"github.com/bytedance/mockey"
	"github.com/knoxgao67/temporal-dispatcher-demo/dal/kafka"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/testsuite"
	"testing"
	"time"
)

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}

func (s *UnitTestSuite) SetupTest() {
}

func (s *UnitTestSuite) Test_MessageConsumer() {
	mockey.PatchConvey("Normal", s.T(), func() {
		uuidVal := uuid.New()
		mockey.Mock(mockey.GetMethod(kafka.Consumer{}, "ReadMessage")).Return(uuidVal).Build()

		env := s.NewTestWorkflowEnvironment()
		env.SetStartWorkflowOptions(client.StartWorkflowOptions{
			WorkflowExecutionTimeout: time.Minute,
		})
		env.RegisterWorkflow(MessageConsumerWorkflow)
		env.RegisterActivity(HandleMessage)

		env.ExecuteWorkflow(MessageConsumerWorkflow, MessageConsumerParam{})
		s.True(env.IsWorkflowCompleted())

		s.NoError(env.GetWorkflowError())
		var result MessageConsumerResult
		s.NoError(env.GetWorkflowResult(&result))
		s.Equal("Handled#"+uuidVal, result.Message)
	})
}
