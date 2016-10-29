package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jmcfarlane/httprwd"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	addr = flag.String("addr", ":8080", "The address to listen on for HTTP requests.")

	handlerDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "example",
		Name:      "handler_duration_histogram_seconds",
		Help:      "Histogram of http call duration",
	}, []string{"route", "method", "code"})
)

func measure(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d := &httprwd.ResponseWriterDelegate{ResponseWriter: w}
		start := time.Now()
		handler.ServeHTTP(d, r)
		handlerDuration.WithLabelValues(
			r.URL.Path,
			r.Method,
			strconv.Itoa(d.Code),
		).Observe(time.Since(start).Seconds())
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}

func init() {
	prometheus.MustRegister(handlerDuration)
}

func main() {
	flag.Parse()
	http.Handle("/", measure(index))
	http.Handle("/metrics", prometheus.Handler())
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Panic(err)
	}
}
