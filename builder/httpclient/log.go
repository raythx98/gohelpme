package httpclient

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/raythx98/gohelpme/tool/httphelper"
	"github.com/raythx98/gohelpme/tool/slogger"
)

func init() {
	slogger.Init()
}

type LogRoundTripper struct{}

func (t *LogRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	startAt := time.Now()
	reqLogGroup := t.createRequestLogGroup(req, startAt)

	resp, err := http.DefaultTransport.RoundTrip(req)

	slog.LogAttrs(
		req.Context(), slogger.GetLogLevel(err), "outgoing http execution",
		slog.String("time taken", time.Since(startAt).String()),
		reqLogGroup,
		t.createResponseLogGroup(resp),
	)

	return resp, err
}

func (t *LogRoundTripper) createRequestLogGroup(req *http.Request, startAt time.Time) slog.Attr {
	return slog.Group("request",
		slog.String("endpoint", fmt.Sprintf("%s %s", req.Method, req.URL.String())),
		slog.Any("headers", req.Header),
		slog.String("body", httphelper.CopyRequestBody(req)),
		slog.Time("started at", startAt.Truncate(time.Second)),
	)
}

func (t *LogRoundTripper) createResponseLogGroup(resp *http.Response) slog.Attr {
	if resp == nil {
		return slog.Attr{}
	}

	return slog.Group("response",
		slog.Time("completed at", time.Now().Truncate(time.Second)),
		slog.Int("status code", resp.StatusCode),
		slog.String("status", resp.Status),
		slog.Any("response body", httphelper.CopyResponseBody(resp)),
	)
}
