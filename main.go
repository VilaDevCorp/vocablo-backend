package main

import (
	"vocablo/api"
	"vocablo/conf"
	"vocablo/db"
	"vocablo/svc"
	"vocablo/svc/mail"

	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("Starting application")
	err := conf.Setup()
	if err != nil {
		log.Fatal().Err(err).Msg("Fatal error in configuration setup")
		return
	}
	err = db.Setup()
	if err != nil {
		log.Fatal().Err(err).Msg("Fatal error in db setup")
		return
	}
	svc.Setup(db.GetClient(), &mail.MailSvcImpl{})
	api.Start()
}
