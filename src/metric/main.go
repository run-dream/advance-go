package main

import (
	"math/rand"
	"net/http"
	"time"

	"metric/pkg/metric"
)

func main() {
	mux := http.NewServeMux()
	var exchange = metric.NewMetricExchange("test")
	mux.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(exchange.Stats()))
	})
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		exchange.Events <- &metric.Event{Type: "start", StartTime: startTime}
		rand.Seed(time.Now().UnixNano())
		cost := rand.Intn(10)
		time.Sleep(time.Duration(cost) * time.Second)
		duration := time.Since(startTime)
		eventType := "success"
		if cost%2 == 1 {
			eventType = "fail"
		}
		exchange.Events <- &metric.Event{Type: eventType, StartTime: startTime, RunDuration: duration}
		w.Write([]byte("test"))
	})
	server := http.Server{
		Handler: mux,
		Addr:    ":8001",
	}
	server.ListenAndServe()
}
