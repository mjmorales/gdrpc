package rpc

import (
	"encoding/json"
	"go/ast"
	"os"
	"reflect"
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

func (s *ResourceSpec) AddFields(typeSpec *ast.TypeSpec, structType *ast.StructType) (*ResourceSpec, error) {
	for _, field := range structType.Fields.List {
		// Skip unexported fields (lowercase)
		if field.Names[0].Name[0] < 65 || field.Names[0].Name[0] > 90 {
			continue
		}

		tag := reflect.StructTag(field.Tag.Value)
		typeName, tagExists := tag.Lookup("rpcType")
		if !tagExists {
			typeName = field.Type.(*ast.Ident).Name
		}

		fieldSpec := Field{Name: field.Names[0].Name}
		fieldSpec.Type = typeName
		s.Fields = append(s.Fields, fieldSpec)
	}

	return s, nil
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
	Services  []ServiceSpec  `json:"services"`
	Resources []ResourceSpec `json:"resources"`
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
