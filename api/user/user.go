package user

import (
	"net/http"
	"vocablo/svc"
	"vocablo/svc/user"
	"vocablo/utils"

	"github.com/gin-gonic/gin"
)

func Create(c *gin.Context) {
	var form user.CreateForm
	err := c.ShouldBind(&form)
	if err != nil {
		res := utils.ErrorResponse(http.StatusBadRequest, utils.GetStringPointer(err.Error()), nil)
		c.AbortWithStatusJSON(res.Status, res.Body)
		return
	}
	svc := svc.Get()
	result, err := svc.User.Create(c.Request.Context(), form)
	if err != nil {
		res := utils.InternalError(err)
		c.AbortWithStatusJSON(res.Status, res.Body)
		return
	}
	res := utils.SuccessResponse(result)
	c.JSON(res.Status, res.Body)
}

// func Get(c *gin.Context) {
// 	unparsedId, _ := c.Params.Get("id")
// 	parsedId, err := uuid.Parse(unparsedId)
// 	if err != nil {
// 		res := utils.ErrorResponse(http.StatusBadRequest, utils.GetStringPointer(err.Error()), nil)
// 		c.AbortWithStatusJSON(res.Status, res.Body)
// 		return
// 	}
// 	svc := svc.Get()
// 	result, err := svc.User.Get(c.Request.Context(), parsedId)
// 	if err != nil {
// 		res := utils.InternalError(err)
// 		c.AbortWithStatusJSON(res.Status, res.Body)
// 		return
// 	}
// 	res := utils.SuccessResponse(result)
// 	c.JSON(res.Status, res.Body)
// }

// func Search(c *gin.Context) {
// 	var form user.SearchForm
// 	err := c.ShouldBind(&form)
// 	if err != nil {
// 		res := utils.ErrorResponse(http.StatusBadRequest, utils.GetStringPointer(err.Error()), nil)
// 		c.AbortWithStatusJSON(res.Status, res.Body)
// 		return
// 	}
// 	svc := svc.Get()
// 	result, err := svc.User.Search(c.Request.Context(), form)
// 	if err != nil {
// 		res := utils.InternalError(err)
// 		c.AbortWithStatusJSON(res.Status, res.Body)
// 		return
// 	}
// 	res := utils.SuccessResponse(result)
// 	c.JSON(res.Status, res.Body)
// }

// func Delete(c *gin.Context) {
// 	unparsedId, _ := c.Params.Get("id")
// 	parsedId, err := uuid.Parse(unparsedId)
// 	if err != nil {
// 		res := utils.ErrorResponse(http.StatusBadRequest, utils.GetStringPointer(err.Error()), nil)
// 		c.AbortWithStatusJSON(res.Status, res.Body)
// 		return
// 	}
// 	svc := svc.Get()
// 	err = svc.User.Delete(c.Request.Context(), parsedId)
// 	if err != nil {
// 		res := utils.InternalError(err)
// 		c.AbortWithStatusJSON(res.Status, res.Body)
// 		return
// 	}
// 	res := utils.SuccessResponse(nil)
// 	c.JSON(res.Status, res.Body)
// }
