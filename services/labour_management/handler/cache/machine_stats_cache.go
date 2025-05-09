package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/datatypes"
	"sort"
)

type MachineStatisticsInfo struct {
	EventId                    int      `json:"eventId"`
	ProductionOrderId          int      `json:"productionOrderId"`
	CurrentStatus              string   `json:"currentStatus"`
	PartId                     int      `json:"partId"`
	ScheduleStartTime          string   `json:"scheduleStartTime"`
	ScheduleEndTime            string   `json:"scheduleEndTime"`
	ActualStartTime            string   `json:"actualStartTime"`
	ActualEndTime              string   `json:"actualEndTime"`
	EstimatedEndTime           string   `json:"estimatedEndTime"`
	Oee                        int      `json:"oee"`
	Availability               int      `json:"availability"`
	Performance                int      `json:"performance"`
	Quality                    int      `json:"quality"`
	PlannedQuality             int      `json:"plannedQuality"`
	Completed                  int      `json:"completed"`
	Rejects                    int      `json:"rejects"`
	CompletedPercentage        float32  `json:"completedPercentage"`
	OverallCompletedPercentage float32  `json:"overallCompletedPercentage"`
	Daily                      int      `json:"daily"`
	Actual                     int      `json:"actual"`
	ExpectedProductQyt         int      `json:"expectedProductQyt"`
	OverallRejectedQty         int      `json:"overallRejectedQty"`
	ProgressPercentage         float64  `json:"progressPercentage"`
	DailyPlannedQty            int      `json:"dailyPlannedQty"`
	DownTime                   int      `json:"downTime"`
	Remark                     string   `json:"remark"`
	WarningMessage             []string `json:"warningMessage"`
}

// MachineStatsCache represents the information stored in the cache.
type MachineStatsCache struct {
	CacheData   map[int64]MachineStatisticsInfo
	EventIdData map[int][]int64
}

// NewMachineStatsCache initializes and returns a new MachineStatsCache.
func NewMachineStatsCache() *MachineStatsCache {
	return &MachineStatsCache{
		CacheData:   make(map[int64]MachineStatisticsInfo),
		EventIdData: make(map[int][]int64),
	}
}

// Insert adds a new entry to the cache. Returns an error if the key already exists.
func (msc *MachineStatsCache) Insert(ts int64, info datatypes.JSON) error {
	var machineStatisticsInfo MachineStatisticsInfo
	err := json.Unmarshal(info, &machineStatisticsInfo)
	if err == nil {
		if _, exists := msc.CacheData[ts]; exists {
			return fmt.Errorf("entry with timestamp %d already exists", ts)
		}
		msc.CacheData[ts] = machineStatisticsInfo
		msc.EventIdData[machineStatisticsInfo.EventId] = append(msc.EventIdData[machineStatisticsInfo.EventId], ts)
		fmt.Println("inserting data", ts, " event ID ", machineStatisticsInfo.EventId, " actual value ", machineStatisticsInfo.Actual)
	}
	return err
}

// GetActualValueFromEventId retrieves the first entry where the timestamp is less than the given `ts` for the specified `eventId`.
func (msc *MachineStatsCache) GetActualValueFromEventId(eventId int, ts int64) (*MachineStatisticsInfo, error) {
	// Get the list of timestamps for the given eventId.
	timestamps, exists := msc.EventIdData[eventId]
	if !exists {
		return &MachineStatisticsInfo{}, errors.New("no timestamps found for the given eventId")
	}

	// Sort the timestamps in descending order.
	sort.Slice(timestamps, func(i, j int) bool {
		return timestamps[i] > timestamps[j]
	})

	// Find the first timestamp less than the given ts.
	for _, t := range timestamps {
		if t < ts {
			if info, exists := msc.CacheData[t]; exists {
				return &info, nil
			}
		}
	}

	return &MachineStatisticsInfo{}, errors.New("no matching entry found")
}

// Update modifies an existing entry in the cache. Returns an error if the key does not exist.
func (msc *MachineStatsCache) Update(ts int64, info datatypes.JSON) error {
	var machineStatisticsInfo MachineStatisticsInfo
	err := json.Unmarshal(info, &machineStatisticsInfo)
	if err == nil {
		if _, exists := msc.CacheData[ts]; !exists {
			return fmt.Errorf("entry with timestamp %d does not exist", ts)
		}
		msc.CacheData[ts] = machineStatisticsInfo
		msc.EventIdData[machineStatisticsInfo.EventId] = append(msc.EventIdData[machineStatisticsInfo.EventId], ts)
	}

	return err
}

// ClearCache clears all data from the cache.
func (msc *MachineStatsCache) ClearCache() {
	msc.CacheData = make(map[int64]MachineStatisticsInfo)
}

// PrintCache prints the current state of the cache for debugging.
func (msc *MachineStatsCache) PrintCache() {
	for ts, info := range msc.CacheData {
		fmt.Printf("Timestamp: %d, Info: %+v\n", ts, info)
	}
}
