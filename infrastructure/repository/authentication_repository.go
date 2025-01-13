package repository

import (
	"context"
	"log"
	"time"

	domainRepo "github.com/To-ge/gr_backend_go/domain/repository"
	"github.com/To-ge/gr_backend_go/infrastructure/database"
)

type authenticationRepository struct {
	dbc  *database.DBConnector
	rdbc *database.RedisConnector
}

func NewAuthenticationRepository(dbc *database.DBConnector, rdbc *database.RedisConnector) domainRepo.IAuthenticationRepository {
	return &authenticationRepository{
		dbc:  dbc,
		rdbc: rdbc,
	}
}

func (ar *authenticationRepository) SignIn(key string) error {
	// Store session in Redis
	if err := ar.rdbc.Conn.Set(context.Background(), key, time.Now().Format(time.RFC3339), 0).Err(); err != nil {
		log.Printf("Redis Error: %s\n", err.Error())
		return err
	}
	return nil
}

func (ar *authenticationRepository) SignOut(key string) error {
	if err := ar.rdbc.Conn.Del(context.Background(), key).Err(); err != nil {
		log.Printf("Redis Error: %s\n", err.Error())
		return err
	}
	return nil
}
