package envconfig

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

const (
	awsLocalConfigProfileNameVar = "AWS_LOCAL_CONFIG_PROFILE_NAME"
	awsDefaultRegionVar          = "AWS_REGION" // Used if the Secrets region is not explicitly set
	awsSecretsRegionVar          = "AWS_SECRETS_REGION"
	awsSecretManagerArnPrefix    = "arn:aws:secretsmanager:"
)

var (
	awsLocalConfigProfileName string
	awsGenCtx                 = context.Background()
	awsSecretsManagerClient   *secretsmanager.Client
)
