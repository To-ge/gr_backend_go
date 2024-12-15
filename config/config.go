package config

import (
	"os"
	"strconv"
)

type appConfig struct {
	RestInfo   *RestInfo
	GrpcInfo   *GrpcInfo
	DBInfo     *DBInfo
	DomainInfo *DomainInfo
}

type RestInfo struct {
	Address string
}

type GrpcInfo struct {
	Address string
}

type DBInfo struct {
	User     string
	Password string
	Address  string
	DBName   string
	DBPort   string
}

type DomainInfo struct {
	TimerMinutes int
}

func LoadConfig() *appConfig {
	address := ":" + os.Getenv("PUBLIC_PORT")

	restInfo := &RestInfo{
		Address: address,
	}

	address = os.Getenv("PRIVATE_HOST")
	grpcInfo := &GrpcInfo{
		Address: address,
	}

	dbInfo := &DBInfo{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Address:  os.Getenv("DB_ADDRESS"),
		DBName:   os.Getenv("DB_NAME"),
		DBPort:   os.Getenv("DB_PORT"),
	}

	i, _ := strconv.Atoi(os.Getenv("TIMER_MINUTES"))
	domainInfo := &DomainInfo{
		TimerMinutes: i,
	}

	appConfig := appConfig{
		RestInfo:   restInfo,
		GrpcInfo:   grpcInfo,
		DBInfo:     dbInfo,
		DomainInfo: domainInfo,
	}

	return &appConfig
}
