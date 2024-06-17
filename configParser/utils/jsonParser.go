package utils

import (
	"encoding/json"

	cp "github.com/nassorc/gandalf/configParser"
)

func NewJsonConfigParser() cp.IConfigParser {
	return &JsonConfigParser{}
}

type JsonConfigParser struct{}

func (p JsonConfigParser) ParseConfig(data []byte) (*cp.Config, error) {
	var out = cp.Config{}
	err := json.Unmarshal(data, &out)

	if err != nil {
		return &out, err
	}

	return &out, nil
}
