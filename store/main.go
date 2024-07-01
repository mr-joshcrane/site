package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func main() {
	// Load AWS configuration from environment variables, shared credentials, or AWS config files
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	// Create a new DynamoDB client using the loaded AWS configuration

	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String("http://localhost:8000")
	})
	out, err := client.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		TableName: aws.String("TestTable"),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("Artist"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("SongTitle"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("Artist"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("SongTitle"),
				KeyType:       types.KeyTypeRange,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {
		panic("failed to create table, " + err.Error())
	}
	fmt.Println(out)
	// Call the ListTables API operation to list all the tables in DynamoDB
	result, err := client.ListTables(context.TODO(), &dynamodb.ListTablesInput{})
	if err != nil {
		panic("failed to list tables, " + err.Error())
	}

	// Print the table names
	fmt.Println("Tables:")
	for _, tableName := range result.TableNames {
		fmt.Println(tableName)
	}
}
