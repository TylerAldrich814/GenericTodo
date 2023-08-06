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

// todoCreateHabler: AWS Lambda function for Handling the creation of a new Todo.
//   $ context -> context.Context
//   $ request -> DatabaseRequest, contains the body data from the Client.
//
//  if successful
//      return DatabaseResponse(200), nil
//  else
//      return DatabaseResponse(400), Error
func todoCreateHandler(
  context context.Context,
  request DatabaseRequest,
)( DatabaseResponse, error ){
  resp := Response{}
  db := Database{}
  db.NewSession()

  var newTodo Todo
  if err := json.Unmarshal([]byte(request.Body), &newTodo); err != nil {
    return resp.
      AddMessage(
        fmt.Sprintf("Failed to Unmarshal new Todo from Request Body: %v", err),
      ).
      Respond(400)
  }
  if err := db.PutItemInTable(newTodo, TODOTABLENAME); err != nil {
    resp.Message = fmt.Sprintf(
      " --> Failed to place to Todo in the %v Table:\n --> Error: %v",
      TODOTABLENAME,
      err.Error(),
    )
    resp.Respond(400)
  }
  return resp.Respond(200)
}

func main(){
  lambda.Start(todoCreateHandler)
}
