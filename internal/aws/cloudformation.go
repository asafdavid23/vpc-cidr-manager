package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

func InitializeCFNClient(cfg aws.Config) (*cloudformation.Client, error) {
	cfnClient := cloudformation.NewFromConfig(cfg)

	return cfnClient, nil

}

func CreateCFNStack(ctx context.Context, client *cloudformation.Client, stackName string, templateBody string) (*cloudformation.CreateStackOutput, error) {
	output, err := client.CreateStack(ctx, &cloudformation.CreateStackInput{
		StackName:    aws.String(stackName),
		TemplateBody: aws.String(templateBody),
		Capabilities: []types.Capability{
			types.CapabilityCapabilityIam,
			types.CapabilityCapabilityNamedIam,
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create stack: %v", err)
	}

	return output, nil
}

// waitForStackToBeCreated waits for the CloudFormation stack to be created successfully.
func WaitForStackToBeCreated(ctx context.Context, cfnClient *cloudformation.Client, stackName string) error {
	// Wait for the stack to reach the 'CREATE_COMPLETE' status
	waiterCfg := cloudformation.NewStackCreateCompleteWaiter(cfnClient)
	describeStacksInput := &cloudformation.DescribeStacksInput{
		StackName: aws.String(stackName),
	}

	// The waiter will keep checking the stack status until it's either complete or failed
	err := waiterCfg.Wait(ctx, describeStacksInput, 5*time.Minute) // You can change the timeout value as needed
	if err != nil {
		return fmt.Errorf("failed to wait for stack creation to complete: %v", err)
	}

	return nil
}

func ValidateCFNStackTemplate(ctx context.Context, client *cloudformation.Client, templateBody string) error {
	_, err := client.ValidateTemplate(ctx, &cloudformation.ValidateTemplateInput{
		TemplateBody: aws.String(templateBody),
	})

	if err != nil {
		return fmt.Errorf("failed to validate template: %v", err)
	}

	return nil
}
