package agent

import (
	"context"
	"fmt"
)

type Agent struct {
	name      string
	action    func(input interface{}) (output interface{}, err error)
	pubsub    *PubSub // Reference to the PubSub instance
	globalSub <-chan string
	topicSubs map[string]<-chan string
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewAgent(name string, action func(input interface{}) (interface{}, error), pubsub *PubSub, topics ...string) *Agent {
	ctx, cancel := context.WithCancel(context.Background())
	agent := &Agent{
		name:      name,
		action:    action,
		pubsub:    pubsub,
		globalSub: pubsub.SubscribeGlobal(),
		topicSubs: make(map[string]<-chan string),
		ctx:       ctx,
		cancel:    cancel,
	}
	for _, topic := range topics {
		agent.topicSubs[topic] = pubsub.SubscribeTopic(topic)
	}
	return agent
}

func (a *Agent) Run() {
	for {
		select {
		case <-a.ctx.Done():
			return
		case encryptedMsg := <-a.globalSub:
			a.processMessage(encryptedMsg)
		default:
			for _, sub := range a.topicSubs {
				select {
				case encryptedMsg := <-sub:
					a.processMessage(encryptedMsg)
				default:
				}
			}
		}
	}
}

func (a *Agent) processMessage(msg interface{}) {
	_, err := a.action(msg)
	if err != nil {
		fmt.Printf("Error processing message in %s: %v", a.name, err)
	}
}
