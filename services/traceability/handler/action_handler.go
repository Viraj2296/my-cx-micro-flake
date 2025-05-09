package handler

import (
	"github.com/gin-gonic/gin"
)

type Assignment struct {
	Event    int `json:"int"`
	Id       int `json:"id"`
	Resource int `json:"resource"`
}

type OrderScheduledEventUpdateRequest struct {
	EndDate    string      `json:"endDate"`
	StartDate  string      `json:"startDate"`
	Assignment *Assignment `json:"assignment"`
}

type ServiceRequestAction struct {
	Remark string `json:"remark"`
}

func (v *TraceabilityService) recordPOSTActionHandler(ctx *gin.Context) {

}
