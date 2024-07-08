package main

import (
    "fmt"
    "net/http"
    "sync"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    stressor = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "stressor_metric",
        Help: "A made-up metric to demonstrate custom autoscaling",
    })
    mu sync.Mutex
)

func init() {
    prometheus.MustRegister(stressor)
}

func main() {
    http.HandleFunc("/", handler)
    http.Handle("/metrics", promhttp.Handler())
    fmt.Println("Serving on :8080")
    http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
    mu.Lock()
    defer mu.Unlock()

    if r.Method == http.MethodPost {
        if r.URL.Query().Get("action") == "increase" {
            stressor.Inc()
        } else if r.URL.Query().Get("action") == "decrease" {
            stressor.Dec()
        }
    }

    fmt.Fprintf(w, `<html>
        <head><title>Stressor Metric</title></head>
        <body>
            <h1>Stressor Metric: %v</h1>
            <form method="post" action="?action=increase"><button type="submit">Increase</button></form>
            <form method="post" action="?action=decrease"><button type="submit">Decrease</button></form>
        </body>
    </html>`, stressor)
}
