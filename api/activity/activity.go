package activity

// import (
// 	"appname/svc"
// 	"appname/svc/activity"
// 	"appname/utils"
// 	"fmt"

// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// )

// func Create(c *gin.Context) {
// 	var form activity.CreateForm
// 	err := c.ShouldBind(&form)
// 	if err != nil {
// 		res := utils.BadRequest(form, err)
// 		c.AbortWithStatusJSON(res.Status, res.Body)
// 		return
// 	}
// 	svc := svc.Get()
// 	result, err := svc.Activity.Create(c.Request.Context(), form)
// 	if err != nil {
// 		res := utils.InternalError(err)
// 		c.AbortWithStatusJSON(res.Status, res.Body)
// 		return
// 	}
// 	res := utils.OkCreated(result)
// 	c.JSON(res.Status, res.Body)
// }

// func Update(c *gin.Context) {
// 	var form activity.UpdateForm
// 	err := c.ShouldBind(&form)
// 	if err != nil {
// 		res := utils.BadRequest(form, err)
// 		c.AbortWithStatusJSON(res.Status, res.Body)
// 		return
// 	}
// 	svc := svc.Get()
// 	result, err := svc.Activity.Update(c.Request.Context(), form)
// 	if err != nil {
// 		res := utils.InternalError(err)
// 		c.AbortWithStatusJSON(res.Status, res.Body)
// 		return
// 	}
// 	res := utils.OkUpdated(result)
// 	c.JSON(res.Status, res.Body)
// }

// func Get(c *gin.Context) {
// 	unparsedId, _ := c.Params.Get("id")
// 	parsedId, err := uuid.Parse(unparsedId)
// 	if err != nil {
// 		res := utils.BadRequest(unparsedId, err)
// 		c.AbortWithStatusJSON(res.Status, res.Body)
// 		return
// 	}
// 	svc := svc.Get()
// 	result, err := svc.Activity.Get(c.Request.Context(), parsedId)
// 	if err != nil {
// 		res := utils.InternalError(err)
// 		c.AbortWithStatusJSON(res.Status, res.Body)
// 		return
// 	}
// 	res := utils.OkGet(result)
// 	c.JSON(res.Status, res.Body)
// }

// func Search(c *gin.Context) {
// 	var form activity.SearchForm
// 	err := c.ShouldBind(&form)
// 	if err != nil {
// 		res := utils.BadRequest(form, err)
// 		c.AbortWithStatusJSON(res.Status, res.Body)
// 		return
// 	}
// 	svc := svc.Get()
// 	result, err := svc.Activity.Search(c.Request.Context(), form)
// 	if err != nil {
// 		res := utils.InternalError(err)
// 		fmt.Print(res)
// 		c.AbortWithStatusJSON(res.Status, res.Body)
// 		return
// 	}
// 	res := utils.OkGet(result)
// 	c.JSON(res.Status, res.Body)
// }

// func Delete(c *gin.Context) {
// 	var form activity.DeleteForm
// 	err := c.ShouldBind(&form)

// 	if err != nil {
// 		res := utils.BadRequest(form, err)
// 		c.AbortWithStatusJSON(res.Status, res.Body)
// 		return
// 	}
// 	activityIds := form.ActivityIds

// 	svc := svc.Get()
// 	err = svc.Activity.Delete(c.Request.Context(), activityIds)
// 	if err != nil {
// 		res := utils.InternalError(err)
// 		c.AbortWithStatusJSON(res.Status, res.Body)
// 		return
// 	}
// 	res := utils.OkDeleted()
// 	c.JSON(res.Status, res.Body)
// }
