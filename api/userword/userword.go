package userword

import (
	"encoding/json"
	"net/http"
	"vocablo/customerrors"
	"vocablo/schema"
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
		createdWordJson, err := json.Marshal(createdWord)
		if err != nil {
			res = utils.InternalError(err)
		}
		res = utils.SuccessResponse(createdWordJson)
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
		case customerrors.NotAllowedResourceError:
			res = utils.ErrorResponse(http.StatusForbidden,
				utils.GetStringPointer("You dont have permissions to edit this element"), utils.GetStringPointer(customerrors.NOT_ALLOWED_RESOURCE))
		default:
			res = utils.InternalError(err)
		}
	} else {
		res = utils.SuccessResponse(updatedWord)
	}
	c.JSON(res.Status, res.Body)
}

func Get(c *gin.Context) {
	id, present := c.Params.Get("id")
	if !present {
		res := utils.ErrorResponse(http.StatusBadRequest, utils.GetStringPointer("Mandatory fields not present"), nil)
		c.AbortWithStatusJSON(res.Status, res.Body)
		return
	}
	svc := svc.Get()
	word, err := svc.UserWord.Get(c.Request.Context(), id)
	var res utils.HttpResponse
	if err != nil {
		switch err.(type) {
		case customerrors.NotFoundError:
			res = utils.ErrorResponse(http.StatusNotFound, utils.GetStringPointer("Word not found"), nil)
		default:
			res = utils.InternalError(err)
		}
	} else {
		res = utils.SuccessResponse(word)
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
	//we add to the form the user id of the logged user
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
	id, present := c.Params.Get("id")
	if !present {
		res := utils.ErrorResponse(http.StatusBadRequest, utils.GetStringPointer("Mandatory fields not present"), nil)
		c.AbortWithStatusJSON(res.Status, res.Body)
		return
	}
	svc := svc.Get()
	err := svc.UserWord.Delete(c.Request.Context(), id)
	var res utils.HttpResponse
	if err != nil {
		switch err.(type) {
		case customerrors.NotAllowedResourceError:
			res = utils.ErrorResponse(http.StatusForbidden, utils.GetStringPointer("You dont have permissions to delete this element"),
				utils.GetStringPointer(customerrors.NOT_ALLOWED_RESOURCE))
		default:
			res = utils.InternalError(err)
		}
	} else {
		res = utils.SuccessResponse(nil)
	}
	c.JSON(res.Status, res.Body)
}

func UserProgress(c *gin.Context) {
	svc := svc.Get()

	//We search the total words of the user
	searchForm := userword.SearchForm{Count: true}

	totalWordsPage, err := svc.UserWord.Search(c.Request.Context(), searchForm)
	if err != nil {
		res := utils.InternalError(err)
		c.JSON(res.Status, res.Body)
		return
	}

	//We search the learned words of the user
	searchForm.Learned = utils.GetBoolPointer(true)
	learnedWordsPage, err := svc.UserWord.Search(c.Request.Context(), searchForm)
	if err != nil {
		res := utils.InternalError(err)
		c.JSON(res.Status, res.Body)
		return
	}

	//We search the not learned words of the user
	searchForm.Learned = utils.GetBoolPointer(false)
	notLearnedWordsPage, err := svc.UserWord.Search(c.Request.Context(), searchForm)
	if err != nil {
		res := utils.InternalError(err)
		c.JSON(res.Status, res.Body)
		return
	}

	var userProgress schema.UserWordProgress = schema.UserWordProgress{
		TotalWords:     totalWordsPage.NElements,
		LearnedWords:   learnedWordsPage.NElements,
		UnlearnedWords: notLearnedWordsPage.NElements}

	res := utils.SuccessResponse(userProgress)
	c.JSON(res.Status, res.Body)
}
