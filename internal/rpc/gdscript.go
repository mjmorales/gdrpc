package rpc

import (
	_ "embed"
	"fmt"
	"os"
	"strings"
)

const (
	FileNameRPCSpec         = "rpcSpec.json"
	FileNameGDScriptService = "RpcService.gd"
)

var (
	//go:embed templates/rpcService.gd
	rpcService string
)

func GenerateSuperClasses(outputDir string) error {
	fileName := fmt.Sprintf("%s/%s", outputDir, FileNameGDScriptService)
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(rpcService)
	if err != nil {
		return err
	}

	return nil
}

func GenerateGDScript(spec *RPCSpec, outputDir string) error {
	fileName := fmt.Sprintf("%s.gd", spec.Service)
	filePath := fmt.Sprintf("%s/%s", outputDir, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	s := strings.Builder{}
	s.WriteString("extends RpcService")
	s.WriteString("\n\n")
	s.WriteString(fmt.Sprintf("class_name %sInterface", spec.Service))
	s.WriteString("\n\n")

	for _, method := range spec.Methods {
		paramsList := make([]string, 0)
		messageParamsList := make([]string, 0)
		for _, param := range method.Parameters {
			paramsList = append(paramsList, fmt.Sprintf("%s: %s", param.Name, param.Type))
			messageParamsList = append(messageParamsList, fmt.Sprintf("\"%s\": %s", param.Name, param.Name))
		}

		s.WriteString(
			fmt.Sprintf(
				"func %s(%s):\n",
				method.Name,
				strings.Join(paramsList, ", "),
			),
		)

		s.WriteString(
			fmt.Sprintf(
				"    var message = {\n        \"method\": \"%s\",\n        \"params\": {%s}\n    }\n",
				method.Name,
				strings.Join(messageParamsList, ", "),
			),
		)

		s.WriteString(
			fmt.Sprintf(
				"    send_data(message)\n",
			),
		)

		s.WriteRune('\n')
	}

	_, err = file.WriteString(s.String())
	if err != nil {
		return err
	}

	return nil
}
