package teams

import (
	"context"
	client2 "github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/team"
)

type (
	Config struct {
		ApiKey string
	}

	Client struct {
		client *team.Client
	}
)

func NewClient(config Config) (*Client, error) {

	client, err := team.NewClient(&client2.Config{
		ApiKey: config.ApiKey,
	})

	if err != nil {
		return nil, err
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) GetCountOfAllTeams() (float64, error) {

	res, err := c.client.List(context.Background(), &team.ListTeamRequest{})
	if err != nil {
		return 0, nil
	}

	return float64(len(res.Teams)), nil
}
