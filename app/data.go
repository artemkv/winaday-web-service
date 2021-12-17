package app

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	log "github.com/sirupsen/logrus"
)

var (
	WIN_TABLE_NAME                  string = "winaday"
	WIN_TABLE_KEY                   string = "Key"
	WIN_TABLE_SORT_KEY              string = "SortKey"
	WIN_TABLE_TEXT_ATTR             string = "text"
	WIN_TABLE_OVERALL_ATTR          string = "overall"
	WIN_TABLE_USER_ID_ATTR          string = "userId"
	WIN_TABLE_USER_EMAIL_ATTR       string = "email"
	WIN_TABLE_LAST_ACCESSED_AT_ATTR string = "accessedAt"
)

type winItem struct {
	SortKey string
	Text    string
	Overall string
}

func updateUserProfile(userId string, email string, lastAccessed string) error {
	// get service
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return logAndConvertError(err)
	}
	svc := dynamodb.NewFromConfig(cfg)

	// define keys
	hashKey := "USER"
	sortKey := email

	// query input
	input := &dynamodb.PutItemInput{
		TableName: aws.String(WIN_TABLE_NAME),
		Item: map[string]types.AttributeValue{
			WIN_TABLE_KEY:                   &types.AttributeValueMemberS{Value: hashKey},
			WIN_TABLE_SORT_KEY:              &types.AttributeValueMemberS{Value: sortKey},
			WIN_TABLE_USER_ID_ATTR:          &types.AttributeValueMemberS{Value: userId},
			WIN_TABLE_LAST_ACCESSED_AT_ATTR: &types.AttributeValueMemberS{Value: lastAccessed},
		},
		ReturnValues: types.ReturnValueNone,
	}

	// run query
	_, err = svc.PutItem(context.TODO(), input)
	if err != nil {
		return logAndConvertError(err)
	}

	// done
	return nil
}

func updateWin(userId string, date string, win winData) error {
	// get service
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return logAndConvertError(err)
	}
	svc := dynamodb.NewFromConfig(cfg)

	// define keys
	hashKey := fmt.Sprintf("WIN#%s", userId)
	sortKey := date

	// encode data
	text := base64.StdEncoding.EncodeToString([]byte(win.Text))
	overallResult := strconv.Itoa(win.OverallResult)

	// query input
	input := &dynamodb.PutItemInput{
		TableName: aws.String(WIN_TABLE_NAME),
		Item: map[string]types.AttributeValue{
			WIN_TABLE_KEY:          &types.AttributeValueMemberS{Value: hashKey},
			WIN_TABLE_SORT_KEY:     &types.AttributeValueMemberS{Value: sortKey},
			WIN_TABLE_TEXT_ATTR:    &types.AttributeValueMemberS{Value: text},
			WIN_TABLE_OVERALL_ATTR: &types.AttributeValueMemberN{Value: overallResult},
		},
		ReturnValues: types.ReturnValueNone,
	}

	// run query
	_, err = svc.PutItem(context.TODO(), input)
	if err != nil {
		return logAndConvertError(err)
	}

	// done
	return nil
}

func getWin(userId string, date string) (*winData, error) {
	// get service
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, logAndConvertError(err)
	}
	svc := dynamodb.NewFromConfig(cfg)

	// define keys
	hashKey := fmt.Sprintf("WIN#%s", userId)
	sortKey := date

	// query expression
	projection := expression.NamesList(
		expression.Name(WIN_TABLE_SORT_KEY),
		expression.Name(WIN_TABLE_TEXT_ATTR),
		expression.Name(WIN_TABLE_OVERALL_ATTR))
	expr, err := expression.NewBuilder().WithProjection(projection).Build()
	if err != nil {
		return nil, logAndConvertError(err)
	}

	// query input
	input := &dynamodb.GetItemInput{
		TableName: aws.String(WIN_TABLE_NAME),
		Key: map[string]types.AttributeValue{
			WIN_TABLE_KEY:      &types.AttributeValueMemberS{Value: hashKey},
			WIN_TABLE_SORT_KEY: &types.AttributeValueMemberS{Value: sortKey},
		},
		ExpressionAttributeNames: expr.Names(),
		ProjectionExpression:     expr.Projection(),
	}

	// run query
	result, err := svc.GetItem(context.TODO(), input)
	if err != nil {
		return nil, logAndConvertError(err)
	}

	// re-pack the results
	if result.Item == nil {
		return nil, nil
	}
	item := winItem{}
	err = attributevalue.UnmarshalMap(result.Item, &item)
	if err != nil {
		return nil, logAndConvertError(err)
	}
	overallResult, err := strconv.Atoi(item.Overall)
	if err != nil {
		return nil, logAndConvertError(err)
	}
	textBytes, err := base64.StdEncoding.DecodeString(item.Text)
	if err != nil {
		return nil, logAndConvertError(err)
	}
	win := winData{
		Text:          string(textBytes),
		OverallResult: overallResult,
	}

	return &win, nil
}

func logAndConvertError(err error) error {
	log.Printf("%v", err)
	return fmt.Errorf("service unavailable")
}

func generateTimestamp() string {
	return time.Now().Format(time.RFC3339)
}
