package httpclient

import (
	"fmt"
	"github.com/raythx98/gohelpme/tool/logger"
	"net/http"
	"time"

	"github.com/raythx98/gohelpme/tool/httphelper"
)

// LogRoundTripper is an http.RoundTripper that logs requests and responses.
type LogRoundTripper struct {
	log logger.ILogger
}

// NewLogRoundTripper creates a new LogRoundTripper.
func NewLogRoundTripper(log logger.ILogger) *LogRoundTripper {
	return &LogRoundTripper{log: log}
}

// RoundTrip executes a single HTTP transaction, returning a Response for the provided Request.
func (t *LogRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	startAt := time.Now()
	reqLogGroup := t.createRequestLogGroup(req, startAt)

	resp, err := http.DefaultTransport.RoundTrip(req)

	message := fmt.Sprintf("[out-http] %s %s in %s", req.Method, req.URL.String(), time.Since(startAt).String())
	message += formatMessageSuffix(resp, err)

	t.log.Info(req.Context(), message,
		logger.WithField("request", reqLogGroup),
		logger.WithField("response", t.createResponseLogGroup(resp)),
	)

	return resp, err
}

func formatMessageSuffix(resp *http.Response, err error) string {
	if err != nil {
		return fmt.Sprintf(": error: %s", err.Error())
	}
	if resp == nil {
		return ": error, nil response"
	}
	return fmt.Sprintf(": %d", resp.StatusCode)
}

func (t *LogRoundTripper) createRequestLogGroup(req *http.Request, startAt time.Time) map[string]interface{} {
	return map[string]interface{}{
		"endpoint":   fmt.Sprintf("%s %s", req.Method, req.URL.String()),
		"method":     req.Method,
		"headers":    req.Header,
		"body":       httphelper.CopyRequestBody(req),
		"started at": startAt.Truncate(time.Second),
	}
}

func (t *LogRoundTripper) createResponseLogGroup(resp *http.Response) map[string]interface{} {
	if resp == nil {
		return nil
	}

	return map[string]interface{}{
		"status code":  resp.StatusCode,
		"status":       resp.Status,
		"body":         httphelper.CopyResponseBody(resp),
		"completed at": time.Now().Truncate(time.Second),
	}
}
