package actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/services/energy_management/source/consts"
	"cx-micro-flake/services/energy_management/source/dto"
	"cx-micro-flake/services/energy_management/source/models"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.cerex.io/transcendflow/logging"
	time_series "go.cerex.io/transcendflow/time-series"
	"go.cerex.io/transcendflow/util"
	"go.uber.org/zap"
)

// GetMachineEnergyUsageDisplay handles the endpoint for getting energy usage
func (v *Actions) GetMachineEnergyUsageDisplay(ctx *gin.Context) {
	v.Logger.Info("getting the machine energy usage..")
	energyManagementDisplayRequest := dto.EnergyManagementDisplayRequest{}
	if err := ctx.ShouldBindBodyWith(&energyManagementDisplayRequest, binding.JSON); err != nil {
		v.Logger.Error("invalid request", logging.Error(err))
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		v.Logger.Error("error sending response", zap.Error(err))
		return
	}
	var overviewResponse = make(map[string]interface{})

	machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
	_, machineMasterData := machineService.GetMachineInfoById(consts.ProjectID, energyManagementDisplayRequest.MachineId)

	var machineMasterInfo = models.GetMachineMasterInfo(machineMasterData.ObjectInfo)

	var machineConnectStatus string
	var maintenanceColorCode string
	if machineMasterInfo.MachineConnectStatus == consts.MachineConnectStatusMaintenance {
		machineConnectStatus = "Maintenance"
		maintenanceColorCode = "#A31041"
	} else {
		if machineMasterInfo.MachineConnectStatus == consts.MachineConnectStatusLive {
			machineConnectStatus = "Running"
			maintenanceColorCode = "#32602e"
		} else {
			machineConnectStatus = "Stopped"
			maintenanceColorCode = "#E60E18"
		}
	}
	mouldInterface := common.GetService("moulds_module").ServiceInterface.(common.MouldInterface)
	_, brandName := mouldInterface.GetBrandName(consts.ProjectID, machineMasterInfo.Brand)
	overviewResponse["brand"] = brandName
	overviewResponse["machineImage"] = machineMasterInfo.MachineImage
	overviewResponse["model"] = machineMasterInfo.Model
	overviewResponse["newMachineId"] = machineMasterInfo.NewMachineId
	overviewResponse["currentStatus"] = machineConnectStatus
	overviewResponse["colorCode"] = maintenanceColorCode

	now := time.Now()
	var currentDistributionStartDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Format(time.RFC3339Nano)
	startOfDayStr := v.EnergyStartDate
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 0, now.Location())
	var endDateStr = endOfDay.Format(time.RFC3339Nano)
	totalPowerUsageResult, _ := v.getTotalPowerDistribution(currentDistributionStartDate, endDateStr)
	totalIdlePowerUse, _ := v.getTotalIdlePowerUsage(startOfDayStr, endDateStr, energyManagementDisplayRequest.MachineId, consts.ProjectID)
	var totalPowerSinceBeginning = v.getTotalEnergyUsage()
	var energyUsageTs []string
	var energyUsageData []float64

	for _, valueMap := range totalPowerUsageResult {
		var tsValue = util.InterfaceToString(valueMap["timestamp"])
		energyUsageTs = append(energyUsageTs, tsValue)
		var enerValue = util.InterfaceToFloat(valueMap["value"])
		energyUsageData = append(energyUsageData, enerValue)
	}

	overviewResponse["totalPowerUsage"] = totalPowerSinceBeginning
	overviewResponse["totalIdlePowerUse"] = totalIdlePowerUse
	var energyChart = dto.LineChartResponse{}
	energyChart.Chart.Type = "spline"
	energyChart.Title.Text = "Hourly Power Usage"
	energyChart.XAxis.Title.Text = "Time (Hour)"
	energyChart.XAxis.Categories = energyUsageTs
	energyChart.YAxis.Title.Text = "Average Watts/Hour"
	energyChart.YAxis.Min = 0
	energyChart.Credits.Enabled = false
	energyChart.Series = append(energyChart.Series, dto.Series{Name: "Power", Data: energyUsageData, Type: "spline"})

	overviewResponse["energyUsage"] = energyChart

	v.Logger.Info("current distribution start date", logging.String("startDate", currentDistributionStartDate), logging.String("endDate", endDateStr))
	phase1Series, phase2Series, phase3Series, ts := v.getCurrentDistribution(currentDistributionStartDate, endDateStr)

	var currentChart = dto.LineChartResponse{}
	currentChart.Chart.Type = "spline"
	currentChart.Title.Text = "Current Usage"
	currentChart.XAxis.Title.Text = "Time (min)"
	currentChart.XAxis.Categories = ts
	currentChart.YAxis.Title.Text = "amperes (A)"
	currentChart.YAxis.Min = 0
	currentChart.Credits.Enabled = false
	currentChart.Series = append(currentChart.Series, dto.Series{Name: "Phase 1", Data: phase1Series, Type: "spline"})
	currentChart.Series = append(currentChart.Series, dto.Series{Name: "Phase 2", Data: phase2Series, Type: "spline"})
	currentChart.Series = append(currentChart.Series, dto.Series{Name: "Phase 3", Data: phase3Series, Type: "spline"})

	overviewResponse["energyUsage"] = energyChart
	overviewResponse["currentChart"] = currentChart

	ctx.JSON(http.StatusOK, overviewResponse)

}

func (v *Actions) getTotalEnergyUsage() string {
	var query = `from(bucket: "fuyu_iot_sensors_data")
  |> range(start: 0)
  |> filter(fn: (r) => 
       r["_measurement"] == "murata_current_sensors_data" and
       (r["_field"] == "current1" or r["_field"] == "current2") and
       (r["sensor_node"] == "A7FD" or r["sensor_node"] == "D7D4")
     )
  |> aggregateWindow(every: 1h, fn: mean, createEmpty: false)
  |> pivot(rowKey:["_time", "sensor_node"], columnKey: ["_field"], valueColumn: "_value")
  |> map(fn: (r) => ({
       _time: r._time,
       sensor_node: r.sensor_node,
       power: if r.sensor_node == "A7FD" then
                400.0 * (if exists r.current1 then r.current1 else 0.0) +
                400.0 * (if exists r.current2 then r.current2 else 0.0)
              else if r.sensor_node == "D7D4" then
                400.0 * (if exists r.current1 then r.current1 else 0.0)
              else 0.0
     }))
  |> group()
  |> sum(column: "power")
  |> rename(columns: {power: "total_power"})  // This value is in kWh
  |> yield(name: "total_power")
`
	queryResults, err := v.RealtimeDBManager.QueryMessagesRange(query)
	if err != nil {
		v.Logger.Error("query to get the today total power usage", logging.Error(err))
		return "0 kW"
	}

	if len(queryResults) == 0 {
		v.Logger.Info("query has return empty results")
		return "0 kW"
	}
	var totalPower float64
	for _, result := range queryResults {

		for key, value := range result.Values {
			v.Logger.Info("results", logging.String("key", key), logging.Any("value", value))

			if key == "total_power" {
				totalPower = util.InterfaceToFloat(value)
			}
		}

	}
	formattedTotalPower := fmt.Sprintf("%.2f", totalPower/1000)
	return formattedTotalPower
}
func (v *Actions) getTotalPowerDistribution(startDate string, endDate string) ([]map[string]interface{}, error) {
	var err error
	var results = make([]map[string]interface{}, 0)
	var queryResults = make([]time_series.QueryResult, 0)

	var query = fmt.Sprintf(`
from(bucket: "fuyu_iot_sensors_data")
  |> range(start: time(v: "%s"), stop: time(v: "%s"))
  |> filter(fn: (r) => 
       r["_measurement"] == "murata_current_sensors_data" and
       (r["_field"] == "current1" or r["_field"] == "current2") and
       (r["sensor_node"] == "A7FD" or r["sensor_node"] == "D7D4")
     )
  |> aggregateWindow(every: 30m, fn: mean, createEmpty: false)
  |> pivot(rowKey:["_time", "sensor_node"], columnKey: ["_field"], valueColumn: "_value")
  |> map(fn: (r) => ({
       _time: r._time,
       sensor_node: r.sensor_node,
       power: if r.sensor_node == "A7FD" then
                400.0 * (if exists r.current1 then r.current1 else 0.0) +
                400.0 * (if exists r.current2 then r.current2 else 0.0)
              else if r.sensor_node == "D7D4" then
                400.0 * (if exists r.current1 then r.current1 else 0.0)
              else 0.0
     }))
  |> group(columns: ["_time"])
  |> sum(column: "power")
  |> rename(columns: {power: "total_power"})
    `, startDate, endDate)

	v.Logger.Info("running query to get the today total power usage", logging.String("query", query))
	queryResults, err = v.RealtimeDBManager.QueryMessagesRange(query)
	if err != nil {
		v.Logger.Error("query to get the today total power usage", logging.Error(err))
		return results, err
	}

	// Print results
	if len(queryResults) == 0 {
		v.Logger.Info("query has return empty results")
		return results, err
	}

	for _, result := range queryResults {

		var tsValue string
		var totalPower float64
		for key, value := range result.Values {
			v.Logger.Info("results", logging.String("key", key), logging.Any("value", value))
			if key == "_time" {
				var timeValue = value.(time.Time)
				tsValue = timeValue.String()
			}
			if key == "total_power" {
				totalPower = util.InterfaceToFloat(value)
			}
		}

		if tsValue != "" {

			var resultMap = make(map[string]interface{})
			// Use a layout matching the input format
			layout := "2006-01-02 15:04:05 -0700 MST"
			t, err := time.Parse(layout, tsValue)
			singaporeTime := t.Add(8 * time.Hour)
			if err != nil {
				v.Logger.Error("error converting timestamp format", logging.String("tsValue", tsValue))
			} else {
				// Format to "HH:MM:SS"
				formatted := singaporeTime.Format("15:04")
				resultMap["timestamp"] = formatted
				resultMap["value"] = totalPower
				results = append(results, resultMap)
			}
		}
	}
	v.Logger.Info("got the results for getTotalPowerDistribution ", logging.Any("results", len(results)))
	return results, nil
}

// func (v *Actions) getTodayIdlePowerUsage(machineRecordId int, projectId string) ([]map[string]interface{}, float64, error) {
// 	var err error
// 	var totalIdlePowerUsage float64
// 	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
// 	err, scheduledEventObject := productionOrderInterface.GetCurrentAssemblyScheduledEvent(projectId, machineRecordId)
// 	var intervalList = make([]map[string]interface{}, 0)

// 	if scheduledEventObject.Id == 0 {
// 		intervalList, _ = CreateNonOverlappingChunks("", "")
// 	} else {
// 		var scheduledEventData = make(map[string]interface{})
// 		json.Unmarshal(scheduledEventObject.ObjectInfo, &scheduledEventData)

// 		var startDate = util.InterfaceToString(scheduledEventData["startDate"])
// 		var endDate = util.InterfaceToString(scheduledEventData["endDate"])

// 		intervalList, _ = CreateNonOverlappingChunks(startDate, endDate)
// 	}

// 	var finalResults = make([]map[string]interface{}, 0)
// 	for _, interval := range intervalList {
// 		var startDate = util.InterfaceToString(interval["start"])
// 		var endDate = util.InterfaceToString(interval["end"])

// 		var results = make([]map[string]interface{}, 0)
// 		var queryResults = make([]time_series.QueryResult, 0)
// 		var query = fmt.Sprintf(`
// 		from(bucket: "fuyu_iot_sensor_data")
// 		  |> range(start: %q, stop: %q)
// 		  |> filter(fn: (r) => r["_measurement"] == "murata_current_sensors_data")
// 		  |> filter(fn: (r) => r["_field"] == "battery" or r["_field"] == "current1" or r["_field"] == "current2")
// 		  |> filter(fn: (r) => r["sensor_node"] == "A7FD" or r["sensor_node"] == "D7D4")
// 		  |> aggregateWindow(every: 1h, fn: mean, createEmpty: false)
// 		  |> yield(name: "mean")
// 		`, startDate, endDate)

// 		queryResults, err = v.RealtimeDBManager.QueryMessagesRange(query)
// 		if err != nil {
// 			return results, totalIdlePowerUsage, err
// 		}

// 		// Print results
// 		if len(queryResults) == 0 {
// 			fmt.Println("No results returned.")
// 			return results, totalIdlePowerUsage, err
// 		}

// 		for _, result := range queryResults {
// 			var tsValue string
// 			var tsValueFloat float64
// 			for key, value := range result.Values {
// 				if key == "_time" {
// 					tsValue = util.InterfaceToString(value)
// 				}

// 				if key == "_value" {
// 					tsValueFloat = util.InterfaceToFloat(value) * 240
// 				}
// 				totalIdlePowerUsage += tsValueFloat
// 			}
// 			var resultMap = make(map[string]interface{})
// 			resultMap["timestamp"] = tsValue
// 			resultMap["value"] = tsValueFloat
// 			results = append(results, resultMap)

// 		}
// 		finalResults = append(finalResults, results...)
// 	}

// 	return finalResults, totalIdlePowerUsage, err
// }

// func CreateNonOverlappingChunks(startInterval, endInterval string) ([]map[string]interface{}, error) {
// 	var result = make([]map[string]interface{}, 0)
// 	// Get the current time
// 	now := time.Now()

// 	// Get the start of the day (00:00:00)
// 	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

// 	// Convert the start of the day to a Unix timestamp
// 	startOfDayStr := startOfDay.Format(time.RFC3339Nano)

// 	// Get the end of the day (23:59:59)
// 	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

// 	var endDateStr = endOfDay.Format(time.RFC3339Nano)
// 	if startInterval == "" || endInterval == "" {
// 		var defaultDate = map[string]interface{}{
// 			"start": startOfDayStr,
// 			"end":   endDateStr,
// 		}
// 		return []map[string]interface{}{defaultDate}, nil
// 	}

// 	var firstChunkEndDate = startInterval
// 	var secondChunkStartDate = endInterval

// 	var firstChunk = map[string]interface{}{
// 		"start": startOfDayStr,
// 		"end":   firstChunkEndDate,
// 	}
// 	var secondChunk = map[string]interface{}{
// 		"start": secondChunkStartDate,
// 		"end":   endDateStr,
// 	}
// 	result = append(result, firstChunk)
// 	result = append(result, secondChunk)

// 	return result, nil
// }

func (v *Actions) getCurrentDistribution(startDate, endDate string) ([]float64, []float64, []float64, []string) {
	var query = fmt.Sprintf(`from(bucket: "fuyu_iot_sensors_data")
 |> range(start: time(v: "%s"), stop: time(v: "%s"))
  |> filter(fn: (r) => 
    r["_measurement"] == "murata_current_sensors_data" and
    (
      (r["sensor_node"] == "A7FD" and (r["_field"] == "current1" or r["_field"] == "current2")) or
      (r["sensor_node"] == "D7D4" and r["_field"] == "current1")
    )
  )
  |> aggregateWindow(every: 5m, fn: mean, createEmpty: false)
  |> map(fn: (r) => ({
    _time: r._time,
    _value: r._value,
    _field: if r["sensor_node"] == "D7D4" and r["_field"] == "current1" then "current3" else r["_field"],
    sensor_node: r["sensor_node"]
  }))`, startDate, endDate)
	v.Logger.Info("running query to get the current distribution", logging.String("query", query))
	var phase1currentDataSeries []float64 // current1
	var phase2currentDataSeries []float64 // current2
	var phase3currentDataSeries []float64 // current3
	var timestamp []string

	queryResults, err := v.RealtimeDBManager.QueryMessagesRange(query)
	if err != nil {
		v.Logger.Error("query to get current distribution failed", logging.Error(err))
		return phase1currentDataSeries, phase2currentDataSeries, phase3currentDataSeries, timestamp
	}

	if len(queryResults) == 0 {
		v.Logger.Info("query returned empty results")
		return phase1currentDataSeries, phase2currentDataSeries, phase3currentDataSeries, timestamp
	}
	var seenTimestamps = make(map[string]bool)
	for _, result := range queryResults {
		var timeStr string
		var value float64
		var field string

		for key, val := range result.Values {
			switch key {
			case "_time":
				if t, ok := val.(time.Time); ok {
					// Convert to Singapore time (UTC+8)
					singaporeTime := t.Add(8 * time.Hour)
					timeStr = singaporeTime.Format("15:04:05")
				}
			case "_value":
				switch v := val.(type) {
				case float64:
					value = v
				case int:
					value = float64(v)
				}
			case "_field":
				if s, ok := val.(string); ok {
					field = s
				}
			}
		}

		if _, exists := seenTimestamps[timeStr]; !exists {
			timestamp = append(timestamp, timeStr)
			seenTimestamps[timeStr] = true
		}

		switch field {
		case "current1":
			phase1currentDataSeries = append(phase1currentDataSeries, value)
		case "current2":
			phase2currentDataSeries = append(phase2currentDataSeries, value)
		case "current3":
			phase3currentDataSeries = append(phase3currentDataSeries, value)
		}
	}
	return phase1currentDataSeries, phase2currentDataSeries, phase3currentDataSeries, timestamp
}

// This function is used to find idle times in a given master interval based on production intervals.
type Interval struct {
	Start time.Time
	End   time.Time
}

func ParseInterval(start, end string) (Interval, error) {
	layout := time.RFC3339Nano
	s, err := time.Parse(layout, start)
	if err != nil {
		return Interval{}, err
	}
	e, err := time.Parse(layout, end)
	if err != nil {
		return Interval{}, err
	}
	return Interval{Start: s, End: e}, nil
}

func FindIdleTimes(master Interval, productions []Interval) []Interval {
	if len(productions) == 0 {
		return []Interval{master}
	}

	// Sort production intervals by start time
	sort.Slice(productions, func(i, j int) bool {
		return productions[i].Start.Before(productions[j].Start)
	})

	idleTimes := []Interval{}
	curr := master.Start

	for _, p := range productions {
		// Ignore out-of-bound production intervals
		if p.End.Before(master.Start) || p.Start.After(master.End) {
			continue
		}

		// Clamp production interval to master interval
		start := maxTime(p.Start, master.Start)
		end := minTime(p.End, master.End)

		if curr.Before(start) {
			idleTimes = append(idleTimes, Interval{Start: curr, End: start})
		}

		if end.After(curr) {
			curr = end
		}
	}

	if curr.Before(master.End) {
		idleTimes = append(idleTimes, Interval{Start: curr, End: master.End})
	}

	return idleTimes
}

func minTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func (v Actions) getTotalIdlePowerUsage(startDate string, endDate string, machineRecordId int, projectId string) (string, error) {
	var err error
	var totalIdlePowerUsage float64
	var formattedTotalPower = "0"

	// startDate and end date 2022-12-21T14:28:46.805Z
	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	err, listOfScheduledEventObject := productionOrderInterface.GetGivenTimeIntervalScheduledEvent(projectId, machineRecordId, startDate)

	var productionRaw = make([]struct{ Start, End string }, 0)
	v.Logger.Info("Master Date: %s -> %s\n", logging.String("start", startDate), logging.String("end", endDate))
	for _, scheduledEventObject := range listOfScheduledEventObject {
		var scheduledEventData = make(map[string]interface{})
		json.Unmarshal(scheduledEventObject.ObjectInfo, &scheduledEventData)
		var startDate = util.InterfaceToString(scheduledEventData["startDate"])
		var endDate = util.InterfaceToString(scheduledEventData["endDate"])

		var intervalData = struct{ Start, End string }{startDate, endDate}

		productionRaw = append(productionRaw, intervalData)
		v.Logger.Info("Schedule Date: %s -> %s\n", logging.String("start", startDate), logging.String("end", endDate))
	}

	master, _ := ParseInterval(startDate, endDate)

	// productionRaw := []struct{ Start, End string }{
	// 	{"2022-12-21T02:28:46.805Z", "2022-12-21T14:28:46.805Z"},
	// 	{"2022-12-22T02:28:46.805Z", "2022-12-22T14:28:46.805Z"},
	// }

	var productions []Interval
	for _, p := range productionRaw {
		interval, _ := ParseInterval(p.Start, p.End)
		productions = append(productions, interval)
	}

	idles := FindIdleTimes(master, productions)

	totalIdlePowerUsage, err = v.getIdlePowerUsage(idles, machineRecordId)
	if err != nil {
		v.Logger.Error("error getting idle power usage", logging.Error(err))
	}

	if totalIdlePowerUsage > 0 {
		formattedTotalPower = fmt.Sprintf("%.2f", totalIdlePowerUsage/1000)
	}

	return formattedTotalPower, err

}

func (v *Actions) getIdlePowerUsage(listOfIdleIntervals []Interval, machineId int) (float64, error) {
	var err error
	var totalIdlePowerUsage float64

	var cachedData = make([]CachedIdleData, 0)
	if val, ok := v.IdlePowerCache[machineId]; ok {
		cachedData = val
	}
	for _, idleInterval := range listOfIdleIntervals {
		var cachedData = GetCacheData(idleInterval, cachedData)
		if cachedData != nil {
			totalIdlePowerUsage += cachedData.IntervalPowerUseage
			v.Logger.Info("cached data", logging.Any("cachedData", cachedData))
			v.Logger.Info("cached data interval", logging.Any("cachedInterval", idleInterval))
		} else {
			var intervalStartDate = idleInterval.Start.Format(time.RFC3339Nano)
			var intervalEndDate = idleInterval.End.Format(time.RFC3339Nano)
			v.Logger.Info("Idle: %s -> %s\n", logging.String("start", idleInterval.Start.Format(time.RFC3339Nano)), logging.String("end", idleInterval.End.Format(time.RFC3339Nano)))
			var query = fmt.Sprintf(`
				from(bucket: "fuyu_iot_sensors_data")
				|> range(start: time(v: "%s"), stop: time(v: "%s"))
				|> filter(fn: (r) => 
					r["_measurement"] == "murata_current_sensors_data" and
					(r["_field"] == "current1" or r["_field"] == "current2") and
					(r["sensor_node"] == "A7FD" or r["sensor_node"] == "D7D4")
					)
				|> aggregateWindow(every: 30m, fn: mean, createEmpty: false)
				|> pivot(rowKey:["_time", "sensor_node"], columnKey: ["_field"], valueColumn: "_value")
				|> map(fn: (r) => ({
					_time: r._time,
					sensor_node: r.sensor_node,
					power: if r.sensor_node == "A7FD" then
								400.0 * (if exists r.current1 then r.current1 else 0.0) +
								400.0 * (if exists r.current2 then r.current2 else 0.0)
							else if r.sensor_node == "D7D4" then
								400.0 * (if exists r.current1 then r.current1 else 0.0)
							else 0.0
					}))
				|> group(columns: ["_time"])
				|> sum(column: "power")
				|> rename(columns: {power: "total_power"})
    	`, intervalStartDate, intervalEndDate)

			v.Logger.Info("running query to get the idle power usage", logging.String("query", query))
			var queryResults = make([]time_series.QueryResult, 0)
			queryResults, err = v.RealtimeDBManager.QueryMessagesRange(query)
			if err != nil {
				v.Logger.Error("query to get the today total power usage", logging.Error(err))
				continue
			}

			// Print results
			if len(queryResults) == 0 {
				v.Logger.Info("query has return empty results")
				continue
			}

			var powerValue float64
			for _, result := range queryResults {
				for key, value := range result.Values {
					if key == "total_power" {
						totalPower := util.InterfaceToFloat(value)
						powerValue = totalPower
						v.Logger.Info("results", logging.String("key", key), logging.Any("value", value))
						totalIdlePowerUsage += totalPower
					}
				}
			}

			// Update the cache with the new interval data
			cachedInterval := CachedInterval{
				Start:               idleInterval.Start,
				End:                 idleInterval.End,
				IntervalPowerUseage: powerValue,
			}
			v.Set(machineId, cachedInterval)
			v.Logger.Info("cached data", logging.Any("cachedData", cachedInterval))
		}
	}

	return totalIdlePowerUsage, err

}

func (c *Actions) Set(machineID int, interval CachedInterval) {
	if _, exists := c.IdlePowerCache[machineID]; !exists {
		c.IdlePowerCache[machineID] = make([]CachedIdleData, 0)
	}

	c.IdlePowerCache[machineID] = append(c.IdlePowerCache[machineID], CachedIdleData{
		Interval: interval,
	})
}

func GetCacheData(interval Interval, cachedIntervals []CachedIdleData) *CachedInterval {

	for _, item := range cachedIntervals {
		if item.Interval.Start.Equal(interval.Start) && item.Interval.End.Equal(interval.End) {
			return &item.Interval
		}
	}
	return nil
}

type CachedInterval struct {
	Start               time.Time
	End                 time.Time
	IntervalPowerUseage float64
}

// CacheKeyedByMachine stores intervals and results per machine
type IdlePowerUsageCache struct {
	Data map[int][]CachedIdleData
}

type CachedIdleData struct {
	Interval CachedInterval
}
