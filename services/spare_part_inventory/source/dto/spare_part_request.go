package dto

import "cx-micro-flake/services/spare_part_inventory/source/models"

type SparePartRequest struct {
	EquipmentId string `json:"equipmentId"`
}
type PartList struct {
	SparePartId   int    `json:"sparePartId"`
	SparePartName string `json:"sparePartName"`
	Quantity      int    `json:"quantity"`
}

type SparePartRepairData struct {
	JobId           string     `json:"jobId"`
	RequestStatus   string     `json:"requestStatus"`
	IsNeedSparePart bool       `json:"isNeedSparePart"`
	Created_by      string     `json:"created_by"`
	Created_at      string     `json:"created_at"`
	SpareParts      []PartList `json:"spareParts"`
	MachineName     string     `json:"machineName"`
	Id              int        `json:"id"`
}
type ApprovePartList struct {
	SparePartList []models.SparePartList `json:"sparePartList"`
}
