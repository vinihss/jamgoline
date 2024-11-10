package agent

import (
	"errors"
	"sync"
)

// AgentTemplate define um template padrão de agente
type AgentTemplate struct {
	Name   string
	Topics []string
	Action func(input interface{}) (output interface{}, err error)
}

// Registry para armazenar e gerenciar templates
var (
	templates      = make(map[string]AgentTemplate)
	templatesMutex sync.RWMutex
)

// RegisterTemplate registra um novo template de agente
func RegisterTemplate(name string, topics []string, action func(input interface{}) (interface{}, error)) error {
	templatesMutex.Lock()
	defer templatesMutex.Unlock()

	if _, exists := templates[name]; exists {
		return errors.New("template already exists")
	}

	templates[name] = AgentTemplate{
		Name:   name,
		Topics: topics,
		Action: action,
	}

	return nil
}

// GetTemplate busca um template registrado pelo nome
func GetTemplate(name string) (AgentTemplate, error) {
	templatesMutex.RLock()
	defer templatesMutex.RUnlock()

	template, exists := templates[name]
	if !exists {
		return AgentTemplate{}, errors.New("template not found")
	}

	return template, nil
}

// ListTemplates retorna a lista de templates disponíveis
func ListTemplates() []string {
	templatesMutex.RLock()
	defer templatesMutex.RUnlock()

	var names []string
	for name := range templates {
		names = append(names, name)
	}
	return names
}
