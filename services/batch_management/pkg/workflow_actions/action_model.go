package workflow_actions

import (
	"cx-micro-flake/services/batch_management/handler/const_util"
	"cx-micro-flake/services/batch_management/handler/database"
	"encoding/json"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type Job struct {
	Id        int       `json:"id"`
	JobType   string    `json:"jobType"`
	Label     string    `json:"label"`
	Location  int       `json:"location"`
	Timestamp time.Time `json:"timestamp"`
}

func (v *Job) Serialised() ([]byte, error) {
	return json.Marshal(v)
}

func GetRawMaterialPrintJob(dbConnection *gorm.DB, jobId int) (datatypes.JSON, error) {
	err, c := database.Get(dbConnection, const_util.BatchManagementRawMaterialTable, jobId)
	if err == nil {
		rawMaterialBatchInfo := database.GeRawMaterialBatchInfo(c.ObjectInfo)
		job := Job{}
		job.Id = jobId
		job.JobType = "raw_material"
		job.Label = rawMaterialBatchInfo.LabelImage
		job.Location = rawMaterialBatchInfo.Location
		job.Timestamp = time.Now()
		return job.Serialised()
	} else {
		return datatypes.JSON{}, err
	}
}
