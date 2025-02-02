package config

import (
	"context"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

// Configuration represents the application configuration with the default values.
type Configuration struct {
	Environment          string `env:"ENVIRONMENT,required"`
	LogLevel             string `env:"LOG_LEVEL,default=INFO"`
	Port                 string `env:"PORT,default=8000"`
	ConsumerMode         string `env:"CONSUMER_MODE,default=QUEUE"`
	AwsEndpoint          string `env:"AWS_ENDPOINT"`
	AwsAccessKeyID       string `env:"AWS_ACCESS_KEY_ID"`
	AwsSecretAccessKey   string `env:"AWS_SECRET_ACCESS_KEY"`
	AwsRegion            string `env:"AWS_REGION"`
	SQSUrl               string `env:"SQS_URL"`
	InfluxUrl            string `env:"INFLUX_URL"`
	InfluxToken          string `env:"INFLUX_TOKEN"`
	InfluxOrganization   string `env:"INFLUX_ORGANIZATION"`
	InfluxBucketInfinite string `env:"INFLUX_BUCKET_INFINITE"`
	InfluxBucket30Days   string `env:"INFLUX_BUCKET_30_DAYS"`
	InfluxBucket24Hours  string `env:"INFLUX_BUCKET_24_HOURS"`
	MongodbURI           string `env:"MONGODB_URI,required"`
	MongodbDatabase      string `env:"MONGODB_DATABASE,required"`
	PprofEnabled         bool   `env:"PPROF_ENABLED,default=false"`
	P2pNetwork           string `env:"P2P_NETWORK,required"`
	CacheURL             string `env:"CACHE_URL,required"`
	CachePrefix          string `env:"CACHE_PREFIX,required"`
	CacheChannel         string `env:"CACHE_CHANNEL,required"`
}

// New creates a configuration with the values from .env file and environment variables.
func New(ctx context.Context) (*Configuration, error) {
	_ = godotenv.Load(".env", "../.env")

	var configuration Configuration
	if err := envconfig.Process(ctx, &configuration); err != nil {
		return nil, err
	}

	return &configuration, nil
}

// IsQueueConsumer check if consumer mode is QUEUE.
func (c *Configuration) IsQueueConsumer() bool {
	return c.ConsumerMode == "QUEUE"
}
