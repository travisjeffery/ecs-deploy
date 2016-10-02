package main

import (
	"fmt"
	"log"
	"os"

	"github.com/travisjeffery/ecs-deploy/client"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	service = kingpin.Flag("service", "Name of Service to update.").Short("s").Required().String()
	image   = kingpin.Flag("image", "Name of Docker image to run.").Short("i").String()
	tag     = kingpin.Flag("tag", "Tag of Docker image to run.").Short("T").String()
	cluster = kingpin.Flag("cluster", "Name of ECS cluster.").Short("C").Default("default").String()
	task    = kingpin.Flag("task", "Name of task definition. Defaults to service name").Short("t").String()
	region  = kingpin.Flag("region", "Name of AWS region.").Short("r").Default("us-east-1").OverrideDefaultFromEnvar("AWS_DEFAULT_REGION").String()
	count   = kingpin.Flag("count", "Desired count of instantiations to place and run in service. Defaults to existing running count.").Short("c").Default("-1").Int64()
	nowait  = kingpin.Flag("nowait", "Disable waiting for all task definitions to start running").Short("n").Bool()
)

func main() {
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version("1.0.4").Author("Travis Jeffery")
	kingpin.CommandLine.Help = "Update ECS service."
	kingpin.Parse()

	if *task == "" {
		task = service
	}

	prefix := fmt.Sprintf("%s/%s ", *cluster, *service)
	logger := log.New(os.Stderr, prefix, log.LstdFlags)
	c := client.New(region, logger)

	arn := ""
	var err error

	if image != nil {
		arn, err = c.RegisterTaskDefinition(task, image, tag)
		if err != nil {
			logger.Printf("[error] register task definition: %s\n", err)
			return
		}
	}

	err = c.UpdateService(cluster, service, count, &arn)
	if err != nil {
		logger.Printf("[error] update service: %s\n", err)
		return
	}

	if *nowait == false {
		err := c.Wait(cluster, service, &arn)
		if err != nil {
			logger.Printf("[error] wait: %s\n", err)
			return
		}
	}

	logger.Printf("[info] update service success")
}
