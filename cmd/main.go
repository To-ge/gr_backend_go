package main

import (
	"log"
	"net/http"
	"time"

	"github.com/To-ge/gr_backend_go/adapter/grpc"
	"github.com/To-ge/gr_backend_go/adapter/rest"
	"github.com/To-ge/gr_backend_go/config"
	"github.com/To-ge/gr_backend_go/pkg"

	"github.com/joho/godotenv"
)

func main() {
	pkg.InitLogger()
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if err := grpc.InitRouter(); err != nil {
		log.Fatalf("Failed to initialize grpc router, %s", err.Error())
	}

	addr := config.LoadConfig().RestInfo.Address
	router, err := rest.InitRouter()
	if err != nil {
		log.Fatalf("router can't initialize, %s", err.Error())
		return
	}
	srv := &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("server is running! addr: ", addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to listen and serve: %+v", err)
	}

	log.Println("Server exiting")
}
