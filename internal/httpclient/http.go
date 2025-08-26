package httpclient

import (
	"context"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/mooship/blokilo/internal/models"
)

type noopLogger struct{}

func (noopLogger) Errorf(format string, v ...interface{}) {}
func (noopLogger) Warnf(format string, v ...interface{})  {}
func (noopLogger) Debugf(format string, v ...interface{}) {}

func CheckHTTPConnectivity(ctx context.Context, domain string, timeout time.Duration, retries int) models.TestResult {
	client := resty.New().SetTimeout(timeout).SetRetryCount(retries).SetRedirectPolicy(resty.FlexibleRedirectPolicy(10))
	client.SetLogger(noopLogger{})
	url := "https://" + domain
	start := time.Now()
	resp, err := client.R().SetContext(ctx).EnableTrace().Get(url)
	dur := time.Since(start)

	result := models.TestResult{
		Domain:         domain,
		ResponseTime:   dur,
		HTTPStatusCode: 0,
	}

	if err != nil {
		if ctx.Err() != nil {
			result.Status = models.StatusBlocked
			result.Err = ctx.Err()
			return result
		}
		result.Status = models.StatusError
		result.Err = err
		return result
	}

	if resp != nil {
		result.HTTPStatusCode = resp.StatusCode()

		code := resp.StatusCode()
		switch {
		case code == 401 || code == 403 || code == 404 || code == 429 || code == 451 || code == 503:
			result.Status = models.StatusBlocked
		case code >= 200 && code < 400:
			result.Status = models.StatusResolved
		default:
			result.Status = models.StatusBlocked
		}
	} else {
		result.Status = models.StatusError
		if result.Err == nil {
			result.Err = err
		}
	}
	return result
}
