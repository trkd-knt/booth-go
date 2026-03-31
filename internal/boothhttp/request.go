package boothhttp

import (
	"context"
	"net/http"

	"golang.org/x/time/rate"
)

type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

func ExecuteGET(ctx context.Context, client HTTPDoer, limiter *rate.Limiter, userAgent, lang, rawURL string) (*http.Response, error) {
	if err := limiter.Wait(ctx); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept-Language", lang)

	return client.Do(req)
}
