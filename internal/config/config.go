package config

import (
    "gopkg.in/yaml.v2"
    "io/ioutil"
)

type Config struct {
    Agents   []AgentConfig   `yaml:"agents"`
    Pipeline PipelineConfig  `yaml:"pipeline"`
}

type AgentConfig struct {
    Template string   `yaml:"template"`
    Name     string   `yaml:"name"`
    Topics   []string `yaml:"topics"`
}

type PipelineConfig struct {
    Sequence []string `yaml:"sequence"`
}

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
