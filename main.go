package main

import (
    "fmt"
    "net/http"
    "sync"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    dto "github.com/prometheus/client_model/go"
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

    var metric dto.Metric
    stressor.Write(&metric)
    value := metric.GetGauge().GetValue()
    percentage := (value / 30) * 100

    fmt.Fprintf(w, `
    <html>
        <head>
            <title>Stressor Metric</title>
            <style>
                body {
                    display: flex;
                    justify-content: center;
                    align-items: center;
                    height: 100vh;
                    margin: 0;
                    font-family: 'Hack', monospace;
                    background-color: #333;
                    color: #fff;
                }
                .container {
                    text-align: center;
                }
                .gauge {
                    width: 100%;
                    max-width: 300px;
                    height: auto;
                    display: block;
                    margin: 0 auto 20px auto;
                }
                .button {
                    font-size: 16px;
                    padding: 10px 20px;
                    margin: 10px;
                    border: none;
                    border-radius: 5px;
                    background-color: #007bff;
                    color: white;
                    cursor: pointer;
                }
                .button:hover {
                    background-color: #0056b3;
                }
            </style>
        </head>
        <body>
            <div class="container">
                <h1>Stressor Metric</h1>
                <svg class="gauge" viewBox="0 0 100 50">
                    <defs>
                        <linearGradient id="gradient" x1="0%" y1="0%" x2="100%" y2="0%">
                            <stop offset="0%" style="stop-color:green;stop-opacity:1" />
                            <stop offset="50%" style="stop-color:yellow;stop-opacity:1" />
                            <stop offset="100%" style="stop-color:red;stop-opacity:1" />
                        </linearGradient>
                    </defs>
                    <path d="M10,40 A30,30 0 0,1 90,40" fill="none" stroke="url(#gradient)" stroke-width="10" />
                    <line x1="50" y1="40" x2="50" y2="20" stroke="black" stroke-width="2" transform="rotate(%[1]f,50,40)" />
                </svg>
                <form method="post" action="?action=increase">
                    <button class="button" type="submit">Increase</button>
                </form>
                <form method="post" action="?action=decrease">
                    <button class="button" type="submit">Decrease</button>
                </form>
                <p>Current Value: %v</p>
            </div>
        </body>
    </html>
    `, percentage*180/100-90, value) // Calculate the rotation angle for the gauge line
}
