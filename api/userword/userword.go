package userword

import (
	"net/http"
	"vocablo/customerrors"
	"vocablo/svc"
	"vocablo/svc/userword"
	"vocablo/utils"

	"github.com/gin-gonic/gin"
)

func Create(c *gin.Context) {
	var form userword.CreateForm
	err := c.ShouldBind(&form)
	if err != nil {
		res := utils.ErrorResponse(http.StatusBadRequest, utils.GetStringPointer(err.Error()), nil)
		c.AbortWithStatusJSON(res.Status, res.Body)
		return
	}
	svc := svc.Get()
	createdWord, err := svc.UserWord.Create(c.Request.Context(), form)
	var res utils.HttpResponse
	if err != nil {
		switch err.(type) {
		case customerrors.EmptyFormFieldsError:
			res = utils.ErrorResponse(http.StatusBadRequest, utils.GetStringPointer("Mandatory fields not present"), nil)
		default:
			res = utils.InternalError(err)
		}
	} else {

		res = utils.SuccessResponse(createdWord)
	}
	c.JSON(res.Status, res.Body)
}

func Update(c *gin.Context) {
	var form userword.UpdateForm
	err := c.ShouldBind(&form)
	if err != nil {
		res := utils.ErrorResponse(http.StatusBadRequest, utils.GetStringPointer(err.Error()), nil)
		c.AbortWithStatusJSON(res.Status, res.Body)
		return
	}
	svc := svc.Get()
	updatedWord, err := svc.UserWord.Update(c.Request.Context(), form)
	var res utils.HttpResponse
	if err != nil {
		switch err.(type) {
		case customerrors.EmptyFormFieldsError:
			res = utils.ErrorResponse(http.StatusBadRequest, utils.GetStringPointer("Mandatory fields not present"), nil)
		default:
			res = utils.InternalError(err)
		}
	} else {
		res = utils.SuccessResponse(updatedWord)
	}
	c.JSON(res.Status, res.Body)
}

func Search(c *gin.Context) {
	var form userword.SearchForm
	err := c.ShouldBind(&form)
	if err != nil {
		res := utils.ErrorResponse(http.StatusBadRequest, utils.GetStringPointer(err.Error()), nil)
		c.AbortWithStatusJSON(res.Status, res.Body)
		return
	}
	svc := svc.Get()
	page, err := svc.UserWord.Search(c.Request.Context(), form)
	var res utils.HttpResponse
	if err != nil {
		res = utils.InternalError(err)
	} else {
		res = utils.SuccessResponse(page)
	}
	c.JSON(res.Status, res.Body)
}

func Delete(c *gin.Context) {
	id, error := c.Params.Get("id")
	if error {
		res := utils.ErrorResponse(http.StatusBadRequest, utils.GetStringPointer("Mandatory fields not present"), nil)
		c.AbortWithStatusJSON(res.Status, res.Body)
		return
	}
	svc := svc.Get()
	err := svc.UserWord.Delete(c.Request.Context(), id)
	var res utils.HttpResponse
	if err != nil {
		res = utils.InternalError(err)
	} else {
		res = utils.SuccessResponse(nil)
	}
	c.JSON(res.Status, res.Body)
}
