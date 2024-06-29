package workflow_entities

import (
	"context"
	"time"

	"github.com/knoxgao67/temporal-dispatcher-demo/dal/kafka"
	"github.com/pborman/uuid"
	"go.temporal.io/sdk/workflow"
)

type MessageProducerParam struct {
	Message string
}

type MessageProducerResult struct {
	Message string
}

func MessageProducerWorkflow(ctx workflow.Context, param MessageProducerParam) (*MessageProducerResult, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	genMessage := workflow.SideEffect(ctx, func(ctx workflow.Context) interface{} {
		return uuid.New()
	})

	var message string
	err := genMessage.Get(&message)
	if err != nil {
		return nil, err
	}

	param.Message = message

	result := &MessageProducerResult{}
	err = workflow.ExecuteActivity(ctx, SendMessage, param).Get(ctx, result)
	return result, err
}

func SendMessage(ctx context.Context, param MessageProducerParam) (*MessageProducerResult, error) {
	err := kafka.GetProducer().SendMessage(ctx, param.Message)
	if err != nil {
		return nil, err
	}
	return &MessageProducerResult{
		Message: "Completed#" + param.Message,
	}, nil
}
