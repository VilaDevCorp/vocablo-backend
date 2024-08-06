package api

import (
	"fmt"
	"net/http"
	"vocablo/api/auth"
	"vocablo/api/userword"
	"vocablo/api/word"
	"vocablo/conf"
	"vocablo/middleware"

	ginzerolog "github.com/dn365/gin-zerolog"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func Start() error {
	conf := conf.Get()
	r := GetRouter()
	log.Info().Msg(fmt.Sprintf(conf.Port))
	err := r.Run(fmt.Sprintf("%s:%s", conf.IP, conf.Port))
	return err
}

func GetRouter() *gin.Engine {
	api := gin.Default()
	api.Use(gin.Recovery())
	api.Use(ginzerolog.Logger("gin"))
	api.Use(middleware.Cors())
	pub := api.Group("/api/public")
	pub.GET("/health", health)
	pub.POST("/login", auth.Login)
	pub.POST("/register", auth.SignUp)
	pub.POST("/validate/:username/:code", auth.ValidateAccount)
	pub.POST("/validate/:username/resend", auth.ResendValidationCode)
	pub.POST("/forgotten-password/:username", auth.SendForgottenPasswordCode)
	pub.POST("/reset-password/:username/:code", auth.ResetPassword)
	priv := api.Group("/api")
	priv.Use(middleware.Authentication())
	priv.GET("/self", auth.Self)
	priv.POST("/userword", userword.Create)
	priv.PUT("/userword", userword.Update)
	priv.GET("/userword/:id", userword.Get)
	priv.POST("/userword/search", userword.Search)
	priv.DELETE("/userword/:id", userword.Delete)
	priv.POST("/word/search", word.Search)
	return api
}

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Everything is FINE"})
}
