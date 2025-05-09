package dto

type CheckOutRequest struct {
	FaultType int    `json:"faultType"`
	FaultCode int    `json:"faultCode"`
	Remarks   string `json:"remarks"` // base64 string due to unicode
}
