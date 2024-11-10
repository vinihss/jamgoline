package config

import (
	"io/ioutil"
)

// Config define a estrutura do arquivo YAML de configuração
type Config struct {
	Agents   []AgentConfig  `yaml:"agents"`
	Pipeline PipelineConfig `yaml:"pipeline"`
}

type AgentConfig struct {
	Template string   `yaml:"template"`
	Name     string   `yaml:"name"`
	Topics   []string `yaml:"topics"`
}

type PipelineConfig struct {
	Sequence []string `yaml:"sequence"`
}

// LoadConfig carrega as configurações do arquivo YAML
func LoadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
