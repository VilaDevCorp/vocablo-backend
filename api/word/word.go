package word

import (
	"net/http"
	"vocablo/customerrors"
	"vocablo/svc"
	"vocablo/svc/word"
	"vocablo/utils"

	"github.com/gin-gonic/gin"
)

func Search(c *gin.Context) {
	var form word.SearchForm
	err := c.ShouldBind(&form)
	if err != nil {
		res := utils.ErrorResponse(http.StatusBadRequest, utils.GetStringPointer("Mandatory fields are empty"), nil)
		c.JSON(res.Status, res.Body)
		return
	}
	if form.Term == "" || form.Lang == "" {
		res := utils.ErrorResponse(http.StatusBadRequest, utils.GetStringPointer("Mandatory fields are empty"), nil)
		c.JSON(res.Status, res.Body)
		return
	}

	svc := svc.Get()
	words, err := svc.Word.Search(c.Request.Context(), form.Lang, form.Term)
	var res utils.HttpResponse
	if err != nil {
		switch err.(type) {
		case customerrors.EmptyFormFieldsError:
			res = utils.ErrorResponse(http.StatusBadRequest, utils.GetStringPointer("Mandatory fields are empty"), nil)
		default:
			res = utils.InternalError(err)
		}
	} else {
		res = utils.SuccessResponse(words)
	}
	c.JSON(res.Status, res.Body)
}
