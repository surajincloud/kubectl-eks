package kube

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/surajincloud/kubectl-eks/pkg/kube"
)

func GetAWSConfig(ctx context.Context, region string) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return aws.Config{}, err
	}
	if cfg.Region == "" {
		// get region
		region, err = kube.GetRegion(region)
		if err != nil {
			return aws.Config{}, err
		}
		cfg.Region = region
	}
	return cfg, nil
}
