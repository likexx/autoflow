package autoflow

import (
	"reflect"
	"log"
)

type ServerAction struct {

}

func (s *ServerAction) Verify(clientParams map[string]interface{}, params ActionParameter) bool {
	clientResult := clientParams["result"].(string)
	expect := params["expect"]

	log.Printf("[Server action] client result: %s, expect: %s", clientResult, expect)

	return clientResult == expect
}

func invokeServerAction(actionName string, clientParams map[string]interface{}, serverParams ActionParameter) bool {
	a := &ServerAction{}
	obj := reflect.ValueOf(a)

	log.Printf("invoke server action: %s\n", actionName)
	t1 := obj.MethodByName(actionName)
	params := make([]reflect.Value, 2)
	params[0] = reflect.ValueOf(clientParams)
	params[1] = reflect.ValueOf(serverParams)
	result := t1.Call(params)

	return result[0].Interface().(bool)
}