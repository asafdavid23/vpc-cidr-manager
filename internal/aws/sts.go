package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func GetStsClient() (*sts.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		return nil, fmt.Errorf("failed to load configuration, %v", err)
	}

	stsClient := sts.NewFromConfig(cfg)

	return stsClient, nil
}

func AssumeRole(ctx context.Context, cfg aws.Config, stsClient *sts.Client, roleArn string) (aws.Config, error) {
	// Create a new config with the assumed role credentials
	creds := stscreds.NewAssumeRoleProvider(stsClient, roleArn)

	newCfg := cfg.Copy()
	newCfg.Credentials = aws.NewCredentialsCache(creds)

	return newCfg, nil
}
