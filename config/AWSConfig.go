package config

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func InitAwsS3Client(conf *Config) (*s3.Client, error) {
    creds := credentials.NewStaticCredentialsProvider(conf.AWS.AccessKeyID, conf.AWS.SecretKey, "")
    cfg, err := config.LoadDefaultConfig(
        context.TODO(), 
        config.WithCredentialsProvider(creds),
        config.WithRegion(conf.AWS.Region),
    )
    if err != nil { return nil, err }

    return s3.NewFromConfig(cfg), nil
}
