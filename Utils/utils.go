package utils

import (
	"encoding/json"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

type DatabaseResponse = events.APIGatewayProxyResponse
type DatabaseRequest = events.APIGatewayProxyRequest
const TODOTABLENAME = "Todos"

// Response: Struct for handling our events.APIGatewayProxyResponse's.
//    Every Response will contain A Message and a StatusCode
//  $ Message -> string :: Our message to the Client about the outcome of their
//                         request.
//  $ Body    -> []byte :: Only supplied when the Client is expecting Data within
//                         their request Response.
//  $ Code    -> int    :: The StatusCode of their Response Outcome. Used
//                         client-side for how to handle their Reqponse.
//  Usage:
//    var resp Response()
//    if successful{
//      resp.Message = fmt.Spintf("Failed to proccess Request: Error: %v", err)
//      return resp.Respond(400)
//    } else {
//      resp.Message = "Request Was successful"
//      resp.Body = dataBytes
//      resp.Respond(200)
//    }
type Response struct {
  Message string
  Body    []byte
  Code    int
}

func(resp *Response)AddMessage(msg string) *Response{
  resp.Message = msg

  return resp
}
func(resp *Response)AddBody(body []byte) *Response{
  resp.Body = body
  return resp
}

// Respond:: Abstracting away handling Request's. Geared towards AWS Lambda's
//           which returns client.APIGatewayProxyResponse, error. Here, we take
//           the provided Response structure, and we Marshal to Respond struct
//           variables, and turn it into a client.APIGatewayProxyResponse.
//  $ statusCode -> int :: Provided by a downstream function, which tells the
//                         client whether or not their Request passed or faied.
//  if successful
//    returns DatabaseResponse(statusCode), nil
//  else
//    returns DatabaseRequest(400), Error
func(resp *Response)Respond(statusCode int)( DatabaseResponse, error ){
  jsonResponse, err := json.Marshal(resp)
  if err != nil {
    return DatabaseResponse{
      StatusCode: 400,
      Body: "Failed to Marshal Response Error.",
    }, err
  }

  return DatabaseResponse{
    StatusCode: statusCode,
    Body: string(jsonResponse),
  }, nil
}

// DBItem: An Interface for Allowing our Database functions to take in and process
//         any struct that of which implements this interface.
//  $ GetID() string :: Every Database table requires an ID, here we simply return
//                      that value.
type DBItem interface {
  GetID() string
}

// Todo: Our Main datastructure for our Generic Todo List Applciation.
//
//   $ ID      -> string    :: For both sorting our Todo's, and used as the key
//                             in our Database.
//   $ Title   -> string    :: The headline for our Todo.
//   $ Body    -> string    :: The Information for our Tofo.
//   $ DueDate -> time.Time :: When our Todo is due.
type Todo struct {
  ID      string       `json:"id"`
  Title   string    `json:"title"`
  Body    string    `json:"body"`
  DueDate time.Time `json:"due_date"`
}

func( t Todo )GetID() string {
  return t.ID
}

// MarshalJson: Here, we are implementing the json.Marshaler interface, that wau
//              we can Marshal our DueDate for our Todo List.
// returns []byte, error
func(todo *Todo)MarshalJson()( []byte, error ){
  type Alias Todo

  return json.Marshal(&struct {
    DueDate string `json:"due_date"`
    *Alias
  }{
    DueDate: todo.DueDate.Format("2020-01-01 15:15:15"),
    Alias:   (*Alias)(todo),
  })
}

// MarshalJson: Here, we are implementing the json.Marshaler interface, that wau
//              we can Marshal our DueDate for our Todo List.
// $ data -> []byte :: Our incomming Marshaled data, to be converted back into
//                     a Todo Object.
// returns error, if unsuccssul
func(todo *Todo)UnmarshalJson(data []byte) error {
  type Alias Todo

  aux := &struct{
    DueDate string `json:"due_date"`
    *Alias
  }{
    Alias: (*Alias)(todo),
  }
  if err := json.Unmarshal(data, &aux); err != nil {
    return err
  }
  parsedTime, err := time.Parse("2020-01-01 15:15:15", aux.DueDate)
  if err != nil {
    return err
  }

  todo.DueDate = parsedTime
  return nil
}

