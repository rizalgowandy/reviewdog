package bitbucket

import (
	"context"
	"fmt"
	"net/url"

	insights "github.com/reva2/bitbucket-insights-api"
)

// BuildServerAPIContext builds context.Context used to call Bitbucket Server Code Insights API
func BuildServerAPIContext(ctx context.Context, bbURL, user, password, token string) (context.Context, error) {
	var err error
	ctx, err = withServerVariables(ctx, bbURL)
	if err != nil {
		return ctx, err
	}

	if user != "" && password != "" {
		ctx = withServerBasicAuth(ctx, user, password)
	}

	if token != "" {
		ctx = withServerAccessToken(ctx, token)
	}

	return ctx, nil
}

// WithServerBasicAuth adds basic auth credentials to context
func withServerBasicAuth(ctx context.Context, user, password string) context.Context {
	return context.WithValue(ctx, insights.ContextBasicAuth, insights.BasicAuth{
		UserName: user,
		Password: password,
	})
}

// WithServerAccessToken adds basic auth credentials to context
func withServerAccessToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, insights.ContextAccessToken, token)
}

// WithServerVariables adds server variable to context
func withServerVariables(ctx context.Context, bbURL string) (context.Context, error) {
	parsed, err := url.Parse(bbURL)
	if err != nil {
		return ctx, fmt.Errorf("failed to parse Bitbucket Server URL: %w", err)
	}

	if parsed.Scheme == "" {
		return ctx, fmt.Errorf("unable to determine scheme of Bitbucket Server URL: %w", err)
	}

	if parsed.Host == "" {
		return ctx, fmt.Errorf("unable to determine host of Bitbucket Server URL: %w", err)
	}

	return context.WithValue(
		ctx,
		insights.ContextServerVariables,
		map[string]string{
			"protocol":        parsed.Scheme,
			"bitbucketDomain": parsed.Host + parsed.Path,
		},
	), nil
}
