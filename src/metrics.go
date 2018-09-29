package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"sync"
	"time"
)

var (
)

// Create and setup metrics and collection
func setupMetricsCollection() {

	for metricName := range opts.configuration.Metrics {
		metric := opts.configuration.Metrics[metricName]
		metric.Init()

		if metric.Help == "" {
			metric.Help = metricName
		}

		switch metric.Type {
		case "gauge":
			vec := prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: metricName,
					Help: metric.Help,
				},
				metric.Labels,
			)

			metric.prometheus.gauge = vec
			prometheus.MustRegister(vec)
		case "summary":
			vec := prometheus.NewSummaryVec(
				prometheus.SummaryOpts{
					Name: metricName,
					Help: metric.Help,
				},
				metric.Labels,
			)

			metric.prometheus.summary = vec
			prometheus.MustRegister(vec)

		case "histogram":
			vec := prometheus.NewHistogramVec(
				prometheus.HistogramOpts{
					Name: metricName,
					Help: metric.Help,
				},
				metric.Labels,
			)

			metric.prometheus.histogram = vec
			prometheus.MustRegister(vec)
		default:
			panic(fmt.Sprintf("metric type \"\" not defined", metric.Type))
		}

		opts.configuration.Metrics[metricName] = metric
	}
}

// Start backgrounded metrics collection
func startMetricsCollection() {
	go func() {
		for {
			go func() {
				runMetricsCollection()
			}()
			time.Sleep(opts.ScrapeTime)
		}
	}()
}

// Metrics run
func runMetricsCollection() {
	var wg sync.WaitGroup

	callbackChannel := make(chan func())


	for metricName := range opts.configuration.Metrics {
		metric := opts.configuration.Metrics[metricName]

		for _, metricItem := range metric.Items {
			switch metric.Type {
			case "gauge":
				metric.prometheus.gauge.With(metricItem.Labels).Set(metricItem.GenerateValue())
			case "summary":
				metric.prometheus.summary.With(metricItem.Labels).Observe(metricItem.GenerateValue())
			case "histogram":
				metric.prometheus.histogram.With(metricItem.Labels).Observe(metricItem.GenerateValue())
			}
		}
	}


	go func() {
		var callbackList []func()
		for callback := range callbackChannel {
			callbackList = append(callbackList, callback)
		}

		for _, callback := range callbackList {
			callback()
		}

		Logger.Messsage("run: finished")
	}()

	// wait for all funcs
	wg.Wait()
	close(callbackChannel)
}

