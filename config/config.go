package config

import "os"

type appConfig struct {
	RestInfo *RestInfo
	DBInfo   *DBInfo
}

type RestInfo struct {
	Address string
}

type DBInfo struct {
	User     string
	Password string
	Address  string
	DBName   string
	DBPort   string
}

func LoadConfig() *appConfig {
	address := ":" + os.Getenv("PUBLIC_PORT")

	restInfo := &RestInfo{
		Address: address,
	}

	dbInfo := &DBInfo{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Address:  os.Getenv("DB_ADDRESS"),
		DBName:   os.Getenv("DB_NAME"),
		DBPort:   os.Getenv("DB_PORT"),
	}

	appConfig := appConfig{
		RestInfo: restInfo,
		DBInfo:   dbInfo,
	}

	return &appConfig
}
