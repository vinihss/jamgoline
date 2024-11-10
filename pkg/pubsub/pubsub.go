package pubsub

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"sync"
)

// PubSub manages encrypted messaging channels and subscriptions.
type PubSub struct {
	globalChannel chan string
	subscribers   map[string][]chan string
	secretKey     []byte
	mu            sync.RWMutex
}

// NewPubSub initializes a new PubSub instance with AES encryption for messaging.
func NewPubSub(secretKey []byte) (*PubSub, error) {
	if len(secretKey) != 16 && len(secretKey) != 24 && len(secretKey) != 32 {
		return nil, errors.New("AES key must be 16, 24, or 32 bytes")
	}
	return &PubSub{
		globalChannel: make(chan string, 10),
		subscribers:   make(map[string][]chan string),
		secretKey:     secretKey,
	}, nil
}

// SubscribeGlobal returns a read-only channel for global message subscription.
func (ps *PubSub) SubscribeGlobal() <-chan string {
	return ps.globalChannel
}

// SubscribeTopic returns a new channel for subscribing to a specific topic.
func (ps *PubSub) SubscribeTopic(topic string) <-chan string {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ch := make(chan string, 10)
	ps.subscribers[topic] = append(ps.subscribers[topic], ch)
	return ch
}

// PublishGlobal encrypts and sends a message to the global channel.
func (ps *PubSub) PublishGlobal(msg string) error {
	encryptedMsg, err := ps.encryptMessage([]byte(msg))
	if err != nil {
		return err
	}
	go func() {
		ps.globalChannel <- encryptedMsg
	}()
	return nil
}

// PublishTopic encrypts and sends a message to a specific topic.
func (ps *PubSub) PublishTopic(topic string, msg string) error {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	encryptedMsg, err := ps.encryptMessage([]byte(msg))
	if err != nil {
		return err
	}

	for _, ch := range ps.subscribers[topic] {
		go func(c chan string) {
			c <- encryptedMsg
		}(ch)
	}
	return nil
}

// encryptMessage uses AES encryption to encrypt messages.
func (ps *PubSub) encryptMessage(message []byte) (string, error) {
	block, err := aes.NewCipher(ps.secretKey)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(message))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], message)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// decryptMessage uses AES to decrypt encrypted messages.
func (ps *PubSub) decryptMessage(encryptedMessage string) ([]byte, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(encryptedMessage)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(ps.secretKey)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}
