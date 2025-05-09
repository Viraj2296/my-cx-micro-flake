package common

import (
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/config"
	"cx-micro-flake/pkg/config/source/file"
	"cx-micro-flake/pkg/orm"
	"fmt"
	"reflect"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AuthInterface interface {
	GetUserInfoById(userRecordId int) UserBasicInfo
	GetJobRoleName(jobRoleId int) string
	GetJobRoleHierarchy(jobRoleId int) int
	GetUserInfoByEmployeeId(employeeId string) UserBasicInfo
	GetUserType(userId int) string
	GetUserTimezone(userId int) string
	GetUserInfoFromGroupId(groupId int) []UserBasicInfo
	GetBotInfo() UserBasicInfo // this will return the user id, user name, and avatar url
	Authenticate(userName, password string) (string, error)
	AddNotificationIds(userId, notificationId int) error
	AddViewNotificationIds(userId, notificationId int) error
	AddPushNotificationIds(userId, notificationId int) error
	AddViewPushNotificationIds(userId, notificationId int) error
	GetNotificationList(userId int) []NotificationMetaInfo
	GetViewNotificationList(userId int) []NotificationMetaInfo
	GetPushNotificationList(userId int) []NotificationMetaInfo
	GetViewPushNotificationList(userId int) []NotificationMetaInfo
	GetUserList() []UserBasicInfo
	IsAllowed(userId int, projectId, moduleName, componentName string, action string, resourceId string, method string, path string) (bool, bool, string)
	IsEmailExist(email string) bool
	IsAnyDepartmentHeadExist(departmentId int) bool
	GetHeadOfDepartment(departmentId int) UserBasicInfo
	EmailToUserId(email string) int
	GetHeadOfDepartments(userId int) []UserBasicInfo
	GetHeadOfSections(userId int) []UserBasicInfo
	GetSectionUsers(userId int) ([]int, bool)
	GetDepartmentUsers(userId int) (bool, []UserBasicInfo)
	GetUsersByDepartment(departmentId int) []UserBasicInfo
	ConvertToUserTimezoneToISO(userId int, datetime string) string
	GetAllUserBasicInfo2Table() component.TableDataResponse
	GetBasicInfo2TableFromUserId(userId int) component.TableDataResponse
	GetBasicInfo2TableFromUserList(userList reflect.Value) component.TableDataResponse
	UpdateComponentResource(rawComponentResourceInfo datatypes.JSON) error
	GetAllUserBasicInfo2QueryResults() []datatypes.JSON
	GetComponents() []datatypes.JSON
	UpdateComponent(componentName, targetTable string, serialisedObject datatypes.JSON) error
	CreateComponent(serialisedObject datatypes.JSON) (int, error)
	GetAllDeviceTokenFromJobRole(jobRole int) []string
	GetUsersDeviceTokens(listOfUsers []int) []string
	GetUsersFromJobRoleId(jobRole int) []int
	GetUserOneSignalSubscriptionIds(userId int) []string
}

type LMSInterface interface {
	GetComponents() []datatypes.JSON
	UpdateComponent(componentName, targetTable string, serialisedObject datatypes.JSON) error
	CreateComponent(serialisedObject datatypes.JSON) (int, error)
	GetAuthorizedBudgetIds(userId int) []map[string]interface{}
}

type ContentInterface interface {
	GetComponents() []datatypes.JSON
	UpdateComponent(componentName, targetTable string, serialisedObject datatypes.JSON) error
	CreateComponent(serialisedObject datatypes.JSON) (int, error)
}

type WorkHourAssignmentInterface interface {
	GetComponents() []datatypes.JSON
	UpdateComponent(componentName, targetTable string, serialisedObject datatypes.JSON) error
	CreateComponent(serialisedObject datatypes.JSON) (int, error)
}

type ToolingInterface interface {
	GetComponents() []datatypes.JSON
	UpdateComponent(componentName, targetTable string, serialisedObject datatypes.JSON) error
	CreateComponent(serialisedObject datatypes.JSON) (int, error)
}

type BackupAndRecoveryInterface interface {
	GetComponents() []datatypes.JSON
	UpdateComponent(componentName, targetTable string, serialisedObject datatypes.JSON) error
	CreateComponent(serialisedObject datatypes.JSON) (int, error)
}

type SchedulingInterface interface {
	GetComponents() []datatypes.JSON
	UpdateComponent(componentName, targetTable string, serialisedObject datatypes.JSON) error
	CreateComponent(serialisedObject datatypes.JSON) (int, error)
}

type ManufacturingInterface interface {
	GetComponents() []datatypes.JSON
	UpdateComponent(componentName, targetTable string, serialisedObject datatypes.JSON) error
	CreateComponent(serialisedObject datatypes.JSON) (int, error)
	GetAssemblyPartInfo(projectId string, partId int) (error, component.GeneralObject)
}

type ITServiceInterface interface {
	GetComponents() []datatypes.JSON
	UpdateComponent(componentName, targetTable string, serialisedObject datatypes.JSON) error
	CreateComponent(serialisedObject datatypes.JSON) (int, error)
}

type IIoTServiceInterface interface {
	GetComponents() []datatypes.JSON
	UpdateComponent(componentName, targetTable string, serialisedObject datatypes.JSON) error
	CreateComponent(serialisedObject datatypes.JSON) (int, error)
}

type FacilityServiceInterface interface {
	GetComponents() []datatypes.JSON
	UpdateComponent(componentName, targetTable string, serialisedObject datatypes.JSON) error
	CreateComponent(serialisedObject datatypes.JSON) (int, error)
}

type TicketsServiceInterface interface {
	GetComponents() []datatypes.JSON
	UpdateComponent(componentName, targetTable string, serialisedObject datatypes.JSON) error
	CreateComponent(serialisedObject datatypes.JSON) (int, error)
}

type IncidentServiceInterface interface {
	GetComponents() []datatypes.JSON
	UpdateComponent(componentName, targetTable string, serialisedObject datatypes.JSON) error
	CreateComponent(serialisedObject datatypes.JSON) (int, error)
}

type FixedAssetsInterface interface {
	GetComponents() []datatypes.JSON
	UpdateComponent(componentName, targetTable string, serialisedObject datatypes.JSON) error
	CreateComponent(serialisedObject datatypes.JSON) (int, error)
}

type FactoryServiceInterface interface {
	GetDepartmentName(departmentId int) string
	GetSiteName(siteId int) string
	GetComponents() []datatypes.JSON
	UpdateComponent(componentName, targetTable string, serialisedObject datatypes.JSON) error
	CreateComponent(serialisedObject datatypes.JSON) (int, error)
	GetSections(departmentId int) ([]int, error)
	GetBuildingInfo() *[]component.GeneralObject
	IsFactoryBuildingExist(recordId int) bool
}

type AnalyticsServiceInterface interface {
	GetComponents() []datatypes.JSON
	UpdateComponent(componentName, targetTable string, serialisedObject datatypes.JSON) error
	CreateComponent(serialisedObject datatypes.JSON) (int, error)
}
type NotificationInterface interface {
	CreateMessages(projectId string, messages []Message) error
	CreateSystemNotification(projectId string, notificationMessage datatypes.JSON) error
	CreatePushNotification(projectId string, notificationMessage datatypes.JSON) (int, error)
	GetNotificationHistory(listOfId []int) (error, []component.GeneralObject) // Updated signature
	GetOneSignalSubscriptionId(deviceToken string) []string
}

type ProductionOrderInterface interface {
	GetOrderStatusIdFromPreferenceLevel(projectId string, preferenceLevel int) int
	OrderStatusId2String(projectId string, orderStatusId int) string
	GetCurrentScheduledEvent(projectId string, machineId int) (error, component.GeneralObject)
	GetCurrentAssemblyScheduledEvent(projectId string, machineId int) (error, component.GeneralObject)
	GetCurrentToolingScheduledEvent(projectId string, machineId int) (error, component.GeneralObject)
	GetNextToCurrentScheduledEvent(projectId string, machineId, scheduledOrderId int) (error, component.GeneralObject)
	GetNextToCurrentAssemblyScheduledEvent(projectId string, machineId, scheduledOrderId int) (error, component.GeneralObject)
	GetNextToCurrentToolingScheduledEvent(projectId string, machineId, scheduledOrderId int) (error, component.GeneralObject)
	GetMachineProductionOrderInfo(projectId string, productionOrderId int, machineId int) (error, component.GeneralObject)
	GetCompletedQuantity(projectId string, productionOrderId int) int
	GetAssemblyProductionOrderInfo(projectId string, productionOrderId int, machineId int) (error, component.GeneralObject)
	GetAssemblyProductionOrderById(productionOrderId int) (error, component.GeneralObject)
	GetToolingProductionOrderInfo(projectId string, productionOrderId int, machineId int) (error, component.GeneralObject)
	GetProductionOrderStatus(projectId string, statusId int) (error, component.GeneralObject)
	GetScheduledEvents(projectId string) (error, *[]component.GeneralObject)
	GetToolingScheduledEvents(projectId string) (error, *[]component.GeneralObject)
	GetToolingScheduledEventsByProductionId(projectId string, productionOrderId int) (error, *[]component.GeneralObject)
	GetAssemblyScheduledEvents(projectId string) (error, *[]component.GeneralObject)
	UpdateOrderPreferenceLevel(projectId string, userId int, eventId int, preferenceLevel int) error
	UpdateAssemblyOrderPreferenceLevel(projectId string, userId int, eventId int, preferenceLevel int) error
	UpdateToolingOrderPreferenceLevel(projectId string, userId int, eventId int, preferenceLevel int) error
	GetChildEventsOfProductionOrder(projectId string, productionOrderId int) []int
	GetChildEventsOfAssemblyProductionOrder(projectId string, productionOrderId int) []int
	GetChildEventsOfToolingProductionOrder(projectId string, productionOrderId int) []int
	GetAllSchedulerEventForScheduler(projectId string) (error, *[]component.GeneralObject)
	GetAllToolingEventForScheduler(projectId string) (error, *[]component.GeneralObject)
	GetAllAssemblyEventForScheduler(projectId string) (error, *[]component.GeneralObject)
	UpdateScheduledOrderFields(projectId string, eventId int, updatingData datatypes.JSON) error
	UpdateAssemblyScheduledOrderFields(projectId string, eventId int, updatingData datatypes.JSON) error
	GetAssemblyEventHistorySummary() datatypes.JSON
	UpdateToolingScheduledOrderFields(projectId string, eventId int, updatingData datatypes.JSON) error
	GetRejectQtyByProductionOrder(projectId string, productionOrderId int) int
	GetNextMachineScheduledOrderEvent(projectId string, machineId, eventId int) (error, *[]component.GeneralObject)
	GetNextToolingMachineScheduledOrderEvent(projectId string, machineId, eventId int) (error, *[]component.GeneralObject)
	GetNextAssemblyMachineScheduledOrderEvent(projectId string, machineId, eventId int) (error, *[]component.GeneralObject)
	GetPreviousMachineScheduledOrderEvent(projectId string, machineId, eventId int) (error, *[]component.GeneralObject)
	GetPreviousAssemblyMachineScheduledOrderEvent(projectId string, machineId, eventId int) (error, *[]component.GeneralObject)
	GetPreviousToolingMachineScheduledOrderEvent(projectId string, machineId, eventId int) (error, *[]component.GeneralObject)
	GetFirstScheduledOrderEvent(projectId string, machineId int) (error, *[]component.GeneralObject)
	GetLastScheduledOrderEvent(projectId string, machineId int) (error, *[]component.GeneralObject)
	GetFirstAssemblyScheduledOrderEvent(projectId string, machineId int) (error, *[]component.GeneralObject)
	GetFirstToolingScheduledOrderEvent(projectId string, machineId int) (error, *[]component.GeneralObject)
	GetLastAssemblyScheduledOrderEvent(projectId string, machineId int) (error, *[]component.GeneralObject)
	GetLastToolingScheduledOrderEvent(projectId string, machineId int) (error, *[]component.GeneralObject)
	GetScheduledOrderInfo(projectId string, eventId int) (error, component.GeneralObject)
	GetAssemblyScheduledOrderInfo(projectId string, eventId int) (error, component.GeneralObject)
	GetToolingScheduledOrderInfo(projectId string, eventId int) (error, component.GeneralObject)
	GetSchedulerUpdatedDate(projectId string, mouldId int) (error, string)
	GetComponents() []datatypes.JSON
	UpdateComponent(componentName, targetTable string, serialisedObject datatypes.JSON) error
	CreateComponent(serialisedObject datatypes.JSON) (int, error)
	GetToolingPartById(projectId string, partId int) (error, component.GeneralObject)
	GetBomInfo(projectId string, orderId int) (error, component.GeneralObject)
	GetAssemblyScheduledEventByProductionId(projectId string, productionOrderId int) (error, *[]component.GeneralObject)
	GetScheduledEventByProductionId(projectId string, productionOrderId int) (error, *[]component.GeneralObject)
	UpdateMouldScheduleOrderEventMouldBatchId(projectId string, scheduledOrderEventId int, mouldBatchResourceId int) error
	GetReleasedAssemblyOrderEventsBetween(projectId string, startDateTime string, endDateTime string) (error, *[]component.GeneralObject)
	UpdateAssemblyManualOrderCompletedQuantity(serialisedData datatypes.JSON) (interface{}, error)
	GetAssemblyManualOrderHistoryFromEventId(serialisedData datatypes.JSON) (interface{}, error)
	GetNumberOfOdersByStatus(projectId string) (error, *[]component.GeneralObject)
	GetTotalCompletedQuantity(projectId string, mouldId int, sincetime string) int64
	GetCurrentScheduledEventByMachineId(projectId string, machineId int, testStartDate string, testEndDate string) (error, *[]component.GeneralObject)
	GetCurrentScheduledEventByMouldId(projectId string, mouldId int) (error, *[]component.GeneralObject)
	GetGivenTimeIntervalScheduledEvent(projectId string, machineId int, energyStartDate string) (error, []component.GeneralObject)
}

type MachineInterface interface {
	CreateMachineParam(projectId string, testId int, userId int) (error, int)
	GetListOfMachines(projectId string) (error, []component.GeneralObject)
	GetListOfToolingMachines(projectId string) (error, []component.GeneralObject)
	GetListOfAssemblyMachines(projectId string) (error, []component.GeneralObject)
	GetListOfAssemblyLines(projectId string) (error, []component.GeneralObject)
	GetAssemblyLineFromId(projectId string, resourceId int) (error, component.GeneralObject)
	MoveMachineLiveStatusToActive(projectId string, machineId int) error
	MoveMachineToMaintenance(projectId string, machineId int) error
	MoveMachineLiveStatusToMaintenance(projectId string, machineId int) error
	MoveAssemblyMachineToMaintenance(projectId string, machineId int) error
	MoveToolingMachineToMaintenance(projectId string, machineId int) error
	MoveMachineToActive(projectId string, machineId int) error
	MoveAssemblyMachineToActive(projectId string, machineId int) error
	MoveToolingMachineToActive(projectId string, machineId int) error
	GetListOfMachineSubCategory(projectId string) (error, []component.GeneralObject)
	GetDepartmentDisplayEnabledMachines(projectId string, listOfDepartments []int, allowedDepartment []component.OrderedData) datatypes.JSON
	IsMachineAlreadyUnderMaintenance(projectId string, machineId int) (error, bool)
	GetMachineInfoById(projectId string, mouldId int) (error, component.GeneralObject)
	GetAssemblyMachineInfoById(machineId int) (error, component.GeneralObject)
	GetComponents() []datatypes.JSON
	UpdateComponent(componentName, targetTable string, serialisedObject datatypes.JSON) error
	CreateComponent(serialisedObject datatypes.JSON) (int, error)
	AddAssemblyForceStop(projectId string, machineId int, eventId int) error
	AddToolingForceStop(projectId string, machineId int, eventId int) error
	AddMouldingForceStop(projectId string, machineId int, eventId int) error
	GetListOfAllowedAssemblyLines(projectId string, userId int) []int
	CreateAssemblyHmiEntry(projectId string, machineHmiInfo map[string]interface{}) error
	GetAssemblyMachineDefaultManpower(machineId int) int
	ArchivedAssemblyHmiEntry(projectId string, eventId int) error
	GetActualValueFromAssemblyStatistic(projectId string, eventId int, timeEpoch int64) int
	GetAssemblyMachineInfoFromEquipmentId(projectId string, equipmentId string) (error, int)
	GetListOfMachinesNeededHelp(timestamp int64) (error, []datatypes.JSON) // Added method signature
	GetLastNHoursAssemblyStats(hours int) datatypes.JSON
	UpdateMachineParamEditStatus(projectId string, canEdit bool, paramId int) error
	GetListOfEnergyManagementMachines(projectId string) []int
}

type ProjectInterface interface {
	GetProjectDatasourceInfo() []ProjectDatasourceConfig
}

type BatchManagementInterface interface {
	GenerateMouldBatchLabel(projectId string, resourceId int) error
	GetListOfMouldBatch(mouldBatchId string) []component.GeneralObject
	GetBatchRawMaterial(rawMaterialId int) component.GeneralObject
	GetMouldBatch(recordId int) component.GeneralObject
}
type LabourManagementInterface interface {
	GetMobileAllowedJobRoles() []int
	GetLabourManagementShiftTemplate(templateId []int) []component.GeneralObject
}
type TraceabilityInterface interface {
	CreateTraceabilityResource(mouldBatchId string, qaStartTime, qaCompleteTime string, qaPerson int)
}

type QAInterface interface {
	CreateQAResource(serialisedObject datatypes.JSON) error
}
type MouldInterface interface {
	GetNoOfCavity(projectId string, mouldId int) (error, int)
	GetListOfMoulds(projectId string) (error, []component.GeneralObject)
	GetPartInfo(projectId string, partId int) (error, component.GeneralObject)
	GetMouldTestEventsForScheduler(projectId string) (error, *[]component.GeneralObject)
	GetMouldInfoById(projectId string, mouldId int) (error, component.GeneralObject)
	PutToMaintenanceMode(projectId string, mouldId int, userId int) error
	PutToActiveMode(projectId string, mouldId int, userId int) error
	PutToRepairMode(projectId string, mouldId int, userId int) error
	UpdateShotCount(projectId string, mouldId int, shotCount int) error
	GetComponents() []datatypes.JSON
	UpdateComponent(componentName, targetTable string, serialisedObject datatypes.JSON) error
	CreateComponent(serialisedObject datatypes.JSON) (int, error)
	GetMouldsByPartNo(projectId string, partNo string) (error, []component.GeneralObject)
	GetMouldMachineTestParam(machineId int, mouldId int) int
	GetBrandName(projectId string, brandId int) (error, string)
	GetListOfMouldCategory(projectId string) (error, []component.GeneralObject)
	MoveMouldToRepair(projectId string, mouldId int) error
	GetMouldShotCountViewForNotification() (error, *[]component.GeneralObject)
	GetMouldSettingById(projectId string) (error, component.GeneralObject)
	GetMouldNameByMachineId(machineId int) (error, string)
}

type MaintenanceInterface interface {
	GetMaintenanceEventsForScheduler(projectId string) (error, *[]component.GeneralObject)
	GetMaintenanceEventsForAssemblyScheduler(projectId string) (error, *[]component.GeneralObject)
	GetMaintenanceEventsForToolingScheduler(projectId string) (error, *[]component.GeneralObject)
	GetMaintenanceOrderForScheduler(projectId string) (error, *[]component.GeneralObject)
	GetMaintenanceCorrectiveOrderForScheduler(projectId string) (error, *[]component.GeneralObject)
	GetComponents() []datatypes.JSON
	UpdateComponent(componentName, targetTable string, serialisedObject datatypes.JSON) error
	CreateComponent(serialisedObject datatypes.JSON) (int, error)

	GetMouldPreventiveMaintenanceOrderForScheduler(projectId string) (error, *[]component.GeneralObject)
	GetMouldCorrectiveMaintenanceOrderForScheduler(projectId string) (error, *[]component.GeneralObject)
	GetMaintenanceWorkOrderByMachineId(machineId int, table string) (error, string)
}

var globalService = make(map[string]*component.ServiceConfig, 100)

func RegisterService(service *component.ServiceConfig) {
	globalService[service.ServiceName] = service
}

func GetService(serviceType string) *component.ServiceConfig {
	return globalService[serviceType]
}

type BaseService struct {
	ReferenceDatabase *gorm.DB
	Logger            *zap.Logger
	ServiceDatabases  map[string]*gorm.DB
}

func (bs *BaseService) GetDatabase(projectId string) *gorm.DB {
	if database, ok := bs.ServiceDatabases[projectId]; ok {
		return database
	} else {
		return nil
	}
}
func New() *BaseService {
	return &BaseService{}
}

func (bs *BaseService) CrateProjectDatabaseConnection(projectId string, name string) error {
	err := bs.ReferenceDatabase.Exec("CREATE DATABASE " + name).Error
	if err != nil {
		return err
	}
	// now connect db and add to map
	//dbConnection, err := bs.databaseConfig.GetDbConnectionFromDbName(name)
	//if err != nil {
	//	bs.ReferenceDatabase.Exec("DROP DATABASE " + name)
	//	return err
	//}
	//
	//bs.ServiceDatabases[projectId] = dbConnection
	return nil
}

func NewSugaredLogger(logConfig *zap.Config) (*zap.SugaredLogger, error) {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.OutputPaths = logConfig.OutputPaths
	zapLogger, err := cfg.Build()
	if err != nil {
		return nil, err
	} else {
		return zapLogger.Sugar(), nil
	}
}
func NewLogger(logConfig *zap.Config) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.OutputPaths = logConfig.OutputPaths
	zapLogger, err := cfg.Build()
	if err != nil {
		return nil, err
	} else {
		return zapLogger, nil
	}
}
func (bs *BaseService) Init(configFile string, listOfProjectConfig []ProjectDatasourceConfig) error {
	// load the config from a file source
	if err := config.Load(file.NewSource(
		file.WithPath(configFile),
	)); err != nil {
		fmt.Errorf("failed to load the service configuration [%s]", err.Error())
		return err
	}

	var databaseConfig orm.DatabaseConfig
	var logConfig zap.Config
	// read a database host
	if err := config.Get("hosts", "database").Scan(&databaseConfig); err != nil {
		fmt.Errorf("failed to get the database configuration [%s] ", err.Error())
		return err
	}
	// read a database host
	if err := config.Get("log").Scan(&logConfig); err != nil {
		fmt.Errorf("failed to get the log configuration [%s]", err.Error())
		return err
	}

	// read the service config
	var serviceConfig component.ServiceConfig
	if err := config.Get("service").Scan(&serviceConfig); err != nil {
		fmt.Errorf("failed to get the service configuration [%s] ", err.Error())
		return err
	}

	logger, _ := NewLogger(&logConfig)
	database, err := databaseConfig.NewConnection()
	if err != nil {
		logger.Error("Connecting database has failed", zap.String("error", err.Error()))
	}
	bs.ServiceDatabases = make(map[string]*gorm.DB, len(listOfProjectConfig))
	for _, projectConfig := range listOfProjectConfig {
		dbConnection, err := projectConfig.DatasourceConfig.NewConnection()
		if err != nil {
			logger.Error("Reference database has failed", zap.String("error", err.Error()))
			return err
		}
		bs.ServiceDatabases[projectConfig.ProjectId] = dbConnection
	}

	fmt.Println("database connection :", database)
	bs.ReferenceDatabase = database
	bs.Logger = logger

	return nil
}
