package core

import (
	"encoding/json"
	"fmt"
)

func ErrorResponse(status int) MainResponseArgs {
	return MainResponseArgs{
		StatusCode: status,
	}
}

// NewLoggedError generates a new error and logs it to stdout
func NewLoggedError(format string, a ...interface{}) error {
	err := fmt.Errorf(format, a...)

	return err
}

func MainArgsFromMap(args *map[string]interface{}) MainRequestArgs {
	jsonString, _ := json.Marshal(args)
	s := MainRequestArgs{}
	json.Unmarshal(jsonString, &s)
	return s
}

func MainArgsToMap(args *MainResponseArgs) map[string]interface{} {
	var newMap map[string]interface{}

	data, _ := json.Marshal(args)
	json.Unmarshal(data, &newMap)

	return newMap
}
