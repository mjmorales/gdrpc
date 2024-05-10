package gdscript

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/mjmorales/gdrpc/internal/rpc"
)

const (
	FileNameRPCSpec         = "rpcSpec.json"
	FileNameGDScriptService = "RpcService.gd"
)

var (
	//go:embed templates/rpcService.gd
	rpcService string
	//go:embed templates/rpcResource.gd
	rpcResource string
)

func GenerateSuperClasses(outputDir string) error {
	rpcSuperClasses := map[string]string{
		"RpcService":  rpcService,
		"RpcResource": rpcResource,
	}

	for className, superClass := range rpcSuperClasses {
		fileName := fmt.Sprintf("%s/%s.gd", outputDir, className)
		file, err := os.Create(fileName)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = file.WriteString(superClass)
		if err != nil {
			return err
		}
	}

	return nil
}

func GenerateService(spec rpc.ServiceSpec, outputDir string) error {
	fileName := fmt.Sprintf("%s.gd", spec.Name)
	filePath := fmt.Sprintf("%s/%s", outputDir, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	s := strings.Builder{}
	s.WriteString("extends RpcService")
	s.WriteString("\n\n")
	s.WriteString(fmt.Sprintf("class_name %sInterface", spec.Name))
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
