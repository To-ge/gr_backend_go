package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/sessions"
)

type Mode int

const (
	Usually Mode = iota
	Demo
)

type appConfig struct {
	Mode       Mode
	Deployment string
	FEUrl      string
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
	UseTLS      bool
	Location    TestLocation
}

type TestLocation struct {
	Latitude  float64
	Longitude float64
	Altitude  float32
}

func LoadConfig() *appConfig {
	var mode Mode
	if os.Getenv("MODE") == "demo" {
		mode = Demo
	} else {
		mode = Usually
	}
	deployment := os.Getenv("DEPLOYMENT")
	feUrl := func() string {
		switch deployment {
		case "local":
			return os.Getenv("FE_URL_LOCAL")
		case "develop":
			return os.Getenv("FE_URL_DEV")
		default:
			return os.Getenv("FE_URL_DEFAULT")
		}
	}()
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
	useTLS := !strings.Contains(testAddress, "localhost")
	latitude, _ := strconv.ParseFloat(os.Getenv("TEST_LATITUDE"), 64)
	longitude, _ := strconv.ParseFloat(os.Getenv("TEST_LONGITUDE"), 64)
	altitude, _ := strconv.ParseFloat(os.Getenv("TEST_ALTITUDE"), 32)
	testInfo := &TestInfo{
		GrpcAddress: testAddress,
		UseTLS:      useTLS,
		Location: TestLocation{
			Latitude:  latitude,
			Longitude: longitude,
			Altitude:  float32(altitude),
		},
	}

	appConfig := appConfig{
		Mode:       mode,
		Deployment: deployment,
		FEUrl:      feUrl,
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
