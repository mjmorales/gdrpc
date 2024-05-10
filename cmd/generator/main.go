package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/mjmorales/gdrpc/internal/rpc"
	"github.com/mjmorales/gdrpc/internal/rpc/gdscript"
)

const (
	GenerateRpcFromGolang = iota
	GenerateGdscriptFromRpc
)

func generateRpcSpec(dirPath, outputPath string) error {
	spec, err := rpc.ParseDir(dirPath)
	if err != nil {
		return fmt.Errorf("error parsing directory: %w", err)
	}

	data, err := json.MarshalIndent(spec, "", "    ")
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	return nil
}

func generateGdscript(specPath, outputPath string) error {
	gdscript.GenerateSuperClasses(outputPath)

	var spec rpc.RPCSpec
	jsonFile, err := os.Open(specPath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}

	if err := json.NewDecoder(jsonFile).Decode(&spec); err != nil {
		return fmt.Errorf("error decoding JSON: %w", err)
	}

	for _, service := range spec.Services {
		err = gdscript.GenerateService(service, outputPath)
		if err != nil {
			return fmt.Errorf("error generating GDScript: %w", err)
		}
	}

	return nil
}

func main() {
	genType := flag.Int("type", GenerateRpcFromGolang, "Type of generation to perform")
	inputDir := flag.String("path", ".", "Path to the Go project directory")
	outputDir := flag.String("output", "rpc_spec.json", "Output file path for the RPC spec")
	flag.Parse()

	switch *genType {
	case GenerateRpcFromGolang:
		if err := generateRpcSpec(*inputDir, *outputDir); err != nil {
			fmt.Printf("Error generating RPC spec: %v\n", err)
			os.Exit(1)
		}
	case GenerateGdscriptFromRpc:
		if err := generateGdscript(*inputDir, *outputDir); err != nil {
			fmt.Printf("Error generating GDScript: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Println("Invalid generation type specified.")
		os.Exit(1)
	}

	fmt.Println("RPC spec generated successfully.")
}
