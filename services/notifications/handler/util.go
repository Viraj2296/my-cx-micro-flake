package handler

import "cx-micro-flake/pkg/common"

func IsNotificationIdExist(metaInfos []common.NotificationMetaInfo, id int) bool {
	for _, info := range metaInfos {
		if info.Id == id {
			return true
		}
	}
	return false
}
