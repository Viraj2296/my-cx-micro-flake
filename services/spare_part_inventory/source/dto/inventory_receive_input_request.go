package dto

type InventoryReceiveInputRequest struct {
	ResourceId int `json:"resourceId"`
	Quantity   int `json:"quantity"`
}
