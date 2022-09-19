package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"github.com/axxyhtrx/kenetic-prometheus-exporter/zyx"
)

// Define a struct for you collector that contains pointers
// to prometheus descriptors for each metric you wish to expose.
// Note you can also include fields of other types if they provide utility
// but we just won't be exposing them as metrics.
type fooCollector struct {
	cpuUsage *prometheus.Desc
	memFree  *prometheus.Desc
	memUsed  *prometheus.Desc
	rxSpeed  *prometheus.Desc
	txSpeed  *prometheus.Desc
}

// You must create a constructor for you collector that
// initializes every descriptor and returns a pointer to the collector
func newFooCollector() *fooCollector {
	return &fooCollector{
		cpuUsage: prometheus.NewDesc("cpu_usage",
			"Shows router cpu loading",
			nil, nil,
		),
		memFree: prometheus.NewDesc("mem_free",
			"Shows free memory in a router(total - used)",
			nil, nil,
		),
		memUsed: prometheus.NewDesc("mem_used",
			"Shows used memory in a router",
			nil, nil,
		),
		rxSpeed: prometheus.NewDesc("rx_speed",
			"rx speed of uplink interface",
			nil, nil,
		),
		txSpeed: prometheus.NewDesc("tx_speed",
			"tx speed of uplink interface",
			nil, nil,
		),
	}
}

// Each and every collector must implement the Describe function.
// It essentially writes all descriptors to the prometheus desc channel.
func (collector *fooCollector) Describe(ch chan<- *prometheus.Desc) {

	//Update this section with the each metric you create for a given collector
	ch <- collector.cpuUsage
	ch <- collector.memFree
	ch <- collector.memUsed
	ch <- collector.rxSpeed
	ch <- collector.txSpeed

}

// Collect implements required collect function for all promehteus collectors
func (collector *fooCollector) Collect(ch chan<- prometheus.Metric) {

	//Write latest value for each metric in the prometheus metric channel.
	//Note that you can pass CounterValue, GaugeValue, or UntypedValue types here.
	cpu := prometheus.MustNewConstMetric(collector.cpuUsage, prometheus.CounterValue, float64(zyx.Zmon.CPU))
	fmem := prometheus.MustNewConstMetric(collector.memFree, prometheus.CounterValue, float64(zyx.Zmon.FreeRam))
	umem := prometheus.MustNewConstMetric(collector.memUsed, prometheus.CounterValue, float64(zyx.Zmon.UsedRam))
	rxSpeed := prometheus.MustNewConstMetric(collector.rxSpeed, prometheus.CounterValue, float64(zyx.Zmon.RXSpeed))
	txSpeed := prometheus.MustNewConstMetric(collector.txSpeed, prometheus.CounterValue, float64(zyx.Zmon.TXSpeed))
	ch <- cpu
	ch <- fmem
	ch <- umem
	ch <- rxSpeed
	ch <- txSpeed

}

func main() {
	foo := newFooCollector()
	prometheus.MustRegister(foo)
	go zyx.Poller()
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9101", nil))
}
