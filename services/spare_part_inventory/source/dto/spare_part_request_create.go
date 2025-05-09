package dto

import "cx-micro-flake/services/spare_part_inventory/source/models"

type SparePartList struct {
	Quantity    int `json:"quantity"`
	SparePartId int `json:"sparePartId"`
}
type CreateSparePartRequest struct {
	IsNeedSparePart bool                   `json:"isNeedSparePart"`
	SpareParts      []models.SparePartList `json:"spareParts"`
	JobId           string                 `json:"jobId"`
	MachineId       int                    `json:"machineId"`
}
