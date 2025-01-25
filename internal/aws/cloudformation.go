package aws

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
)

func InitializeCFNClient(cfg aws.Config) (*cloudformation.Client, error) {
	cfnClient := cloudformation.NewFromConfig(cfg)

	return cfnClient, nil

}

func CreateCFStack(stackName string, templateBody string) (*cloudformation.CreateStackOutput, error) {
	// Lo
}
