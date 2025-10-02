package envconfig

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// AWS Environment Variables in the scope below must be defined in the Config prior to any variables that use AWS Secrets,
// except AWS_REGION that will be pulled from the OS Env directly if not present in Config

const (
	awsLocalConfigProfileNameVar = "AWS_LOCAL_CONFIG_PROFILE_NAME" // Optional, used primarily for Devstack type of setup

	// AWS_REGION is normally present in the Env of ECS Tasks at runtime
	awsDefaultRegionVar = "AWS_REGION"         // Used if the Secrets region is not explicitly set
	awsSecretsRegionVar = "AWS_SECRETS_REGION" // Required if the AWS Default region is not set in the Env or is different from the Secrets region

	awsSecretManagerArnPrefix = "arn:aws:secretsmanager:"
)

var (
	awsLocalConfigProfileName string
	awsSecretsRegion          string
	awsGenCtx                 = context.Background()
	awsSecretsManagerClient   *secretsmanager.Client
)

func getAwsSecretsManagerClient() (*secretsmanager.Client, error) {
	if awsSecretsManagerClient == nil {
		// TODO: init the AWS Config and the Secrets Manager Client
		// If failed to init, return a descriptive error
	}
	return awsSecretsManagerClient, nil
}
