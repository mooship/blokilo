package httptest

import (
	"context"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/mooship/blokilo/internal/models"
)

func CheckHTTPConnectivity(ctx context.Context, domain string, timeout time.Duration, retries int) models.TestResult {
	client := resty.New().SetTimeout(timeout).SetRetryCount(retries)
	url := "https://" + domain
	start := time.Now()
	resp, err := client.R().SetContext(ctx).Get(url)
	dur := time.Since(start)

	if err != nil {
		if ctx.Err() != nil {
			return models.TestResult{Status: models.StatusBlocked, ResponseTime: dur, Err: ctx.Err()}
		}
		return models.TestResult{Status: models.StatusError, ResponseTime: dur, Err: err}
	}

	if resp.StatusCode() >= 200 && resp.StatusCode() < 400 {
		return models.TestResult{Status: models.StatusResolved, ResponseTime: dur, Domain: domain}
	}
	return models.TestResult{Status: models.StatusBlocked, ResponseTime: dur, Domain: domain}
}
