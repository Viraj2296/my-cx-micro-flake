package workflow_actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/scheduler/handler/const_util"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
)

type EventResourceRequest struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

func (v *ActionService) GetEventsByResources(ctx *gin.Context) {
	projectId := ctx.Param("projectId")
	componentName := ctx.Param("componentName")
	v.Logger.Info("get event by resource called")

	eventResourceRequest := EventResourceRequest{}
	if err := ctx.ShouldBindBodyWith(&eventResourceRequest, binding.JSON); err != nil {
		err = ctx.AbortWithError(http.StatusBadRequest, err)
		if err != nil {
			v.Logger.Error("error sending response", zap.Error(err))
		}
		return
	}
	objectRequestStartDate, errStart := time.Parse("2006-01-02", eventResourceRequest.StartDate)
	objectRequestEndDate, errEnd := time.Parse("2006-01-02", eventResourceRequest.EndDate)
	if errStart != nil || errEnd != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format in request"})
		return
	}
	productionOrderService := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
	maintenanceService := common.GetService("maintenance_module").ServiceInterface.(common.MaintenanceInterface)
	mouldService := common.GetService("moulds_module").ServiceInterface.(common.MouldInterface)
	manufacturingService := common.GetService("manufacturing_module").ServiceInterface.(common.ManufacturingInterface)

	var listfEvents []interface{}

	if componentName == const_util.MouldingSchedulerComponent {
		v.Logger.Info("Processing Moulding Scheduler Component")

		var listOfScheduledOrderEvents *[]component.GeneralObject
		_, listOfScheduledOrderEvents = productionOrderService.GetAllSchedulerEventForScheduler(projectId)

		_, listOfWorkOrderEvents := maintenanceService.GetMaintenanceEventsForScheduler(projectId)
		//_, listOfWorkOrderTaskEvents := maintenanceService.GetMaintenanceOrderTaskEventsForScheduler(projectId)
		_, listOfMouldTestEvents := mouldService.GetMouldTestEventsForScheduler(projectId)

		//*listOfEvents = append(*listOfEvents, *listOfWorkOrderEvents...)
		*listOfScheduledOrderEvents = append(*listOfScheduledOrderEvents, *listOfWorkOrderEvents...)
		*listOfScheduledOrderEvents = append(*listOfScheduledOrderEvents, *listOfMouldTestEvents...)
		_, listOfResources := machineService.GetListOfMachines(projectId)

		assigmentId := 1

		for _, resource := range listOfResources {
			var arrayOfEvents []interface{}
			var resourceObjects = make(map[string]interface{})
			json.Unmarshal(resource.ObjectInfo, &resourceObjects)
			schedulerResponse := component.SchedulerEventResponse{}
			schedulerResponse.MachineId = util.InterfaceToString(resourceObjects["newMachineId"])

			for _, event := range *listOfScheduledOrderEvents {
				var objects = make(map[string]interface{})
				json.Unmarshal(event.ObjectInfo, &objects)

				if resource.Id == util.InterfaceToInt(objects["machineId"]) {
					if objects["objectStatus"].(string) != "Archived" {
						objects["id"] = assigmentId
						objects["eventId"] = event.Id
						startDate := objects["startDate"].(string)
						endDate := objects["endDate"].(string)
						objectStartDate := util.ConvertTimeToTimeZonCorrected("Asia/Singapore", startDate)
						objectEndDate := util.ConvertTimeToTimeZonCorrected("Asia/Singapore", endDate)
						objects["startDate"] = objectStartDate
						objects["endDate"] = objectEndDate

						objectStartDateTime, _ := time.Parse("2006-01-02T15:04:05", objectStartDate)
						objectEndDateTime, _ := time.Parse("2006-01-02T15:04:05", objectEndDate)

						if objectStartDateTime.Before(objectRequestEndDate.Add(24*time.Hour)) &&
							objectEndDateTime.After(objectRequestStartDate) {

							if objects["eventType"] == "production_schedule" {
								machineId := util.InterfaceToInt(objects["machineId"])
								productionOrderId := util.InterfaceToInt(objects["eventSourceId"])

								_, productionGeneralObject := productionOrderService.GetMachineProductionOrderInfo(projectId, productionOrderId, machineId)

								var productionObjectInfo = make(map[string]interface{})
								json.Unmarshal(productionGeneralObject.ObjectInfo, &productionObjectInfo)
								partId := util.InterfaceToInt(productionObjectInfo["partNumber"])

								_, partMasterObject := mouldService.GetPartInfo(projectId, partId)
								var partMasterInfo = make(map[string]interface{})
								json.Unmarshal(partMasterObject.ObjectInfo, &partMasterInfo)

								partName := util.InterfaceToString(partMasterInfo["partNumber"])
								objects["partNo"] = partName
							}
							objects["permissions"] = common.GetPermissions(util.InterfaceToInt(objects["eventStatus"]), util.InterfaceToBool(objects["canComplete"]), util.InterfaceToBool(objects["canForceStop"]), util.InterfaceToBool(objects["canHold"]), util.InterfaceToBool(objects["canRelease"]))

							keysToRemove := []string{"canComplete", "canForceStop", "canHold", "canRelease"}
							objects = common.RemoveKeys(objects, keysToRemove)

							arrayOfEvents = append(arrayOfEvents, objects)
							// machineId := util.InterfaceToInt(objects["machineId"])
							assigmentId += 1

						}
					}

				} else {
					continue
				}

			}
			if len(arrayOfEvents) > 0 {
				schedulerResponse.Events = arrayOfEvents
				listfEvents = append(listfEvents, schedulerResponse)
			}
		}
		v.Logger.Info("Completed processing Moulding Scheduler Component", zap.Int("totalEvents", len(listfEvents)))

	} else if componentName == const_util.AssemblySchedulerComponent {

		_, listOfEvents := productionOrderService.GetAllAssemblyEventForScheduler(projectId)

		_, listOfWorkOrderEvents := maintenanceService.GetMaintenanceEventsForAssemblyScheduler(projectId)
		_, listOfResources := machineService.GetListOfAssemblyMachines(projectId)

		*listOfEvents = append(*listOfEvents, *listOfWorkOrderEvents...)

		assigmentId := 1

		for _, resource := range listOfResources {
			var arrayOfEvents []interface{}
			var resourceObjects = make(map[string]interface{})
			json.Unmarshal(resource.ObjectInfo, &resourceObjects)
			schedulerResponse := component.SchedulerEventResponse{}
			schedulerResponse.MachineId = util.InterfaceToString(resourceObjects["newMachineId"])
			for _, event := range *listOfEvents {
				var objects = make(map[string]interface{})
				json.Unmarshal(event.ObjectInfo, &objects)

				if resource.Id == util.InterfaceToInt(objects["machineId"]) {

					if objects["objectStatus"].(string) != "Archived" {
						objects["id"] = assigmentId
						objects["eventId"] = event.Id
						startDate := objects["startDate"].(string)
						endDate := objects["endDate"].(string)
						objectStartDate := util.ConvertTimeToTimeZonCorrected("Asia/Singapore", startDate)
						objectEndDate := util.ConvertTimeToTimeZonCorrected("Asia/Singapore", endDate)
						objects["startDate"] = objectStartDate
						objects["endDate"] = objectEndDate

						objectStartDateTime, _ := time.Parse("2006-01-02T15:04:05", objectStartDate)
						objectEndDateTime, _ := time.Parse("2006-01-02T15:04:05", objectEndDate)

						if objectStartDateTime.Before(objectRequestEndDate.Add(24*time.Hour)) &&
							objectEndDateTime.After(objectRequestStartDate) {

							machineId := util.InterfaceToInt(objects["machineId"])
							productionOrderId := util.InterfaceToInt(objects["eventSourceId"])

							_, productionGeneralObject := productionOrderService.GetAssemblyProductionOrderInfo(projectId, productionOrderId, machineId)

							var productionObjectInfo = make(map[string]interface{})
							json.Unmarshal(productionGeneralObject.ObjectInfo, &productionObjectInfo)
							partId := util.InterfaceToInt(productionObjectInfo["partNumber"])

							_, partMasterObject := manufacturingService.GetAssemblyPartInfo(projectId, partId)
							var partMasterInfo = make(map[string]interface{})
							json.Unmarshal(partMasterObject.ObjectInfo, &partMasterInfo)

							partName := util.InterfaceToString(partMasterInfo["partNumber"])
							objects["partNo"] = partName

							objects["name"] = partName + "_" + util.InterfaceToString(objects["name"])
							objects["permissions"] = common.GetPermissionsForAssembly(util.InterfaceToInt(objects["eventStatus"]), util.InterfaceToBool(objects["canComplete"]), util.InterfaceToBool(objects["canForceStop"]), util.InterfaceToBool(objects["canHold"]), util.InterfaceToBool(objects["canRelease"]))

							keysToRemove := []string{"canComplete", "canForceStop", "canHold", "canRelease"}
							objects = common.RemoveKeys(objects, keysToRemove)

							fmt.Println("Passed")
							arrayOfEvents = append(arrayOfEvents, objects)
							assigmentId += 1

						}

					}
				} else {
					continue
				}

			}
			if len(arrayOfEvents) > 0 {
				schedulerResponse.Events = arrayOfEvents
				listfEvents = append(listfEvents, schedulerResponse)
			}
		}
		v.Logger.Info("Completed processing Assembly Scheduler Component", zap.Int("totalEvents", len(listfEvents)))

	} else if componentName == const_util.ToolingSchedulerComponent {
		_, listOfEvents := productionOrderService.GetAllToolingEventForScheduler(projectId)
		_, listOfWorkOrderEvents := maintenanceService.GetMaintenanceEventsForToolingScheduler(projectId)
		_, listOfResources := machineService.GetListOfToolingMachines(projectId)

		*listOfEvents = append(*listOfEvents, *listOfWorkOrderEvents...)

		assigmentId := 1

		for _, resource := range listOfResources {
			var arrayOfEvents []interface{}
			var resourceObjects = make(map[string]interface{})
			json.Unmarshal(resource.ObjectInfo, &resourceObjects)
			schedulerResponse := component.SchedulerEventResponse{}
			schedulerResponse.MachineId = util.InterfaceToString(resourceObjects["newMachineId"])

			for _, event := range *listOfEvents {
				var objects = make(map[string]interface{})
				json.Unmarshal(event.ObjectInfo, &objects)

				if resource.Id == util.InterfaceToInt(objects["machineId"]) {
					if objects["objectStatus"].(string) != "Archived" {
						objects["id"] = assigmentId
						objects["eventId"] = event.Id
						startDate := objects["startDate"].(string)
						endDate := objects["endDate"].(string)
						objectStartDate := util.ConvertTimeToTimeZonCorrected("Asia/Singapore", startDate)
						objectEndDate := util.ConvertTimeToTimeZonCorrected("Asia/Singapore", endDate)
						objects["startDate"] = objectStartDate
						objects["endDate"] = objectEndDate

						objectStartDateTime, _ := time.Parse("2006-01-02T15:04:05", objectStartDate)
						objectEndDateTime, _ := time.Parse("2006-01-02T15:04:05", objectEndDate)

						if objectStartDateTime.Before(objectRequestEndDate.Add(24*time.Hour)) &&
							objectEndDateTime.After(objectRequestStartDate) {
							fmt.Println("Passed")

							partId := util.InterfaceToInt(objects["partId"])
							eventSourceId := util.InterfaceToInt(objects["eventSourceId"])

							_, partGeneralObject := productionOrderService.GetToolingPartById(projectId, partId)

							var partObjectInfo = make(map[string]interface{})
							json.Unmarshal(partGeneralObject.ObjectInfo, &partObjectInfo)
							partName := util.InterfaceToString(partObjectInfo["name"])
							objects["partNo"] = partName

							_, bomGeneralObject := productionOrderService.GetBomInfo(projectId, eventSourceId)
							var bomObjectInfo = make(map[string]interface{})
							json.Unmarshal(bomGeneralObject.ObjectInfo, &bomObjectInfo)
							bomName := util.InterfaceToString(bomObjectInfo["name"])
							objects["bomName"] = bomName

							objects["permissions"] = common.GetPermissions(util.InterfaceToInt(objects["eventStatus"]), util.InterfaceToBool(objects["canComplete"]), util.InterfaceToBool(objects["canForceStop"]), util.InterfaceToBool(objects["canHold"]), util.InterfaceToBool(objects["canRelease"]))

							keysToRemove := []string{"canComplete", "canForceStop", "canHold", "canRelease"}
							objects = common.RemoveKeys(objects, keysToRemove)

							arrayOfEvents = append(arrayOfEvents, objects)
							assigmentId += 1

						}

					}
				} else {
					continue
				}

			}

			if len(arrayOfEvents) > 0 {
				schedulerResponse.Events = arrayOfEvents
				listfEvents = append(listfEvents, schedulerResponse)
			}
		}
		v.Logger.Info("Completed processing Tooling Scheduler Component", zap.Int("totalEvents", len(listfEvents)))

	}
	v.Logger.Info("Sending response", zap.Int("totalEvents", len(listfEvents)))
	ctx.JSON(http.StatusOK, listfEvents)

}
