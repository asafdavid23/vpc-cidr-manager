package aws

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func PushToDynamoDB(client *dynamodb.Client, vpcInfo VPCInfo) error {
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
