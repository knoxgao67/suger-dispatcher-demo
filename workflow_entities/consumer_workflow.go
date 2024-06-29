package workflow_entities

import (
	"context"
	"time"

	"github.com/knoxgao67/temporal-dispatcher-demo/dal/kafka"
	"go.temporal.io/sdk/workflow"
)

type MessageConsumerParam struct {
	Message string
}
type MessageConsumerResult struct {
	Message string
}

func MessageConsumerWorkflow(ctx workflow.Context, param MessageConsumerParam) (*MessageConsumerResult, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	readMsg := workflow.SideEffect(ctx, func(ctx workflow.Context) interface{} {
		return kafka.GetConsumer().ReadMessage(time.Second * 10)
	})

	var msg string
	err := readMsg.Get(&msg)
	if err != nil {
		return nil, err
	}

	param.Message = msg

	result := &MessageConsumerResult{}
	err = workflow.ExecuteActivity(ctx, HandleMessage, param).Get(ctx, &result)
	return result, err
}

// HandleMessage Receive one msg and exit.
func HandleMessage(ctx context.Context, param MessageConsumerParam) (*MessageConsumerResult, error) {
	// handle message
	result := &MessageConsumerResult{Message: "Handled#" + param.Message}
	time.Sleep(time.Second * 5)
	return result, nil
}
