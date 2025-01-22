package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

func AssumeRole(roleArn string) (*sts.AssumeRoleOutput, error) {
	sess := session.Must(session.NewSession())
	stsSvc := sts.New(sess)

	input := &sts.AssumeRoleInput{
		RoleArn:         aws.String(roleArn),
		RoleSessionName: aws.String("vpc-cidr-manager"),
	}

	result, err := stsSvc.AssumeRole(input)

	if err != nil {
		return nil, err
	}

	return result, nil
}
