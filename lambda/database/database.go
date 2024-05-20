package database

import (
	"go-aws/lambda/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)
const (
    TABLE_NAME = "userTable"
)

type UserStore interface {
    DoesUserExist(username string) (bool, error)
    InserUser(user types.RegisterUser) error
}

type DynamoDBClient struct {
    databaseStore *dynamodb.DynamoDB
}

func NewDynamoDBClient() DynamoDBClient {
    dbSession := session.Must(session.NewSession())
    db := dynamodb.New(dbSession)
    return DynamoDBClient{
        databaseStore:db,
    }
}

// Does this user exist

func (u DynamoDBClient) DoesUserExist(username string) (bool, error) {
    result, err := u.databaseStore.GetItem(&dynamodb.GetItemInput{
       TableName: aws.String(TABLE_NAME),
       Key: map[string]*dynamodb.AttributeValue {
            "username" : {
                S: aws.String(username),
            },
       },
    })
    if err != nil {
        return true, err
    }

    if result.Item == nil{
        return false, nil
    }
    return true, nil
}

func (u DynamoDBClient) InserUser(user types.RegisterUser) error {
    //assemble items and insert it
    item := &dynamodb.PutItemInput{
        TableName: aws.String(TABLE_NAME),
        Item: map[string]*dynamodb.AttributeValue{
            "username": {
                S: aws.String(user.Username),
            },
            "password": {
                S: aws.String(user.Password),
            },
        },
    }
    _, err := u.databaseStore.PutItem(item)
    if err != nil {
        return err
    }
    return nil
}






