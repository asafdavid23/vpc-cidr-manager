package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func GetStsClient() (*sts.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		return nil, fmt.Errorf("failed to load configuration, %v", err)
	}

	stsSvc := sts.NewFromConfig(cfg)

	return stsSvc, nil
}

func AssumeRole(roleName string, account string) (*sts.AssumeRoleOutput, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration, %v", err)
	}

	stsSvc := sts.NewFromConfig(cfg)
	roleArn := "arn:aws:iam::" + account + ":role/" + roleName

	input := &sts.AssumeRoleInput{
		RoleArn:         aws.String(roleArn),
		RoleSessionName: aws.String("vpc-cidr-manager"),
	}

	result, err := stsSvc.AssumeRole(context.TODO(), input)

	if err != nil {
		return nil, err
	}

	return result, nil
}
