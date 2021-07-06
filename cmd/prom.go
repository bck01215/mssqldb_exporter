package main

import (
	"errors"
	"net/http"
	"path"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

type error interface {
	Error() string
}

func do_stuff(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("target")
	configs := get_metric_info(target)
	registry := prometheus.NewRegistry()
	for _, metric := range configs.Metrics {
		make_gauges(registry, metric)
	}
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

func main() {
	logrus.Info("Starting exporter")
	http.HandleFunc(path.Join("/probe"), func(w http.ResponseWriter, r *http.Request) {
		do_stuff(w, r)
	})

	logrus.Fatal(http.ListenAndServe(":9101", nil))
}

func make_gauges(reg *prometheus.Registry, metric Metric) {
	var label_vals []string
	any_errs := false
L:
	for _, i := range metric.Labels {
		switch metric.Values {
		case nil:
			logrus.Error("Metric data collection failed. Unable to add metric: " + metric.Name)
			any_errs = true
			break L
		default:
			label_vals = append(label_vals, metric.Values[i].(string))

		}
	}
	if any_errs == false {
		x := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metric.Name,
			Help: metric.Help,
		}, metric.Labels)
		x, err := set_values(x, label_vals, metric)
		if err == nil {
			reg.MustRegister(x)
		}
	}
}

func set_values(x *prometheus.GaugeVec, label_vals []string, metric Metric) (*prometheus.GaugeVec, error) {
	switch metric.Values[metric.Value].(type) {
	case float64:
		x.WithLabelValues(label_vals...).Set(float64(metric.Values[metric.Value].(float64)))
	case int64:
		x.WithLabelValues(label_vals...).Set(float64(metric.Values[metric.Value].(int64)))
	case string:
		float, err := strconv.ParseFloat(metric.Values[metric.Value].(string), 64)
		if err != nil {
			logrus.Error("Could not convert string to float64: ", err)
			break
		}
		x.WithLabelValues(label_vals...).Set(float)
	case nil:
		return nil, errors.New("Metric data collection failed. Unable to add metric: " + metric.Name)
	}
	return x, nil
}
