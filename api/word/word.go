package word

import (
	"net/http"
	"vocablo/customerrors"
	"vocablo/svc"
	"vocablo/utils"

	"github.com/gin-gonic/gin"
)

func Search(c *gin.Context) {
	term, present := c.Params.Get("term")
	if !present {
		res := utils.ErrorResponse(http.StatusBadRequest, utils.GetStringPointer("Mandatory fields are empty"), nil)
		c.JSON(res.Status, res.Body)
		return
	}
	lang, present := c.Params.Get("lang")
	if !present {
		res := utils.ErrorResponse(http.StatusBadRequest, utils.GetStringPointer("Mandatory fields are empty"), nil)
		c.JSON(res.Status, res.Body)
		return
	}

	svc := svc.Get()
	words, err := svc.Word.Search(c.Request.Context(), term, lang)
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
