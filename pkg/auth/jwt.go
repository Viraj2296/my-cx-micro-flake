package auth

import (
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte(os.Getenv("api_secret"))

func createTokenWithClaims(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func IsTokenExpired(r *http.Request) bool {
	tokenString := ExtractToken(r)
	_, _, err := parseToken(tokenString)
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			return ve.Errors&jwt.ValidationErrorExpired != 0
		}
	}
	return false
}
func CreateEmailInvitationToken(userId string) (string, error) {
	claims := jwt.MapClaims{
		"authorized": true,
		"id":         userId,
		"exp":        time.Now().Add(48 * time.Hour).Unix(),
	}
	return createTokenWithClaims(claims)
}

func CreateToken(userId int, defaultTokenExpiry float64) (string, error) {
	claims := jwt.MapClaims{
		"authorized": true,
		"id":         userId,
		"exp":        time.Now().Add(time.Duration(defaultTokenExpiry) * time.Hour).Unix(),
	}
	return createTokenWithClaims(claims)
}

func CreateInfToken(userId int) (string, error) {
	claims := jwt.MapClaims{
		"authorized": true,
		"id":         userId,
		"exp":        time.Now().Add(10 * 365 * 24 * time.Hour).Unix(),
	}
	return createTokenWithClaims(claims)
}

func CreateCustomToken(projectId, zone string, tokenClaims map[string]interface{}) (string, error) {
	claims := jwt.MapClaims{}
	for k, v := range tokenClaims {
		claims[k] = v
	}

	if endDateValue, ok := tokenClaims["endDate"]; ok {
		endTime := util.ConvertUserTimezoneToUTC(zone, util.InterfaceToString(endDateValue))
		expiryHours, err := util.CalculateHoursFromCurrentTime(endTime)
		if err == nil {
			claims["exp"] = time.Now().Add(time.Duration(expiryHours) * time.Hour).Unix()
		} else {
			claims["exp"] = time.Now().Add(3600 * time.Hour).Unix()
		}
	}
	return createTokenWithClaims(claims)
}

func CreateResourceToken(projectId, componentName string, recordId int) (string, error) {
	claims := jwt.MapClaims{
		"authorized":    true,
		"id":            recordId,
		"componentName": componentName,
		"projectId":     projectId,
		"exp":           time.Now().Add(3600 * time.Hour).Unix(),
	}
	return createTokenWithClaims(claims)
}

func CreateResourceInfToken(projectId, componentName string, recordId int) (string, error) {
	claims := jwt.MapClaims{
		"authorized":    true,
		"id":            recordId,
		"componentName": componentName,
		"projectId":     projectId,
		"exp":           time.Now().Add(10 * 365 * 24 * time.Hour).Unix(),
	}
	return createTokenWithClaims(claims)
}

func CreateRefreshToken(userId int, platform string, expiry float64) (string, error) {
	claims := jwt.MapClaims{
		"authorized":       true,
		"is_refresh_token": true,
		"id":               userId,
		"platform":         platform,
		"exp":              time.Now().Add(time.Duration(expiry) * time.Hour).Unix(),
	}
	return createTokenWithClaims(claims)
}

func CreateRefreshInfToken(userId int, platform string) (string, error) {
	claims := jwt.MapClaims{
		"authorized":       true,
		"is_refresh_token": true,
		"id":               userId,
		"platform":         platform,
		"exp":              time.Now().Add(10 * 365 * 24 * time.Hour).Unix(),
	}
	return createTokenWithClaims(claims)
}
func parseToken(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	if err != nil {
		return nil, nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, nil, fmt.Errorf("invalid token")
	}
	return token, claims, nil
}

func ExtractComponentName(tokenString string) (string, error) {
	_, claims, err := parseToken(tokenString)
	if err != nil {
		return "", err
	}
	return util.InterfaceToString(claims["componentName"]), nil
}

func ExtractResourceId(tokenString string) (int, error) {
	_, claims, err := parseToken(tokenString)
	if err != nil {
		return 0, err
	}
	return util.InterfaceToInt(claims["id"]), nil
}

func ExtractProjectId(tokenString string) (string, error) {
	_, claims, err := parseToken(tokenString)
	if err != nil {
		return "", err
	}
	return util.InterfaceToString(claims["projectId"]), nil
}

func ExtractToken(r *http.Request) string {
	if token := r.URL.Query().Get("token"); token != "" {
		return token
	}
	auth := r.Header.Get("Authorization")
	parts := strings.Split(auth, " ")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}

func TokenValid(r *http.Request) (string, error) {
	tokenStr := ExtractToken(r)
	_, _, err := parseToken(tokenStr)
	return tokenStr, err
}

func IsTokenStringValid(tokenString string) (bool, error) {
	_, _, err := parseToken(tokenString)
	if err != nil {
		return false, err
	}
	return true, err
}
func ExtractRefreshTokenID(r *http.Request) (int, error) {
	tokenStr := ExtractToken(r)
	_, claims, err := parseToken(tokenStr)
	if err != nil {
		return 0, err
	}
	if claims["is_refresh_token"] != true {
		return 0, fmt.Errorf("invalid refresh token")
	}
	return util.InterfaceToInt(claims["id"]), nil
}

func ExtractRefreshTokenPlatform(r *http.Request) (string, error) {
	tokenStr := ExtractToken(r)
	_, claims, err := parseToken(tokenStr)
	if err != nil {
		return "", err
	}
	return util.InterfaceToString(claims["platform"]), nil
}

func Pretty(data interface{}) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(b))
}
