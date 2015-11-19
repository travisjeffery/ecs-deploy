package main

import (
	"fmt"
	
	"github.com/travisjeffery/ecs-update/client"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	service = kingpin.Flag("service", "Name of Service to update.").Required().String()
	image = kingpin.Flag("image", "Name of Docker image to run.").Required().String()
	cluster = kingpin.Flag("cluster", "Name of ECS cluster.").Default("default").String()
	region = kingpin.Flag("region", "Name of AWS region.").Default("us-east-1").String()
)

func main() {
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version("1.0").Author("Travis Jeffery")
	kingpin.CommandLine.Help = "Update ECS service."
	kingpin.Parse()
	c := client.New(*region)

	arn, err := c.RegisterTaskDefinition(*service, *image)
	if err != nil {
		fmt.Printf("register task definition error: %s\n", err.Error())
		return
	}

	err = c.UpdateService(*cluster, *service, arn)
	if err != nil {
		fmt.Printf("update service error: %s\n", err.Error())
		return
	}
	fmt.Printf("update service success")
}
