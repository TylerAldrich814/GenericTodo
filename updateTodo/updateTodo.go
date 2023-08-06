package main

import (
	"context"
	"encoding/json"
	"fmt"
	"todoDatabase"
	"utils"

	"github.com/aws/aws-lambda-go/lambda"
)
type Database           = database.Database
type Todo               = utils.Todo
type Response           = utils.Response
type DatabaseResponse   = utils.DatabaseResponse
type DatabaseRequest    = utils.DatabaseRequest
const TODOTABLENAME     = utils.TODOTABLENAME

// todoUpdateHabler: AWS Lambda function for Handling the Update of a
//                   particular Todo. Will fail if todo doesn't exist.
//   $ context -> context.Context
//   $ request -> DatabaseRequest, contains the body data from the Client.
//
//  if successful
//      return DatabaseResponse(200), nil
//  else
//      return DatabaseResponse(400), Error
func todoUpdateHandler(
  context context.Context,
  request DatabaseRequest,
)( DatabaseResponse, error ){
  resp := Response{}
  db := Database{}
  db.NewSession()

  updateTodo := struct{
    keyName   string
    id        string
    todo      Todo
  }{}

  if err := json.Unmarshal([]byte(request.Body), &updateTodo); err != nil {
    return resp.
      AddMessage(
        fmt.Sprintf("Failed to Unmarshal Request Body\nError: %v", err),
      ).
      Respond(400)
  }
  if err := db.UpdateItemInTable(
    TODOTABLENAME,
    updateTodo.keyName,
    updateTodo.id,
    updateTodo.todo,
  ); err != nil {
    return resp.
      AddMessage(
        fmt.Sprintf("Failed to Update Database with '%v'\nError: %v",
          updateTodo.todo.Title, err,
        ),
      ).
      Respond(400)
  }

  return resp.
    AddMessage("Successfully Updated Todo").
    Respond(200)
}

func UpdateTodo(){
  lambda.Start(todoUpdateHandler)
}
