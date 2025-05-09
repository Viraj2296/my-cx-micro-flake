package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

var (
	InvalidRequestFormat = "invalid request format"
)

func ValidateRequest(data interface{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := ctx.ShouldBindBodyWith(&data, binding.JSON); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		ctx.Next()
	}

}
