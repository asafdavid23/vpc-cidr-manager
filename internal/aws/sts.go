package aws

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func GetStsClient(cfg aws.Config) (*sts.Client, error) {
	stsClient := sts.NewFromConfig(cfg)

	return stsClient, nil
}

func AssumeRole(cfg aws.Config, stsClient *sts.Client, roleArn string) (aws.Config, error) {
	// Create a new config with the assumed role credentials
	creds := stscreds.NewAssumeRoleProvider(stsClient, roleArn)

	newCfg := cfg.Copy()
	newCfg.Credentials = aws.NewCredentialsCache(creds)

	return newCfg, nil
}
