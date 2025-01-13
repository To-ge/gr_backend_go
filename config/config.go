package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gorilla/sessions"
)

type appConfig struct {
	RestInfo   *RestInfo
	GrpcInfo   *GrpcInfo
	DBInfo     *DBInfo
	RedisInfo  *RedisInfo
	DomainInfo *DomainInfo
	TestInfo   *TestInfo
}

type RestInfo struct {
	Address      string
	CookieSecret string
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

type RedisInfo struct {
	Address string
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
	cookieSecret := os.Getenv("COOKIE_SECRET")

	restInfo := &RestInfo{
		Address:      address,
		CookieSecret: cookieSecret,
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

	redisInfo := &RedisInfo{
		Address: os.Getenv(("REDIS_ADDRESS")),
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
		RedisInfo:  redisInfo,
		DomainInfo: domainInfo,
		TestInfo:   testInfo,
	}

	secret := os.Getenv("COOKIE_SECRET")
	SessionStore = sessions.NewCookieStore([]byte(secret))

	return &appConfig
}

var (
	SessionStore *sessions.CookieStore
	SessionKey   = "key"
)
