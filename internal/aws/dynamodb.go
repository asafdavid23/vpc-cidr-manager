package aws

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func PushToDynamoDB(client *dynamodb.Client, vpcInfo VPCInfo) error {

	itemExists, err := CheckItemExists(context.TODO(), client, vpcInfo.CIDR)

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
