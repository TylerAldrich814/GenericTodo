package database

import (
	"encoding/json"
	"errors"
	"log"
	"utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)
type DBItem = utils.DBItem

type Database struct {
  db *dynamodb.DynamoDB
}

// NewSession: return a Database struct. with a DynamoDB instance,
//    At the moment, NewSession is hard wired to use use Region "us-east-2".
func(db *Database)NewSession() {
  sess := session.Must(session.NewSession(&aws.Config{
    Region: aws.String("us-east-2"),
  }))
  db.db = dynamodb.New(sess)
}

// PutItemInTable: Takes my DBItem, and a TableName. Places the DBItem in a
//     DynamoDB Table.
//  $ item      -> DBItem :: Any struct that implements DBTable
//  $ tableName -> string :: The Chosen name for the DynamoDB Table.
//
//  if successful
//      return nil
//  else
//      return Error
func(db *Database)PutItemInTable(item DBItem, tableName string) error {
  av, err := dynamodbattribute.MarshalMap(item)
  if err != nil {
    return err
  }
  input := &dynamodb.PutItemInput{
    Item:      av,
    TableName: aws.String(tableName),
  }

  if _, err = db.db.PutItem(input); err != nil {
    return err
  }
  return nil
}

// DeleteItemFromTable: Providing the KeyName and ID, attempts to Delete an Item
//     in the DynamoDB table, via the provided table name.
//  $ tableName    -> string :: The DynamoDB TableName
//  $ keyName      -> string :: The DynamoDB Table Key
//  $ id           -> string :: The search ID in tableName[KEY]
//  $ toBeReturned -> bool   :: If true, return the Deleted Value, if found
//
//  if successful, returns (*DBItem-or-nil), nil.
//  else we failed to delete the item, then we return nil, Error
func( db *Database )DeleteItemFromTable(
  tableName,
  keyName,
  id           string,
  toBeReturned bool,
)( *DBItem,error ){
  dbKey := map[string] *dynamodb.AttributeValue{
    keyName: {
      S: aws.String(id),
    },
  }
  var returnValue string
  if toBeReturned{
    returnValue = "ALL_OLD"
  } else {
    returnValue = "NONE"
  }

  deleteItemInput := &dynamodb.DeleteItemInput{
    TableName:    aws.String(tableName),
    Key:          dbKey,
    ReturnValues: &returnValue,
  }

  item, err := db.db.DeleteItem(deleteItemInput)
  if err != nil {
    log.Printf(" --> Failed to Delete Item from Database: %v", err.Error())
    return nil, err
  }
  if item.Attributes == nil && toBeReturned{
    log.Printf("Item to delete could not be returned")
  }else if item.Attributes == nil {
    return nil, nil
  }
  var dbItem DBItem
  if err = dynamodbattribute.UnmarshalMap(item.Attributes, &dbItem); err != nil {
    log.Printf(" --> Deleted Item, Failed to Unmarshal it into a DBItem")
    return nil, errors.New("Failed to Unmarshal Deleted Item")
  }

  return &dbItem, nil
}

// UpdateItemInTable: With the given DynamoDB Search parameters, and a DBItem,
//     we attempt to Find and update the DynamoDB Table Item.
//  $ tableName    -> string :: The DynamoDB TableName
//  $ keyName      -> string :: The DynamoDB Table Key
//  $ id           -> string :: The search ID in tableName[KEY]
//  $ item         -> DBItem :: The provided DBItem to Update tableName[KEY] with
//
//  if successful
//      return nil
//  else
//      return Error
func(db *Database)UpdateItemInTable(
  tableName,
  keyName,
  id        string,
  item      DBItem,
) error {
  dbKey := map[string]*dynamodb.AttributeValue{
    keyName: {
      S: aws.String(id),
    },
  }
  itemJson, err := json.Marshal(item)
  if err != nil {
    log.Printf(
      " --> Error while Marshaling DBItem for Table Updata: Error: %v",
    err.Error(),
  )
    return err
  }

  updateExpr := "SET #item = :itemValue"
  expressionAttirbuteValues := map[string]*dynamodb.AttributeValue{
    ":itemValue": {
      S: aws.String(string(itemJson)),
    },
  }
  expressionAttirbuteNames := map[string]*string{
    "#item": aws.String("Item"),
  }

  input := &dynamodb.UpdateItemInput{
    TableName:        aws.String(tableName),
    Key:              dbKey,
    UpdateExpression: aws.String(updateExpr),
    ExpressionAttributeValues: expressionAttirbuteValues,
    ExpressionAttributeNames:  expressionAttirbuteNames,
  }

  if _, err = db.db.UpdateItem(input); err != nil{
    return err
  }

  return nil
}

// GetItemFromTable: With the provided DynamoDB search parameters, we attempt to
//     search for and return the corresponding DynamoDB Table Item.
//  $ tableName    -> string :: The DynamoDB TableName
//  $ keyName      -> string :: The DynamoDB Table Key
//  $ id           -> string :: The search ID in tableName[KEY]
//
//  if successful
//     return an *DBTable, nil
//  else
//     return nil, Error
func( db *Database )GetItemFromTable(
  tableName,
  keyName,
  id        string,
)( *DBItem, error ){
  dbKey := map[string]*dynamodb.AttributeValue{
    keyName: {
      S: aws.String(id),
    },
  }

  input := &dynamodb.GetItemInput{
    TableName: aws.String(tableName),
    Key:       dbKey,
  }

  result, err := db.db.GetItem(input)
  if err != nil {
    return nil, err
  }
  if result.Item == nil {
    return nil, nil
  }
  var item DBItem
  if err = dynamodbattribute.UnmarshalMap(result.Item, &item); err != nil {
    log.Printf(" --> Failed to Unmarshal Record: Error: %v", err.Error())
    return nil, err
  }

  return &item, nil
}
