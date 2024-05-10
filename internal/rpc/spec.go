package rpc

import (
	"encoding/json"
	"os"
)

type Parameter struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Method struct {
	Name       string      `json:"name"`
	Parameters []Parameter `json:"parameters"`
	Response   string      `json:"response"`
}

type RPCSpec struct {
	Service string   `json:"service"`
	Methods []Method `json:"methods"`
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
