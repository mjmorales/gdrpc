package rpc

import (
	"encoding/json"
	"reflect"
)

func GenerateRPCSpec(svc interface{}) ([]byte, error) {
	svcType := reflect.TypeOf(svc)
	methods := make([]map[string]interface{}, 0)

	for i := 0; i < svcType.NumMethod(); i++ {
		method := svcType.Method(i)
		if tag, ok := method.Type.Field(0).Tag.Lookup("rpc"); ok && tag == "method" {
			methodName := method.Name
			methodDetails := map[string]interface{}{
				"name":       methodName,
				"parameters": []string{},
				"returns":    []string{},
			}

			for j := 0; j < method.Type.NumIn(); j++ {
				paramType := method.Type.In(j)
				methodDetails["parameters"] = append(methodDetails["parameters"].([]string), paramType.Name())
			}

			for k := 0; k < method.Type.NumOut(); k++ {
				returnType := method.Type.Out(k)
				methodDetails["returns"] = append(methodDetails["returns"].([]string), returnType.Name())
			}

			methods = append(methods, methodDetails)
		}
	}

	return json.MarshalIndent(map[string]interface{}{"methods": methods}, "", "    ")
}
