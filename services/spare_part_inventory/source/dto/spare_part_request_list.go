package dto

type SparePartRequestList struct {
	Id              int    `json:"id"`
	LocationName    string `json:"locationName"`
	OnHandQty       int    `json:"onHandQty"`
	SparePartNumber string `json:"sparePartNumber"`
}
