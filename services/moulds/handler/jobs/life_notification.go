package jobs

import (
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/moulds/const_util"
	"cx-micro-flake/services/moulds/handler/database"
	"encoding/json"
	"time"

	"go.uber.org/zap"
)

func (v *JobService) SendMouldLifeNotification() {
	if v.PoolingInterval == 0 || v.PoolingInterval < 0 {
		v.PoolingInterval = 30
	}
	v.Logger.Info("machine stats pooling is starting up....", zap.Int("pooling_interval", v.PoolingInterval))
	var duration = time.Duration(v.PoolingInterval) * time.Second
	for {
		time.Sleep(duration)
		// mouldService := common.GetService("moulds_module").ServiceInterface.(common.MouldInterface)
		// err, generalObject := mouldService.GetMouldShotCountViewForNotification()
		condition := "object_info ->> '$.isNotificationSent' = 'false'"
		generalObject, _ := database.GetConditionalObjects(v.Database, const_util.MouldShoutCountViewTable, condition, 0)

		for _, mould := range *generalObject {
			shotCountViewObject := database.GetMouldShoutCountViewInfo(mould.ObjectInfo)

			_, mouldInfoObject := database.Get(v.Database, const_util.MouldMasterTable, mould.Id)
			mouldInfo := make(map[string]interface{})
			json.Unmarshal(mouldInfoObject.ObjectInfo, &mouldInfo)

			mouldToolLife := util.InterfaceToInt(mouldInfo["toolLife"])

			_, mouldSettingObject := database.Get(v.Database, const_util.MouldSettingTable, 1)
			mouldSetting := make(map[string]interface{})
			json.Unmarshal(mouldSettingObject.ObjectInfo, &mouldSetting)

			lifeNotificationThreshold := util.InterfaceToInt(mouldSetting["lifeNotificationThreshold"])
			shotCountRatio := (shotCountViewObject.CurrentShotCount * 100) / mouldToolLife
			if shotCountRatio >= lifeNotificationThreshold {
				notificationUsers := util.InterfaceToIntArray(mouldSetting["mouldLifeNotificationGroups"])
				for _, notificationUserId := range notificationUsers {
					err := v.EmailHandler.EmailGenerator(v.Database, const_util.MouldShotCountTemplateType, notificationUserId, const_util.MouldMasterComponent, mould.Id)
					if err == nil {
						shotCountViewObject.IsNotificationSent = true
						var updateObject = make(map[string]interface{})
						updateObject["object_info"] = shotCountViewObject.Serialised()

						err := database.Update(v.Database, const_util.MouldShoutCountViewTable, mould.Id, updateObject)
						if err != nil {
							return
						}

					}
				}
			}

		}
	}
}
