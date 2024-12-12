package main

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
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

	address := config.LoadConfig().GrpcInfo.Address
	fmt.Println("address ", address)
	conn, err := grpc.NewClient(
		fmt.Sprintf("localhost%s", address),
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
			Latitude:  rand.Float64() * 100,
			Longitude: rand.Float64() * 100,
			Altitude:  rand.Float32() * 100,
		}
		stream.Send(message)
		fmt.Printf("send %d message...\n", i)
		time.Sleep(1 * time.Second)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
}
