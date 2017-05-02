package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/opts"
	_"os"
	"os"
)

func main() {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	var managerClient *http.Client
	managerUrl := os.Getenv("MANAGER_URL")
	managerCli, err := client.NewClient(managerUrl,"", managerClient,nil)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	filterOpt := opts.NewFilterOpt()
	filterOpt.Set("event=stop")
	stopEvent, _ := cli.Events(ctx,types.EventsOptions{"", "", filterOpt.Value()})
	for {
		fmt.Println("Observing stop events...")
		fmt.Println(<-stopEvent)
		rescalePendingTasks(managerCli)
	}
}

func rescalePendingTasks(cli *client.Client) {
	fmt.Println("Collecting pending tasks...")
	pendingTasks := getPendingTasks(cli)
	incompleteServices := make(map[string][]swarm.Task)
	for _, task := range pendingTasks {
		incompleteServices[task.ServiceID] = append(incompleteServices[task.ServiceID], task)
	}
	for serviceID, tasks := range incompleteServices {
		ctx := context.Background()
		service, _, err := cli.ServiceInspectWithRaw(ctx, serviceID, types.ServiceInspectOptions{})
		if err != nil {
			panic(err)
		}
		fmt.Printf("Trying to reschedule service: %s \n", service.Spec.Name)
		*service.Spec.Mode.Replicated.Replicas -= uint64(len(tasks))
		_, err = cli.ServiceUpdate(ctx, service.ID, service.Version, service.Spec, types.ServiceUpdateOptions{})
		if err != nil {
			panic(err)
		}
		service, _, err = cli.ServiceInspectWithRaw(ctx, serviceID, types.ServiceInspectOptions{})
		if err != nil {
			panic(err)
		}
		*service.Spec.Mode.Replicated.Replicas += uint64(len(tasks))
		_, err = cli.ServiceUpdate(ctx, service.ID, service.Version, service.Spec, types.ServiceUpdateOptions{})
		if err != nil {
			panic(err)
		}
	}
}

func getPendingTasks(cli *client.Client) ([]swarm.Task) {
	tasks, err := cli.TaskList(context.Background(), types.TaskListOptions{})
	if err != nil {
		panic(err)
	}
	var pendingTasks []swarm.Task
	for _,task := range tasks {
		if task.Status.State == "pending" {
			pendingTasks = append(pendingTasks,task)
		}
	}
	return pendingTasks
}
