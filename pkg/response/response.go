package response

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	ResourceNotFoundErrorCode = 1000
	ResourceAlreadyArchived   = 1001
	InternalSystemErrorCode   = 1002
	ResourceNotFound          = "Resource Not Found"
	InternalSystemError       = "Internal System Error"
	ObjectAlreadyArchived     = "Resource Already Archived"
)

func SendResourceNotFound(ctx *gin.Context) {
	SendDetailedError(ctx, http.StatusBadRequest, errors.New(ResourceNotFound), ResourceNotFoundErrorCode, "The resource processing is not found, this is due to invalid resource id")
}

func SendObjectCreationMessage(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, GeneralResponse{
		Code:    200,
		Message: "The resource is successfully created",
	})
}
func SendInternalSystemError(ctx *gin.Context) {
	SendDetailedError(ctx, http.StatusBadRequest, errors.New(InternalSystemError), InternalSystemErrorCode, "A malfunction has occurred within the API system. Please try again later or contact support if the issue persists")
}
func SendAlreadyArchivedError(ctx *gin.Context) {
	SendDetailedError(ctx, http.StatusBadRequest, errors.New(ObjectAlreadyArchived), ResourceAlreadyArchived, "You are trying to modify archived object, This modification is rejected")
	return
}

func SendSimpleError(ctx *gin.Context, status int, err error, errorCode int) {
	httpError := simpleError{
		Code:    errorCode,
		Message: err.Error(),
	}

	ctx.JSON(status, httpError)
}

func SendDetailedError(ctx *gin.Context, status int, err error, errorCode int, description string) {
	httpError := detailedError{
		Code:        errorCode,
		Message:     err.Error(),
		Description: description,
	}

	ctx.JSON(status, httpError)
}

func DispatchDetailedError(ctx *gin.Context, errorCode int, errorDetail *DetailedError) {
	httpError := detailedError{
		Code:        errorCode,
		Message:     errorDetail.Header,
		Description: errorDetail.Description,
	}

	ctx.JSON(http.StatusBadRequest, httpError)
}
func SendValidationError(ctx *gin.Context, status int, err error, errorCode int, description string) {
	httpError := authorisationError{
		Code:             errorCode,
		Message:          err.Error(),
		Description:      description,
		IsRouteEnabled:   true,
		RoutingComponent: "module_selection",
	}

	ctx.JSON(status, httpError)
}

func AbortWithAccessDeniedError(ctx *gin.Context, status int, err error, errorCode int, description string, isRouteEnabled bool, routingComponent string) {
	httpError := authorisationError{
		Code:                errorCode,
		Message:             err.Error(),
		Description:         description,
		IsRouteEnabled:      isRouteEnabled,
		RoutingComponent:    routingComponent,
		AuthorisationFailed: true,
	}
	ctx.AbortWithStatusJSON(status, httpError)
}

func AbortWithError(ctx *gin.Context, status int, errorCode int, err error) {
	er := simpleError{
		Code:    errorCode,
		Message: err.Error(),
	}
	ctx.AbortWithStatusJSON(status, er)
}

func AbortWithTokenError(ctx *gin.Context, status int, errorCode int, err error) {
	er := tokenExpiredError{
		Code:      errorCode,
		DoRefresh: true,
		Message:   err.Error(),
	}
	ctx.AbortWithStatusJSON(status, er)
}

type simpleError struct {
	Code    int    `json:"code" example:"345"`
	Message string `json:"message" example:"exception happened during upstream parser"`
}

type tokenExpiredError struct {
	Code      int    `json:"code" example:"345"`
	DoRefresh bool   `json:"doRefresh"`
	Message   string `json:"message" example:"exception happened during upstream parser"`
}

type detailedError struct {
	Code        int    `json:"code" example:"345"`
	Message     string `json:"message" example:"exception happened during upstream parser"`
	Description string `json:"description" example:"unable to process this request, please try again later"`
}

type authorisationError struct {
	Code                int    `json:"code" example:"345"`
	Message             string `json:"message" example:"exception happened during upstream parser"`
	IsRouteEnabled      bool   `json:"isRouteEnabled" example:"exception happened during upstream parser"`
	AuthorisationFailed bool   `json:"authorisationFailed"`
	RoutingComponent    string `json:"routingComponent" example:"exception happened during upstream parser"`
	Description         string `json:"description" example:"unable to process this request, please try again later"`
}

type GeneralResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	RecordId  int    `json:"recordId"`
	CanDelete bool   `json:"canDelete"`
}

type ValidationResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	CanDelete bool   `json:"canDelete"`
}

type DetailedError struct {
	Header      string
	Description string
}
