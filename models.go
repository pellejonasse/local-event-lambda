package localeventlambda

import "github.com/aws/aws-lambda-go/lambda"

type Option func(*startConfig)

const (
	defaultPort            = "3000"
	defaultLocalServerPath = "/"
)

type startConfig struct {
	LocalPort       string          // Port to listen on when running locally
	LocalServerPath string          // Path to listen on when running locally (e.g. "/my-lambda")
	lambdaOptions   []lambda.Option // Options passed to lambda.StartWithOptions when running on AWS
}

func WithLocalPort(port string) Option {
	return func(cfg *startConfig) {
		cfg.LocalPort = port
	}
}

func WithLocalServerPath(path string) Option {
	return func(cfg *startConfig) {
		cfg.LocalServerPath = path
	}
}

// WithLambdaOption passes a lambda.Option through to lambda.StartWithOptions when running on AWS.
// Can be called multiple times to add multiple options.
func WithLambdaOption(opt lambda.Option) Option {
	return func(cfg *startConfig) {
		cfg.lambdaOptions = append(cfg.lambdaOptions, opt)
	}
}
