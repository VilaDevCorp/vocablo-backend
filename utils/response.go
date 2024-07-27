package utils

import (
	"fmt"
	"net/http"
)

type ResponseBody struct {
	Data         interface{} `json:"data"`
	ErrorMessage *string     `json:"errorMessage"`
	ErrorCode    *string     `json:"errorCode"`
}

type HttpResponse struct {
	Status int          `json:"status"`
	Body   ResponseBody `json:"body"`
}

func SuccessResponse(data interface{}) HttpResponse {
	return HttpResponse{Status: http.StatusOK, Body: ResponseBody{Data: data}}
}

func ErrorResponse(status int, errorMessage *string, errorCode *string) HttpResponse {
	return HttpResponse{Status: status, Body: ResponseBody{ErrorMessage: errorMessage, ErrorCode: errorCode}}
}

func InternalError(err error) HttpResponse {
	return ErrorResponse(http.StatusInternalServerError, GetStringPointer(err.Error()), nil)
}

func NotFound(entity string, entityId string) HttpResponse {
	return HttpResponse{Status: http.StatusNotFound, Body: ResponseBody{ErrorMessage: GetStringPointer(fmt.Sprintf("%s %s not found", entity, entityId))}}
}
