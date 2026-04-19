package localeventlambda

import (
	"fmt"
	"os"
	"reflect"

	"github.com/aws/aws-lambda-go/lambda"

	fiberproxy "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/gofiber/fiber/v2"
)




// start lambda handler with default options
func Start(handler any) {
	startLambda(handler)
}

// start lambda handler with custom options
func StartWithOptions(handler any, options ...Option) {
	startLambda(handler, options...)
}

func startLambda(handler any, options ...Option) {
	opts := localConfig{
		Port:            defaultPort,
		LocalServerPath: defaultLocalServerPath,
	}
	if len(options) > 0 {
		for _, option := range options {
			option(&opts)
		}
		if opts.Port == "" {
			opts.Port = defaultOptions.Port
		}
		if opts.LocalServerPath == "" {
			opts.LocalServerPath = defaultOptions.LocalServerPath
		}
	}

	switch reflect.TypeOf(handler) {
		// either our lambda is a fiber app handler or a regular event based handler
	case reflect.TypeOf(&fiber.App{}):
		startFiberApp(handler.(*fiber.App), opts)
	
	default:
		startEventBasedHandler(handler, opts)
	}
}

func isLocalLambda() bool {
	if _, ok := os.LookupEnv("AWS_LAMBDA_FUNCTION_NAME"); ok {
		return false
	}
	return true
}

func startFiberApp(app *fiber.App, options localConfig) {
	if isLocalLambda() {
		app.Listen(fmt.Sprintf("127.0.0.1:%s", options.Port))
	} 
	
	lambda.Start(fiberproxy.New(app).ProxyWithContext)
}

func startEventBasedHandler(handler any, options localConfig) {
	if isLocalLambda() {
		app := registerNewFiberApp(handler, options)
		app.Listen(fmt.Sprintf("127.0.0.1:%s", options.Port))
	}

	lambda.Start(handler)
}


func registerNewFiberApp(handler any, options ...localConfig) *fiber.App {
	app := fiber.New()
	lambdaHandler := lambda.NewHandler(handler) // Set a reasonable timeout for local testing
	fiberHandler := func(c *fiber.Ctx) error {
		body := c.Body()
		response, err := lambdaHandler.Invoke(c.Context(), body)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Error invoking lambda: %v", err))
		}
		return c.Status(fiber.StatusOK).Send(response)
	}

	app.Post(options.LocalServerPath, fiberHandler)
	return app
}

