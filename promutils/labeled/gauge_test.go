package labeled

import (
	"context"
	"github.com/lyft/flytestdlib/contextutils"
	"github.com/lyft/flytestdlib/promutils"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)


func TestLabeledGauge(t *testing.T) {
	UnsetMetricKeys()
	assert.NotPanics(t, func() {
		SetMetricKeys(contextutils.ProjectKey, contextutils.DomainKey, contextutils.WorkflowIDKey, contextutils.TaskIDKey, contextutils.LaunchPlanIDKey)
	})

	scope := promutils.NewScope("testscope")
	ctx := context.Background()
	ctx = contextutils.WithProjectDomain(ctx, "flyte", "dev")
	g := NewGauge("unittest", "some desc", scope)
	assert.NotNil(t, g)

	g.Inc(ctx)

	const header = `
		# HELP testscope:unittest some desc
        # TYPE testscope:unittest gauge
	`
	var expected = `
        testscope:unittest{domain="dev",lp="",project="flyte",task="",wf=""} 1
	`
	err := testutil.CollectAndCompare(g.GaugeVec, strings.NewReader(header + expected))
	assert.NoError(t, err)

	g.Set(ctx, 42)
	expected = `
        testscope:unittest{domain="dev",lp="",project="flyte",task="",wf=""} 42
	`
	err = testutil.CollectAndCompare(g.GaugeVec, strings.NewReader(header + expected))
	assert.NoError(t, err)

	g.Add(ctx, 1)
	expected = `
        testscope:unittest{domain="dev",lp="",project="flyte",task="",wf=""} 43
	`
	err = testutil.CollectAndCompare(g.GaugeVec, strings.NewReader(header + expected))
	assert.NoError(t, err)

	g.Dec(ctx)
	expected = `
        testscope:unittest{domain="dev",lp="",project="flyte",task="",wf=""} 42
	`
	err = testutil.CollectAndCompare(g.GaugeVec, strings.NewReader(header + expected))
	assert.NoError(t, err)

	g.Sub(ctx, 1)
	expected = `
        testscope:unittest{domain="dev",lp="",project="flyte",task="",wf=""} 41
	`
	err = testutil.CollectAndCompare(g.GaugeVec, strings.NewReader(header + expected))
	assert.NoError(t, err)

	g.SetToCurrentTime(ctx)
}
