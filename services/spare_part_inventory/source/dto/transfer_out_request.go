package dto

type TransferRequest struct {
	SourceLocation      int    `json:"sourceLocation"`
	DestinationLocation int    `json:"destinationLocation"`
	Quantity            int    `json:"quantity"`
	ServiceNotification string `json:"serviceNotification"`
}
