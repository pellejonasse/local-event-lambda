package main

import (
	"context"
	"fmt"
	"os"

	localeventlambda "local-event-lambda"

	"github.com/aws/aws-lambda-go/events"
)

type envConfig struct {
	QueueName string
}

type Handler struct {
	cfg envConfig
}

func newHandler() *Handler {
	return &Handler{
		cfg: envConfig{
			QueueName: os.Getenv("QUEUE_NAME"),
		},
	}
}

func (h *Handler) Handle(ctx context.Context, event events.SQSEvent) error {
	for _, record := range event.Records {
		fmt.Printf("[%s] Processing message: %s\n", h.cfg.QueueName, record.Body)
	}
	return nil
}

func main() {
	h := newHandler()
	localeventlambda.Start(h.Handle)
}
