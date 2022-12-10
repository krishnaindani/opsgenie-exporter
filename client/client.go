package client

import (
	"github.com/opsgenie/opsgenie-go-sdk-v2/alert"
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
)

type (
	Config struct {
		APIKey string
	}

	Client struct {
		alerts *alert.Client
	}
)

func NewClient(config Config) (*Client, error) {

	alertClient, err := alert.NewClient(&client.Config{
		ApiKey: config.APIKey,
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		alerts: alertClient,
	}, nil

}
