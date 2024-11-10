
# jamgoline

## Overview
This project demonstrates a PubSub system with encrypted messaging, a templated Agent system, and a command-line interface (CLI) for managing agents. 
Agents can be created from pre-defined templates or custom templates, allowing flexible message handling across various topics.

## Project Structure
```
jamgoline/
├── cmd/
│   └── jamgoline/
│       └── main.go                # Executável principal com interface CLI
├── config/
│   └── agents_config.yaml         # Arquivo de configuração para definir o pipeline de execução dos agentes
├── pkg/
│   ├── agent/
│   │   ├── agent.go               # Lógica dos agentes
│   │   ├── agent_template.go      # Gerenciamento e registro de templates de agentes
│   │   ├── agent_factory.go       # Lógica de criação de agentes a partir de templates
│   │   └── agent_test.go
│   ├── pubsub/
│   │   ├── pubsub.go
│   │   └── pubsub_test.go
│   └── crypto/
│       ├── crypto.go
│       └── crypto_test.go
└── internal/
    ├── config/
    │   └── config.go              # Lógica de carregamento e parsing de configurações
    └── logger/
        └── logger.go              # Sistema de log centralizado



```

## Features
- **Agent Templates**: Register agents with unique configurations, topics, and actions.
- **Encrypted PubSub**: Uses AES encryption for secure message transmission in the PubSub channels.
- **Configurable Pipeline**: Define agents and their sequence of execution in `agents_config.yaml`.
- **CLI Interface**: Manage templates and execute the pipeline via terminal.

## Usage

### 1. Running the Project
To compile and run the project, navigate to `cmd/jamgoline` and execute:
```
go run main.go <command>
```

### 2. CLI Commands
- `list`: Lists available agent templates.
- `run`: Executes the configured agent pipeline as specified in `agents_config.yaml`.

### 3. Configuration
Edit `config/agents_config.yaml` to add or modify agents and define the pipeline execution order.

### Adding Custom Templates
To add a custom template:
1. Modify `agent_template.go` to register the new template.
2. Use the template name in `agents_config.yaml` to configure custom agents.

## Example
The sample configuration in `agents_config.yaml` includes two agents, `CustomLoggerAgent` and `CustomAnalyticsAgent`, subscribed to different topics. The pipeline specifies that these agents should run in a defined order.

## Requirements
- Go 1.16 or later
- YAML library (`gopkg.in/yaml.v2`)

## License
This project is open-source and available under the MIT License.
