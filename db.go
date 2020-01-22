package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("us-east-1"))

func getItem(name string) (*book, error) {
    input := &dynamodb.GetItemInput{
        TableName: aws.String("Scratchers"),
        Key: map[string]*dynamodb.AttributeValue{
            "NAME": {
                S: aws.String(name),
            },
        },
    }
    result, err := db.GetItem(input)
    if err != nil {
        return nil, err
    }
    if result.Item == nil {
        return nil, nil
    }
    bk := new(book)
    err = dynamodbattribute.UnmarshalMap(result.Item, bk)
    if err != nil {
        return nil, err
    }
    return bk, nil
}

// Add a book record to DynamoDB.
func putItem(bk *book) error {
    input := &dynamodb.PutItemInput{
        TableName: aws.String("Scratchers"),
        Item: map[string]*dynamodb.AttributeValue{
            "NAME": {
                S: aws.String(bk.NAME),
            },
            "Description": {
                S: aws.String(bk.Description),
            },
            "Size": {
                S: aws.String(bk.Size),
            },
            "Price": {
                S: aws.String(bk.Price),
            },
        },
    }
_, err := db.PutItem(input)
    return err
}
