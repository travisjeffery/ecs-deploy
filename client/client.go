package client

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

type Client struct {
	svc          *ecs.ECS
	logger       *log.Logger
	pollInterval time.Duration
}

func New(region *string, logger *log.Logger) *Client {
	sess := session.New(&aws.Config{Region: region})
	svc := ecs.New(sess)
	return &Client{
		svc:          svc,
		pollInterval: time.Second * 5,
		logger:       logger,
	}
}

// RegisterTaskDefinition updates the existing task definition's image.
func (c *Client) RegisterTaskDefinition(task *string, image, tag *[]string) (string, error) {
	output, err := c.svc.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{
		TaskDefinition: task,
	})
	if err != nil {
		return "", err
	}
	for _, d := range output.TaskDefinition.ContainerDefinitions {
		for j := range *image {
			if strings.HasPrefix(*d.Image, (*image)[j]) {
				i := fmt.Sprintf("%s:%s", (*image)[j], (*tag)[j])
				d.Image = &i
			}
		}
	}
	input := &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: output.TaskDefinition.ContainerDefinitions,
		Family:               output.TaskDefinition.Family,
		NetworkMode:          output.TaskDefinition.NetworkMode,
		PlacementConstraints: output.TaskDefinition.PlacementConstraints,
		TaskRoleArn:          output.TaskDefinition.TaskRoleArn,
		Volumes:              output.TaskDefinition.Volumes,
	}
	resp, err := c.svc.RegisterTaskDefinition(input)
	if err != nil {
		return "", err
	}
	return *resp.TaskDefinition.TaskDefinitionArn, nil
}

// UpdateService updates the service to use the new task definition.
func (c *Client) UpdateService(cluster, service *string, count *int64, arn *string) error {
	input := &ecs.UpdateServiceInput{
		Cluster: cluster,
		Service: service,
	}
	if *count != -1 {
		input.DesiredCount = count
	}
	if arn != nil {
		input.TaskDefinition = arn
	}
	_, err := c.svc.UpdateService(input)
	return err
}

// Wait waits for the service to finish being updated.
func (c *Client) Wait(cluster, service, arn *string) error {
	t := time.NewTicker(c.pollInterval)
	for {
		select {
		case <-t.C:
			s, err := c.GetDeployment(cluster, service, arn)
			if err != nil {
				return err
			}
			c.logger.Printf("[info] --> desired: %d, pending: %d, running: %d", *s.DesiredCount, *s.PendingCount, *s.RunningCount)
			if *s.RunningCount == *s.DesiredCount {
				return nil
			}
		}
	}
}

// GetDeployment gets the deployment for the arn.
func (c *Client) GetDeployment(cluster, service, arn *string) (*ecs.Deployment, error) {
	input := &ecs.DescribeServicesInput{
		Cluster:  cluster,
		Services: []*string{service},
	}
	output, err := c.svc.DescribeServices(input)
	if err != nil {
		return nil, err
	}
	ds := output.Services[0].Deployments
	for _, d := range ds {
		if *d.TaskDefinition == *arn {
			return d, nil
		}
	}
	return nil, nil
}

// GetContainerDefinitions get container definitions of the service.
func (c *Client) GetContainerDefinitions(task *string) ([]*ecs.ContainerDefinition, error) {
	output, err := c.svc.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{
		TaskDefinition: task,
	})
	if err != nil {
		return nil, err
	}
	return output.TaskDefinition.ContainerDefinitions, nil
}
