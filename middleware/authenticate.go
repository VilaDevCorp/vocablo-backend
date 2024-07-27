package middleware

import (
	"appname/customerrors"
	"appname/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		csrf := c.Request.Header.Get("X-API-CSRF")
		if csrf == "" {
			res := utils.ErrorResponse(http.StatusUnauthorized, utils.GetStringPointer("Not CSRF token present"), utils.GetStringPointer(customerrors.NOT_CSRF_TOKEN))
			c.AbortWithStatusJSON(res.Status, res.Body)
			return
		}
		jwt, err := c.Cookie("JWT_TOKEN")
		if err != nil {
			res := utils.ErrorResponse(http.StatusUnauthorized, utils.GetStringPointer("Not JWT present"), utils.GetStringPointer(customerrors.NOT_JWT_TOKEN))
			c.AbortWithStatusJSON(res.Status, res.Body)
			return
		}
		tokenClaims, jwtError := utils.ValidateToken(jwt)
		if jwtError != nil || !utils.CompareHash(csrf, tokenClaims.Csrf) {
			fmt.Print(jwtError)
			res := utils.ErrorResponse(http.StatusUnauthorized, utils.GetStringPointer("Invalid or expired JWT"), utils.GetStringPointer(customerrors.INVALID_TOKEN))
			c.AbortWithStatusJSON(res.Status, res.Body)
			return
		}
		c.Next()
	}
}
