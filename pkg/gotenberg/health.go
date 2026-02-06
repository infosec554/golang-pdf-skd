package gotenberg

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type HealthStatus struct {
	Available    bool
	StatusCode   int
	ResponseTime time.Duration
	Version      string
	Error        error
}

func (g *gotenbergClient) HealthCheck(ctx context.Context) *HealthStatus {
	start := time.Now()
	status := &HealthStatus{}

	req, err := http.NewRequestWithContext(ctx, "GET", g.baseURL+"/health", nil)
	if err != nil {
		status.Error = fmt.Errorf("failed to create request: %w", err)
		return status
	}

	resp, err := g.httpClient.Do(req)
	status.ResponseTime = time.Since(start)

	if err != nil {
		status.Error = fmt.Errorf("health check failed: %w", err)
		return status
	}
	defer resp.Body.Close()

	status.StatusCode = resp.StatusCode
	status.Available = resp.StatusCode == http.StatusOK
	status.Version = resp.Header.Get("Gotenberg-Version")

	return status
}

func (g *gotenbergClient) IsHealthy(ctx context.Context) bool {
	status := g.HealthCheck(ctx)
	return status.Available
}

func (g *gotenbergClient) WaitForReady(ctx context.Context, timeout time.Duration, interval time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if g.IsHealthy(ctx) {
				return nil
			}
			time.Sleep(interval)
		}
	}

	return fmt.Errorf("gotenberg not ready after %v", timeout)
}
