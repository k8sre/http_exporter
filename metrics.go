/**
 * @Author: Sonic Ma
 * @Author: sonic.ma@outlook.com
 * @Date: 2020/3/30 4:20 下午
 * @Desc:
 */

package main

//noinspection GoUnresolvedReference
import (
	"github.com/prometheus/client_golang/prometheus"
)

// PromCollectors has instances of Prometheus Collectors
type PromCollectors struct {
	count      *prometheus.CounterVec
	error      *prometheus.CounterVec
	code       *prometheus.GaugeVec
	duration   *prometheus.HistogramVec
	length     *prometheus.GaugeVec
	hash       *prometheus.GaugeVec
	registerer prometheus.Registerer
}

func (promCollectors *PromCollectors) Count() *prometheus.CounterVec {
	return promCollectors.count
}

// Register registers all collectors
func (promCollectors *PromCollectors) Register(registerer prometheus.Registerer) {

	promCollectors.registerer = registerer

	promCollectors.count = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_count",
		Help: "Total number of performed check",
	}, []string{"url"})
	registerer.MustRegister(promCollectors.count)

	promCollectors.error = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_error",
		Help: "Total number of error",
	}, []string{"url"})
	registerer.MustRegister(promCollectors.error)

	promCollectors.code = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "http_code",
		Help: "Response code",
	}, []string{"url"})
	registerer.MustRegister(promCollectors.code)

	promCollectors.duration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_duration_seconds",
		Help: "Histogram of request duration",
	}, []string{"url"})
	registerer.MustRegister(promCollectors.duration)

	promCollectors.length = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "http_length",
		Help: "Page length",
	}, []string{"url"})
	registerer.MustRegister(promCollectors.length)

	promCollectors.hash = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "http_hash",
		Help: "Page hash",
	}, []string{"url"})
	registerer.MustRegister(promCollectors.hash)

}

// Update copied values from latest measurements to Prometheus collectors
func (promCollectors *PromCollectors) Update(site string, result Result, err error) {

	siteLabels := prometheus.Labels{"url": site}
	promCollectors.count.With(siteLabels).Inc()
	promCollectors.code.With(siteLabels).Set(float64(result.StatusCode))
	promCollectors.duration.With(siteLabels).Observe(result.Duration.Seconds())
	promCollectors.length.With(siteLabels).Set(float64(result.Length))
	promCollectors.hash.With(siteLabels).Set(float64(result.Hash))
	if err != nil {
		promCollectors.error.With(siteLabels).Inc()
	}
}
