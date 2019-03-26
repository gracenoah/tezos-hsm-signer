package watermark

import (
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// DynamoWatermark stores the last-signed level in memory
type DynamoWatermark struct {
	table    string
	dynamodb *dynamodb.DynamoDB
}

// GetDynamoWatermark returns a new dynamo watermark manager
func GetDynamoWatermark(table string) Watermark {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_DEFAULT_REGION")),
	})
	if err != nil {
		log.Fatal("Unable to initialize dynamo watermark")
	}

	return &DynamoWatermark{
		table:    table,
		dynamodb: dynamodb.New(sess),
	}
}

// getDynamoKey used as the hash identifier of each entry
func getDynamoKey(keyHash string, chainID string, opMagicByte uint8) string {
	return fmt.Sprintf("%v-%v-%v", keyHash, chainID, opMagicByte)
}

// getCurrentLevel watermarked in Dynamo
func (mw *DynamoWatermark) getCurrentLevel(keyHash string, chainID string, opMagicByte uint8) (*big.Int, error) {
	// Get Item
	result, err := mw.dynamodb.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(mw.table),
		Key: map[string]*dynamodb.AttributeValue{
			"KeyChainOp": {S: aws.String(getDynamoKey(keyHash, chainID, opMagicByte))},
		},
		ConsistentRead: aws.Bool(true),
	})
	if err != nil {
		// There was an error retrieving the dynamo item
		return nil, err
	}
	if result.Item["Level"] == nil {
		// The key does not exist in dynamo
		return nil, nil
	}
	newInt, _ := new(big.Int).SetString(*result.Item["Level"].S, 10)
	return newInt, nil
}

// putItem for the first time into Dynamo
func (mw *DynamoWatermark) putItem(keyHash string, chainID string, opMagicByte uint8, level *big.Int) error {
	_, err := mw.dynamodb.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(mw.table),
		Item: map[string]*dynamodb.AttributeValue{
			"KeyChainOp": {S: aws.String(getDynamoKey(keyHash, chainID, opMagicByte))},
			"Level":      {S: aws.String(level.String())},
		},
		ConditionExpression: aws.String("attribute_not_exists(KeyChainOp)"),
	})
	return err
}

// updateItem with a new level in dynamo
func (mw *DynamoWatermark) updateItem(keyHash string, chainID string, opMagicByte uint8, currentLevel *big.Int, newLevel *big.Int) error {
	// Update Item
	_, err := mw.dynamodb.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(mw.table),
		Key: map[string]*dynamodb.AttributeValue{
			"KeyChainOp": {S: aws.String(getDynamoKey(keyHash, chainID, opMagicByte))},
		},
		ExpressionAttributeNames: map[string]*string{
			"#Level": aws.String("Level"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":newval":  &dynamodb.AttributeValue{S: aws.String(newLevel.String())},
			":currval": &dynamodb.AttributeValue{S: aws.String(currentLevel.String())},
		},
		UpdateExpression:    aws.String("SET #Level = :newval"),
		ConditionExpression: aws.String("#Level = :currval"),
	})
	return err
}

// IsSafeToSign returns true if the provided (key, chainID, opMagicByte) tuple has
// not yet been signed at this or greater levels
func (mw *DynamoWatermark) IsSafeToSign(keyHash string, chainID string, opMagicByte uint8, level *big.Int) bool {

	currentLevel, err := mw.getCurrentLevel(keyHash, chainID, opMagicByte)
	if err != nil {
		log.Println("Error: Unable to get current level", err)
		return false
	}

	// Create a new item if none currently exists
	if currentLevel == nil {
		err := mw.putItem(keyHash, chainID, opMagicByte, level)
		if err != nil {
			return false
		}
		return true
	}

	// Update existing items
	if level.Cmp(currentLevel) != 1 {
		log.Println("Warning: Attempted to sign at an unsafe level. Will not allow.")
		return false
	} else {
		err := mw.updateItem(keyHash, chainID, opMagicByte, currentLevel, level)
		if err != nil {
			return false
		}
		return true
	}
}
