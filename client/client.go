package client

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

type Client struct {
	svc *ecs.ECS
}

func New(region string) *Client {
	sess := session.New(&aws.Config{Region: aws.String(region)})
	svc := ecs.New(sess)
	return &Client{
		svc: svc,
	}
}

// RegisterTaskDefinition updates the existing task definition's image.
func (c *Client) RegisterTaskDefinition(service, image string) (string, error) {
	defs, err := c.GetContainerDefinitions(service)
	if err != nil {
		return "", err
	}
	defs[0].Image = aws.String("taplytics/app:master")
	input := &ecs.RegisterTaskDefinitionInput{
		Family:               aws.String(service),
		ContainerDefinitions: defs,
	}
	resp, err := c.svc.RegisterTaskDefinition(input)
	if err != nil {
		return "", err
	}
	return *resp.TaskDefinition.TaskDefinitionArn, nil
}

// UpdateService updates the service to use the new task definition.
func (c *Client) UpdateService(cluster, service string, count int64, arn string) error {
	input := &ecs.UpdateServiceInput{
		Cluster:        aws.String(cluster),
		DesiredCount:   aws.Int64(count),
		Service:        aws.String(service),
		TaskDefinition: aws.String(arn),
	}
	_, err := c.svc.UpdateService(input)
	return err
}

// GetContainerDefinitions get container definitions of the service.
func (c *Client) GetContainerDefinitions(service string) ([]*ecs.ContainerDefinition, error) {
	output, err := c.svc.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(service),
	})
	if err != nil {
		return nil, err
	}
	return output.TaskDefinition.ContainerDefinitions, nil
}
