package quiz

import (
	"net/http"
	"vocablo/customerrors"
	"vocablo/svc"
	"vocablo/svc/quiz"
	"vocablo/utils"

	"github.com/gin-gonic/gin"
)

func Create(c *gin.Context) {
	var form quiz.CreateForm
	err := c.ShouldBind(&form)
	if err != nil {
		res := utils.ErrorResponse(http.StatusBadRequest, utils.GetStringPointer(err.Error()), nil)
		c.AbortWithStatusJSON(res.Status, res.Body)
		return
	}

	svc := svc.Get()
	createdQuiz, err := svc.Quiz.Create(c.Request.Context(), form)
	var res utils.HttpResponse
	if err != nil {
		switch err.(type) {
		case customerrors.NotEnoughWordsForQuizError:
			res = utils.ErrorResponse(http.StatusConflict, utils.GetStringPointer(err.Error()),
				utils.GetStringPointer(customerrors.NOT_ENOUGH_WORDS_FOR_QUIZ))
		default:
			res = utils.InternalError(err)
		}
	} else {
		res = utils.SuccessResponse(createdQuiz)
	}
	c.JSON(res.Status, res.Body)
}

func Answer(c *gin.Context) {
	var filledQuiz quiz.Quiz
	err := c.ShouldBind(&filledQuiz)
	if err != nil {
		res := utils.ErrorResponse(http.StatusBadRequest, utils.GetStringPointer(err.Error()), nil)
		c.AbortWithStatusJSON(res.Status, res.Body)
		return
	}

	svc := svc.Get()
	score, err := svc.Quiz.Answer(c.Request.Context(), filledQuiz)
	var res utils.HttpResponse
	if err != nil {
		switch err.(type) {
		default:
			res = utils.InternalError(err)
		}
	} else {
		res = utils.SuccessResponse(score)
	}
	c.JSON(res.Status, res.Body)
}
