package handler

import (
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/iot/pkg/consts"
	"go.uber.org/zap"
	"net/http"
	"os/exec"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func (v *IoTService) overview(ctx *gin.Context) {
	// get the total machines
	dbConnection := v.BaseService.ServiceDatabase
	totalMessagesQuery := "select count(*) as total_messages from message"
	totalMessageSizeQuery := "select sum(char_length(body)) as total_message_length FROM message"
	totalUniqueTopicsQuery := "select count(distinct(topic)) as total_unique_topics  from message"
	lastMessageReceivedTimeQuery := "select ts as last_message_received_time from message order by ts desc limit 1"

	var queryResults map[string]interface{}

	var totalMessages int
	var totalMessageSize int
	var totalUniqueTopics int
	var lastMessageReceivedTime int

	dbConnection.Raw(totalMessagesQuery).Scan(&queryResults)
	totalMessages = util.InterfaceToInt(queryResults["total_messages"])
	dbConnection.Raw(totalMessageSizeQuery).Scan(&queryResults)
	totalMessageSize = util.InterfaceToInt(queryResults["total_message_length"])
	dbConnection.Raw(totalUniqueTopicsQuery).Scan(&queryResults)
	totalUniqueTopics = util.InterfaceToInt(queryResults["total_unique_topics"])
	dbConnection.Raw(lastMessageReceivedTimeQuery).Scan(&queryResults)
	lastMessageReceivedTime = util.InterfaceToInt(queryResults["last_message_received_time"])

	var response = make(map[string]interface{}, 100)
	//[{"nodeName":"PUMA TCP Source", }]
	response["nodeName"] = "PUMA TCP Source"
	response["messageCount"] = totalMessages
	response["messageSize"] = strconv.Itoa(totalMessageSize/(1024*1024)) + " MB"
	response["uniqueTopics"] = totalUniqueTopics
	t := time.UnixMilli(int64(lastMessageReceivedTime))
	response["lastReceivedMessageTime"] = t.Format(consts.TimeLayout)
	// t.Format(time.RFC822)
	var overviewResponse []interface{}
	overviewResponse = append(overviewResponse, response)
	ctx.JSON(http.StatusOK, overviewResponse)

}

func (v *IoTService) startBroker() *int {
	cmd := exec.Command("broker")
	_, err := cmd.Output()

	if err != nil {
		v.BaseService.Logger.Error("error starting broker process", zap.String("error", err.Error()))

		return nil
	}

	pid := cmd.Process.Pid
	return &pid
}

func (v *IoTService) stopBroker(pid int) {
	err := syscall.Kill(pid, syscall.SIGKILL)
	if err != nil {
		v.BaseService.Logger.Error("error in stopping broker process", zap.String("error", err.Error()))

		return
	}
}
