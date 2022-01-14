package app

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	log "github.com/sirupsen/logrus"
)

const (
	OVERALL_DAY_RESULT_NO_WIN_YET           = 0
	OVERALL_DAY_RESULT_GOT_MY_WIN           = 1
	OVERALL_DAY_RESULT_COULD_NOT_GET_MY_WIN = 2
	OVERALL_DAY_RESULT_UNUSED               = 3
	OVERALL_DAY_RESULT_AWESOME_ACHIEVEMENT  = 4
)

const (
	//WIN_TABLE_NAME            string = "winaday-test"
	WIN_TABLE_NAME            string = "winaday"
	WIN_TABLE_KEY             string = "Key"
	WIN_TABLE_SORT_KEY        string = "SortKey"
	WIN_TABLE_TEXT_ATTR       string = "text"
	WIN_TABLE_OVERALL_ATTR    string = "overall"
	WIN_TABLE_PRIORITIES_ATTR string = "priorities"
	WIN_TABLE_ITEMS_ATTR      string = "items"
	WIN_TABLE_UPDATED_AT_ATTR string = "udpatedAt"
)

type winItem struct {
	SortKey    string
	Text       string
	Overall    string
	Priorities []string
}

type prioritiesListItem struct {
	SortKey string
	Items   []priorityItem
}

type priorityItem struct {
	Id        string
	Text      string
	Color     int
	IsDeleted bool
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

	priorities, err := attributevalue.MarshalList(win.Priorities)
	if err != nil {
		return logAndConvertError(err)
	}

	// query input
	input := &dynamodb.PutItemInput{
		TableName: aws.String(WIN_TABLE_NAME),
		Item: map[string]types.AttributeValue{
			WIN_TABLE_KEY:             &types.AttributeValueMemberS{Value: hashKey},
			WIN_TABLE_SORT_KEY:        &types.AttributeValueMemberS{Value: sortKey},
			WIN_TABLE_TEXT_ATTR:       &types.AttributeValueMemberS{Value: text},
			WIN_TABLE_OVERALL_ATTR:    &types.AttributeValueMemberN{Value: overallResult},
			WIN_TABLE_PRIORITIES_ATTR: &types.AttributeValueMemberL{Value: priorities},
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
		expression.Name(WIN_TABLE_OVERALL_ATTR),
		expression.Name(WIN_TABLE_PRIORITIES_ATTR))
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
		Priorities:    item.Priorities,
	}

	return &win, nil
}

func updatePriorities(userId string, priorities priorityListData, updatedAt string) error {
	// get service
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return logAndConvertError(err)
	}
	svc := dynamodb.NewFromConfig(cfg)

	// define keys
	hashKey := "PRIORITIES"
	sortKey := userId

	// encode data
	encodedPriorities, err := attributevalue.MarshalList(encodePriorities(priorities.Items, 100))
	if err != nil {
		return logAndConvertError(err)
	}

	// query input
	input := &dynamodb.PutItemInput{
		TableName: aws.String(WIN_TABLE_NAME),
		Item: map[string]types.AttributeValue{
			WIN_TABLE_KEY:             &types.AttributeValueMemberS{Value: hashKey},
			WIN_TABLE_SORT_KEY:        &types.AttributeValueMemberS{Value: sortKey},
			WIN_TABLE_ITEMS_ATTR:      &types.AttributeValueMemberL{Value: encodedPriorities},
			WIN_TABLE_UPDATED_AT_ATTR: &types.AttributeValueMemberS{Value: updatedAt},
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

func getPriorities(userId string) (*priorityListData, error) {
	// get service
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, logAndConvertError(err)
	}
	svc := dynamodb.NewFromConfig(cfg)

	// define keys
	hashKey := "PRIORITIES"
	sortKey := userId

	// query expression
	projection := expression.NamesList(
		expression.Name(WIN_TABLE_SORT_KEY),
		expression.Name(WIN_TABLE_ITEMS_ATTR))
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
	item := prioritiesListItem{}
	err = attributevalue.UnmarshalMap(result.Item, &item)
	if err != nil {
		return nil, logAndConvertError(err)
	}

	prioritiesToDecode := make([]priorityData, len(item.Items))
	for i, p := range item.Items {
		prioritiesToDecode[i] = priorityData{
			Id:        p.Id,
			Text:      p.Text,
			Color:     p.Color,
			IsDeleted: p.IsDeleted,
		}
	}
	prioritiesDecoded, err := decodePriorities(prioritiesToDecode)
	if err != nil {
		return nil, logAndConvertError(err)
	}

	priorityList := priorityListData{
		Items: prioritiesDecoded,
	}

	return &priorityList, nil
}

func logAndConvertError(err error) error {
	log.Printf("%v", err)
	return fmt.Errorf("service unavailable")
}

// base-64 encodes the priority text
// keeps only maxItems, trying to keep as much non-deleted (active) priorities as possible
func encodePriorities(priorities []priorityData, maxItems int) []priorityData {
	active := 0
	total := 0

	for _, p := range priorities {
		if !p.IsDeleted {
			active++
		}
		total++
	}

	activeAllowed := active
	if activeAllowed > maxItems {
		activeAllowed = maxItems
	}
	deletedAllowed := maxItems - activeAllowed
	if deletedAllowed < 0 {
		deletedAllowed = 0
	}
	totalAllowed := total
	if totalAllowed > maxItems {
		totalAllowed = maxItems
	}

	var encoded = make([]priorityData, totalAllowed)

	pos := 0
	for _, p := range priorities {
		take := false
		if p.IsDeleted {
			if deletedAllowed > 0 {
				take = true
				deletedAllowed--
			}
		} else {
			if activeAllowed > 0 {
				take = true
				activeAllowed--
			}
		}

		if take {
			encoded[pos] = priorityData{
				Id:        p.Id,
				Color:     p.Color,
				Text:      base64.StdEncoding.EncodeToString([]byte(p.Text)),
				IsDeleted: p.IsDeleted,
			}
			pos++
		}
	}

	return encoded
}

func decodePriorities(priorities []priorityData) ([]priorityData, error) {
	var decoded = make([]priorityData, len(priorities))

	for i, p := range priorities {
		textBytes, err := base64.StdEncoding.DecodeString(p.Text)
		if err != nil {
			return nil, logAndConvertError(err)
		}

		decoded[i] = priorityData{
			Id:        p.Id,
			Color:     p.Color,
			Text:      string(textBytes),
			IsDeleted: p.IsDeleted,
		}
	}

	return decoded, nil
}

// Returns wins [from:to]
func getWins(userId string, from string, to string) ([]winOnDayData, error) {
	// get service
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, logAndConvertError(err)
	}
	svc := dynamodb.NewFromConfig(cfg)

	// define keys
	hashKey := fmt.Sprintf("WIN#%s", userId)

	// query expression
	projection := expression.NamesList(
		expression.Name(WIN_TABLE_SORT_KEY),
		expression.Name(WIN_TABLE_TEXT_ATTR),
		expression.Name(WIN_TABLE_OVERALL_ATTR),
		expression.Name(WIN_TABLE_PRIORITIES_ATTR))
	expr, err := expression.NewBuilder().WithKeyCondition(
		expression.KeyAnd(
			expression.Key(WIN_TABLE_KEY).Equal(expression.Value(hashKey)),
			expression.KeyBetween(expression.Key(WIN_TABLE_SORT_KEY), expression.Value(from), expression.Value(to))),
	).WithProjection(projection).Build()
	if err != nil {
		return nil, logAndConvertError(err)
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(WIN_TABLE_NAME),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
	}

	// run query
	result, err := svc.Query(context.TODO(), input)
	if err != nil {
		return nil, logAndConvertError(err)
	}

	// re-pack the results
	wins := make([]winOnDayData, len(result.Items))
	for i, v := range result.Items {
		// TODO: extract into a separate method, to reuse in getWin
		item := winItem{}
		err = attributevalue.UnmarshalMap(v, &item)
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
			Priorities:    item.Priorities,
		}

		wins[i] = winOnDayData{
			Date: item.SortKey,
			Win:  win,
		}
	}

	// done
	return wins, nil
}

// Returns wins [from:to]
func getWinDays(userId string, from string, to string) ([]string, error) {
	// get service
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, logAndConvertError(err)
	}
	svc := dynamodb.NewFromConfig(cfg)

	// define keys
	hashKey := fmt.Sprintf("WIN#%s", userId)

	// query expression
	projection := expression.NamesList(
		expression.Name(WIN_TABLE_SORT_KEY),
		expression.Name(WIN_TABLE_OVERALL_ATTR))
	expr, err := expression.NewBuilder().WithKeyCondition(
		expression.KeyAnd(
			expression.Key(WIN_TABLE_KEY).Equal(expression.Value(hashKey)),
			expression.KeyBetween(expression.Key(WIN_TABLE_SORT_KEY), expression.Value(from), expression.Value(to))),
	).WithProjection(projection).Build()
	if err != nil {
		return nil, logAndConvertError(err)
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(WIN_TABLE_NAME),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
	}

	// run query
	result, err := svc.Query(context.TODO(), input)
	if err != nil {
		return nil, logAndConvertError(err)
	}

	// re-pack the results
	days := make([]string, 0, len(result.Items))
	for _, v := range result.Items {
		item := winItem{}
		err = attributevalue.UnmarshalMap(v, &item)
		if err != nil {
			return nil, logAndConvertError(err)
		}
		overallResult, err := strconv.Atoi(item.Overall)
		if err != nil {
			return nil, logAndConvertError(err)
		}
		if overallResult == OVERALL_DAY_RESULT_GOT_MY_WIN ||
			overallResult == OVERALL_DAY_RESULT_AWESOME_ACHIEVEMENT {
			days = append(days, item.SortKey)
		}
	}

	// done
	return days, nil
}
