package localeventlambda

import "github.com/aws/aws-lambda-go/lambda"

type Option func(*localConfig)

const (
	defaultPort            = "3000"
	defaultLocalServerPath = "/"
)

type localConfig struct {
	Port            string        // Port to listen on when running locally
	LocalServerPath string        // Path to listen on when running locally (e.g. "/my-lambda")
	LambdaOptions   lambda.Option // Options passed to lambda.StartWithOptions when running on AWS
}

func WithPort(port string) Option {
	return func(cfg *localConfig) {
		cfg.Port = port
	}
}

func WithLocalServerPath(path string) Option {
	return func(cfg *localConfig) {
		cfg.LocalServerPath = path
	}
}

func WithLambdaOptions(opts lambda.Option) Option {
	return func(cfg *localConfig) {
		cfg.LambdaOptions = opts
	}
}
