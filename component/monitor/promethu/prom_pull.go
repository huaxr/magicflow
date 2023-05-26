// Author: XinRui Hua
// Time:   2022/4/14 下午5:11
// Git:    huaxr

package promethu

import (
	"context"
	"fmt"
	"net/http"

	"github.com/huaxr/magicflow/component/logx"
	"github.com/huaxr/magicflow/component/monitor/promethu/metric"
	"github.com/huaxr/magicflow/pkg/confutil"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func PromPull(ctx context.Context) {
	metric.RegisterMetrics()
	worker := NewManager(ctx)

	// Since we are dealing with custom Collector implementations, it might
	// be a good idea to try it out with a pedantic registry.
	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(worker)

	gatherers := prometheus.Gatherers{
		prometheus.DefaultGatherer,
		reg,
	}

	h := promhttp.HandlerFor(gatherers,
		promhttp.HandlerOpts{
			ErrorLog:      &promlog{},
			ErrorHandling: promhttp.ContinueOnError,
		})
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})

	port := confutil.GetProm().PullPort
	logx.L().Infof(fmt.Sprintf("Start prom server at http://0.0.0.0:%s/metrics", port))
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		logx.L().Errorf("Error occur when start server %v", err)
	}
}
