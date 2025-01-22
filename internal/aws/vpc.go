package aws

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type VPCInfo struct {
	CIDR       string    `json:"cidrBlock"`
	VpcID      string    `json:"vpcId"`
	VpcName    string    `json:"vpcName"`
	ReservedAt time.Time `json:"reservedAt"`
	ReservedBy string    `json:"reservedBy"`
	Status     string    `json:"status"`
}

func GetEc2Client() (*ec2.Client, error) {
	region := os.Getenv("AWS_REGION")

	if region == "" {
		return nil, fmt.Errorf("AWS_REGION environment variable is not set")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)

	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	client := ec2.NewFromConfig(cfg)

	if client == nil {
		return nil, fmt.Errorf("EC2 client is nil")
	}

	return client, nil
}

func GetVpcInfo(client *ec2.Client, vpcId string) (VPCInfo, error) {
	var vpcInfo VPCInfo
	var sessionName string

	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		return VPCInfo{}, fmt.Errorf("unable to load SDK config, %v", err)
	}

	stsClient := sts.NewFromConfig(cfg)

	output, err := stsClient.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})

	if err != nil {
		return VPCInfo{}, fmt.Errorf("unable to get caller identity, %v", err)
	}

	arnParts := strings.Split(*output.Arn, "/")

	if len(arnParts) > 2 {
		sessionName = arnParts[2]
	} else {
		return VPCInfo{}, fmt.Errorf("unable to parse session name from ARN: %s", *output.Arn)
	}

	input := &ec2.DescribeVpcsInput{
		VpcIds: []string{vpcId},
	}

	result, err := client.DescribeVpcs(context.TODO(), input)

	if err != nil {
		return VPCInfo{}, fmt.Errorf("%v", err)
	}

	if len(result.Vpcs) == 0 {
		return VPCInfo{}, fmt.Errorf("VPC with ID %s not found", vpcId)
	}

	for _, vpc := range result.Vpcs {
		vpcInfo = VPCInfo{
			CIDR:       *vpc.CidrBlock,
			VpcID:      *vpc.VpcId,
			ReservedAt: time.Now(),
			ReservedBy: sessionName,
			Status:     "reserved",
		}

		if vpc.Tags != nil {
			for _, tag := range vpc.Tags {
				if *tag.Key == "Name" {
					vpcInfo.VpcName = *tag.Value
				}
			}
		}
	}

	if err != nil {
		return VPCInfo{}, fmt.Errorf("Error pushing VPC info to DynamoDB: %v", err)
	}

	return vpcInfo, nil
}
