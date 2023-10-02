package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"metrics/api"
	"net"
	"time"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	api.MetricsServer
}

type Metric struct {
	Name      string    // The name of the metric
	Value     float64   // The value of the metric
	Timestamp time.Time // Timestamp of adding a metric
}

var metrics = make(map[string][]Metric)

func (s *GRPCServer) Option(ctx context.Context, req *api.OptionRequest) (*api.OptionResponse, error) {
	option := req.GetX()
	switch option {
	case 1:
		GenerateMetrics()
		return &api.OptionResponse{Result: "Metrics are written in the format metric-No.(from 1 to 10000, for example metric-1, etc) and value"}, nil
	case 2:
		return &api.OptionResponse{Result: "Input the name and value of the metric separated by a space"}, nil
	case 3:
		return &api.OptionResponse{Result: "Input the name of the metric"}, nil
	default:
		return &api.OptionResponse{Result: "Exit"}, nil
	}
}

func (s *GRPCServer) Add(ctx context.Context, req *api.AddRequest) (*api.AddResponse, error) {
	name, value := req.Name, req.Value
	AddMetric(name, float64(value))
	return &api.AddResponse{Result: "Metric added"}, nil

}

func (s *GRPCServer) GetAvgValue(ctx context.Context, req *api.GavRequest) (*api.GavResponse, error) {
	now := time.Now()
	CleanUpMetrics(now)
	result := GetAvgValue(req.Name)
	return &api.GavResponse{Result: float32(result)}, nil
}

func GenerateMetrics() {
	for i := 1; i <= 10000; i++ {
		now := time.Now()
		value := rand.Float64()
		name := fmt.Sprintf("metric-%d", i)
		metric := Metric{Name: name, Value: value, Timestamp: now}
		metrics[name] = append(metrics[name], metric)
	}
}

func AddMetric(metricName string, value float64) {
	now := time.Now()
	metrics[metricName] = append(metrics[metricName], Metric{Timestamp: now, Value: value})
}

func GetAvgValue(metricName string) float64 {
	MetricData := metrics[metricName]
	if len(MetricData) == 0 {
		return 0
	}

	var sum float64
	var count int

	for _, data := range MetricData {
		sum += data.Value
		count++
	}
	return sum / float64(count)
}

func CleanUpMetrics(currentTime time.Time) {
	for metricName, data := range metrics {
		var newData []Metric

		for _, metric := range data {
			if currentTime.Sub(metric.Timestamp) <= time.Minute {
				newData = append(newData, metric)
			}
		}

		metrics[metricName] = newData
	}
}

func main() {
	s := grpc.NewServer()
	srv := &GRPCServer{}
	api.RegisterMetricsServer(s, srv)

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("server listening at %v", l.Addr())

	if err := s.Serve(l); err != nil {
		log.Fatal(err)
	}
}
