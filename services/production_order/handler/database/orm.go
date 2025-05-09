package database

import "time"

// AssemblySchedulerEventHistory view of any tables should be linked view event history
type AssemblySchedulerEventHistory struct {
	Id            int       `gorm:"column:id;not null;primaryKey;autoIncrement" json:"id"`
	EventId       int       `gorm:"column:event_id;not null;index;foreignKey:event_id;references:Id" json:"eventId"` // Foreign key
	EventStatusId int       `gorm:"column:event_status_id" json:"eventStatusId"`
	CreatedAt     time.Time `gorm:"column:created_at;not null;primaryKey;autoCreateTime" json:"createdAt"`
}
