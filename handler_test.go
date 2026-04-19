package localeventlambda

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func postEvent(t *testing.T, app *fiber.App, path string, payload any) *http.Response {
	t.Helper()
	body, err := json.Marshal(payload)
	assert.NoError(t, err)
	req, _ := http.NewRequest(http.MethodPost, path, strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	return resp
}

func TestEventBasedHandlers(t *testing.T) {
	cfg := startConfig{LocalPort: defaultPort, LocalServerPath: defaultLocalServerPath}

	t.Run("SQS", func(t *testing.T) {
		var received events.SQSEvent

		handler := func(ctx context.Context, event events.SQSEvent) error {
			received = event
			return nil
		}

		app := registerNewFiberApp(handler, cfg)
		payload := events.SQSEvent{
			Records: []events.SQSMessage{
				{MessageId: "msg-1", Body: "hello from sqs"},
			},
		}

		resp := postEvent(t, app, "/", payload)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Len(t, received.Records, 1)
		assert.Equal(t, "msg-1", received.Records[0].MessageId)
	})

	t.Run("SNS", func(t *testing.T) {
		var received events.SNSEvent

		handler := func(ctx context.Context, event events.SNSEvent) error {
			received = event
			return nil
		}

		app := registerNewFiberApp(handler, cfg)
		payload := events.SNSEvent{
			Records: []events.SNSEventRecord{
				{SNS: events.SNSEntity{MessageID: "sns-1", Message: "hello from sns"}},
			},
		}

		resp := postEvent(t, app, "/", payload)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Len(t, received.Records, 1)
		assert.Equal(t, "sns-1", received.Records[0].SNS.MessageID)
	})

	t.Run("EventBridge", func(t *testing.T) {
		var received events.CloudWatchEvent

		handler := func(ctx context.Context, event events.CloudWatchEvent) error {
			received = event
			return nil
		}

		app := registerNewFiberApp(handler, cfg)
		payload := events.CloudWatchEvent{
			ID:     "eb-1",
			Source: "com.myapp.orders",
			Detail: json.RawMessage(`{"orderId":"123"}`),
		}

		resp := postEvent(t, app, "/", payload)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "eb-1", received.ID)
		assert.Equal(t, "com.myapp.orders", received.Source)
	})
}

func TestFiberAppHandler(t *testing.T) {
	fiberApp := fiber.New()
	fiberApp.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	resp, err := fiberApp.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
