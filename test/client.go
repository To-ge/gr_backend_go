package main

import (
	"context"
	"fmt"
	"log"
	"time"

	v1 "github.com/To-ge/gr_backend_go/adapter/grpc/api/gen/go/v1"
	"github.com/To-ge/gr_backend_go/config"
	"github.com/To-ge/gr_backend_go/pkg"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	pkg.InitLogger()
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal("Error loading .env file")
	}
	fmt.Println("finished loading...")

	testInfo := config.LoadConfig().TestInfo
	fmt.Println("address: ", testInfo.GrpcAddress)
	conn, err := grpc.NewClient(
		testInfo.GrpcAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal("Connection failed.")
		return
	}
	defer conn.Close()

	client := v1.NewTelemetryServiceClient(conn)

	stream, err := client.SendLocation(context.Background())
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal("Connection failed.")
		return
	}

	for i := 0; i < 5; i++ {
		message := &v1.SendLocationRequest{
			Timestamp: time.Now().Unix(),
			Latitude:  testInfo.Location.Latitude + float64(i)*0.001,
			Longitude: testInfo.Location.Longitude + float64(i)*0.001,
			Altitude:  testInfo.Location.Altitude,
		}
		stream.Send(message)
		fmt.Printf("send %d message...\n", i)
		time.Sleep(2 * time.Second)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
}
