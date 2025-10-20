package envconfig

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"log"
	"strings"
	"time"
)

// AWS Environment Variables in the scope below must be defined in the Config prior to any variables that use AWS Secrets,
// except AWS_REGION that will be pulled from the OS Env directly if not present in Config

const (
	awsLocalConfigProfileNameVar = "AWS_LOCAL_CONFIG_PROFILE_NAME" // Optional, used primarily for Devstack type of setup

	// AWS_REGION is normally present in the Env of ECS Tasks at runtime
	awsDefaultRegionVar = "DEFAULT_AWS_REGION" // Used if the Secrets region is not explicitly set
	awsSecretsRegionVar = "AWS_SECRETS_REGION" // Required  if the AWS Default region is not set in the Env or is different from the Secrets region

	awsSecretManagerArnPrefix = "arn:aws:secretsmanager:"

	awsRetryBackoffInitialDelaySec = 1
	awsRetryMaxAttempts            = 5
	awsRetryMaxBackoffDelaySec     = 30
)

var (
	awsLocalConfigProfileName string
	awsSecretsRegion          string
	awsGenCtx                 = context.Background()
	awsSecretsManagerClient   *secretsmanager.Client
)

func GetAwsSecretsManagerClient() (*secretsmanager.Client, error) {
	if awsSecretsManagerClient == nil {
		// TODO: init the AWS Config and the Secrets Manager Client
		// If failed to init, return a descriptive error
		awsCfg, err := getAwsConfig(awsGenCtx, awsSecretsRegion, awsLocalConfigProfileName)
		if err != nil {
			log.Printf("Error getting aws config: %v", err)
			return nil, err
		}
		if awsCfg == nil {
			return nil, err
		}
		awsSecretsManagerClient = secretsmanager.NewFromConfig(*awsCfg, func(options *secretsmanager.Options) {
			options.Retryer = retry.AddWithMaxBackoffDelay(options.Retryer, time.Duration(awsRetryBackoffInitialDelaySec)*time.Second)
		})
	}
	return awsSecretsManagerClient, nil
}

func getAwsConfig(ctx context.Context, region string, configProfileName string) (*aws.Config, error) {
	var (
		err    error
		cnf    aws.Config
		optFns []func(*awsConfig.LoadOptions) error
	)

	optFns = append(optFns, awsConfig.WithRetryer(func() aws.Retryer {
		retrier := retry.AddWithMaxAttempts(retry.NewStandard(), awsRetryMaxAttempts)
		return retry.AddWithMaxBackoffDelay(retrier, time.Second*time.Duration(awsRetryMaxBackoffDelaySec))
	}))

	if len(region) > 0 {
		optFns = append(optFns, awsConfig.WithRegion(region), awsConfig.WithDefaultRegion(region))
	}

	if len(configProfileName) > 0 {
		log.Printf("Shared Config Profile Name:  %s", configProfileName)
		optFns = append(optFns, awsConfig.WithSharedConfigProfile(configProfileName))
	}

	cnf, err = awsConfig.LoadDefaultConfig(ctx, optFns...)
	if err != nil {
		log.Fatal(err) // Exits after logging
		return nil, err
	}
	return &cnf, nil
}

func RetrieveSecretStringVal(ctx context.Context, client *secretsmanager.Client, secretId string) (string, error) {
	res, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretId),
	})
	if err != nil {
		log.Printf("Error while retrieving secret key for secret  %s", secretId)
		return "", err
	}
	if res.SecretString == nil {
		log.Printf("Secret with id %s has no String value", secretId)
		return "", err
	}

	log.Printf("Retrieved secret: %s", *res.SecretString)
	return *res.SecretString, nil
}

func ParseSecretArn(arn string) (secretId string, secretKey string, err error) {
	parts := strings.Split(arn, ":")
	log.Printf("parts %v, %d", parts, len(parts))
	if len(parts) < 7 {
		err = errors.New("invalid secret arn")
		log.Printf("In ParseSecretArn arn %s", arn)
		return
	}
	secretId = parts[0]
	for i := 1; i < 7; i++ {
		secretId = secretId + ":" + parts[i]
	}
	if len(parts) > 7 {
		secretKey = parts[7]
	}
	return
}
