package workflow_entities

import (
	"fmt"

	"github.com/pborman/uuid"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/workflow"
)

type SetupParams struct {
	ProducerScheduleID, ConsumerScheduleID string
	ProducerNumber, ConsumerNumber         int
}

func BatchProducerWorkflow(ctx workflow.Context, param SetupParams) error {
	return baseMultiChildWorkflow(ctx, param.ProducerNumber, "batch_producer", MessageProducerWorkflow)
}
func BatchConsumerWorkflow(ctx workflow.Context, param SetupParams) error {
	return baseMultiChildWorkflow(ctx, param.ConsumerNumber, "batch_consumer", MessageConsumerWorkflow)
}

func baseMultiChildWorkflow(ctx workflow.Context, cnt int, name string, childWorkflow interface{}) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("multi child workflow started", "name", name)

	wg := workflow.NewWaitGroup(ctx)
	wg.Add(cnt)
	for i := 0; i < cnt; i++ {
		idx := i
		workflow.Go(ctx, func(ctx workflow.Context) {
			childWorkflowID := fmt.Sprintf("%s_child:%d:%s", name, idx, uuid.New())
			childCtx := workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
				WorkflowID:        childWorkflowID,
				ParentClosePolicy: enums.PARENT_CLOSE_POLICY_TERMINATE, // the child workflow will also terminate
			})
			err := workflow.ExecuteChildWorkflow(childCtx, childWorkflow).Get(ctx, nil)
			if err != nil {
				logger.Error("child workflow failed",
					"error", err,
					"workflowId", childWorkflowID)
			}
			wg.Done()
		})
	}
	wg.Wait(ctx)
	return nil
}
