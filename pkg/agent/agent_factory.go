package agent

// NewAgentFromTemplate cria um novo agente com base em um template registrado
func NewAgentFromTemplate(templateName, customName string, topics []string, pubsub *PubSub) (*Agent, error) {
	template, err := GetTemplate(templateName)
	if err != nil {
		return nil, err
	}
	return NewAgent(customName, template.Action, pubsub, topics...), nil
}
