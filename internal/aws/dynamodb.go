package aws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
)

func PushToDynamoDB(ctx context.Context, client *dynamodb.Client, vpcInfo VPCInfo) error {

	itemExists, err := CheckItemExists(ctx, client, vpcInfo.CIDR)

	if err != nil {
		return fmt.Errorf("Got error checking if item exists: %v", err)
	}

	if itemExists {
		return fmt.Errorf("Item already exists in DynamoDB")
	}

	// Covert VPCInfo to JSON and push to DynamoDB
	av, err := attributevalue.MarshalMap(vpcInfo)

	if err != nil {
		return fmt.Errorf("Got error marshalling map: %v", err)
	}

	tableName := os.Getenv("DDB_TABLE_NAME")

	_, err = client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	})

	if err != nil {
		return fmt.Errorf("Got error calling PutItem: %v", err)
	}

	return nil
}

func CheckItemExists(ctx context.Context, client *dynamodb.Client, cidr string) (bool, error) {
	tableName := os.Getenv("DDB_TABLE_NAME")

	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"CIDR": &types.AttributeValueMemberS{Value: cidr},
		},
		ProjectionExpression: aws.String("vpcId"),
	}

	output, err := client.GetItem(ctx, input)

	if err != nil {
		return false, fmt.Errorf("Got error calling GetItem: %v", err)
	}

	if output.Item == nil {
		return false, nil
	}

	return true, nil
}

func GetDynamoDBClient(ctx context.Context, cfg aws.Config) (*dynamodb.Client, error) {
	client := dynamodb.NewFromConfig(cfg)

	if client == nil {
		return nil, errors.New("DynamoDB client is nil")
	}

	return client, nil
}

func checkTableExists(ctx context.Context, client *dynamodb.Client, name string) error {
	// Check if table already exists
	_, err := client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(name),
	})

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func CreateDynamoDBTable(ctx context.Context, client *dynamodb.Client, name string, logger *log.Logger) error {
	err := checkTableExists(ctx, client, name)

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

func ReserveCIDR(ctx context.Context, client *dynamodb.Client, cidr string, vpcID string, vpcName string, logger *log.Logger) error {
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
		sessionName = arnParts[2]
		logger.Debugf("Session name: %s", sessionName)
	} else {
		return fmt.Errorf("failed to parse session name from ARN: %s", *output.Arn)
	}

	err = checkTableExists(ctx, client, tableName)

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
			"VpcId":      &types.AttributeValueMemberS{Value: vpcID},
			"VpcName":    &types.AttributeValueMemberS{Value: vpcName},
			"ReservedAt": &types.AttributeValueMemberS{Value: fmt.Sprintf("%v", time.Now())},
			"ReservedBy": &types.AttributeValueMemberS{Value: sessionName},
			"Status":     &types.AttributeValueMemberS{Value: "reserved"},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to reserve CIDR: %w", err)
	}
	return nil
}

func ReleaseCidr(ctx context.Context, client *dynamodb.Client, cidr string, logger *log.Logger) error {
	tableName := os.Getenv("DDB_TABLE_NAME")

	if tableName == "" {
		return fmt.Errorf("DDB_TABLE_NAME environment variable is not set")
	}

	err := checkTableExists(ctx, client, tableName)

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

func ListCIDRs(ctx context.Context, client *dynamodb.Client, tableName string, outputFormat string) error {
	if tableName == "" {
		return fmt.Errorf("DDB_TABLE_NAME environment variable is not set")
	}

	err := checkTableExists(ctx, client, tableName)

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	output, err := client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	})

	if err != nil {
		return fmt.Errorf("failed to scan table: %w", err)
	}

	if len(output.Items) == 0 {
		return fmt.Errorf("no CIDRs found in table %s", tableName)
	}

	switch outputFormat {
	case "json":
		outputJSON, err := json.MarshalIndent(output.Items, "", "  ")

		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}

		fmt.Println(string(outputJSON))

	case "table":
		headers := []string{}

		if len(output.Items) > 0 {
			for key := range output.Items[0] {
				headers = append(headers, key)
			}
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader(headers)

		for _, item := range output.Items {
			row := []string{}

			for _, header := range headers {
				if attr, ok := item[header].(*types.AttributeValueMemberS); ok {
					row = append(row, attr.Value)
				} else {
					row = append(row, "")
				}
			}

			table.Append(row)
		}

		table.Render()

	default:
		return fmt.Errorf("unsupported output format: %s", outputFormat)
	}

	return nil
}
