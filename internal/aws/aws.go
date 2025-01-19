package aws

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func GetDynamoDBClient() (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}
	return dynamodb.NewFromConfig(cfg), nil
}

func CreateDynamoDBTable(client *dynamodb.Client, name string) error {
	_, err := client.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
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
			{
				AttributeName: aws.String("Status"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("ReservedAt"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	return nil
}

func ReserveCIDR(cidr string) error {
	client, err := GetDynamoDBClient()
	if err != nil {
		return err
	}

	_, err = client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("VpcCidrReservations"),
		Item: map[string]types.AttributeValue{
			"CIDR":       &types.AttributeValueMemberS{Value: cidr},
			"Status":     &types.AttributeValueMemberS{Value: "reserved"},
			"ReservedAt": &types.AttributeValueMemberS{Value: fmt.Sprintf("%v", time.Now())},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to reserve CIDR: %w", err)
	}
	return nil
}

func DeleteCIDR(cidr string) error {
	client, err := GetDynamoDBClient()
	if err != nil {
		return err
	}

	_, err = client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String("VpcCidrReservations"),
		Key: map[string]types.AttributeValue{
			"CIDR": &types.AttributeValueMemberS{Value: cidr},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to delete CIDR: %w", err)
	}
	return nil
}

func ImportCIDRs() error {
	// Logic for importing CIDRs from live VPCs (mocked here)
	// Implement actual logic to pull CIDR blocks from AWS VPCs and insert them into DynamoDB.

	cidrs := []string{"10.0.0.0/24", "10.0.1.0/24"}
	for _, cidr := range cidrs {
		err := ReserveCIDR(cidr)
		if err != nil {
			return fmt.Errorf("failed to import CIDR %s: %w", cidr, err)
		}
	}
	return nil
}
