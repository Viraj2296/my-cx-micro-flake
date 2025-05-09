package workflow_actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/labour_management/handler/const_util"
	"cx-micro-flake/services/labour_management/handler/database"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
)

type LineRequest struct {
	LineId int `json:"lineId"`
}
type DataSeries struct {
	Name string `json:"name"`
	Data []int  `json:"data"`
}
type LineChartResponse struct {
	Chart struct {
		Type string `json:"type"`
	} `json:"chart"`
	Title struct {
		Text string `json:"text"`
	} `json:"title"`
	Subtitle struct {
		Text string `json:"text"`
	} `json:"subtitle"`
	XAxis struct {
		Categories []string `json:"categories"`
	} `json:"xAxis"`
	YAxis struct {
		Title struct {
			Text string `json:"text"`
		} `json:"title"`
	} `json:"yAxis"`
	PlotOptions struct {
		Line struct {
			DataLabels struct {
				Enabled bool `json:"enabled"`
			} `json:"dataLabels"`
			EnableMouseTracking bool `json:"enableMouseTracking"`
		} `json:"line"`
	} `json:"plotOptions"`
	Series []DataSeries `json:"series"`
}
type LineDisplayResponse struct {
	ShiftId                  string            `json:"shiftId"`
	Model                    string            `json:"model"`
	CanDisplay               bool              `json:"canDisplay"`
	Material                 string            `json:"material"`
	ShiftTarget              int               `json:"shiftTarget"`
	ShiftActual              int               `json:"shiftActual"`
	Different                int               `json:"different"`
	PartImage                string            `json:"partImage"`
	ShiftStartDateTime       string            `json:"shiftStartDateTime"`
	ShiftEndDateTime         string            `json:"shiftEndDateTime"`
	ShiftActualStartDateTime string            `json:"shiftActualStartDateTime"`
	Operators                []LineUsers       `json:"operators"`
	ManagementTeam           []ManagementTeam  `json:"managementTeam"`
	ShiftPerformanceChart    LineChartResponse `json:"shiftPerformanceChart"`
}

type LineUsers struct {
	ProfileImage string `json:"profileImage"`
	Name         string `json:"name"`
	RoleName     string `json:"roleName"`
}

type Labour struct {
	Role      string `json:"role"`
	Name      string `json:"name"`
	AvatarUrl string `json:"avatarUrl"`
}

type ManagementTeam struct {
	RoleName       string   `json:"roleName"`
	Labours        []Labour `json:"labours"`
	HierarchyLevel int      `json:"hierarchyLevel"`
}
type assemblyManualOrderHistoryByEventRequest struct {
	EventId int `json:"eventId"`
}
type assemblyManualOrderHistoryByEventResponse struct {
	CreatedAt         string `json:"createdAt"`
	CompletedQuantity int    `json:"completedQuantity"`
}

func (v *ActionService) GetLineDisplay(ctx *gin.Context) {
	v.Logger.Info("getting the line display...")

	lineRequest := LineRequest{}
	if err := ctx.ShouldBindBodyWith(&lineRequest, binding.JSON); err != nil {
		err := ctx.AbortWithError(http.StatusBadRequest, err)
		if err != nil {
			v.Logger.Error("error sending response to client", zap.String("error", err.Error()))
		}
		return
	}

	var conditionString = " object_info->>'$.shiftStatus' = " + strconv.Itoa(const_util.ShiftStatusActive)
	err, listOfActiveShifts := database.GetConditionalObjects(v.Database, const_util.LabourManagementShiftMasterTable, conditionString)
	if err != nil {
		// it is the error send the error description
		v.Logger.Error("error getting shift master", zap.String("error", err.Error()))
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Internal Server Error",
				Description: "Sorry, server couldn't able to complete this action due to internal error, please report this error to system admin",
			})
		return
	}

	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
	manufacturingModuleInterface := common.GetService("manufacturing_module").ServiceInterface.(common.ManufacturingInterface)

	// Identify what shift is running on requested line
	var lineResponse = LineDisplayResponse{}
	lineResponse.ShiftTarget = 0
	lineResponse.Operators = make([]LineUsers, 0)
	lineResponse.ManagementTeam = make([]ManagementTeam, 0) // Changed to slice of ManagementTeam
	lineResponse.Different = 0
	lineResponse.Material = "-"
	lineResponse.ShiftId = "-"
	lineResponse.PartImage = "-"
	lineResponse.ShiftStartDateTime = "-"
	lineResponse.ShiftActualStartDateTime = "-"
	lineResponse.ShiftEndDateTime = "-"
	lineResponse.ShiftActual = 0
	lineResponse.Model = "-"
	lineResponse.CanDisplay = false

	lineChartResponse := LineChartResponse{}
	lineChartResponse.Chart.Type = "line"
	lineChartResponse.Title.Text = "Shift Performance"
	lineChartResponse.Subtitle.Text = "Shift: FYCASSY (KX)SHIFT0180"
	lineChartResponse.XAxis.Categories = make([]string, 0)
	lineChartResponse.YAxis.Title.Text = "Shift Output (Quantity)"
	lineChartResponse.PlotOptions.Line.DataLabels.Enabled = true
	lineChartResponse.PlotOptions.Line.EnableMouseTracking = true
	lineChartResponse.Series = make([]DataSeries, 0)

	var lineIdFound = false
	var foundShiftId = 0

	var assemblyEventId = 0
	var isMachineConnectedDriver = false
	var assemblyManualHistoryListResponse []assemblyManualOrderHistoryByEventResponse
	var shiftActual = 0
	for _, activeShiftInterface := range listOfActiveShifts {
		shiftMasterInfo := database.GetShiftMasterInfo(activeShiftInterface.ObjectInfo)
		var shiftTarget = 0

		for _, schedulerEventId := range shiftMasterInfo.ScheduledOrderEvents {
			v.Logger.Info("processing shift event id", zap.Int("scheduler_event_id", schedulerEventId))
			err, assemblySchedulerOrderInterface := productionOrderInterface.GetAssemblyScheduledOrderInfo(const_util.ProjectID, schedulerEventId)
			if err != nil {
				v.Logger.Error("error getting scheduler order event", zap.String("error", err.Error()))
				continue
			}
			assemblySchedulerInfo := GetAssemblyScheduledOrderEventInfo(assemblySchedulerOrderInterface.ObjectInfo)
			err, assemblyMachineObject := machineService.GetAssemblyMachineInfoById(assemblySchedulerInfo.MachineId)
			if err != nil {
				v.Logger.Error("error getting machine details", zap.String("error", err.Error()))
				continue
			}
			assemblyMachineMasterInfo := GetAssemblyMachineMasterInfo(assemblyMachineObject.ObjectInfo)

			err, productionOrderData := productionOrderInterface.GetAssemblyProductionOrderById(assemblySchedulerInfo.EventSourceId)
			if err != nil {
				v.Logger.Error("error getting production order information", zap.String("error", err.Error()))
				continue
			}
			var productionOrderInfo = GetAssemblyProductionOrderInfo(productionOrderData.ObjectInfo)

			_, partObject := manufacturingModuleInterface.GetAssemblyPartInfo(const_util.ProjectID, productionOrderInfo.PartNumber)
			partInfo := GetPartInfo(partObject.ObjectInfo)
			v.Logger.Info("part details", zap.Any("part_details", partInfo), zap.Any("assembly_master_info", assemblyMachineMasterInfo))

			if assemblyMachineMasterInfo.AssemblyLineOption == lineRequest.LineId {
				// Yes, found the line
				// Get the shift end time
				var shiftEndTime = "-"
				err, shiftTemplateInterface := database.Get(v.Database, const_util.LabourManagementShiftTemplateTable, shiftMasterInfo.ShiftTemplateId)
				if err == nil {
					shiftTemplateInfo := database.GetSShiftTemplateInfo(shiftTemplateInterface.ObjectInfo)
					err, endTime := v.GetShiftTime(shiftTemplateInfo.ShiftStartTime, shiftTemplateInfo.ShiftPeriod)
					if err != nil {
						v.Logger.Error("error getting shift time", zap.String("error", err.Error()))
					} else {
						v.Logger.Info("shift end time", zap.String("end_time", endTime))
						shiftEndTime = addHoursToRFC3339(endTime, 8)
					}
				}

				startDataTime, _ := CombineDateTimeToISO8601(shiftMasterInfo.ShiftStartDate, shiftMasterInfo.ShiftStartTime)
				actualStartDateTime, _ := ConvertToISO8601NoMilliseconds(shiftMasterInfo.CreatedAt)
				shiftTarget = assemblySchedulerInfo.ScheduledQty
				shiftActual = assemblySchedulerInfo.CompletedQty

				lineResponse.Model = assemblyMachineMasterInfo.Model
				lineResponse.PartImage = partInfo.Image
				lineResponse.Material = partInfo.PartNumber
				lineResponse.ShiftId = shiftMasterInfo.ShiftReferenceId
				lineResponse.Different = shiftActual - shiftTarget
				lineResponse.ShiftActual = shiftActual
				lineResponse.ShiftStartDateTime = startDataTime
				lineResponse.ShiftActualStartDateTime = actualStartDateTime
				lineResponse.ShiftEndDateTime = shiftEndTime
				lineResponse.ShiftTarget = shiftTarget
				lineResponse.CanDisplay = true

				assemblyEventId = schedulerEventId
				isMachineConnectedDriver = assemblyMachineMasterInfo.IsMESDriverConfigured

				if !assemblyMachineMasterInfo.IsMESDriverConfigured {
					var historyListRequest assemblyManualOrderHistoryByEventRequest
					historyListRequest.EventId = schedulerEventId
					serialisedData, _ := json.Marshal(historyListRequest)

					historyListResponse, err := productionOrderInterface.GetAssemblyManualOrderHistoryFromEventId(serialisedData)
					if err != nil {
						v.Logger.Error("error getting manual order history", zap.String("error", err.Error()))
					} else {
						err := json.Unmarshal(historyListResponse.([]byte), &assemblyManualHistoryListResponse)
						if err != nil {
							v.Logger.Error("error unmarshalling manual order history", zap.String("error", err.Error()))
						}
					}
				}

				// Get the assembly line name
				err, assemblyLineObject := machineService.GetAssemblyLineFromId(const_util.ProjectID, assemblyMachineMasterInfo.AssemblyLineOption)
				if err == nil {
					lineResponse.Model = database.GetAssemblyMachineLineInfo(assemblyLineObject.ObjectInfo).Name
				}

				lineIdFound = true
				foundShiftId = activeShiftInterface.Id
				break
			} else {
				err, assemblyLineObject := machineService.GetAssemblyLineFromId(const_util.ProjectID, lineRequest.LineId)
				if err == nil {
					lineResponse.Model = database.GetAssemblyMachineLineInfo(assemblyLineObject.ObjectInfo).Name
				}
				v.Logger.Info("no production orders running for the", zap.String("model_name ", lineResponse.Model))
			}

			if lineIdFound {
				break
			}
		}
	}

	if lineIdFound {
		// Attendance should be taken for lines, not for shifts
		var attendanceCondition = " object_info->>'$.shiftResourceId' = " + strconv.Itoa(foundShiftId)
		v.Logger.Error("line found, getting the attendance", zap.Any("attendanceCondition", attendanceCondition))
		err, listOfAttendance := database.GetConditionalObjects(v.Database, const_util.LabourManagementAttendanceTable, attendanceCondition)
		if err != nil {
			// Add the empty array
			v.Logger.Error("error getting operators", zap.String("error", err.Error()))
			lineResponse.Operators = make([]LineUsers, 0)
		}

		var uniqueUserList []int
		for _, info := range listOfAttendance {
			attendanceInfo := database.GetAttendanceInfo(info.ObjectInfo)
			// Only show the line checked in by the operator to the requested lines
			// Remove the users already checked out
			if attendanceInfo.CheckOutDate == "" || attendanceInfo.CheckOutTime == "" {
				if util.HasInt(lineRequest.LineId, attendanceInfo.ManufacturingLines) {
					uniqueUserList = append(uniqueUserList, attendanceInfo.UserResourceId)
				}
			}
		}
		uniqueUserList = util.RemoveDuplicateInt(uniqueUserList)

		// Create a map to group labours by role name
		roleLabourMap := make(map[string]*ManagementTeam)

		// Assuming role is a string and we are getting hierarchy level
		for _, userId := range uniqueUserList {
			userInfo := authService.GetUserInfoById(userId)
			jobRoleName := authService.GetJobRoleName(userInfo.JobRole)

			// Get the hierarchy level
			hierarchyLevel := authService.GetJobRoleHierarchy(userInfo.JobRole)

			labour := Labour{
				Role:      jobRoleName,
				Name:      userInfo.FullName,
				AvatarUrl: userInfo.AvatarUrl,
			}

			if util.HasInt(userInfo.JobRole, v.LabourManagementSettingInfo.LineDisplayRoles) {
				// If the role already exists in the map, append the labour
				if team, exists := roleLabourMap[jobRoleName]; exists {
					team.Labours = append(team.Labours, labour) // Correctly appending to the Labours slice
				} else {
					// Otherwise, create a new entry
					roleLabourMap[jobRoleName] = &ManagementTeam{
						RoleName:       jobRoleName,
						HierarchyLevel: hierarchyLevel,
						Labours:        []Labour{labour}, // Initialize with the first labour
					}
				}
			} else {
				lineResponse.Operators = append(lineResponse.Operators, LineUsers{
					ProfileImage: userInfo.AvatarUrl,
					Name:         userInfo.FullName,
					RoleName:     jobRoleName,
				})
			}
		}

		// Convert the map back to a slice
		for _, team := range roleLabourMap {
			lineResponse.ManagementTeam = append(lineResponse.ManagementTeam, *team) // Correctly dereferencing
		}

		// Sorting ManagementTeam by HierarchyLevel
		sort.Slice(lineResponse.ManagementTeam, func(i, j int) bool {
			return lineResponse.ManagementTeam[i].HierarchyLevel < lineResponse.ManagementTeam[j].HierarchyLevel
		})

	} else {
		v.Logger.Warn("No lines found send the display, so sending the empty response with canDisplay false")
	}
	arrayOfTimeInterval, err := generateHourlyIntervals(lineResponse.ShiftStartDateTime, lineResponse.ShiftEndDateTime)
	//TODO , remove the first element if the array size is more than two. For example if the shift is start time 06:30, we need to remove, and the
	// element should 07:30
	if len(arrayOfTimeInterval) > 0 {
		arrayOfTimeInterval = arrayOfTimeInterval[1:]
	}
	var listOfData []int
	if err == nil {
		lineChartResponse.XAxis.Categories = arrayOfTimeInterval
		var hourShiftTarget = lineResponse.ShiftTarget / len(arrayOfTimeInterval)
		var combinedShiftTarget = hourShiftTarget
		for index, _ := range arrayOfTimeInterval {
			if index == len(arrayOfTimeInterval)-1 {
				listOfData = append(listOfData, lineResponse.ShiftTarget)
			} else {
				listOfData = append(listOfData, combinedShiftTarget)
				combinedShiftTarget += hourShiftTarget
			}
		}
	}
	arrayOfStringTimeInterval, err := generateHourlyIntervalsWithTimestamps(lineResponse.ShiftStartDateTime, lineResponse.ShiftEndDateTime)
	if len(arrayOfStringTimeInterval) > 0 {
		arrayOfStringTimeInterval = arrayOfStringTimeInterval[1:]
	}
	v.Logger.Info("shift dates for the time interval", zap.Any("ShiftStartDateTime", lineResponse.ShiftStartDateTime), zap.Any("ShiftStartDateTime", lineResponse.ShiftEndDateTime), zap.Int("event_id", assemblyEventId))
	var shiftActualSeries = make([]int, 0)
	if err == nil {
		shiftActualSeries = v.getShiftActualHours(assemblyEventId, isMachineConnectedDriver, assemblyManualHistoryListResponse, arrayOfStringTimeInterval)
	}
	lineChartResponse.Series = append(lineChartResponse.Series, DataSeries{
		Name: "Shift Target",
		Data: listOfData,
	})
	lineChartResponse.Series = append(lineChartResponse.Series, DataSeries{
		Name: "Shift Actual",
		Data: shiftActualSeries,
	})
	lineResponse.ShiftPerformanceChart = lineChartResponse
	ctx.JSON(http.StatusOK, lineResponse)
}

func generateHourlyIntervals(startTimeStr, endTimeStr string) ([]string, error) {
	const layout = "2006-01-02T15:04:05Z"

	// Parse the input time strings
	startTime, err := time.Parse(layout, startTimeStr)
	if err != nil {
		return nil, fmt.Errorf("invalid start time format: %v", err)
	}

	endTime, err := time.Parse(layout, endTimeStr)
	if err != nil {
		return nil, fmt.Errorf("invalid end time format: %v", err)
	}

	// Ensure startTime is before endTime
	if !startTime.Before(endTime) {
		return nil, fmt.Errorf("start time must be before end time")
	}

	// Generate the intervals
	var intervals []string
	current := startTime
	for current.Before(endTime) {
		intervals = append(intervals, current.Format("15:04")) // Only include HH:mm
		current = current.Add(time.Hour)
	}
	// Append the endTime if it doesn't align exactly
	if current.Before(endTime) || current.Equal(endTime) {
		intervals = append(intervals, current.Format("15:04"))
	}

	return intervals, nil
}

// TODO write the function to check if any production orders are running in that line, then assign,
func (v *ActionService) isLineStopped(lineId int) bool {
	return true
}
func CombineDateTimeToISO8601(shiftEndDate, shiftEndTime string) (string, error) {
	// Combine the date and time into a single string
	dateTime := fmt.Sprintf("%s %s", shiftEndDate, shiftEndTime)

	// Parse the date and time in MySQL format (YYYY-MM-DD HH:MM:SS)
	t, err := time.Parse("2006-01-02 15:04:05", dateTime)
	if err != nil {
		return "", err
	}

	// Format into ISO 8601 with Z (UTC time)
	return t.UTC().Format("2006-01-02T15:04:05Z"), nil
}
func ConvertToISO8601NoMilliseconds(timestamp string) (string, error) {
	// Parse the timestamp with milliseconds and UTC ('Z')
	parsedTime, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return "", err
	}

	// Add 8 hours to the parsed time
	updatedTime := parsedTime.Add(8 * time.Hour)
	// Format into ISO 8601 without milliseconds and keep UTC ('Z')
	return updatedTime.UTC().Format("2006-01-02T15:04:05Z"), nil
}

func addHoursToRFC3339(datetime string, hours int) string {
	// Parse the input time in RFC3339 format
	parsedTime, err := time.Parse(time.RFC3339, datetime)
	if err != nil {
		fmt.Println("Error parsing date time:", err)
		return ""
	}
	// Add the specified number of hours to the parsed time
	updatedTime := parsedTime.Add(time.Duration(hours) * time.Hour)

	// Format the updated time back to RFC3339 format
	return updatedTime.Format(time.RFC3339)
}

func generateHourlyIntervalsWithTimestamps(startTimeStr, endTimeStr string) ([]string, error) {
	// This will return
	// [2025-01-16T08:30:00Z , 2025-01-16T09:30:00Z ,2025-01-16T10:30:00Z]
	const layout = "2006-01-02T15:04:05Z"

	// Parse the input time strings
	startTime, err := time.Parse(layout, startTimeStr)
	if err != nil {
		return nil, fmt.Errorf("invalid start time format: %v", err)
	}

	endTime, err := time.Parse(layout, endTimeStr)
	if err != nil {
		return nil, fmt.Errorf("invalid end time format: %v", err)
	}

	// Ensure startTime is before endTime
	if !startTime.Before(endTime) {
		return nil, fmt.Errorf("start time must be before end time")
	}

	// Generate the intervals
	var intervals []string
	current := startTime
	for current.Before(endTime) {
		intervals = append(intervals, current.Format(layout)) // Use full ISO 8601 format
		current = current.Add(time.Hour)
	}
	// Append the endTime if it doesn't align exactly
	if current.Before(endTime) || current.Equal(endTime) {
		intervals = append(intervals, current.Format(layout))
	}

	return intervals, nil
}

func (v *ActionService) getShiftActualHours(eventId int, isMESDriverConfigured bool, historyResponse []assemblyManualOrderHistoryByEventResponse, arrayOfTimeInterval []string) []int {
	var result = make([]int, 0)
	v.Logger.Info("time interval list", zap.Any("arrayOfTimeInterval", arrayOfTimeInterval))
	var timeNow = util.GetCurrentTime(const_util.ISOTimeLayout)
	var timeNowEpoch = util.ConvertStringToDateTime(timeNow)
	if !isMESDriverConfigured {
		v.Logger.Info("getting data for driver configured events")
		// Sorting the data by CreatedAt in ascending order
		sort.Slice(historyResponse, func(i, j int) bool {
			// Parse the dates
			timeI, errI := time.Parse(time.RFC3339, historyResponse[i].CreatedAt)
			timeJ, errJ := time.Parse(time.RFC3339, historyResponse[j].CreatedAt)

			// Handle potential parsing errors
			if errI != nil || errJ != nil {
				return false // If parsing fails, maintain the current order
			}

			return timeI.Before(timeJ) // Compare the parsed times
		})
		v.Logger.Info("processing history response ordered", zap.Any("historyResponse", historyResponse))
		var intervalArrayLength = len(arrayOfTimeInterval)
		for intervalIndex, timeValue := range arrayOfTimeInterval {
			var timeValueUTC = util.ConvertSingaporeTimeToUTC(timeValue)
			var timeEpoch = util.ConvertStringToDateTime(timeValueUTC)
			if timeEpoch.DateTimeEpoch < timeNowEpoch.DateTimeEpoch {

				// Add next time interval into function
				var nextTimeIntervalEpoch int64
				if intervalArrayLength >= (intervalIndex + 1) {
					// Edge case. Then we will add one hour to current time epoch
					nextTimeIntervalEpoch = time.Unix(timeEpoch.DateTimeEpoch, 0).Add(time.Hour).Unix()
				} else {
					var nextTimeIntervalUTC = util.ConvertSingaporeTimeToUTC(arrayOfTimeInterval[intervalIndex+1])
					var nextTimeIntervalObject = util.ConvertStringToDateTime(nextTimeIntervalUTC)
					nextTimeIntervalEpoch = nextTimeIntervalObject.DateTimeEpoch

				}
				var actualValue = v.getLatestManualCount(historyResponse, timeEpoch.DateTimeEpoch, nextTimeIntervalEpoch)
				result = append(result, actualValue)
				//var isValueAdded = false
				//for _, historyElement := range historyResponse {
				//	var historyEpoch = util.ConvertStringToDateTime(historyElement.CreatedAt)
				//	v.Logger.Info("preparing the data timeEpoch.DateTimeEpoch ", zap.Any("epcch", timeEpoch.DateTimeEpoch), zap.Any("history", historyEpoch.DateTimeEpoch))
				//	if timeEpoch.DateTimeEpoch < historyEpoch.DateTimeEpoch {
				//		result = append(result, historyElement.CompletedQuantity)
				//		isValueAdded = true
				//		break
				//	}
				//}
				//if !isValueAdded {
				//	result = append(result, 0)
				//}

			}

		}
		return result
	} else {
		for _, timeValue := range arrayOfTimeInterval {
			var timeValueUTC = util.ConvertSingaporeTimeToUTC(timeValue)
			var timeEpoch = util.ConvertStringToDateTime(timeValueUTC)
			if timeEpoch.DateTimeEpoch < timeNowEpoch.DateTimeEpoch {
				statsInfo, err := v.MachineStatsCache.GetActualValueFromEventId(eventId, timeEpoch.DateTimeEpoch)
				if err != nil {
					v.Logger.Error("error getting actual time from cache", zap.Error(err))
				}
				v.Logger.Info("searching event and timestamp", zap.Any("event_id", eventId), zap.Any("timestamp", timeValueUTC), zap.Any("epoch", timeEpoch.DateTimeEpoch), zap.Any("actual_value", statsInfo.Actual))
				result = append(result, statsInfo.Actual)
			} else {
				break
			}
		}
	}

	v.Logger.Info("cached values", zap.Any("cache_data", v.ShiftActualValueCache))
	for _, value := range v.ShiftActualValueCache[eventId] {
		result = append(result, value.ActualValue)
	}
	v.Logger.Info("return value", zap.Any("result", result))
	return result
}

func (v *ActionService) getLatestManualCount(historyResponse []assemblyManualOrderHistoryByEventResponse, startTimeEpoch int64, endTimeEpoch int64) int {
	// This function check latest value in particular time interval
	var result = 0

	if startTimeEpoch == 0 || endTimeEpoch == 0 {
		return result
	}

	for _, historyElement := range historyResponse {
		var historyEpoch = util.ConvertStringToDateTime(historyElement.CreatedAt)

		if startTimeEpoch < historyEpoch.DateTimeEpoch && historyEpoch.DateTimeEpoch < endTimeEpoch {
			result = historyElement.CompletedQuantity

		}

		// Cut the redundent loops
		if historyEpoch.DateTimeEpoch > endTimeEpoch {
			break
		}
	}
	v.Logger.Info("preparing the data timeEpoch.DateTimeEpoch ", zap.Any("startTime", startTimeEpoch), zap.Any("endTime", endTimeEpoch))
	return result
}
