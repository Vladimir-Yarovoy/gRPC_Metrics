package main

import (
	"context"
	"fmt"
	"log"
	"metrics/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	var option int

	fmt.Printf("Input 1 to add random metrics\nInput 2 to add a metric (name, value)\nInput 3 to get the moving average of the metric\n")
	fmt.Scan(&option)

	conn, err := grpc.Dial(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := api.NewMetricsClient(conn)
	resOption, err := c.Option(context.Background(), &api.OptionRequest{X: int32(option)})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resOption.GetResult())

	switch option {

	case 1:
		return

	case 2:
		var name string
		var value float32

		fmt.Scan(&name, &value)

		resAdd, err := c.Add(context.Background(), &api.AddRequest{Name: name, Value: value})
		if err != nil {
			log.Fatal(err)
		}
		log.Println(resAdd)

	case 3:
		var name string
		fmt.Scan(&name)
		resAvgValue, err := c.GetAvgValue(context.Background(), &api.GavRequest{Name: name})
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Result:", resAvgValue.GetResult())
	default:
		return
	}

}
