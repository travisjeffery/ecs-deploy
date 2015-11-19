package client

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/aws/session"
)

type Client struct {
	svc *ecs.ECS
}

func New(region) *Client {
	sess := session.New(&aws.Config{Region: aws.String(region)})
	svc := ecs.New(sess)
	return &Client{
		svc: svc,
	}
}

// Register Task Definition by updating existing task definition's image.
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

// Update service to use new task definition.
func (c *Client) UpdateService(cluster, service, arn string) error {
	input := &ecs.UpdateServiceInput{
		Cluster:        aws.String(cluster),
		DesiredCount:   aws.Int64(1),
		Service:        aws.String(service),
		TaskDefinition: aws.String(arn),
	}
	_, err := c.svc.UpdateService(input)
	return err
}

// Get existing container definitions.
func (c *Client) GetContainerDefinitions(service string) ([]*ecs.ContainerDefinition, error) {
	output, err := c.svc.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(service),
	})
	if err != nil {
		return nil, err
	}
	return output.TaskDefinition.ContainerDefinitions, nil
}
