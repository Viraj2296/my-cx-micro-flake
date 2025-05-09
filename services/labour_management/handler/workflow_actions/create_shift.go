package workflow_actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/labour_management/handler/const_util"
	"cx-micro-flake/services/labour_management/handler/database"
	"cx-micro-flake/services/labour_management/handler/helper"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type createShiftRequest struct {
	ShiftTemplateId      int   `json:"shiftTemplateId"`
	ScheduledOrderEvents []int `json:"scheduledOrderEvents"` // when creating the shift, users can select the list of lines this shift assigned it.
}

// src : [1,2,3]
// dst : [4,5,6]  -> allow to create shift
// dst : [1,2,3] -> won't allow you to create
// dst : [4,5,6,7,8,9,1]  -> not allow

func hasSameSchedulerEvents(srcEvents, dstEvents []int) bool {
	srcSet := make(map[int]struct{})
	for _, event := range srcEvents {
		srcSet[event] = struct{}{}
	}
	for _, event := range dstEvents {
		if _, exists := srcSet[event]; exists {
			return true
		}
	}
	return false
}

func (v *ActionService) CreateShift(ctx *gin.Context) {
	// get the department and site to create the shift ID
	v.Logger.Info("handling create shift request")
	createShiftRequest := createShiftRequest{}
	if err := ctx.ShouldBindBodyWith(&createShiftRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var userId = common.GetUserId(ctx)
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	var basicUserInfo = authService.GetUserInfoById(userId)

	v.Logger.Info("shift is getting created by", zap.Int("user_id", basicUserInfo.UserId))
	// first if we have thing configured for checking roles.
	if len(v.LabourManagementSettingInfo.ShiftCreationRoles) == 0 {
		v.Logger.Warn("no shift creation roles configured.. using default shift supervisor")
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Operation",
				Description: "The system is not configured with shift management roles. Please ask your administrator to configure the roles",
			})
		return
	} else {
		if !util.HasInt(basicUserInfo.JobRole, v.LabourManagementSettingInfo.ShiftCreationRoles) {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Invalid Operation",
					Description: "Only configured job roles are allowed to create shift, and manage shifts, please check your access privileges with system admin",
				})
			return
		}
	}

	// load all the shift master, and check the similar shift is created with same events, then don't create them
	err, listOfShiftMasterInterface := database.GetObjects(v.Database, const_util.LabourManagementShiftMasterTable)
	if err == nil {
		for _, shiftMasterInterface := range listOfShiftMasterInterface {
			shiftMasterInfo := database.GetShiftMasterInfo(shiftMasterInterface.ObjectInfo)
			if shiftMasterInfo.ObjectStatus != component.ObjectStatusArchived {
				if hasSameSchedulerEvents(shiftMasterInfo.ScheduledOrderEvents, createShiftRequest.ScheduledOrderEvents) {
					response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
						&response.DetailedError{
							Header:      "Shift is already created",
							Description: "You have already created the shift using following scheduled orders, please go to active shift and use the shift",
						})
					return
				}
			}

		}
	}

	err, shiftTemplateInterface := database.Get(v.Database, const_util.LabourManagementShiftTemplateTable, createShiftRequest.ShiftTemplateId)
	if err != nil {
		v.Logger.Error("error getting shift template", zap.String("error", err.Error()))
		response.SendResourceNotFound(ctx)
	}
	var shiftTemplateInfo = database.GetSShiftTemplateInfo(shiftTemplateInterface.ObjectInfo)
	err, shiftStartModifiedTime := v.GetShiftTime(shiftTemplateInfo.ShiftStartTime, 0)
	if err != nil {
		response.SendInternalSystemError(ctx)
	}
	err, shiftEndTime := v.GetShiftTime(shiftTemplateInfo.ShiftStartTime, shiftTemplateInfo.ShiftPeriod)
	if err != nil {
		response.SendInternalSystemError(ctx)
	}
	if len(basicUserInfo.Department) > 1 {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Department Setting",
				Description: "You are seeing this error because you are either not assigned to any department or assigned to multiple departments. Please report this issue to the system administrator",
			})
		return
	}

	shiftEndDate, shiftEndTime := getDateAndTime(shiftEndTime)
	shiftStartDate, shiftStartTime := getDateAndTime(shiftStartModifiedTime)

	var shiftMasterInfo = database.ShiftMasterInfo{
		SiteId:                shiftTemplateInfo.SiteId,
		CreatedAt:             util.GetCurrentTime("2006-01-02T15:04:05.000Z"),
		CreatedBy:             userId,
		CanCheckIn:            false,
		ShiftStatus:           const_util.ShiftStatusPending,
		CanShiftStop:          false,
		DepartmentId:          basicUserInfo.Department[0],
		ObjectStatus:          component.ObjectStatusActive,
		ShiftEndDate:          shiftEndDate,
		ShiftEndTime:          shiftEndTime,
		CanShiftStart:         true,
		LastUpdatedAt:         util.GetCurrentTime("2006-01-02T15:04:05.000Z"),
		LastUpdatedBy:         userId,
		ActualManPower:        0,
		ShiftStartDate:        shiftStartDate,
		ShiftStartTime:        shiftStartTime,
		ShiftSupervisor:       userId,
		ScheduledOrderEvents:  createShiftRequest.ScheduledOrderEvents,
		IsSupervisorCheckedIn: false,
		ShiftTemplateId:       createShiftRequest.ShiftTemplateId,
		ActionRemarks:         make([]database.ActionRemarks, 0),
		CanRollBack:           true,
	}
	object := component.GeneralObject{
		ObjectInfo: shiftMasterInfo.Serialised(),
	}

	err, createdResourceId := database.CreateFromGeneralObject(v.Database, const_util.LabourManagementShiftMasterTable, object)
	if err != nil {
		response.SendDetailedError(ctx, http.StatusBadRequest, const_util.GetError("Error in Resource Creation "), const_util.ErrorCreatingObjectInformation, err.Error())
		return
	}
	factoryInterface := common.GetService("factory_module").ServiceInterface.(common.FactoryServiceInterface)
	var siteName = factoryInterface.GetSiteName(shiftTemplateInfo.SiteId)
	var departmentName = factoryInterface.GetDepartmentName(basicUserInfo.Department[0])
	updatingData := make(map[string]interface{})
	var composedShiftId = siteName + departmentName + helper.GenerateShiftID(createdResourceId)
	shiftMasterInfo.ShiftReferenceId = composedShiftId
	updatingData["object_info"] = shiftMasterInfo.Serialised()
	err = database.Update(v.Database, const_util.LabourManagementShiftMasterTable, createdResourceId, updatingData)
	v.Logger.Info("shift is created successfully", zap.Any("record_id", createdResourceId))

	// all success, now create the shift master production one by one
	v.createShiftMasterProduction(createdResourceId, userId, createShiftRequest.ScheduledOrderEvents)
	ctx.JSON(http.StatusCreated, response.GeneralResponse{
		Code:    0,
		Message: "Great, you have successfully created the shift. Now you can start the shift and begin checking in",
	})
}

func (v *ActionService) createShiftMasterProduction(shiftResourceId int, userId int, ScheduledOrderEvents []int) {
	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	for _, scheduledOrderEventId := range ScheduledOrderEvents {
		// each events get the details
		err, scheduledOrderInterface := productionOrderInterface.GetAssemblyScheduledOrderInfo(const_util.ProjectID, scheduledOrderEventId)
		if err != nil {
			v.Logger.Error("error getting scheduler order information", zap.String("error", err.Error()))
			continue
		}
		scheduledOrderInfo := GetAssemblyScheduledOrderEventInfo(scheduledOrderInterface.ObjectInfo)
		shiftMasterProductionOrderInfo := database.ShiftMasterProductionInfo{
			ShiftId:                    shiftResourceId,
			ScheduledEventId:           scheduledOrderEventId,
			ActualManpower:             0,
			ActualManHourPartTimer:     0,
			ShiftTargetOutputPartTimer: 0,
			ShiftActualOutputPartTimer: 0,
			CreatedAt:                  util.GetCurrentTime("2006-01-02T15:04:05.000Z"),
			LastUpdatedAt:              util.GetCurrentTime("2006-01-02T15:04:05.000Z"),
			CreatedBy:                  userId,
			LastUpdatedBy:              userId,
			MachineId:                  scheduledOrderInfo.MachineId,
			CanUpdateManpower:          false,
		}
		generalObject := component.GeneralObject{ObjectInfo: shiftMasterProductionOrderInfo.Serialised()}
		err, resourceId := database.CreateFromGeneralObject(v.Database, const_util.LabourManagementShiftProductionTable, generalObject)
		if err == nil {
			v.Logger.Info("new shift production order is created", zap.Any("record_id", resourceId))
		} else {
			v.Logger.Error("error creating shift production record", zap.Any("error", err.Error()))
		}

	}

}
func getDateAndTime(datetime string) (string, string) {
	parsedTime, err := time.Parse(time.RFC3339, datetime)
	if err != nil {
		fmt.Println("Error parsing date time:", err)
		return "", ""
	}

	// Add 8 hours to the parsed time
	updatedTime := parsedTime.Add(8 * time.Hour)

	// Extract the date and time components from the updated time
	extractedDate := updatedTime.Format("2006-01-02")
	extractedTime := updatedTime.Format("15:04:05")
	return extractedDate, extractedTime
}
