package middlewares

import (
	"bytes"
	"cx-micro-flake/pkg/auth"
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func logRequestDetails(r *http.Request) {
	fmt.Println("----- HTTP Request Details -----")
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL.String())

	// Print headers
	fmt.Println("Headers:")
	for name, values := range r.Header {
		for _, value := range values {
			fmt.Printf("  %s: %s\n", name, value)
		}
	}

	// Print query params
	fmt.Println("Query Parameters:")
	for key, values := range r.URL.Query() {
		for _, value := range values {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}

	// Optionally print body (if small and not already read)
	if r.Body != nil && r.ContentLength < 1024*10 { // limit to 10KB
		bodyBytes, err := io.ReadAll(r.Body)
		if err == nil {
			fmt.Println("Body:")
			fmt.Println(string(bodyBytes))
			// Rewind body for next handlers
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
	}
	fmt.Println("----- End Request Details -----")
}
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, DELETE,OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func TokenAuthMiddleware() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		tokenString, err := auth.TokenValid(ctx.Request)
		if err != nil {
			fmt.Println("token validation has failed ", err.Error(), "token string value ", tokenString)
			logRequestDetails(ctx.Request)
			response.AbortWithTokenError(ctx, http.StatusUnauthorized, common.InvalidTokenOrTokenMissing, errors.New("token validation has failed ["+err.Error()+"]"))
			return
		}
		userId, err := auth.ExtractResourceId(tokenString)
		if err != nil {
			fmt.Println("invalid token format, extracting user id had failed ", err.Error(), "token string value ", tokenString, "user id ", userId)
			response.AbortWithError(ctx, http.StatusUnauthorized, common.ExtractingUserIdHadFailedFromToken, errors.New("invalid token format, extracting user id had failed"))
			return
		}
		ctx.Set("id", userId)
		ctx.Set("token", tokenString)
	}
}

func TokenComponentResourceMiddleware() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		tokenString, err := auth.TokenValid(ctx.Request)
		if err != nil {
			response.AbortWithError(ctx, http.StatusUnauthorized, common.InvalidTokenOrTokenMissing, errors.New("invalid token or token missing"))
			return
		}

		projectId, err := auth.ExtractProjectId(tokenString)
		if err != nil {
			response.AbortWithError(ctx, http.StatusUnauthorized, common.ExtractingUserIdHadFailedFromToken, errors.New("invalid token format, extracting project id had failed"))
			return
		}
		componentName, err := auth.ExtractComponentName(tokenString)
		if err != nil {
			response.AbortWithError(ctx, http.StatusUnauthorized, common.ExtractingUserIdHadFailedFromToken, errors.New("invalid token format, extracting component  had failed"))
			return
		}
		resourceId, err := auth.ExtractResourceId(tokenString)
		if err != nil {
			response.AbortWithError(ctx, http.StatusUnauthorized, common.ExtractingUserIdHadFailedFromToken, errors.New("invalid token format, extracting resource id had failed"))
			return
		}
		ctx.Set("projectId", projectId)
		ctx.Set("componentName", componentName)
		ctx.Set("id", resourceId)

	}

}
func PermissionMiddleware(moduleName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method
		userId := common.GetUserId(ctx)
		authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
		userType := authService.GetUserType(userId)
		if userType == "system_admin" {
			return
		}
		componentName := util.GetComponentName(ctx)
		recordId := util.GetOriginalRecordId(ctx)
		action := util.GetActionName(ctx)
		projectId := util.GetProjectId(ctx)
		var isAllowed bool
		urlPath := ctx.Request.URL.Path
		fmt.Println("component name ", componentName, " action :", action, "module name :", moduleName, " project id :", projectId, " record id ", recordId, " urlPath :", urlPath)
		isAllowed, isRouteEnabled, routingComponent := authService.IsAllowed(userId, projectId, moduleName, componentName, action, recordId, method, urlPath)
		fmt.Println("isAllowed", isAllowed, " isRouteEnabled :", isRouteEnabled, " routingComponent :", routingComponent)
		return
		//if !isAllowed {
		//	if method == "GET" {
		//		response.AbortWithAccessDeniedError(ctx, http.StatusBadRequest, errors.New("Access Denied"), common.AccessDenied, common.AccessDeniedDescriptionForGetResources, isRouteEnabled, routingComponent)
		//	}
		//	if method == "PUT" {
		//		response.AbortWithAccessDeniedError(ctx, http.StatusBadRequest, errors.New("Access Denied"), common.AccessDenied, common.AccessDeniedDescriptionForUpdateResources, isRouteEnabled, routingComponent)
		//	}
		//	if method == "POST" {
		//		if action == "" {
		//			response.AbortWithAccessDeniedError(ctx, http.StatusBadRequest, errors.New("Access Denied"), common.AccessDenied, common.AccessDeniedDescriptionForCreateResources, isRouteEnabled, routingComponent)
		//		} else {
		//			response.AbortWithAccessDeniedError(ctx, http.StatusBadRequest, errors.New("Access Denied"), common.AccessDenied, common.AccessDeniedDescriptionForActionResources, isRouteEnabled, routingComponent)
		//		}
		//
		//	}
		//	if method == "DELETE" {
		//		response.AbortWithAccessDeniedError(ctx, http.StatusBadRequest, errors.New("Access Denied"), common.AccessDenied, common.AccessDeniedDescriptionForDeleteResources, isRouteEnabled, routingComponent)
		//	}
		//	return
		//}
	}
}
