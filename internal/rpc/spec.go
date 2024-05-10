package rpc

import (
	"encoding/json"
	"go/ast"
	"os"
)

type Parameter struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type MethodSpec struct {
	Name       string      `json:"name"`
	Parameters []Parameter `json:"parameters"`
	Returns    []string    `json:"returns"`
}

type Field struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type ResourceSpec struct {
	Name   string  `json:"name"`
	Fields []Field `json:"fields"`
}

type ServiceSpec struct {
	Name    string       `json:"service"`
	Methods []MethodSpec `json:"methods"`
}

func (s *ServiceSpec) AddMethod(typeSpec *ast.TypeSpec, interfaceType *ast.InterfaceType) (*ServiceSpec, error) {
	for _, method := range interfaceType.Methods.List {
		methodSpec := MethodSpec{Name: method.Names[0].Name}
		funcType, ok := method.Type.(*ast.FuncType)
		if !ok {
			continue
		}

		methodSpec.Parameters = ParseParameters(funcType.Params)
		methodSpec.Returns = ParseResults(funcType.Results)
		s.Methods = append(s.Methods, methodSpec)
	}

	return s, nil
}

type RPCSpec struct {
	Services   []ServiceSpec  `json:"services"`
	Ressources []ResourceSpec `json:"resources"`
}

func ReadSpec(filename string) (*RPCSpec, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var spec RPCSpec
	if err := json.Unmarshal(data, &spec); err != nil {
		return nil, err
	}

	return &spec, nil
}
