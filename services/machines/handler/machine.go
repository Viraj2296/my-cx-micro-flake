package handler

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"

	"gorm.io/datatypes"
)

type ActionRemarks struct {
	ExecutedTime  string `json:"executedTime"`
	Status        string `json:"status"`
	UserId        int    `json:"userId"`
	Remarks       string `json:"remarks"`
	ProcessedTime string `json:"processedTime"`
}

type Assignment struct {
	Event    int `json:"int"`
	Id       int `json:"id"`
	Resource int `json:"resource"`
}

type MachinesRecordTrail struct {
	Id            int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	RecordId      int            `json:"recordId" gorm:"primary_key;not_null"`
	ComponentName string         `json:"componentName" gorm:"primary_key;not_null"`
	ObjectInfo    datatypes.JSON `json:"objectInfo"`
}

type MachineOverview struct {
	SnapshotTime datatypes.Time `json:"snapshotTime"`
	//MachineSnapshot
}

type MachineFilter struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type MachineMaster struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type MachineConnectStatus struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type PartMaster struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type AssemblyLineType struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}
type AssemblyEquipmentName struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type AssemblyEquipmentTypeMaster struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

func (mm *MachineMaster) getMachineMasterInfo() *MachineMasterInfo {
	machineMasterInfo := MachineMasterInfo{}
	json.Unmarshal(mm.ObjectInfo, &machineMasterInfo)
	return &machineMasterInfo
}

func (mm *AssemblyMachineMaster) getAssemblyMachineMasterInfo() *AssemblyMachineMasterInfo {
	machineMasterInfo := AssemblyMachineMasterInfo{}
	json.Unmarshal(mm.ObjectInfo, &machineMasterInfo)
	return &machineMasterInfo
}

func (mm *ToolingMachineMaster) getToolingMachineMasterInfo() *ToolingMachineMasterInfo {
	machineMasterInfo := ToolingMachineMasterInfo{}
	json.Unmarshal(mm.ObjectInfo, &machineMasterInfo)
	return &machineMasterInfo
}

type DisplaySettingInfo struct {
	LastUpdatedAt   string `json:"lastUpdatedAt"`
	LastUpdatedBy   int    `json:"lastUpdatedBy"`
	DisplayEnabled  bool   `json:"displayEnabled"`
	DisplayInterval int    `json:"displayInterval"`
	CreatedBy       int    `json:"createdBy"`
	CreatedAt       string `json:"createdAt"`
	ObjectStatus    string `json:"objectStatus"`
}
type MachineMasterInfo struct {
	Site                         string      `json:"site"`
	Brand                        int         `json:"brand"`
	Model                        string      `json:"model"`
	Plant                        string      `json:"plant"`
	Tonnage                      string      `json:"tonnage"`
	Category                     int         `json:"category"`
	Location                     string      `json:"location"`
	Supplier                     string      `json:"supplier"`
	Department                   int         `json:"department"`
	SubCategory                  int         `json:"subCategory"`
	MachineImage                 string      `json:"machineImage"`
	NewMachineId                 string      `json:"newMachineId"`
	OldMachineId                 string      `json:"oldMachineId"`
	SerialNumber                 string      `json:"serialNumber"`
	MachineStatus                int         `json:"machineStatus"`
	CommissionedDate             interface{} `json:"commissionedDate"`
	FixedAssetNumber             string      `json:"fixedAssetNumber"`
	MachineDescription           string      `json:"machineDescription"`
	MachineConnectStatus         int         `json:"machineConnectStatus"`
	ObjectStatus                 string      `json:"objectStatus"`
	CreatedBy                    int         `json:"createdBy"`
	CreatedAt                    string      `json:"createdAt"`
	LastUpdatedBy                int         `json:"lastUpdatedBy"`
	LastUpdatedAt                string      `json:"lastUpdatedAt"`
	LastUpdatedMachineLiveStatus string      `json:"lastUpdatedMachineLiveStatus"`
	CanCreateWorkOrder           bool        `json:"canCreateWorkOrder"`
	CanCreateCorrectiveWorkOrder bool        `json:"canCreateCorrectiveWorkOrder"`
	DelayStatus                  string      `json:"delayStatus"`
	DelayPeriod                  int64       `json:"delayPeriod"`
	CurrentCycleCount            int         `json:"currentCycleCount"`
	EnableProductionOrder        bool        `json:"enableProductionOrder"`
}
type MachineParameter struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type AssemblyLinesInfo struct {
	Name                   string    `json:"name"`
	TvUsers                []int     `json:"tvUsers"`
	CreatedAt              time.Time `json:"createdAt"`
	CreatedBy              int       `json:"createdBy"`
	Description            string    `json:"description"`
	ObjectStatus           string    `json:"objectStatus"`
	LastUpdatedAt          time.Time `json:"lastUpdatedAt"`
	LastUpdatedBy          int       `json:"lastUpdatedBy"`
	DefaultPlannedManpower int       `json:"defaultPlannedManpower"`
}
type AssemblyMachineLines struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

func getAssemblyLinesInfo(objectInfo datatypes.JSON) *AssemblyLinesInfo {
	assemblyLinesInfo := AssemblyLinesInfo{}
	err := json.Unmarshal(objectInfo, &assemblyLinesInfo)
	if err != nil {
		return &AssemblyLinesInfo{}
	}
	return &assemblyLinesInfo
}
func (v *MachineParameter) getMachineParamInfo() *MachineParamInfo {
	machineParamInfo := MachineParamInfo{}
	json.Unmarshal(v.ObjectInfo, &machineParamInfo)
	return &machineParamInfo
}

type MachineParamInfo struct {
	TestId                   int         `json:"testId"`
	CreatedAt                string      `json:"createdAt"`
	CreatedBy                int         `json:"createdBy"`
	LastUpdatedAt            string      `json:"lastUpdatedAt"`
	LastUpdatedBy            int         `json:"lastUpdatedBy"`
	ObjectStatus             string      `json:"objectStatus"`
	PartNo                   int         `json:"partNo"`
	RevNo                    float64     `json:"revNo"`
	McNo                     string      `json:"mcNo"`
	McTonnage                string      `json:"mcTonnage"`
	FeedstockNo              string      `json:"feedstockNo"`
	ProgramNo                string      `json:"programNo"`
	DateInjected             string      `json:"dateInjected"`
	ScrewDiameter            int         `json:"screwDiameter"`
	GreenPartWeightRange     string      `json:"greenPartWeightRange"`
	MachineParamReferenceId  string      `json:"machineParamReferenceId"`
	ForceSpringLoadedMold    interface{} `json:"forceSpringLoadedMold"`
	StrokeSpringLoadedMold   interface{} `json:"strokeSpringLoadedMold"`
	VelocitySpringLoadedMold interface{} `json:"velocitySpringLoadedMold"`
	CanEdit                  bool        `json:"canEdit"`
	// Injection

	Injection1StInjectionPressure float64 `json:"injection1stInjectionPressure"`
	InjectionPosition1            float64 `json:"injectionPosition1"`
	InjectionPosition2            float64 `json:"injectionPosition2"`
	InjectionPosition3            float64 `json:"injectionPosition3"`
	InjectionPosition4            float64 `json:"injectionPosition4"`
	InjectionPosition5            float64 `json:"injectionPosition5"`
	InjectionSpeed1               float64 `json:"injectionSpeed1"`
	InjectionSpeed2               float64 `json:"injectionSpeed2"`
	InjectionSpeed3               float64 `json:"injectionSpeed3"`
	InjectionSpeed4               float64 `json:"injectionSpeed4"`
	InjectionSpeed5               float64 `json:"injectionSpeed5"`

	// Holding Pressure
	HoldingPressure1     float64 `json:"holdingPressure1"`
	HoldingPressure2     float64 `json:"holdingPressure2"`
	HoldingPressure3     float64 `json:"holdingPressure3"`
	HoldingPressure4     float64 `json:"holdingPressure4"`
	HoldingPressure5     float64 `json:"holdingPressure5"`
	HoldingPressureTime1 float64 `json:"holdingPressureTime1"`
	HoldingPressureTime2 float64 `json:"holdingPressureTime2"`
	HoldingPressureTime3 float64 `json:"holdingPressureTime3"`
	HoldingPressureTime4 float64 `json:"holdingPressureTime4"`
	HoldingPressureTime5 float64 `json:"holdingPressureTime5"`

	ClampingUnitOpeningV1 int         `json:"clampingUnitOpeningV1"`
	ClampingUnitOpeningV2 int         `json:"clampingUnitOpeningV2"`
	ClampingUnitOpeningV3 float64     `json:"clampingUnitOpeningV3"`
	ClampingUnitOpeningV4 float64     `json:"clampingUnitOpeningV4"`
	ClampingUnitOpeningV5 int         `json:"clampingUnitOpeningV5"`
	ClampingUnitOpeningS1 float64     `json:"clampingUnitOpeningS1"`
	ClampingUnitOpeningS2 int         `json:"clampingUnitOpeningS2"`
	ClampingUnitOpeningS3 float64     `json:"clampingUnitOpeningS3"`
	ClampingUnitOpeningS4 float64     `json:"clampingUnitOpeningS4"`
	ClampingUnitOpeningS5 int         `json:"clampingUnitOpeningS5"`
	ClampingUnitOpeningF1 int         `json:"clampingUnitOpeningF1"`
	ClampingUnitOpeningF2 int         `json:"clampingUnitOpeningF2"`
	ClampingUnitOpeningF3 float64     `json:"clampingUnitOpeningF3"`
	ClampingUnitOpeningF4 interface{} `json:"clampingUnitOpeningF4"`
	ClampingUnitOpeningF5 int         `json:"clampingUnitOpeningF5"`

	ClampingUnitClosingV1 int         `json:"clampingUnitClosingV1"`
	ClampingUnitClosingV2 interface{} `json:"clampingUnitClosingV2"`
	ClampingUnitClosingV3 interface{} `json:"clampingUnitClosingV3"`
	ClampingUnitClosingV4 float64     `json:"clampingUnitClosingV4"`
	ClampingUnitClosingV5 int         `json:"clampingUnitClosingV5"`
	ClampingUnitClosingS1 int         `json:"clampingUnitClosingS1"`
	ClampingUnitClosingS2 interface{} `json:"clampingUnitClosingS2"`
	ClampingUnitClosingS3 interface{} `json:"clampingUnitClosingS3"`
	ClampingUnitClosingS4 float64     `json:"clampingUnitClosingS4"`
	ClampingUnitClosingS5 float64     `json:"clampingUnitClosingS5"`
	ClampingUnitClosingF1 int         `json:"clampingUnitClosingF1"`
	ClampingUnitClosingF2 interface{} `json:"clampingUnitClosingF2"`
	ClampingUnitClosingF3 interface{} `json:"clampingUnitClosingF3"`
	ClampingUnitClosingF4 float64     `json:"clampingUnitClosingF4"`
	ClampingUnitClosingF5 int         `json:"clampingUnitClosingF5"`

	//Barrel Temperature
	BarrelTemperatureFeedZone        float64 `json:"barrelTemperatureFeedZone"`
	BarrelTemperatureCompressionZone float64 `json:"barrelTemperatureCompressionZone"`
	BarrelTemperatureMeteringZone    float64 `json:"barrelTemperatureMeteringZone"`
	BarrelTemperatureFrontZone       float64 `json:"barrelTemperatureFrontZone"`
	BarrelTemperatureNozzle          float64 `json:"barrelTemperatureNozzle"`
	BarrelTemperatureAdditionalZone1 float64 `json:"barrelTemperatureAdditionalZone1"`
	BarrelTemperatureAdditionalZone2 float64 `json:"barrelTemperatureAdditionalZone2"`

	//Ejector
	EjectorDeMoldingForwardTime      float64 `json:"ejectorDeMoldingForwardTime"`
	EjectorDeMoldingForwardDistance  float64 `json:"ejectorDeMoldingForwardDistance"`
	EjectorDeMoldingForwardForce     float64 `json:"ejectorDeMoldingForwardForce"`
	EjectorDeMoldingBackwardTime     float64 `json:"ejectorDeMoldingBackwardTime"`
	EjectorDeMoldingBackwardDistance float64 `json:"ejectorDeMoldingBackwardDistance"`
	EjectorDeMoldingBackwardForce    float64 `json:"ejectorDeMoldingBackwardForce"`

	MeteringDec1Q float64     `json:"meteringDec1Q"`
	MeteringDec1V float64     `json:"meteringDec1V"`
	MeteringDec2Q float64     `json:"meteringDec2Q"`
	MeteringDec2V float64     `json:"meteringDec2V"`
	MeteringV1Ccm interface{} `json:"meteringV1Ccm"`
	MeteringV1Rpm int         `json:"meteringV1Rpm"`
	MeteringV2Ccm float64     `json:"meteringV2Ccm"`
	MeteringV2Rpm int         `json:"meteringV2Rpm"`
	MeteringP1    int         `json:"meteringP1"`
	MeteringP2    int         `json:"meteringP2"`

	MultipleEject        int         `json:"multipleEject"`
	VibratingStroke      interface{} `json:"vibratingStroke"` // Type retained from `MachineParamInfo`
	EjectorBackwardDelay float64     `json:"ejectorBackwardDelay"`
	EjectorForwardDelay  float64     `json:"ejectorForwardDelay"`

	StartTimeSpecific       float64 `json:"startTimeSpecific"`
	StartVolumeSpecific     float64 `json:"startVolumeSpecific"`
	MouldTempCavitySide     float64 `json:"mouldTempCavitySide"`
	MouldTempCoreSide       float64 `json:"mouldTempCoreSide"`
	MouldProtectionStrokes  int     `json:"mouldProtectionStrokes"`
	MouldProtectionTime     float64 `json:"mouldProtectionTime"`
	MouldSafetyDeviceTravel int     `json:"mouldSafetyDeviceTravel"`

	//Mold Temperature
	MoldTemperatureCoolingMedium   int     `json:"moldTemperatureCoolingMedium"` // Heater,Chiller,Water
	MoldTemperatureFixedHalf       float64 `json:"moldTemperatureFixedHalf"`
	MoldTemperatureMovingHalf      float64 `json:"moldTemperatureMovingHalf"`
	MoldTemperatureSliderTopBottom float64 `json:"moldTemperatureSliderTopBottom"`
	MoldTemperatureSliderSide      float64 `json:"moldTemperatureSliderSide"`

	// InjectionMaster
	FillingTime           string  `json:"fillingTime"`
	PeakInjPress          float64 `json:"peakInjPress"`
	MinCushion            float64 `json:"minCushion"`
	PlasticizingTime      float64 `json:"plasticizingTime"`
	TotalWeightShotWeight float64 `json:"totalWeightShotWeight"`
	RunnerWeight          float64 `json:"runnerWeight"`
	TransferPosition      float64 `json:"transferPosition"`
	Cushion               float64 `json:"cushion"`
	ScrewDecompression    float64 `json:"screwDecompression"`
	FillTime              float64 `json:"fillTime"`
	CycleTime             float64 `json:"cycleTime"`
	CoolingTime           float64 `json:"coolingTime"`
	ClampingForce         float64 `json:"clampingForce"`
	ScrewSpeed            float64 `json:"screwSpeed"`
	RecoveryTime          float64 `json:"recoveryTime"`
	BackPressure          float64 `json:"backPressure"`

	//Hot Runner
	NoOfTips       int     `json:"noOfTips"`
	Temperature1   float64 `json:"temperature1"`
	Temperature2   float64 `json:"temperature2"`
	Temperature3   float64 `json:"temperature3"`
	Temperature4   float64 `json:"temperature4"`
	Temperature5   float64 `json:"temperature5"`
	Temperature6   float64 `json:"temperature6"`
	Temperature7   float64 `json:"temperature7"`
	Temperature8   float64 `json:"temperature8"`
	Temperature9   float64 `json:"temperature9"`
	Temperature10  float64 `json:"temperature10"`
	Temperature11  float64 `json:"temperature11"`
	Temperature12  float64 `json:"temperature12"`
	Temperature13  float64 `json:"temperature13"`
	Temperature14  float64 `json:"temperature14"`
	Temperature15  float64 `json:"temperature15"`
	Temperature16  float64 `json:"temperature16"`
	Temperature17  float64 `json:"temperature17"`
	Temperature18  float64 `json:"temperature18"`
	Temperature19  float64 `json:"temperature19"`
	Temperature20  float64 `json:"temperature20"`
	Temperature21  float64 `json:"temperature21"`
	Temperature22  float64 `json:"temperature22"`
	Temperature23  float64 `json:"temperature23"`
	Temperature24  float64 `json:"temperature24"`
	Temperature25  float64 `json:"temperature25"`
	Temperature26  float64 `json:"temperature26"`
	Temperature27  float64 `json:"temperature27"`
	Temperature28  float64 `json:"temperature28"`
	Temperature29  float64 `json:"temperature29"`
	Temperature30  float64 `json:"temperature30"`
	Temperature31  float64 `json:"temperature31"`
	Temperature32  float64 `json:"temperature32"`
	Temperature33  float64 `json:"temperature33"`
	Temperature34  float64 `json:"temperature34"`
	Temperature35  float64 `json:"temperature35"`
	Temperature36  float64 `json:"temperature36"`
	Temperature37  float64 `json:"temperature37"`
	Temperature38  float64 `json:"temperature38"`
	Temperature39  float64 `json:"temperature39"`
	Temperature40  float64 `json:"temperature40"`
	Temperature41  float64 `json:"temperature41"`
	Temperature42  float64 `json:"temperature42"`
	Temperature43  float64 `json:"temperature43"`
	Temperature44  float64 `json:"temperature44"`
	Temperature45  float64 `json:"temperature45"`
	Temperature46  float64 `json:"temperature46"`
	Temperature47  float64 `json:"temperature47"`
	Temperature48  float64 `json:"temperature48"`
	Temperature49  float64 `json:"temperature49"`
	Temperature50  float64 `json:"temperature50"`
	Temperature51  float64 `json:"temperature51"`
	Temperature52  float64 `json:"temperature52"`
	Temperature53  float64 `json:"temperature53"`
	Temperature54  float64 `json:"temperature54"`
	Temperature55  float64 `json:"temperature55"`
	Temperature56  float64 `json:"temperature56"`
	Temperature57  float64 `json:"temperature57"`
	Temperature58  float64 `json:"temperature58"`
	Temperature59  float64 `json:"temperature59"`
	Temperature60  float64 `json:"temperature60"`
	Temperature61  float64 `json:"temperature61"`
	Temperature62  float64 `json:"temperature62"`
	Temperature63  float64 `json:"temperature63"`
	Temperature64  float64 `json:"temperature64"`
	Temperature65  float64 `json:"temperature65"`
	Temperature66  float64 `json:"temperature66"`
	Temperature67  float64 `json:"temperature67"`
	Temperature68  float64 `json:"temperature68"`
	Temperature69  float64 `json:"temperature69"`
	Temperature70  float64 `json:"temperature70"`
	Temperature71  float64 `json:"temperature71"`
	Temperature72  float64 `json:"temperature72"`
	Temperature73  float64 `json:"temperature73"`
	Temperature74  float64 `json:"temperature74"`
	Temperature75  float64 `json:"temperature75"`
	Temperature76  float64 `json:"temperature76"`
	Temperature77  float64 `json:"temperature77"`
	Temperature78  float64 `json:"temperature78"`
	Temperature79  float64 `json:"temperature79"`
	Temperature80  float64 `json:"temperature80"`
	Temperature81  float64 `json:"temperature81"`
	Temperature82  float64 `json:"temperature82"`
	Temperature83  float64 `json:"temperature83"`
	Temperature84  float64 `json:"temperature84"`
	Temperature85  float64 `json:"temperature85"`
	Temperature86  float64 `json:"temperature86"`
	Temperature87  float64 `json:"temperature87"`
	Temperature88  float64 `json:"temperature88"`
	Temperature89  float64 `json:"temperature89"`
	Temperature90  float64 `json:"temperature90"`
	Temperature91  float64 `json:"temperature91"`
	Temperature92  float64 `json:"temperature92"`
	Temperature93  float64 `json:"temperature93"`
	Temperature94  float64 `json:"temperature94"`
	Temperature95  float64 `json:"temperature95"`
	Temperature96  float64 `json:"temperature96"`
	Temperature97  float64 `json:"temperature97"`
	Temperature98  float64 `json:"temperature98"`
	Temperature99  float64 `json:"temperature99"`
	Temperature100 float64 `json:"temperature100"`
	Temperature101 float64 `json:"temperature101"`
	Temperature102 float64 `json:"temperature102"`
	Temperature103 float64 `json:"temperature103"`
	Temperature104 float64 `json:"temperature104"`
	Temperature105 float64 `json:"temperature105"`
	Temperature106 float64 `json:"temperature106"`
	Temperature107 float64 `json:"temperature107"`
	Temperature108 float64 `json:"temperature108"`
	Temperature109 float64 `json:"temperature109"`
	Temperature110 float64 `json:"temperature110"`
	Temperature111 float64 `json:"temperature111"`
	Temperature112 float64 `json:"temperature112"`
	Temperature113 float64 `json:"temperature113"`
	Temperature114 float64 `json:"temperature114"`
	Temperature115 float64 `json:"temperature115"`
	Temperature116 float64 `json:"temperature116"`
	Temperature117 float64 `json:"temperature117"`
	Temperature118 float64 `json:"temperature118"`
	Temperature119 float64 `json:"temperature119"`
	Temperature120 float64 `json:"temperature120"`
}
type MachineStatus struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type MachineDisplaySetting struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type AssemblyMachineDisplaySetting struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type MachineModuleSetting struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

func (mhs *MachineDisplaySetting) getDisplaySettingInfoInfo() *DisplaySettingInfo {
	displaySettingInfo := DisplaySettingInfo{}
	json.Unmarshal(mhs.ObjectInfo, &displaySettingInfo)
	return &displaySettingInfo
}

type MachineModuleSettingInfo struct {
	HmiRefreshInterval             int       `json:"hmiRefreshInterval"`
	DefaultScheduleView            bool      `json:"defaultScheduleView"`
	SchedulerViewEndDate           time.Time `json:"schedulerViewEndDate"`
	DisplayRotateInterval          int       `json:"displayRotateInterval"`
	SchedulerViewStartDate         time.Time `json:"schedulerViewStartDate"`
	SchedulerDateSelectionInterval string    `json:"schedulerDateSelectionInterval"`
}

func (mhs *MachineModuleSetting) getMachineModuleSettingInfo() *MachineModuleSettingInfo {
	machineModuleSettingInfo := MachineModuleSettingInfo{}
	json.Unmarshal(mhs.ObjectInfo, &machineModuleSettingInfo)
	return &machineModuleSettingInfo
}

type AssemblyMachineModuleSettingInfo struct {
	HmiRefreshInterval             int       `json:"hmiRefreshInterval"`
	DefaultScheduleView            bool      `json:"defaultScheduleView"`
	SchedulerViewEndDate           time.Time `json:"schedulerViewEndDate"`
	DisplayRotateInterval          int       `json:"displayRotateInterval"`
	SchedulerViewStartDate         time.Time `json:"schedulerViewStartDate"`
	SchedulerDateSelectionInterval string    `json:"schedulerDateSelectionInterval"`
}

func (mhs *AssemblyMachineModuleSetting) getAssemblyMachineModuleSettingInfo() *AssemblyMachineModuleSettingInfo {
	machineModuleSettingInfo := AssemblyMachineModuleSettingInfo{}
	json.Unmarshal(mhs.ObjectInfo, &machineModuleSettingInfo)
	return &machineModuleSettingInfo
}

type MachineHMISetting struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

func (mhs *MachineHMISetting) getHMISettingInfo() *HmiSettingInfo {
	hmiSettingInfo := HmiSettingInfo{}
	json.Unmarshal(mhs.ObjectInfo, &hmiSettingInfo)
	return &hmiSettingInfo
}

type HmiSettingInfo struct {
	CreatedAt                      string   `json:"createdAt"`
	CreatedBy                      int      `json:"createdBy"`
	HmiOperators                   []int    `json:"hmiOperators"`
	ObjectStatus                   string   `json:"objectStatus"`
	LastUpdatedAt                  string   `json:"lastUpdatedAt"`
	LastUpdatedBy                  int      `json:"lastUpdatedBy"`
	HmiStopReasons                 []int    `json:"hmiStopReasons"`
	HmiAutoStopPeriod              string   `json:"hmiAutoStopPeriod"`
	MachineLiveDetectionInterval   string   `json:"machineLiveDetectionInterval"`
	WarningMessageGenerationPeriod string   `json:"warningMessageGenerationPeriod"`
	WarningTargetEmailId           string   `json:"warningTargetEmailId"`
	WarningCCTargetEmailId         []string `json:"warningCCTargetEmailId"`
	ResetOperators                 []int    `json:"resetOperators"`
}

type MachineEmailTemplate struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type GeneralObject struct {
	Id         int             `json:"id"`
	ObjectInfo json.RawMessage `json:"object_info"`
}

type MachineEmailMasterFields struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type MachineCategory struct {
	Id         int            `json:"id"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type CategoryInfo struct {
	Name          string    `json:"name"`
	Site          int       `json:"site"`
	CreatedAt     time.Time `json:"createdAt"`
	CreatedBy     int       `json:"createdBy"`
	Description   string    `json:"description"`
	ObjectStatus  string    `json:"objectStatus"`
	LastUpdatedAt time.Time `json:"lastUpdatedAt"`
	LastUpdatedBy int       `json:"lastUpdatedBy"`
}

func (mc *MachineCategory) getCategoryInfo() *CategoryInfo {
	categoryInfo := CategoryInfo{}
	json.Unmarshal(mc.ObjectInfo, &categoryInfo)
	return &categoryInfo
}

type MachineSubCategory struct {
	Id         int            `json:"id"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type SubCategoryInfo struct {
	Name          string    `json:"name"`
	Category      int       `json:"category"`
	CreatedAt     time.Time `json:"createdAt"`
	CreatedBy     int       `json:"createdBy"`
	Description   string    `json:"description"`
	LastUpdatedAt time.Time `json:"lastUpdatedAt"`
	LastUpdatedBy int       `json:"lastUpdatedBy"`
}

func (mc *MachineSubCategory) getSubCategoryInfo() *SubCategoryInfo {
	subCategoryInfo := SubCategoryInfo{}
	json.Unmarshal(mc.ObjectInfo, &subCategoryInfo)
	return &subCategoryInfo
}

type MachineStatistics struct {
	MachineId int            `json:"machineId"`
	TS        int64          `json:"ts"`
	StatsInfo datatypes.JSON `json:"statsInfo"`
}

type AssemblyMachineStatistics struct {
	MachineId int            `json:"machineId"`
	TS        int64          `json:"ts"`
	StatsInfo datatypes.JSON `json:"statsInfo"`
}

type ToolingMachineStatistics struct {
	MachineId int            `json:"machineId"`
	TS        int64          `json:"ts"`
	StatsInfo datatypes.JSON `json:"statsInfo"`
}

type MachineStatisticsInfo struct {
	EventId                    int      `json:"eventId"`
	ProductionOrderId          int      `json:"productionOrderId"`
	CurrentStatus              string   `json:"currentStatus"`
	PartId                     int      `json:"partId"`
	ScheduleStartTime          string   `json:"scheduleStartTime"`
	ScheduleEndTime            string   `json:"scheduleEndTime"`
	ActualStartTime            string   `json:"actualStartTime"`
	ActualEndTime              string   `json:"actualEndTime"`
	EstimatedEndTime           string   `json:"estimatedEndTime"`
	Oee                        int      `json:"oee"`
	Availability               int      `json:"availability"`
	Performance                int      `json:"performance"`
	Quality                    int      `json:"quality"`
	PlannedQuality             int      `json:"plannedQuality"`
	Completed                  int      `json:"completed"`
	Rejects                    int      `json:"rejects"`
	CompletedPercentage        float32  `json:"completedPercentage"`
	OverallCompletedPercentage float32  `json:"overallCompletedPercentage"`
	Daily                      int      `json:"daily"`
	Actual                     int      `json:"actual"`
	ExpectedProductQyt         int      `json:"expectedProductQyt"`
	OverallRejectedQty         int      `json:"overallRejectedQty"`
	ProgressPercentage         float64  `json:"progressPercentage"`
	DailyPlannedQty            int      `json:"dailyPlannedQty"`
	DownTime                   int      `json:"downTime"`
	Remark                     string   `json:"remark"`
	WarningMessage             []string `json:"warningMessage"`
}

type MachineHMI struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type MachineHMIInfo struct {
	CreatedAt        string  `json:"createdAt"`
	EventId          int     `json:"eventId"`
	Operator         int     `json:"operator"`
	ReasonId         int     `json:"reasonId"`
	MachineId        int     `json:"machineId"`
	Status           string  `json:"status"`
	SetupTime        *string `json:"setupTime"`
	RejectedQuantity *int    `json:"rejectedQuantity"`
	HMIStatus        string  `json:"hmiStatus"`
	Remark           string  `json:"remark"`
}

func (mhi *MachineHMI) getMachineHMIInfo() *MachineHMIInfo {
	machineHMIInfo := MachineHMIInfo{}
	json.Unmarshal(mhi.ObjectInfo, &machineHMIInfo)
	return &machineHMIInfo
}

type StopReasons struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type MachineScheduler struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type MachineComponent struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type MachineWidget struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type MachineDashboard struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type MachineTimelineEvent struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type TimelineEventInfoUpdateRequest struct {
	EndDate     string      `json:"endDate"`
	StartDate   string      `json:"startDate"`
	EventStatus string      `json:"eventStatus"`
	Assignment  *Assignment `json:"assignment"`
}

type MachineTimelineResource struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type MachineTimelineAssignment struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type MachineTimelineAssignmentInfo struct {
	Event        int    `json:"event"`
	Resource     int    `json:"resource"`
	ObjectStatus string `json:"objectStatus"`
}

type MachineHMIStopReason struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

func (mhsr *MachineHMIStopReason) getMachineHMIStopReasonInfo() *MachineHMIStopReasonInfo {
	machineHMIStopReasonInfo := MachineHMIStopReasonInfo{}
	json.Unmarshal(mhsr.ObjectInfo, &machineHMIStopReasonInfo)
	return &machineHMIStopReasonInfo
}

type MachineHMIStopReasonInfo struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	ObjectStatus string `json:"objectStatus"`
	IconCls      string `json:"iconCls"`
}

type TimelineAssignmentInfo struct {
	Event    int `json:"event"`
	Resource int `json:"resource"`
}

type MachineData struct {
	CycleCount              int `json:"cycle_count"`
	Auto                    int `json:"auto"`
	MaintenanceDoorAlarm    int `json:"maintenance_door_alarm"`
	SafetyDoorAlarm         int `json:"safety_door_alarm"`
	EstopAlarm              int `json:"estop_alarm"`
	TowerLightRed           int `json:"tower_light_red"`
	MachineConnectionStatus int `json:"machine_connection_status"`
	InTimestamp             int `json:"in_timestamp"`
}

type Message struct {
	Topic                      string         `json:"topic" gorm:"primary_key;not_null"`
	TS                         int64          `json:"ts" gorm:"primary_key;not_null"`
	Id                         int32          `json:"id" gorm:"primary_key;not_null"`
	DriverMessageGeneratedTime int64          `json:"driverMessageGeneratedTime"`
	Body                       datatypes.JSON `json:"body"`
}

type ProductionOrderMaster struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type RefreshExportRequest struct {
	SelectedId []string `json:"selectEdId"`
}

type AssemblyMachineMaster struct {
	Id         int            `json:"id"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type AssemblyMachineMasterInfo struct {
	NewMachineId                 string `json:"newMachineId"`
	MessageFlag                  string `json:"messageFlag"`
	Area                         string `json:"area"`
	Level                        string `json:"level"`
	HelpButtonStationNo          string `json:"helpButtonStationNo"`
	Department                   int    `json:"department"`
	CreatedBy                    int    `json:"createdBy"`
	CreatedAt                    string `json:"createdAt"`
	ObjectStatus                 string `json:"objectStatus"`
	MachineImage                 string `json:"machineImage"`
	Description                  string `json:"description"`
	MachineStatus                int    `json:"machineStatus"`
	MachineConnectStatus         int    `json:"machineConnectStatus"`
	LastUpdatedAt                string `json:"lastUpdatedAt"`
	DelayStatus                  string `json:"delayStatus"`
	DelayPeriod                  int64  `json:"delayPeriod"`
	CurrentCycleCount            int    `json:"currentCycleCount"`
	CanCreateWorkOrder           bool   `json:"canCreateWorkOrder"`
	LastUpdatedMachineLiveStatus string `json:"lastUpdatedMachineLiveStatus"`
	AssemblyLineOption           int    `json:"assemblyLineOption"`
	Model                        string `json:"model"` // This was added as part of labour management module
	IsMESDriverConfigured        bool   `json:"isMESDriverConfigured"`
	IsEnabled                    bool   `json:"isEnabled"`
	AutoRejectHistoryGeneration  bool   `json:"AutoRejectHistoryGeneration"`
	LineMappingMessageFlag       string `json:"lineMappingMessageFlag"`
	EquipmentId                  string `json:"equipmentId"`   // This was added to check-in the machine using the mobile using qr code
	EquipmentName                string `json:"equipmentName"` // This was added as part of machine downtime module
}

func GetAssemblyMachineLineInfo(serialisedData datatypes.JSON) *AssemblyMachineMasterInfo {
	assemblyMachineLineInfo := AssemblyMachineMasterInfo{}
	err := json.Unmarshal(serialisedData, &assemblyMachineLineInfo)
	if err != nil {
		return &AssemblyMachineMasterInfo{}
	}
	return &assemblyMachineLineInfo
}
func GetAssemblyMachineFromId(dbConnection *gorm.DB, machineId int) *AssemblyMachineMasterInfo {
	err, generalObject := Get(dbConnection, AssemblyMachineMasterTable, machineId)
	if err != nil {
		return nil
	} else {
		return GetAssemblyMachineLineInfo(generalObject.ObjectInfo)
	}
}

type AssemblyMachineHmi struct {
	Id         int            `json:"id"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type AssemblyMachineHmiSetting struct {
	Id         int            `json:"id"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type AssemblyMachineModuleSetting struct {
	Id         int            `json:"id"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type DepartmentDisplayMachines struct {
	DepartmentId   int    `json:"department"`
	Name           string `json:"name"`
	ListOfMachines []int  `json:"listOfMachines"`
}

type ToolingMachineMaster struct {
	Id         int            `json:"id"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type ToolingMachineHMI struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type ToolingMachineDisplaySetting struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type ToolingMachineHmiSetting struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type ToolingMachineMasterInfo struct {
	NewMachineId                 string `json:"newMachineId"`
	TopicFlag                    string `json:"topicFlag"`
	Department                   int    `json:"department"`
	CreatedBy                    int    `json:"createdBy"`
	CreatedAt                    string `json:"createdAt"`
	MachineStatus                int    `json:"machineStatus"`
	MachineConnectStatus         int    `json:"machineConnectStatus"`
	LastUpdatedAt                string `json:"lastUpdatedAt"`
	DelayStatus                  string `json:"delayStatus"`
	DelayPeriod                  int64  `json:"delayPeriod"`
	CurrentCycleCount            int    `json:"currentCycleCount"`
	LastUpdatedMachineLiveStatus string `json:"lastUpdatedMachineLiveStatus"`
	ObjectStatus                 string `json:"objectStatus"`
	MachineImage                 string `json:"machineImage"`
	Description                  string `json:"description"`
	Model                        string `json:"model"`
	CanCreateWorkOrder           bool   `json:"canCreateWorkOrder"`
	IsMachiningTimeAvailable     bool   `json:"isMachiningTimeAvailable"`
}
type MouldMachineBrand struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type MouldMachineSetting struct {
	Id         int            `json:"id" gorm:"primary_key;not_null"` //machine Id
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}
