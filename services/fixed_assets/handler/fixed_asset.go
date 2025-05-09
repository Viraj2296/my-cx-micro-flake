package handler

import (
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"time"

	"gorm.io/datatypes"
)

type ActionRemarks struct {
	ExecutedTime  string `json:"executedTime"`
	Status        string `json:"status"`
	UserId        int    `json:"userId"`
	Remarks       string `json:"remarks"`
	ProcessedTime string `json:"processedTime"`
}

type FixedAssetRecordTrail struct {
	Id            int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	RecordId      int            `json:"recordId" gorm:"primary_key;not_null"`
	ComponentName string         `json:"componentName" gorm:"primary_key;not_null"`
	ObjectInfo    datatypes.JSON `json:"objectInfo"`
}

type FixedAssetsComponent struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type FixedAssetMaster struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type FixedAssetContract struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type FixedAssetClass struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type FixedAssetStatus struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type FixedAssetContractStatus struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type FixedAssetSetupMaster struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type FixedAssetCategory struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type ITAssetCategory struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type FixedAssetSubCategory struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type FixedAssetDisposal struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type FixedAssetTransfer struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type FixedAssetTransferStatus struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type FixedAssetDisposalStatus struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type FixedAssetLocation struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type FixedAssetDynamicField struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type FixedAssetDynamicFieldConfiguration struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type FixedAssetEmailTemplate struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type FixedAssetEmailTemplateField struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

func (v *FixedAssetDynamicFieldConfiguration) getFixedAssetDynamicFieldInfo() *FixedAssetDynamicFieldInfo {
	fixedAssetDynamicFieldInfo := FixedAssetDynamicFieldInfo{}
	json.Unmarshal(v.ObjectInfo, &fixedAssetDynamicFieldInfo)
	return &fixedAssetDynamicFieldInfo
}

type DynamicFields struct {
	Property   string `json:"property"`
	Type       string `json:"type"`
	GridSystem string `json:"gridSystem"`
	Label      string `json:"label"`
}

type FixedAssetDynamicFieldInfo struct {
	ComponentName string                   `json:"componentName"`
	Id            int                      `json:"id"`
	DynamicFields []component.RecordSchema `json:"dynamicFields"`
}

type FixedAssetTransferRequest struct {
	AssetId             int             `json:"assetId"`
	Name                string          `json:"name"`
	Description         string          `json:"description"`
	DetailedDescription string          `json:"detailedDescription"`
	ActionRemarks       []ActionRemarks `json:"actionRemarks"`
	WorkflowLevel       int             `json:"workflowLevel"`
	TransferLocation    int             `json:"transferLocation"`
	TransferStatus      int             `json:"disposalStatus"`
	ActionStatus        string          `json:"actionStatus"`
	CanSubmit           bool            `json:"canSubmit"`
	CanHODApprove       bool            `json:"canHODApprove"`
	CanHODReject        bool            `json:"canHODReject"`
	CanCEOApprove       bool            `json:"canCEOApprove"`
	CanCEOReject        bool            `json:"canCEOReject"`
	Labels              []string        `json:"labels"`
	CreatedAt           string          `json:"createdAt"`
	CreatedBy           int             `json:"createdBy"`
	LastUpdatedAt       string          `json:"lastUpdatedAt"`
	LastUpdatedBy       int             `json:"lastUpdatedBy"`
	ObjectStatus        string          `json:"objectStatus"`
}

func (v *FixedAssetTransfer) getTransferRequest() *FixedAssetTransferRequest {
	assetDisposalRequest := FixedAssetTransferRequest{}
	json.Unmarshal(v.ObjectInfo, &assetDisposalRequest)
	return &assetDisposalRequest
}

func (v *FixedAssetTransferRequest) Serialize() []byte {
	serialisedObject, _ := json.Marshal(v)
	return serialisedObject
}

func (v *FixedAssetTransferRequest) DatabaseSerialize(userId int) map[string]interface{} {
	updatingData := make(map[string]interface{})
	v.LastUpdatedBy = userId
	v.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	updatingData["object_info"] = v.Serialize()
	return updatingData
}

type FixedAssetDisposalRequest struct {
	AssetId             int             `json:"assetId"`
	Name                string          `json:"name"`
	Description         string          `json:"description"`
	DetailedDescription string          `json:"detailedDescription"`
	ActionRemarks       []ActionRemarks `json:"actionRemarks"`
	WorkflowLevel       int             `json:"workflowLevel"`
	DisposalStatus      int             `json:"disposalStatus"`
	ActionStatus        string          `json:"actionStatus"`
	CanSubmit           bool            `json:"canSubmit"`
	CanHODApprove       bool            `json:"canHODApprove"`
	CanHODReject        bool            `json:"canHODReject"`
	CanCEOApprove       bool            `json:"canCEOApprove"`
	CanCEOReject        bool            `json:"canCEOReject"`
	Labels              []string        `json:"labels"`
	CreatedAt           string          `json:"createdAt"`
	CreatedBy           int             `json:"createdBy"`
	LastUpdatedAt       string          `json:"lastUpdatedAt"`
	LastUpdatedBy       int             `json:"lastUpdatedBy"`
	ObjectStatus        string          `json:"objectStatus"`
	CanUserAcknowledge  bool            `json:"canUserAcknowledge"`
}

func (v *FixedAssetDisposal) getDisposalRequest() *FixedAssetDisposalRequest {
	assetDisposalRequest := FixedAssetDisposalRequest{}
	json.Unmarshal(v.ObjectInfo, &assetDisposalRequest)
	return &assetDisposalRequest
}

func (v *FixedAssetDisposalRequest) Serialize() []byte {
	serialisedObject, _ := json.Marshal(v)
	return serialisedObject
}

func (v *FixedAssetDisposalRequest) DatabaseSerialize(userId int) map[string]interface{} {
	updatingData := make(map[string]interface{})
	v.LastUpdatedBy = userId
	v.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	updatingData["object_info"] = v.Serialize()
	return updatingData
}

type FixedAssetMasterInfo struct {
	Cost                   int                `json:"cost"`
	Name                   string             `json:"name"`
	Image                  string             `json:"image"`
	Model                  string             `json:"model"`
	Labels                 []string           `json:"labels"`
	Vendor                 string             `json:"vendor"`
	Category               int                `json:"category"`
	CreatedAt              time.Time          `json:"createdAt"`
	CreatedBy              int                `json:"createdBy"`
	ScrapDate              string             `json:"scrapDate"`
	Technician             int                `json:"technician"`
	AssetNumber            string             `json:"assetNumber"`
	AssetStatus            int                `json:"assetStatus"`
	Description            string             `json:"description"`
	AssignedDate           string             `json:"assignedDate"`
	SerialNumber           string             `json:"serialNumber"`
	LastUpdatedAt          time.Time          `json:"lastUpdatedAt"`
	LastUpdatedBy          int                `json:"lastUpdatedBy"`
	MaintenanceTeam        []int              `json:"maintenanceTeam"`
	UsedInLocation         string             `json:"usedInLocation"`
	VendorReference        string             `json:"vendorReference"`
	WarrantyExpirationDate string             `json:"warrantyExpirationDate"`
	IsWarrantyEnabled      bool               `json:"isWarrantyEnabled"`
	QRCodeUrl              string             `json:"qrCodeUrl"`
	AssetConfiguration     AssetConfiguration `json:"assetConfiguration"`
	ObjectStatus           string             `json:"objectStatus"`
	MaintenaceTeam         []int              `json:"maintenaceTeam"`
	DetailedDescription    string             `json:"detailedDescription"`
	HasAttachment          bool               `json:"hasAttachment"`
	Attachment             string             `json:"attachment"`
	PurchaseDate           string             `json:"purchaseDate"`
	SapAssetId             string             `json:"sapAssetId"`
	VendorEmail            string             `json:"vendorEmail"`
}

type AssetConfiguration struct {
	OsTypes        int    `json:"osTypes"`
	UserName       string `json:"userName"`
	OfficeVersions int    `json:"officeVersions"`
}

func (v *FixedAssetMaster) getFixedAssetMasterInfo() *FixedAssetMasterInfo {
	itServiceRequestInfo := FixedAssetMasterInfo{}
	json.Unmarshal(v.ObjectInfo, &itServiceRequestInfo)
	return &itServiceRequestInfo
}

type FileUploadResponse struct {
	Url              string        `json:"url"`
	Name             string        `json:"name"`
	Size             string        `json:"size"`
	IsFile           bool          `json:"isFile"`
	ChainReference   string        `json:"chainReference"`
	Path             string        `json:"path"`
	MimeType         string        `json:"mimeType"`
	CreatedAt        time.Time     `json:"createdAt"`
	CreatedBy        int           `json:"createdBy"`
	LastUpdatedAt    time.Time     `json:"lastUpdatedAt"`
	LastUpdatedBy    int           `json:"lastUpdatedBy"`
	ObjectStatus     string        `json:"objectStatus"`
	IsShared         bool          `json:"isShared"`
	Share            []interface{} `json:"share"`
	Tags             interface{}   `json:"tags"`
	FavoriteList     interface{}   `json:"favoriteList"`
	FileTypeIcon     string        `json:"fileTypeIcon"`
	FilePreviewImage string        `json:"filePreviewImage"`
	Description      string        `json:"description"`
}
