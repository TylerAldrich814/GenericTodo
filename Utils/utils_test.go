package utils

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)



func TestMarshalingAndUnmarshaling(t *testing.T){
  tests := struct{
    todos []Todo
  }{
    todos: []Todo{
      Todo {
          ID:      "1",
          Title:   "Comeplete the Project",
          DueDate: time.Date(2023, time.September, 15, 12, 0, 0, 0, time.UTC),
        },
      Todo {
          ID:      "2",
          Title:   "Land your first Backend Engineering Position",
          DueDate: time.Date(2023, time.September, 20, 11, 20, 10, 0, time.UTC),
        },
      Todo {
          ID:      "3",
          Title:   "Take over the workd",
          DueDate: time.Date(2023, time.December, 12, 23, 58, 59, 59, time.UTC),
        },
    },
  }
  for _, tt := range tests.todos {
    jsonData, err := json.Marshal(&tt)
    if err != nil {
      t.Errorf("Failed to Marshal Todo: %v", err)
      continue
    }
    var got Todo
    err = json.Unmarshal(jsonData, &got)
    if err != nil {
      t.Errorf("Failed to Unmarshal Todo: %v", err)
      continue
    }
    if !reflect.DeepEqual(got, tt){
      t.Errorf("Failed: Got '%v' But extected '%v'", got, tt)
    }
  }
}
