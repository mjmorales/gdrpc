package rpc

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"slices"
	"strings"
)

const (
	DocTagService = "@rpc:service"
)

type MethodSpec struct {
	Name       string      `json:"name"`
	Parameters []Parameter `json:"parameters"`
	Returns    []string    `json:"returns"`
}

type ServiceSpec struct {
	Name    string       `json:"service"`
	Methods []MethodSpec `json:"methods"`
}

func ParseDir(path string) ([]ServiceSpec, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var specs []ServiceSpec
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if !ok || genDecl.Doc == nil {
					continue
				}

				for _, spec := range genDecl.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}

					interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
					if !ok {
						continue
					}

					isService := slices.ContainsFunc(genDecl.Doc.List, func(c *ast.Comment) bool {
						return strings.Contains(c.Text, DocTagService)
					})

					if !isService {
						continue
					}

					service := ServiceSpec{Name: typeSpec.Name.Name}
					for _, method := range interfaceType.Methods.List {
						methodSpec := MethodSpec{Name: method.Names[0].Name}
						funcType, ok := method.Type.(*ast.FuncType)
						if !ok {
							continue
						}

						methodSpec.Parameters = ParseParameters(funcType.Params)
						methodSpec.Returns = ParseResults(funcType.Results)
						service.Methods = append(service.Methods, methodSpec)
					}
					if len(service.Methods) > 0 {
						specs = append(specs, service)
					}
				}
			}
		}
	}
	return specs, nil
}

func ParseParameters(fields *ast.FieldList) []Parameter {
	if fields == nil {
		return nil
	}

	var params []Parameter
	for _, field := range fields.List {
		for _, name := range field.Names {
			params = append(params, Parameter{Name: name.Name, Type: fmt.Sprint(field.Type)})
		}
	}

	return params
}

func ParseResults(fields *ast.FieldList) []string {
	if fields == nil {
		return nil
	}
	var results []string
	for _, field := range fields.List {
		typeStr := fmt.Sprint(field.Type)
		if len(field.Names) > 0 {
			for _, name := range field.Names {
				results = append(results, fmt.Sprintf("%s %s", name, typeStr))
			}
		} else {
			results = append(results, typeStr)
		}
	}
	return results
}
