package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Zones struct {
		Excludes []string `yaml:"excludes"`
	} `yaml:"zones"`

	ResourceTypes struct {
		Excludes []string `yaml:"excludes"`
	} `yaml:"resource-types"`

	ResourceIDs struct {
		Excludes []ResourceIDFilter `yaml:"excludes"`
	} `yaml:"resource-ids"`
}

type ResourceIDFilter struct {
	ResourceType string `yaml:"resourceType"`
	ID           string `yaml:"id"`
}

func LoadConfig(path string) (*Config, error) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling YAML: %w", err)
	}

	return &config, nil
}
