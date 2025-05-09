package dto

type CancelJobRequest struct {
	Remarks string `json:"remarks"` // base64 string due to unicode
}
