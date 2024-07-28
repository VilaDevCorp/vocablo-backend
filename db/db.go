package db

import (
	"context"
	"fmt"
	"vocablo/conf"
	"vocablo/ent"
	"vocablo/ent/migrate"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

var db *ent.Client

func GetClient() *ent.Client {
	return db
}

func Setup() error {
	err := StartConnection()
	if err != nil {
		log.Debug().Err(err).Msg("Failed opening connection to DB")
		return err
	}
	ctx := context.Background()

	err = db.Schema.Create(
		ctx,
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	)
	if err != nil {
		log.Fatal().Msgf("failed creating schema resources: %v", err)
	}

	return nil
}

func getConnection() string {
	conf := conf.Get()
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", conf.DB.User, conf.DB.Pass,
		conf.DB.Host, conf.DB.Port, conf.DB.DB)
	return connString
}

func StartConnection() error {
	connString := getConnection()
	var err error
	db, err = ent.Open("postgres", connString)
	if err != nil {
		return err
	}
	return nil
}
