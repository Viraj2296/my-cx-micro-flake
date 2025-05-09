package workflow_actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/labour_management/handler/cache"
	"cx-micro-flake/services/labour_management/handler/const_util"
	"cx-micro-flake/services/labour_management/handler/database"
	"cx-micro-flake/services/labour_management/handler/notification"
	"encoding/json"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AssemblyMachineMasterInfo struct {
	Area                         string    `json:"area"`
	Level                        string    `json:"level"`
	CreatedAt                    time.Time `json:"createdAt"`
	CreatedBy                    int       `json:"createdBy"`
	Department                   int       `json:"department"`
	DelayPeriod                  int       `json:"delayPeriod"`
	DelayStatus                  string    `json:"delayStatus"`
	Description                  string    `json:"description"`
	MessageFlag                  string    `json:"messageFlag"`
	MachineImage                 string    `json:"machineImage"`
	NewMachineId                 string    `json:"newMachineId"`
	ObjectStatus                 string    `json:"objectStatus"`
	LastUpdatedAt                time.Time `json:"lastUpdatedAt"`
	LastUpdatedBy                int       `json:"lastUpdatedBy"`
	MachineStatus                int       `json:"machineStatus"`
	AssemblyLineOption           int       `json:"assemblyLineOption"`
	AssemblyLineTypeOption       int       `json:"assemblyLineTypeOption"`
	CanCreateWorkOrder           bool      `json:"canCreateWorkOrder"`
	HelpButtonStationNo          string    `json:"helpButtonStationNo"`
	MachineConnectStatus         int       `json:"machineConnectStatus"`
	LastUpdatedMachineLiveStatus time.Time `json:"lastUpdatedMachineLiveStatus"`
	Model                        string    `json:"model"`
	IsMESDriverConfigured        bool      `json:"isMESDriverConfigured"`
}

func GetAssemblyMachineMasterInfo(eventObject datatypes.JSON) *AssemblyMachineMasterInfo {
	assemblyMachineMasterInfo := AssemblyMachineMasterInfo{}
	json.Unmarshal(eventObject, &assemblyMachineMasterInfo)
	return &assemblyMachineMasterInfo
}

type AssemblyScheduledOrderEventInfo struct {
	Name            string    `json:"name"`
	EndDate         time.Time `json:"endDate"`
	IconCls         string    `json:"iconCls"`
	IsUpdate        bool      `json:"isUpdate"`
	Draggable       bool      `json:"draggable"`
	EventType       string    `json:"eventType"`
	MachineId       int       `json:"machineId"`
	StartDate       time.Time `json:"startDate"`
	EventColor      string    `json:"eventColor"`
	EventStatus     int       `json:"eventStatus"`
	PercentDone     int       `json:"percentDone"`
	RejectedQty     int       `json:"rejectedQty"`
	CanForceStop    bool      `json:"canForceStop"`
	CompletedQty    int       `json:"completedQty"`
	ObjectStatus    string    `json:"objectStatus"`
	ScheduledQty    int       `json:"scheduledQty"`
	EventSourceId   int       `json:"eventSourceId"`
	LastUpdatedAt   time.Time `json:"lastUpdatedAt"`
	LastUpdatedBy   int       `json:"lastUpdatedBy"`
	IsAbortEnabled  bool      `json:"isAbortEnabled"`
	ProductionOrder string    `json:"productionOrder"`
	PlannedManPower int       `json:"plannedManPower"`
	PriorityLevel   int       `json:"priorityLevel"`
	Line            string    `json:"line"`
	Remarks         string    `json:"remarks"`
}

func GetAssemblyScheduledOrderEventInfo(serialisedData datatypes.JSON) *AssemblyScheduledOrderEventInfo {
	assemblyProductionOrderInfo := AssemblyScheduledOrderEventInfo{}
	err := json.Unmarshal(serialisedData, &assemblyProductionOrderInfo)
	if err != nil {
		return &AssemblyScheduledOrderEventInfo{}
	}
	return &assemblyProductionOrderInfo
}

type AssemblyProductionOrderInfo struct {
	OrderStatus                      int     `json:"orderStatus"`
	Balance                          int     `json:"balance"`
	MouldId                          int     `json:"mouldId"`
	ProdQty                          int     `json:"prodQty"`
	Remarks                          string  `json:"remarks"`
	OrderQty                         int     `json:"orderQty"`
	PartName                         string  `json:"partName"`
	CycleTime                        float64 `json:"cycleTime"`
	DailyRate                        int     `json:"dailyRate"`
	MachineId                        int     `json:"machineId"`
	ProdOrder                        string  `json:"prodOrder"`
	PartNumber                       int     `json:"partNumber"`
	WorkCenter                       string  `json:"workCenter"`
	RecycleMaterial                  string  `json:"recycleMaterial"`
	BomTextPercentageOfMaterialUsage float64 `json:"bomTextPercentageOfMaterialUsage"`
	RemainingScheduledQty            int     `json:"remainingScheduledQty"`
	ObjectStatus                     string  `json:"objectStatus"`
	CreatedAt                        string  `json:"createdAt"`
	CreatedBy                        int     `json:"createdBy"`
	LastUpdatedAt                    string  `json:"lastUpdatedAt"`
	LastUpdatedBy                    int     `json:"lastUpdatedBy"`
}

func GetAssemblyProductionOrderInfo(serialisedData datatypes.JSON) *AssemblyProductionOrderInfo {
	assemblyProductionOrderInfo := AssemblyProductionOrderInfo{}
	err := json.Unmarshal(serialisedData, &assemblyProductionOrderInfo)
	if err != nil {
		return &AssemblyProductionOrderInfo{}
	}
	return &assemblyProductionOrderInfo
}

type PartInfo struct {
	Image          string    `json:"image"`
	PartNumber     string    `json:"partNumber"`
	Description    string    `json:"description"`
	LastUpdatedAt  time.Time `json:"lastUpdatedAt"`
	LastUpdatedBy  int       `json:"lastUpdatedBy"`
	IsBatchManaged bool      `json:"isBatchManaged"`
}

func GetPartInfo(eventObject datatypes.JSON) *PartInfo {
	partInfo := PartInfo{}
	err := json.Unmarshal(eventObject, &partInfo)
	if err != nil {
		return &PartInfo{}
	}
	return &partInfo
}

type ActionService struct {
	Logger                      *zap.Logger
	Database                    *gorm.DB
	EmailHandler                *notification.EmailHandler
	LabourManagementSettingInfo *database.LabourManagementSettingInfo
	ShiftActualValueCache       map[int][]database.ActualShiftValueCache
	MachineStatsCache           *cache.MachineStatsCache
}

func (v *ActionService) Init(logger *zap.Logger) {
	v.Logger = logger
}
func getTimeDifference(dst string) string {
	currentTime := util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	var difference = util.ConvertStringToDateTime(currentTime).DateTimeEpoch - util.ConvertStringToDateTime(dst).DateTimeEpoch
	if difference < 60 {
		// this is seconds
		return strconv.Itoa(int(difference)) + "  seconds"
	} else if difference < 3600 {
		minutes := difference / 60
		return strconv.Itoa(int(minutes)) + "  minutes"
	} else {
		minutes := difference / 3600
		return strconv.Itoa(int(minutes)) + "  hour"
	}
}
func (v *ActionService) getBasicInfo(ctx *gin.Context) (int, *gorm.DB, common.UserBasicInfo, *database.ShiftMasterInfo) {
	recordId := util.GetRecordId(ctx)
	userId := common.GetUserId(ctx)
	v.Logger.Info("getting labour management shift master for record", zap.Any("record_id", recordId))
	_, serviceRequestGeneralObject := database.Get(v.Database, const_util.LabourManagementShiftMasterTable, recordId)
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	basicUserInfo := authService.GetUserInfoById(userId)

	return recordId, v.Database, basicUserInfo, database.GetShiftMasterInfo(serviceRequestGeneralObject.ObjectInfo)
}

func (v *ActionService) getShiftMasterInfo(ctx *gin.Context, recordId int) (int, *gorm.DB, common.UserBasicInfo, *database.ShiftMasterInfo) {
	userId := common.GetUserId(ctx)
	v.Logger.Info("getting labour management shift master for record", zap.Any("record_id", recordId))
	_, serviceRequestGeneralObject := database.Get(v.Database, const_util.LabourManagementShiftMasterTable, recordId)
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	basicUserInfo := authService.GetUserInfoById(userId)

	return recordId, v.Database, basicUserInfo, database.GetShiftMasterInfo(serviceRequestGeneralObject.ObjectInfo)
}
