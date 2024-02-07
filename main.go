package main

import (
	"context"
	"log"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

// HelloWorldWorkflow is a workflow definition.
// It should have at least one argument - a context, and return an error.
// The workflow can also return one (or two, if error is returned) additional values.
func HelloWorldWorkflow(ctx workflow.Context, name string) (string, error) {
	// Create options for the activity.
	ao := workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	name += "hey "

	var result string
	err := workflow.ExecuteActivity(ctx, HelloWorldActivity, name).Get(ctx, &result)
	if err != nil {
		return "", err
	}
	return result, nil
}

// HelloWorldActivity is an activity definition.
func HelloWorldActivity(ctx context.Context, name string) (string, error) {
	return "Hello, " + name + "!", nil
}

func main() {
	c, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("Error creating Temporal client", err)
	}
	defer c.Close()

	// This worker hosts both Workflow and Activity functions
	w := worker.New(c, "hello-world-task-queue", worker.Options{})

	w.RegisterWorkflow(HelloWorldWorkflow)
	w.RegisterActivity(HelloWorldActivity)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
