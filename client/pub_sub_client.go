package client

import (
	"github.com/hirakiuc/jobworker-go/models"
)

// PubSubClient describe a CloudPubSub client to fetch messages.
type PubSubClient struct {
}

// NewPubSubClient return a PubSubClient instance.
func NewPubSubClient() (*PubSubClient, error) {
	return &PubSubClient{}, nil
}

// GetMessage return the new message which received from CloudPubSub
func (f *PubSubClient) GetMessage() (*models.Message, error) {
	return nil, nil
}
