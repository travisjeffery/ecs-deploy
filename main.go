package main

import (
	"fmt"
	"log"
	"os"

	"github.com/travisjeffery/ecs-deploy/client"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	service = kingpin.Flag("service", "Name of Service to update.").Required().String()
	image   = kingpin.Flag("image", "Name of Docker image to run.").Strings()
	tag     = kingpin.Flag("tag", "Tag of Docker image to run.").Strings()
	cluster = kingpin.Flag("cluster", "Name of ECS cluster.").Default("default").String()
	task    = kingpin.Flag("task", "Name of task definition. Defaults to service name").String()
	region  = kingpin.Flag("region", "Name of AWS region.").Default("us-east-1").OverrideDefaultFromEnvar("AWS_DEFAULT_REGION").String()
	count   = kingpin.Flag("count", "Desired count of instantiations to place and run in service. Defaults to existing running count.").Default("-1").Int64()
	nowait  = kingpin.Flag("nowait", "Disable waiting for all task definitions to start running").Bool()
	// VERSION is set via ldflag
	VERSION = "0.0.0"
)

func main() {
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version(VERSION).Author("Travis Jeffery")
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
			os.Exit(1)
		}
	}

	err = c.UpdateService(cluster, service, count, &arn)
	if err != nil {
		logger.Printf("[error] update service: %s\n", err)
		os.Exit(1)
	}

	if *nowait == false {
		err := c.Wait(cluster, service, &arn)
		if err != nil {
			logger.Printf("[error] wait: %s\n", err)
			os.Exit(1)
		}
	}

	logger.Printf("[info] update service success")
}
