package jobs

import (
	"cx-micro-flake/pkg/common"
	"encoding/json"
	"go.uber.org/zap"
	"gorm.io/datatypes"
	"time"
)

type machineStatsData struct {
	Ts        int64          `json:"ts"`
	StatsInfo datatypes.JSON `json:"statsInfo"`
}

func (v *JobService) LoadMachineStats() {
	if v.PoolingInterval == 0 || v.PoolingInterval < 0 {
		v.PoolingInterval = 30
	}
	v.Logger.Info("machine stats pooling is starting up....", zap.Int("pooling_interval", v.PoolingInterval))
	var duration = time.Duration(v.PoolingInterval) * time.Second
	for {
		machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
		var last12HoursStatsData = machineService.GetLastNHoursAssemblyStats(12)
		var machineStats []machineStatsData
		err := json.Unmarshal(last12HoursStatsData, &machineStats)
		if err == nil {
			// clear the cache before load it
			v.MachineStatsCache.ClearCache()
			for _, machineStat := range machineStats {
				err := v.MachineStatsCache.Insert(machineStat.Ts, machineStat.StatsInfo)
				if err != nil {
					v.Logger.Error("error inserting element to cache", zap.Error(err))
				} else {
					v.Logger.Info("inserting the cache element", zap.Any("ts", machineStat.Ts), zap.Any("stats_info", machineStat.StatsInfo))
				}
			}

		} else {
			v.Logger.Error("error unmarshalling json", zap.Error(err))
		}
		time.Sleep(duration)
	}
}
