package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go/aws"
)

func GetIAMClient() (*iam.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	client := iam.NewFromConfig(cfg)

	return client, nil
}

func CreateAssumableRole(client *iam.Client, roleName string, policyFile, trustFile string) error {
	ctx := context.Background()

	// Load the trust relationship policy
	trustRelationship, err := os.ReadFile(trustFile)

	if err != nil {
		return fmt.Errorf("failed to read trust relationship file: %v", err)
	}

	// Create the IAM role with the trust relationship
	createRoleInput := &iam.CreateRoleInput{
		RoleName:                 aws.String(roleName),
		AssumeRolePolicyDocument: aws.String(string(trustRelationship)),
		Description:              aws.String("Assumable role for VPC CIDR Manager"),
	}

	_, err = client.CreateRole(ctx, createRoleInput)

	if err != nil {
		return fmt.Errorf("failed to create role: %v", err)
	}

	// Load the policy document
	policyDocument, err := os.ReadFile(policyFile)

	if err != nil {
		return fmt.Errorf("failed to read policy file: %v", err)
	}

	// Validate policy JSON
	var policyJSON map[string]interface{}

	if err = json.Unmarshal(policyDocument, &policyJSON); err != nil {
		return fmt.Errorf("failed to unmarshal policy document: %v", err)
	}

	putPolicyInput := &iam.PutRolePolicyInput{
		RoleName:       aws.String(roleName),
		PolicyName:     aws.String(roleName + "-policy"),
		PolicyDocument: aws.String(string(policyDocument)),
	}

	_, err = client.PutRolePolicy(ctx, putPolicyInput)

	if err != nil {
		return fmt.Errorf("failed to attach policy to role: %w", err)
	}

	return nil
}
