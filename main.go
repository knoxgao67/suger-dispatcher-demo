package main

import (
	"context"
	"flag"
	"log"
	"sync"

	"github.com/knoxgao67/temporal-dispatcher-demo/common"
	"github.com/knoxgao67/temporal-dispatcher-demo/config"
	"github.com/knoxgao67/temporal-dispatcher-demo/dal/kafka"
	"github.com/knoxgao67/temporal-dispatcher-demo/dal/temporal_client"
	"github.com/knoxgao67/temporal-dispatcher-demo/dal/zk"
	"github.com/knoxgao67/temporal-dispatcher-demo/workflow_entities"
	"github.com/pborman/uuid"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.uber.org/multierr"
)

var cleanFn []func()
var mutex sync.Mutex

var cfgFile = flag.String("f", "./config/config.yaml", "config filename")

func main() {
	flag.Parse()
	config.Init(*cfgFile)

	zk.Init()
	kafka.Init()
	temporal_client.Init()
	defer finish()

	addCleanFn(func() {
		if err := multierr.Append(kafka.GetProducer().Close(), kafka.GetConsumer().Close()); err != nil {
			log.Println("unable to close kafka client", err)
		}
	})
	addCleanFn(func() {
		zk.Close()
		temporal_client.Close()
	})

	cfg := config.GetConfig()

	go startSchedule(context.Background(), workflow_entities.SetupParams{
		ProducerScheduleID: "ProducerSchedule",
		ConsumerScheduleID: "ConsumerSchedule",
		ProducerNumber:     cfg.ProducerNumber,
		ConsumerNumber:     cfg.ConsumerNumber,
	})
	startWorker(context.Background())
}

func startSchedule(ctx context.Context, param workflow_entities.SetupParams) {
	workflowID := "schedule_workflow_" + uuid.New()
	// set global lock by zk, only one master can start the schedule
	unlock := zk.Lock(common.GlobalLockPath)
	addCleanFn(unlock)

	c := temporal_client.GetClient()
	start := func(cronExpression string, scheduleID string, workflowEntity interface{}) {
		scheduleHandler := c.ScheduleClient().GetHandle(ctx, scheduleID)
		_ = scheduleHandler.Delete(ctx)
		_, err := c.ScheduleClient().Create(ctx, client.ScheduleOptions{
			ID: scheduleID,

			Spec: client.ScheduleSpec{
				CronExpressions: []string{cronExpression},
			},
			Overlap: enums.SCHEDULE_OVERLAP_POLICY_TERMINATE_OTHER, // terminate old workflow if exist
			Action: &client.ScheduleWorkflowAction{
				ID:        workflowID,
				Workflow:  workflowEntity,
				Args:      []interface{}{param},
				TaskQueue: common.DispatcherQueue, // TODO 单独设置queue，便于更加灵活的控制？
			},
		})
		if err != nil {
			log.Fatalln("Unable to create schedule for module", err)
		}
	}

	start("*/1 * * * *", param.ProducerScheduleID, workflow_entities.BatchProducerWorkflow)
	start("*/2 * * * *", param.ConsumerScheduleID, workflow_entities.BatchConsumerWorkflow)
}

func startWorker(ctx context.Context) {
	w := worker.New(temporal_client.GetClient(), common.DispatcherQueue, worker.Options{})

	// This worker hosts both Workflow and Activity functions.
	w.RegisterWorkflow(workflow_entities.MessageProducerWorkflow)
	w.RegisterWorkflow(workflow_entities.MessageConsumerWorkflow)
	w.RegisterWorkflow(workflow_entities.BatchProducerWorkflow)
	w.RegisterWorkflow(workflow_entities.BatchConsumerWorkflow)

	w.RegisterActivity(workflow_entities.SendMessage)
	w.RegisterActivity(workflow_entities.HandleMessage)

	// Start listening to the Task Queue.
	err := w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}

func finish() {
	for _, fn := range cleanFn {
		fn()
	}
}

func addCleanFn(fn func()) {
	mutex.Lock()
	defer mutex.Unlock()
	cleanFn = append(cleanFn, fn)
}
