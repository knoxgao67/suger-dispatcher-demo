package workflow_entities

import (
	"github.com/bytedance/mockey"
	"github.com/knoxgao67/temporal-dispatcher-demo/dal/kafka"
	"go.temporal.io/sdk/client"
	"strings"
	"time"
)

func (s *UnitTestSuite) Test_Producer() {
	mockey.PatchConvey("Normal", s.T(), func() {
		mockey.Mock(mockey.GetMethod(kafka.Producer{}, "SendMessage")).Return(nil).Build()

		env := s.NewTestWorkflowEnvironment()
		env.SetStartWorkflowOptions(client.StartWorkflowOptions{
			WorkflowExecutionTimeout: time.Minute,
		})
		env.RegisterWorkflow(MessageProducerWorkflow)
		env.RegisterActivity(SendMessage)

		env.ExecuteWorkflow(MessageProducerWorkflow, MessageProducerParam{})
		s.True(env.IsWorkflowCompleted())

		s.NoError(env.GetWorkflowError())
		var result MessageProducerResult
		s.NoError(env.GetWorkflowResult(&result))
		s.True(strings.Contains(result.Message, "Completed#"))
	})

}
