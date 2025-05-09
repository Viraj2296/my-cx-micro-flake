package common

type EmailSettingRequest struct {
	Host           string   `json:"host"`
	User           string   `json:"user"`
	Password       string   `json:"password"`
	CertFile       string   `json:"certFile"`
	KeyFile        string   `json:"keyFile"`
	SkipVerify     bool     `json:"skipVerify"`
	FromAddress    string   `json:"fromAddress"`
	ToAddress      string   `json:"toAddress"`
	FromName       string   `json:"fromName"`
	EhloIdentity   string   `json:"ehloIdentity"`
	StartTLSPolicy string   `json:"startTLSPolicy"`
	ContentTypes   []string `json:"contentTypes"`
	SampleContent  string   `json:"sampleContent"`
	InstanceName   string   `json:"instanceName"`
}
