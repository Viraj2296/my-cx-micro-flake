package views

import (
	"cx-micro-flake/services/machines/handler/database"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ViewManager struct {
	DbConn *gorm.DB
	Logger *zap.Logger
}

// UpdateAssemblyView update the assembly view using machine ID
func (v *ViewManager) UpdateAssemblyView(machineId int, viewData map[string]interface{}) error {
	if err := v.DbConn.Model(&database.AssemblyMachineView{Id: machineId}).Updates(viewData).Error; err != nil {
		return err
	}
	return nil
}

func (v *ViewManager) CreateNewView(machineId int) error {
	if err := v.DbConn.Save(&database.AssemblyMachineView{Id: machineId, DelayStatus: "Delayed Feed"}).Error; err != nil {
		return err
	}
	return nil
}

func (v *ViewManager) CreateOrUpdateMouldingView(machineId int, viewData map[string]interface{}) error {
	var machineView database.MouldingMachineView

	// Check if the record exists
	err := v.DbConn.First(&machineView, "id = ?", machineId).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err {
			// Record does not exist, create a new one
			machineView = database.MouldingMachineView{
				Id:          machineId,
				DelayStatus: "Delayed Feed",
			}
			if err := v.DbConn.Create(&machineView).Error; err != nil {
				return err
			}
		} else {
			// Another error occurred
			return err
		}
	} else {
		// Record exists, update it
		if err := v.DbConn.Model(&database.MouldingMachineView{Id: machineId}).Updates(viewData).Error; err != nil {
			return err
		}
	}

	return nil
}

func (v *ViewManager) CreateOrUpdateToolingView(machineId int, viewData map[string]interface{}) error {
	var machineView database.ToolingMachineView

	// Check if the record exists
	err := v.DbConn.First(&machineView, "id = ?", machineId).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err {
			// Record does not exist, create a new one
			machineView = database.ToolingMachineView{
				Id:          machineId,
				DelayStatus: "Delayed Feed",
			}
			if err := v.DbConn.Create(&machineView).Error; err != nil {
				return err
			}
		} else {
			// Another error occurred
			return err
		}
	} else {
		// Record exists, update it
		if err := v.DbConn.Model(&database.ToolingMachineView{Id: machineId}).Updates(viewData).Error; err != nil {
			return err
		}
	}

	return nil
}

func (v *ViewManager) CreateOrUpdateAssemblyView(machineId int, viewData map[string]interface{}) error {
	var machineView database.AssemblyMachineView

	// Check if the record exists
	err := v.DbConn.First(&machineView, "id = ?", machineId).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err {
			// Record does not exist, create a new one
			machineView = database.AssemblyMachineView{
				Id:          machineId,
				DelayStatus: "Delayed Feed",
			}
			if err := v.DbConn.Create(&machineView).Error; err != nil {
				return err
			}
		} else {
			// Another error occurred
			return err
		}
	} else {
		// Record exists, update it
		if err := v.DbConn.Model(&database.AssemblyMachineView{Id: machineId}).Updates(viewData).Error; err != nil {
			return err
		}
	}

	return nil
}
