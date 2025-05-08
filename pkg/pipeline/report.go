package pipeline

import (
	"context"
	"fmt"
	"github.com/cncf-tags/green-reviews-tooling/internal/dagger"
	"github.com/cncf-tags/green-reviews-tooling/pkg/cmd"
	"github.com/cncf-tags/green-reviews-tooling/pkg/monitoring"
)

// / pipeline/report.go
func (p *Pipeline) report(
	ctx context.Context,
	cncfProject, config, version, benchmarkJobURL string,
	benchmarkJobDurationMins int,
) (*dagger.Container, error) {

	// Show progress using Dagger's logging
	if _, err := p.echo(ctx, "🟢 Starting metrics collection..."); err != nil {
		return nil, err
	}

	// function that gets pod name

	results, err := monitoring.ComputeBenchmarkingResults(ctx)
	if err != nil {
		// Log error through Dagger
		_, _ = p.echo(ctx, fmt.Sprintf("❌ Metrics collection failed: %v", err))
		return nil, fmt.Errorf("metrics collection failed: %w", err)
	}

	// Send results to Dagger logs
	for _, res := range results {
		msg := fmt.Sprintf("📊 %s\n   Type: %s\n   Value: %s",
			res.Query,
			res.Type,
			res.Value,
		)
		if _, err := p.echo(ctx, msg); err != nil {
			return nil, err
		}
	}

	// Add final success message
	if _, err := p.echo(ctx, "✅ Benchmarking completed successfully"); err != nil {
		return nil, err
	}

	return p.container, nil
}

// Todo generalize with any potential name in the future
func (p *Pipeline) getFalcoPodName(ctx context.Context, client *dagger.Client) (string, error) {
	// container := client.Container().
	// 	From("bitnami/kubectl").
	// 	WithMountedFile("/root/.kube/config", client.Host().File("./.kube/config")).
	// 	WithEnvVariable("KUBECONFIG", "/root/.kube/config").
	// 	WithExec([]string{"kubectl", "get", "pods", "-n", "benchmark", "-l", "app=falco", "-o", "name"})

	// todo avoid harcoded variables

	output, err := p.exec(ctx, cmd.GetPodName("benchmark", "falco"))
	if err != nil {
		return "", err
	}

	p.echo(ctx, "############## output ##################")
	p.echo(ctx, output)

	// lines := strings.Split(output, "\n")
	// for _, line := range lines {
	// 	if strings.Contains(line, "falco-driver-modern-ebpf") {
	// 		parts := strings.Split(line, "/")
	// 		return parts[len(parts)-1], nil
	// 	}
	// }

	return "", fmt.Errorf("Falco pod not found")
}
