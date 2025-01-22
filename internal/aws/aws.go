package aws

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func GetDynamoDBClient() (*dynamodb.Client, error) {
	region := os.Getenv("AWS_REGION")

	if region == "" {
		return nil, errors.New("AWS_REGION environment variable is not set")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("AWS_REGION")),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}
	client := dynamodb.NewFromConfig(cfg)

	if client == nil {
		return nil, errors.New("DynamoDB client is nil")
	}

	return client, nil
}

func checkTableExists(client *dynamodb.Client, name string) error {
	// Check if table already exists
	_, err := client.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String(name),
	})

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func CreateDynamoDBTable(client *dynamodb.Client, name string, logger *log.Logger) error {
	err := checkTableExists(client, name)

	if err != nil {
		var notFound *types.ResourceNotFoundException
		var inUse *types.ResourceInUseException

		if errors.As(err, &notFound) {
			logger.Debugf("Table %s does not exist, creating it now.\n", name)

			// Create table
			_, err = client.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
				TableName: aws.String(name),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("CIDR"),
						KeyType:       types.KeyTypeHash,
					},
				},
				AttributeDefinitions: []types.AttributeDefinition{
					{
						AttributeName: aws.String("CIDR"),
						AttributeType: types.ScalarAttributeTypeS,
					},
				},
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(5),
					WriteCapacityUnits: aws.Int64(5),
				},
			})

			if err != nil {
				return fmt.Errorf("%w", err)
			}
		} else if errors.As(err, &inUse) {
			logger.Printf("Table %s is already in use.\n", name)
		} else {
			return fmt.Errorf("%w", err)
		}
	} else {
		logger.Printf("Table %s already exists.\n", name)
	}

	return nil
}

func ReserveCIDR(client *dynamodb.Client, cidr string, vpcID string, vpcName string, logger *log.Logger) error {
	var sessionName string
	tableName := os.Getenv("DDB_TABLE_NAME")

	if tableName == "" {
		return fmt.Errorf("DDB_TABLE_NAME environment variable is not set")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		return fmt.Errorf("failed to load SDK config: %w", err)
	}

	stsClient := sts.NewFromConfig(cfg)

	output, err := stsClient.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})

	if err != nil {
		return fmt.Errorf("failed to get caller identity: %w", err)
	}

	arnParts := strings.Split(*output.Arn, "/")

	if len(arnParts) > 2 {
		sessionName := arnParts[2]
		logger.Debugf("Session name: %s", sessionName)
	} else {
		return fmt.Errorf("failed to parse session name from ARN: %s", *output.Arn)
	}

	err = checkTableExists(client, tableName)

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	// Validate the input CIDR
	_, newCIDR, err := net.ParseCIDR(cidr)

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	// Fetch all reserved CIDRs
	existingCIDRs, err := fetchExistingCIDRs(client, tableName)

	if err != nil {
		return fmt.Errorf("failed to fetch existing CIDRs: %w", err)
	}

	// Check if the new CIDR overlaps with any existing CIDRs

	for _, existingCIDR := range existingCIDRs {
		_, existingCIDRBlock, err := net.ParseCIDR(existingCIDR)

		if err != nil {
			continue
		}

		if cidrOverlaps(existingCIDRBlock, newCIDR) {
			return fmt.Errorf("CIDR %s overlaps with existing CIDR %s", cidr, existingCIDR)
		}
	}

	// Reserve the new CIDR
	_, err = client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: map[string]types.AttributeValue{
			"CIDR":       &types.AttributeValueMemberS{Value: cidr},
			"Status":     &types.AttributeValueMemberS{Value: "reserved"},
			"ReservedAt": &types.AttributeValueMemberS{Value: fmt.Sprintf("%v", time.Now())},
			"ReservedBy": &types.AttributeValueMemberS{Value: sessionName},
			"VpcId":      &types.AttributeValueMemberS{Value: vpcID},
			"VpcName":    &types.AttributeValueMemberS{Value: vpcName},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to reserve CIDR: %w", err)
	}
	return nil
}

func ReleaseCidr(client *dynamodb.Client, cidr string, logger *log.Logger) error {
	tableName := os.Getenv("DDB_TABLE_NAME")

	if tableName == "" {
		return fmt.Errorf("DDB_TABLE_NAME environment variable is not set")
	}

	err := checkTableExists(client, tableName)

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	_, err = client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"CIDR": &types.AttributeValueMemberS{Value: cidr},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to delete CIDR: %w", err)
	}
	return nil
}

// func ImportCIDRs() error {
// 	// Logic for importing CIDRs from live VPCs (mocked here)
// 	// Implement actual logic to pull CIDR blocks from AWS VPCs and insert them into DynamoDB.

// 	cidrs := []string{"10.0.0.0/24", "10.0.1.0/24"}
// 	for _, cidr := range cidrs {
// 		err := ReserveCIDR(cidr)
// 		if err != nil {
// 			return fmt.Errorf("failed to import CIDR %s: %w", cidr, err)
// 		}
// 	}
// 	return nil
// }

// fetchExistingCIDRs retrieves all reserved CIDRs from the DynamoDB table
func fetchExistingCIDRs(client *dynamodb.Client, tableName string) ([]string, error) {
	output, err := client.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName:            aws.String(tableName),
		ProjectionExpression: aws.String("CIDR"),
	})
	if err != nil {
		return nil, err
	}

	var cidrs []string
	for _, item := range output.Items {
		if cidrAttr, ok := item["CIDR"].(*types.AttributeValueMemberS); ok {
			cidrs = append(cidrs, cidrAttr.Value)
		}
	}
	return cidrs, nil
}

// cidrOverlaps checks if two CIDRs overlap
func cidrOverlaps(cidr1, cidr2 *net.IPNet) bool {
	return cidr1.Contains(cidr2.IP) || cidr2.Contains(cidr1.IP)
}
