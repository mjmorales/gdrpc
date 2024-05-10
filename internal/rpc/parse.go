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

func ParseDir(path string) (RPCSpec, error) {
	rpcSpec := RPCSpec{}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
	if err != nil {
		return rpcSpec, err
	}

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

					isService := slices.ContainsFunc(genDecl.Doc.List, func(c *ast.Comment) bool {
						return strings.Contains(c.Text, DocTagService)
					})

					if isService {
						interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
						if !ok {
							continue
						}

						serviceSpec := &ServiceSpec{Name: typeSpec.Name.Name}
						serviceSpec, err = serviceSpec.AddMethod(typeSpec, interfaceType)
						if err != nil {
							return rpcSpec, err
						}

						rpcSpec.Services = append(rpcSpec.Services, *serviceSpec)
					}
				}
			}
		}
	}

	return rpcSpec, nil
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
