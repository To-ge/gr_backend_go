package config

import (
	"fmt"
	"os"
	"strconv"
)

type appConfig struct {
	RestInfo   *RestInfo
	GrpcInfo   *GrpcInfo
	DBInfo     *DBInfo
	DomainInfo *DomainInfo
	TestInfo   *TestInfo
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

type TestInfo struct {
	GrpcAddress string
	Location    TestLocation
}

type TestLocation struct {
	Latitude  float64
	Longitude float64
	Altitude  float32
}

func LoadConfig() *appConfig {
	address := ":" + os.Getenv("PUBLIC_PORT")

	restInfo := &RestInfo{
		Address: address,
	}

	address = fmt.Sprintf(":%s", os.Getenv("PRIVATE_PORT"))
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

	testAddress := os.Getenv("TEST_PRIVATE_ADDRESS")
	latitude, _ := strconv.ParseFloat(os.Getenv("TEST_LATITUDE"), 64)
	longitude, _ := strconv.ParseFloat(os.Getenv("TEST_LONGITUDE"), 64)
	altitude, _ := strconv.ParseFloat(os.Getenv("TEST_ALTITUDE"), 32)
	testInfo := &TestInfo{
		GrpcAddress: testAddress,
		Location: TestLocation{
			Latitude:  latitude,
			Longitude: longitude,
			Altitude:  float32(altitude),
		},
	}

	appConfig := appConfig{
		RestInfo:   restInfo,
		GrpcInfo:   grpcInfo,
		DBInfo:     dbInfo,
		DomainInfo: domainInfo,
		TestInfo:   testInfo,
	}

	return &appConfig
}
