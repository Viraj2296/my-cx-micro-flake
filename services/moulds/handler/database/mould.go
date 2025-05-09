package database

import (
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"time"

	"gorm.io/datatypes"
)

type MouldOverview struct {
	SnapshotTime datatypes.Time `json:"snapshotTime"`
}

type MouldsRecordTrail struct {
	Id            int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	RecordId      int            `json:"recordId" gorm:"primary_key;not_null"`
	ComponentName string         `json:"componentName" gorm:"primary_key;not_null"`
	ObjectInfo    datatypes.JSON `json:"objectInfo"`
}

type MouldMaster struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

func (mm *MouldMaster) GetMouldMasterInfo() *MouldMasterInfo {
	mouldMasterInfo := MouldMasterInfo{}
	json.Unmarshal(mm.ObjectInfo, &mouldMasterInfo)
	return &mouldMasterInfo
}

type MouldMasterInfo struct {
	A1                  bool        `json:"a1"`
	A2                  bool        `json:"a2"`
	B1                  bool        `json:"b1"`
	B2                  bool        `json:"b2"`
	Pro                 bool        `json:"pro"`
	CH18                bool        `json:"cH18"`
	CH20                bool        `json:"cH20"`
	CH21                bool        `json:"cH21"`
	CH24                bool        `json:"cH24"`
	CH27                bool        `json:"cH27"`
	CH30                bool        `json:"cH30"`
	Site                string      `json:"site"`
	Type                string      `json:"type"`
	Plant               string      `json:"plant"`
	RevNo               bool        `json:"revNo"`
	Width               int         `json:"width"`
	Plate               bool        `json:"2Plate"`
	Plate1              bool        `json:"3Plate"`
	Height              int         `json:"height"`
	Length              int         `json:"length"`
	Others              bool        `json:"others"`
	PartNo              int         `json:"partNo"`
	Polish              bool        `json:"polish"`
	Spring              bool        `json:"spring"`
	ToolNo              string      `json:"toolNo"`
	Vendor              string      `json:"vendor"`
	FanGate             bool        `json:"fanGate"`
	NoOfCav             int         `json:"noOfCav"`
	Remarks             string      `json:"remarks"`
	Reverse             bool        `json:"reverse"`
	TabGate             bool        `json:"tabGate"`
	Category            int         `json:"category"`
	CavityNo            bool        `json:"cavityNo"`
	Customer            string      `json:"customer"`
	EdgeGate            bool        `json:"edgeGate"`
	Location            string      `json:"location"`
	MoonGate            bool        `json:"moonGate"`
	NoOfPart            int         `json:"noOfPart"`
	PinPoint            bool        `json:"pinPoint"`
	SubToPin            bool        `json:"subToPin"`
	ToolLife            int         `json:"toolLife"`
	MTHYRCode           bool        `json:"MTHYRCode"`
	CreatedAt           string      `json:"createdAt"`
	CreatedBy           int         `json:"createdBy"`
	HalfRound           interface{} `json:"halfRound"`
	OldToolId           string      `json:"oldToolId"`
	ReturnPin           bool        `json:"returnPin"`
	Trapzodal           bool        `json:"trapzodal"`
	ValveGate           bool        `json:"valveGate"`
	Controller          string      `json:"controller"`
	EdmTexture          bool        `json:"edmTexture"`
	ExportTool          bool        `json:"exportTool"`
	Interblock          bool        `json:"interblock"`
	MouldImage          string      `json:"mouldImage"`
	PinEjector          bool        `json:"pinEjector"`
	Production          bool        `json:"production"`
	RepeatTool          bool        `json:"repeatTool"`
	ResignCode          bool        `json:"resignCode"`
	DirectSprue         bool        `json:"directSprue"`
	EarlyReturn         bool        `json:"earlyReturn"`
	FixtureTool         bool        `json:"fixtureTool"`
	Maintenance         bool        `json:"maintenance"`
	MouldStatus         int         `json:"mouldStatus"`
	SubCategory         int         `json:"subCategory"`
	BladeEjector        bool        `json:"bladeEjector"`
	DirectHotTip        bool        `json:"directHotTip"`
	ObjectStatus        string      `json:"objectStatus"`
	SafetyDevice        bool        `json:"safetyDevice"`
	StripperRing        bool        `json:"stripperRing"`
	SupportPlate        bool        `json:"supportPlate"`
	CoolingOfCors       bool        `json:"coolingOfCors"`
	DoubleEjector       bool        `json:"doubleEjector"`
	EjectorSleeve       bool        `json:"ejectorSleeve"`
	HotSprueBrush       bool        `json:"hotSprueBrush"`
	LastUpdatedAt       string      `json:"lastUpdatedAt"`
	LastUpdatedBy       int         `json:"lastUpdatedBy"`
	RoundedRunner       bool        `json:"roundedRunner"`
	StripperPlate       bool        `json:"stripperPlate"`
	SubmarineGate       bool        `json:"submarineGate"`
	TextureFinish       bool        `json:"textureFinish"`
	DirectHotSprue      bool        `json:"directHotSprue"`
	InsertsMolding      bool        `json:"insertsMolding"`
	InsultingPlate      bool        `json:"insultingPlate"`
	StripperLedges      bool        `json:"stripperLedges"`
	CoolingOfSlides     bool        `json:"coolingOfSlides"`
	InReliefOnMould     bool        `json:"inReliefOnMould"`
	PartDescription     bool        `json:"partDescription"`
	RecessedInMould     bool        `json:"recessedInMould"`
	ToolMaintenance     string      `json:"toolMaintenance"`
	CommissionedDate    string      `json:"commissionedDate"`
	CoolingOfCavities   bool        `json:"coolingOfCavities"`
	MarkConnectionsBy   bool        `json:"markConnectionsBy"`
	CanCreateWorkOrder  bool        `json:"canCreateWorkOrder"`
	CanCustomerApprove  bool        `json:"canCustomerApprove"`
	CoolingOfCorePlate  bool        `json:"coolingOfCorePlate"`
	CavNoMarkedOnInsert bool        `json:"cavNoMarkedOnInsert"`
	StripperAtFixedHalf bool        `json:"stripperAtFixedHalf"`
	TemplateTableFields []struct {
		Id         int         `json:"id"`
		H13        int         `json:"h13"`
		P20        int         `json:"p20"`
		Ssoc       int         `json:"ssoc"`
		Brass      int         `json:"brass"`
		Nak80      int         `json:"nak80"`
		SStar      int         `json:"sStar"`
		Sdk61      int         `json:"sdk61"`
		DievAr     int         `json:"dievAr"`
		Hrc        int         `json:"4448hrc"`
		Hrc1       int         `json:"4852hrc"`
		Hrc2       int         `json:"5254hrc"`
		Assab760   interface{} `json:"assab760"`
		Category   string      `json:"category"`
		Assab8402  int         `json:"assab8402"`
		Hardening  int         `json:"hardening"`
		Nitriding  int         `json:"nitriding"`
		Assab718Hh int         `json:"assab718hh"`
	} `json:"templateTableFields"`
	CanSubmitTestRequest          bool          `json:"canSubmitTestRequest"`
	CoolingOfCavityPlate          bool          `json:"coolingOfCavityPlate"`
	EjectorChamberClosed          bool          `json:"ejectorChamberClosed"`
	EstimatedMouldWeight          int           `json:"estimatedMouldWeight"`
	LimitSwitchSafetyPin          bool          `json:"limitSwitchSafetyPin"`
	StripperAtMovingHalf          bool          `json:"stripperAtMovingHalf"`
	HydraulicallyOperated         bool          `json:"hydraulicallyOperated"`
	SpringOnCentralEjector        bool          `json:"springOnCentralEjector"`
	CoolingOfTopClampingPlate     bool          `json:"coolingOfTopClampingPlate"`
	MechanicallyOperatedTheBars   bool          `json:"mechanicallyOperatedTheBars"`
	CoolingOfBottomClampingPlate  bool          `json:"coolingOfBottomClampingPlate"`
	EjectorMechanismInFrontMould  bool          `json:"ejectorMechanismInFrontMould"`
	CavNoResignCodeMarkedOnRunner bool          `json:"cavNoResignCodeMarkedOnRunner"`
	Description                   string        `json:"description"`
	ShotCount                     int           `json:"shotCount"`
	MouldLifeNotification         []int         `json:"mouldLifeNotification"`
	IsNotificationSend            bool          `json:"isNotificationSend"`
	PartListArray                 []interface{} `json:"partListArray"`
	ProjectName                   string        `json:"projectName"`
	CycleTime                     string        `json:"cycleTime"`
	Pwt                           string        `json:"pwt"`
	Rwt                           string        `json:"rwt"`
	CanModify                     bool          `json:"canModify"`
	ModificationCount             int           `json:"modificationCount"`
}

func (v *MouldMasterInfo) DatabaseSerialize(userId int) map[string]interface{} {
	updatingData := make(map[string]interface{})
	v.LastUpdatedBy = userId
	v.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	updatingData["object_info"] = v.Serialize()
	return updatingData
}

func (v *MouldMasterInfo) Serialize() []byte {
	rawData, _ := json.Marshal(v)
	return rawData
}

type MouldComponent struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type PartMaster struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type MouldStatus struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

func (ms *MouldStatus) GetMouldStatusInfo() *MouldStatusInfo {
	mouldStatusInfo := MouldStatusInfo{}
	json.Unmarshal(ms.ObjectInfo, &mouldStatusInfo)
	return &mouldStatusInfo
}

type MouldTestStatus struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type MouldTestStatusInfo struct {
	Status               string `json:"status"`
	ColorCode            string `json:"colorCode"`
	Description          string `json:"description"`
	NotificationGroup    []int  `json:"notificationGroup"`
	NotificationTemplate string `json:"notificationTemplate"`
}

func (ms *MouldTestStatus) GetMouldTestStatusInfo() *MouldTestStatusInfo {
	mouldTestStatusInfo := MouldTestStatusInfo{}
	json.Unmarshal(ms.ObjectInfo, &mouldTestStatusInfo)
	return &mouldTestStatusInfo
}

type MouldStatusInfo struct {
	Status               string `json:"status"`
	ColorCode            string `json:"colorCode"`
	Preference           int    `json:"preference"`
	Description          string `json:"description"`
	NotificationGroup    []int  `json:"notificationGroup"`
	NotificationTemplate string `json:"notificationTemplate"`
}
type MouldCategory struct {
	Id         int            `json:"Id"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type MouldSubCategory struct {
	Id         int            `json:"Id"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type MouldTestRequest struct {
	Id         int            `json:"Id"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type MouldEmailTemplate struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type MouldEmailTemplateField struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type MouldStatusOverview struct {
	GroupByField string      `json:"groupByField"`
	DisplayField string      `json:"displayField"`
	Cards        interface{} `json:"cards"`
}

func (ms *MouldTestRequest) GetMouldTestRequestInfo() *MouldTestRequestInfo {
	mouldTestRequestInfo := MouldTestRequestInfo{}
	json.Unmarshal(ms.ObjectInfo, &mouldTestRequestInfo)
	return &mouldTestRequestInfo
}

type ActionRemarks struct {
	ExecutedTime  string `json:"executedTime"`
	Status        string `json:"status"`
	UserId        int    `json:"userId"`
	Remarks       string `json:"remarks"`
	ProcessedTime string `json:"processedTime"`
}
type MouldTestRequestInfo struct {
	MouldId                int             `json:"mouldId"`
	OffTool                int             `json:"offTool"`
	TestRequestReferenceId string          `json:"testRequestReferenceId"`
	MachineParamId         int             `json:"machineParamId"`
	TonFixed               int             `json:"tonFixed"`
	MachineId              int             `json:"machineId"`
	CostCentre             int             `json:"costCentre"`
	ShotRemarks            string          `json:"shotRemarks"`
	TonRangeMax            int             `json:"tonRangeMax"`
	TonRangeMin            int             `json:"tonRangeMin"`
	ObjectStatus           string          `json:"objectStatus"`
	MouldTestStatus        int             `json:"mouldTestStatus"`
	MouldStatusRemarks     string          `json:"mouldStatusRemarks"`
	OverallTestRemarks     string          `json:"overallTestRemarks"`
	RequestTestEndDate     string          `json:"requestTestEndDate"`
	RequestShotQuantity    string          `json:"requestShotQuantity"`
	RequestTestStartDate   string          `json:"requestTestStartDate"`
	TestedBy               int             `json:"testedBy"`
	ApprovedBy             int             `json:"approvedBy"`
	CreatedBy              int             `json:"createdBy"`
	CreatedAt              string          `json:"createdAt"`
	LastUpdatedBy          int             `json:"lastUpdatedBy"`
	LastUpdatedAt          string          `json:"lastUpdatedAt"`
	CanCheckIn             bool            `json:"canCheckIn"`
	CanContinueTest        bool            `json:"canContinueTest"`
	CanCheckOut            bool            `json:"canCheckOut"`
	CanApprove             bool            `json:"canApprove"`
	CanForceStop           bool            `json:"canForceStop"`
	IsAbortEnabled         bool            `json:"isAbortEnabled"`
	ActionStatus           string          `json:"actionStatus"`
	CanComplete            bool            `json:"canComplete"`
	CanView                bool            `json:"canView"`
	CanRelease             bool            `json:"canRelease"`
	CanUnRelease           bool            `json:"canUnRelease"`
	IsUpdate               bool            `json:"isUpdate"`
	Draggable              bool            `json:"draggable"`
	ActionRemarks          []ActionRemarks `json:"actionRemarks"`
	StrfNumber             string          `json:"strfNumber"`
}

func GetMouldTestRequestInfo(serialisedData datatypes.JSON) MouldTestRequestInfo {
	mouldTestRequestInfo := MouldTestRequestInfo{}
	err := json.Unmarshal(serialisedData, &mouldTestRequestInfo)
	if err != nil {
		return MouldTestRequestInfo{}
	}
	return mouldTestRequestInfo
}
func (v *MouldTestRequestInfo) Serialize() []byte {
	rawData, _ := json.Marshal(v)
	return rawData
}

func (v *MouldTestRequestInfo) DatabaseSerialize(userId int) map[string]interface{} {
	updatingData := make(map[string]interface{})
	v.LastUpdatedBy = userId
	v.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	updatingData["object_info"] = v.Serialize()
	return updatingData
}

type MouldManualShotCount struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type MouldMachineMaster struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type MouldSetting struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}
type MouldSettingInfo struct {
	MouldTestingGroup           int       `json:"mouldTestingGroup"`
	CreatedAt                   time.Time `json:"createdAt"`
	CreatedBy                   int       `json:"createdBy"`
	Description                 string    `json:"description"`
	ObjectStatus                string    `json:"objectStatus"`
	LastUpdatedAt               time.Time `json:"lastUpdatedAt"`
	LastUpdatedBy               int       `json:"lastUpdatedBy"`
	MouldLifeNotificationGroups []int     `json:"mouldLifeNotificationGroups"`
	LifeNotificationThreshold   int       `json:"lifeNotificationThreshold"`
}

func (v *MouldSettingInfo) Serialised() datatypes.JSON {
	serialisedData, _ := json.Marshal(v)
	return serialisedData
}

func GetMouldSettingInfo(serialisedData datatypes.JSON) *MouldSettingInfo {
	mouldSettingInfo := MouldSettingInfo{}
	err := json.Unmarshal(serialisedData, &mouldSettingInfo)
	if err != nil {
		return &MouldSettingInfo{}
	}
	return &mouldSettingInfo
}

type MouldManualShotCountInfo struct {
	MouldId             string `json:"mouldId"`
	CreatedAt           string `json:"createdAt"`
	CreatedBy           int    `json:"createdBy"`
	ShotCount           int    `json:"shotCount"`
	ObjectStatus        string `json:"objectStatus"`
	LastUpdatedAt       string `json:"lastUpdatedAt"`
	LastUpdatedBy       int    `json:"lastUpdatedBy"`
	ProductionOrder     string `json:"productionOrder"`
	ProductionOrderDate string `json:"productionOrderDate"`
}

func (ms *MouldManualShotCount) GetMouldManualShotCountInfo() *MouldManualShotCountInfo {
	mouldTestRequestInfo := MouldManualShotCountInfo{}
	json.Unmarshal(ms.ObjectInfo, &mouldTestRequestInfo)
	return &mouldTestRequestInfo
}

func GetMouldManualShotCountInfo(serialisedData datatypes.JSON) MouldManualShotCountInfo {
	mouldTestRequestInfo := MouldManualShotCountInfo{}
	err := json.Unmarshal(serialisedData, &mouldTestRequestInfo)
	if err != nil {
		return MouldManualShotCountInfo{}
	}
	return mouldTestRequestInfo
}

type MouldShotCountView struct {
	MouldId    int            `json:"mould_id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}
type MouldShotCountViewInfo struct {
	CreatedAt          string `json:"createdAt"`
	UpdatedAt          string `json:"updatedAt"`
	CreatedBy          int    `json:"createdBy"`
	UpdatedBy          int    `json:"updatedBy"`
	CurrentShotCount   int    `json:"currentShotCount"`
	IsNotificationSent bool   `json:"isNotificationSent"`
}

func (v *MouldShotCountViewInfo) Serialised() datatypes.JSON {
	serialisedData, _ := json.Marshal(v)
	return serialisedData
}

func GetMouldShoutCountViewInfo(serialisedData datatypes.JSON) *MouldShotCountViewInfo {
	mouldShotCountViewInfo := MouldShotCountViewInfo{}
	err := json.Unmarshal(serialisedData, &mouldShotCountViewInfo)
	if err != nil {
		return &MouldShotCountViewInfo{}
	}
	return &mouldShotCountViewInfo
}

type MouldModificationHistory struct {
	Id         int            `json:"Id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}
type MouldModificationHistoryInfo struct {
	CreatedAt          string `json:"createdAt"`
	UpdatedAt          string `json:"updatedAt"`
	CreatedBy          int    `json:"createdBy"`
	UpdatedBy          int    `json:"updatedBy"`
	ModifyCount        int    `json:"modifyCount"`
	MouldId            int    `json:"mouldId"`
	MouldTestRequestId int    `json:"mouldTestRequestId"`
}

func (v *MouldModificationHistoryInfo) Serialised() datatypes.JSON {
	serialisedData, _ := json.Marshal(v)
	return serialisedData
}

func GetMouldModificationHistoryInfo(serialisedData datatypes.JSON) *MouldModificationHistoryInfo {
	mouldModificationHistoryInfo := MouldModificationHistoryInfo{}
	err := json.Unmarshal(serialisedData, &mouldModificationHistoryInfo)
	if err != nil {
		return &MouldModificationHistoryInfo{}
	}
	return &mouldModificationHistoryInfo
}
