package views

import (
	"cx-micro-flake/services/production_order/handler/database"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ViewManager struct {
	DbConn *gorm.DB
	Logger *zap.Logger
}

func (v *ViewManager) CreateLabourManagementShiftLinesHistory(eventId int, eventStatusId int) error {
	if err := v.DbConn.Save(&database.AssemblySchedulerEventHistory{
		EventId:       eventId,
		EventStatusId: eventStatusId,
	}).Error; err != nil {
		return err
	}
	return nil
}

type EventStatusCount struct {
	EventID       int `json:"eventId"`
	EventStatusID int `json:"EventStatusId"`
	TotalCount    int `json:"totalCount"`
}

func (v *ViewManager) GetEventStatusCounts() ([]EventStatusCount, error) {
	var results []EventStatusCount

	err := v.DbConn.Table("assembly_scheduler_event_history").
		Select("event_id, event_status_id, COUNT(*) AS total_count").
		Where("DATE(created_at) = CURDATE()").
		Group("event_id, event_status_id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}
