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

// todoGetHabler: AWS Lambda function for Querying and Returning a request Todo.
//   $ context -> context.Context
//   $ request -> DatabaseRequest, contains the body data from the Client.
//
//  if successful
//      return DatabaseResponse(200).Body(todoJSON), nil
//  else
//      return DatabaseResponse(400), Error
func todoGetHandler(
  context context.Context,
  request DatabaseRequest,
)( DatabaseResponse, error ){
  resp := Response{}
  db := Database{}
  db.NewSession()

  todoSearch := struct{
    keyName string
    id      string
  }{}
  if err := json.Unmarshal([]byte(request.Body), &todoSearch); err != nil {
    return resp.
      AddMessage(
        fmt.Sprintf("Failed to Unmarshal Request Body\nError: %v", err.Error()),
      ).
      Respond(400)
  }
  todo, err := db.GetItemFromTable(
    TODOTABLENAME,
    todoSearch.keyName,
    todoSearch.id,
  )
  if err != nil {
    return resp.
      AddMessage(
        fmt.Sprintf(
          "Failed to Find %v in Database Table\nError: %v",
          todoSearch.id, err.Error()),
      ).
      Respond(400)
  }
  todoJson, err := json.Marshal(todo)
  if err != nil {
    return resp.
      AddMessage(fmt.Sprintf("Failed to Marshal Todo\nError: %v", err)).
      Respond(400)
  }

  return resp.
    AddMessage("Successfully Obtained Todo Item").
    AddBody(todoJson).
    Respond(200)
}

func main(){
  lambda.Start(todoGetHandler)
}
