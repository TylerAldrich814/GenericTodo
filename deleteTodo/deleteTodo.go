package main

import (
	"context"
	"encoding/json"
	"fmt"
	"utils"

	"github.com/aws/aws-lambda-go/lambda"

	"todoDatabase"
)

type Database           = database.Database
type Todo               = utils.Todo
type Response           = utils.Response
type DatabaseRequest  = utils.DatabaseRequest
type DatabaseResponse = utils.DatabaseResponse
const TODOTABLENAME     = utils.TODOTABLENAME

// todoDeleteHabler: AWS Lambda function for Handling the Deletion of a new Todo.
//   $ context -> context.Context
//   $ request -> DatabaseRequest, contains the body data from the Client.
//
//  if successful && request.Body contains 'toBeReturned = True'
//      return DatabaseResponse(200).Body(DELETEDITEM), nil
//  else if successful:
//      return DatabaseResponse(200), nil
//  else
//      return DatabaseResponse(400), Error
func todoDeleteHandler(
  context context.Context,
  request DatabaseRequest,
)( DatabaseResponse, error ){
  resp := Response{}
  db   := Database{}
  db.NewSession()

  toBeDeleted := struct{
    tableName   string
    keyName     string
    id          string
    toBeReturned bool
  }{}
  if err := json.Unmarshal([]byte(request.Body), &toBeDeleted); err != nil {
    return resp.
      AddMessage(
        fmt.Sprintf(
          "Failed to Unmarshal 'toBeDeleted' Object from Request Body: %v",
        err.Error(),
        ),
      ).
      Respond(200)
  }
  item, err := db.DeleteItemFromTable(
    toBeDeleted.tableName,
    toBeDeleted.keyName,
    toBeDeleted.id,
    toBeDeleted.toBeReturned,
  )
  if err != nil {
    return resp.
      AddMessage(
        fmt.Sprintf(
          "Failed to Delete %v from the %v\n -Error: %v",
          toBeDeleted.id, toBeDeleted.tableName, err,
        ),
      ).Respond(400)
  }
  var itemJson []byte
  if toBeDeleted.toBeReturned {
    itemJson, err = json.Marshal(&item)
    if err != nil {
      return resp.AddMessage(
        fmt.Sprintf(
          "Failed to Marshal the returned Deleted Value\n - Error: %v",
          err,
        ),
      ).
      Respond(400)
    }
  }

  return resp.
    AddMessage("Successfully Deleted Todo from Database Table").
    AddBody(itemJson).
    Respond(200)
}

func DeleteTodo(){
  lambda.Start(todoDeleteHandler)
}
