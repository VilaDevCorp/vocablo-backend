package conf

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Conf struct {
	Env    string
	Dev    EnvConf
	Prod   EnvConf
	IP     string
	Port   string
	DB     DatabaseConf
	Mail   MailConf
	JwtKey string
}

type EnvConf struct {
	CookieHost     string
	CookieSecure   bool
	CookieHttpOnly bool
	CorsHost       string
	FrontUrl       string
}

type DatabaseConf struct {
	Host string
	Port string
	User string
	Pass string
	DB   string
}

type MailConf struct {
	User     string
	Pass     string
	SmtpHost string
	SmtpPort string
	Texts    MailTexts
}

type MailTexts struct {
	ValidateEmail MailText
	ResetPassword MailText
}

type MailText struct {
	Subject string
	Body    string
}

var conf Conf

func Setup() error {
	_, currentPath, _, _ := runtime.Caller(0)

	log.Info().Msg("Initializing conf setup")
	viper.SetConfigName("configuration")
	viper.SetConfigType("yml")
	viper.AddConfigPath(filepath.Join(filepath.Dir(currentPath), ".."))
	err := readAndUnmarshal()

	if err != nil {
		log.Error().Err(err).Msg("Error reading conf file")
		return errors.New("error reading conf file")
	}
	if os.Getenv("DB_PASS") != "" {
		conf.DB.Pass = os.Getenv("DB_PASS")
	}
	if os.Getenv("JWT_SECRET") != "" {
		conf.JwtKey = os.Getenv("JWT_SECRET")
	}
	if os.Getenv("MAIL_PASS") != "" {
		conf.Mail.Pass = os.Getenv("MAIL_PASS")
	}
	log.Info().Msg("Conf file loaded succesfully")
	return nil
}

func readAndUnmarshal() error {
	err := viper.ReadInConfig()
	if err != nil {
		log.Error().Err(err)
		return err
	}
	err = viper.Unmarshal(&conf)
	if err != nil {
		log.Debug().Err(err)
		return err
	}
	return nil
}

func Get() *Conf {
	return &conf
}
