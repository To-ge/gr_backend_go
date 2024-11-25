package main

import (
	"log"
	"net/http"
	"time"

	"github.com/To-ge/gr_backend_go/adapter/rest"
	"github.com/To-ge/gr_backend_go/config"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
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

	// graceful shutdown
	// https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/notify-without-context/server.go
	log.Println("server is running! addr: ", addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to listen and serve: %+v", err)
	}

	log.Println("Server exiting")
}
