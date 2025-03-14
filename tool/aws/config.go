package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

type IConfig interface {
	GetInstance() aws.Config
}

type Config struct {
	awsConfig aws.Config
}

func NewConfig(ctx context.Context, cfg ConfigProvider) (*Config, error) {
	awsConfig, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(cfg.GetAwsRegion()),
		config.WithCredentialsProvider(aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(cfg.GetAwsAccessKey(), cfg.GetAwsSecretKey(), ""))),
	)
	if err != nil {
		return nil, err
	}

	return &Config{awsConfig: awsConfig}, nil
}

func (c *Config) GetInstance() aws.Config {
	return c.awsConfig
}
