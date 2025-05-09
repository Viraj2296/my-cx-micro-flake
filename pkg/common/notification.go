package common

// AttachedFile struct represents email attached files.
type AttachedFile struct {
	Name    string
	Content []byte
}

// Message is representation of the email message.
type Message struct {
	To            []string
	SingleEmail   bool
	From          string
	Subject       string
	Body          map[string]string
	Info          string
	ReplyTo       []string
	EmbeddedFiles []string
	AttachedFiles []*AttachedFile
}

type SystemNotification struct {
	Name               string `json:"name"`
	IconCls            string `json:"iconCls"`
	RecordId           int    `json:"recordId"`
	ColorCode          string `json:"colorCode"`
	Component          string `json:"component"`
	GeneratedTime      string `json:"generatedTime"`
	Description        string `json:"description"`
	RouteLinkComponent string `json:"routeLinkComponent"`
	TargetUsers        []int  `json:"targetUsers"`
}

type PushNotificationMessage struct {
	IncludePlayerIDs []string               `json:"IncludePlayerIDs"`
	Headings         map[string]string      `json:"headings"`
	Contents         map[string]string      `json:"contents"`
	Data             map[string]interface{} `json:"data,omitempty"`
	DeliveryStatus   string                 `json:"deliveryStatus"`
	RetryCount       int                    `json:"retryCount"`
	ReferenceData    map[string]interface{} `json:"referenceData,omitempty"`
}
