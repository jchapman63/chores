package sns

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type SNSClient struct {
	Client *sns.Client
}

func NewSNSClient(ctx context.Context) (*SNSClient, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	// Create a new SNS client
	client := sns.NewFromConfig(cfg)

	return &SNSClient{Client: client}, nil
}
