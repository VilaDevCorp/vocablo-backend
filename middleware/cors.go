package middleware

import (
	"appname/conf"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		conf := conf.Get()
		if conf.Env == "prod" {
			c.Header("Access-Control-Allow-Origin", conf.Prod.CorsHost)
		} else {
			c.Header("Access-Control-Allow-Origin", conf.Dev.CorsHost)
		}
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "Accept, Content-Type, credentials, Content-Length, Accept-Encoding, X-API-CSRF, Authorization")
		c.Header("Access-Control-Expose-Headers", "Set-Cookie")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
