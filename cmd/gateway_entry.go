package main

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/services/analytics"
	"cx-micro-flake/services/auth"
	"cx-micro-flake/services/batch_management"
	"cx-micro-flake/services/content"
	energy_management "cx-micro-flake/services/energy_management/source"
	"cx-micro-flake/services/facility"
	"cx-micro-flake/services/factory"
	"cx-micro-flake/services/fixed_assets"
	"cx-micro-flake/services/incident"
	"cx-micro-flake/services/iot/cmd"
	"cx-micro-flake/services/it"
	"cx-micro-flake/services/labour_management"
	machine_downtime "cx-micro-flake/services/machine_downtime/source"
	"cx-micro-flake/services/machines"
	"cx-micro-flake/services/maintenance"
	"cx-micro-flake/services/manufacturing"
	"cx-micro-flake/services/moulds"
	"cx-micro-flake/services/notifications"
	"cx-micro-flake/services/production_order"
	"cx-micro-flake/services/project"
	"cx-micro-flake/services/qa"
	"cx-micro-flake/services/scheduler"
	spare_part_inventory "cx-micro-flake/services/spare_part_inventory/source"
	"cx-micro-flake/services/tickets"
	"cx-micro-flake/services/tooling"
	"cx-micro-flake/services/traceability"
	"cx-micro-flake/services/work_hour_assignment"
	"fmt"
	"go.cerex.io/transcendflow/base_modules"
	"go.cerex.io/transcendflow/config"
	"go.cerex.io/transcendflow/logging"
	"go.cerex.io/transcendflow/service"
	"go.cerex.io/transcendflow/service_zmq"
	"go.cerex.io/transcendflow/util/color"
	"net/http"
	"os"
	"strconv"
)

const (
	AuthService               = "auth"
	ContentService            = "content"
	MouldsService             = "moulds"
	MaintenanceService        = "maintenance"
	ProjectService            = "project"
	Machines                  = "machines"
	Notifications             = "notifications"
	ProductionOrderService    = "production_order"
	IoTService                = "iot"
	AnalyticsService          = "analytics"
	SchedulerService          = "scheduler"
	ITService                 = "it"
	FixedAssetService         = "fixed_assets"
	FactoryService            = "factory"
	IncidentService           = "incident"
	FacilityService           = "facility"
	ToolingService            = "tooling"
	WorkHourAssignmentService = "work_hour_assignment"
	TicketsService            = "tickets"
	ManufacturingService      = "manufacturing"
	BatchManagementService    = "batch_management"
	QAService                 = "qa"
	TraceabilityService       = "traceability"
	LabourManagementService   = "labour_management"
	MachineDownTimeService    = "machine_downtime"
	SparePartInventoryService = "spare_part_inventory"
	EnergyManagementService   = "energy_management"
)

type AppConfig struct {
	Server         ServerConfig `mapstructure:"server"`
	LoadComponent  bool         `mapstructure:"loadComponent"`
	DeploymentMode string       `mapstructure:"deploymentMode"`
	BaseServices   []string     `mapstructure:"baseServices"`
	CustomServices []string     `mapstructure:"customServices"`
	ServicePort    int          `mapstructure:"servicePort"`
}
type ServerConfig struct {
	Address string `mapstructure:"address"`
	Port    int    `mapstructure:"port"`
}

func (v *AppConfig) getAddress() string {
	return v.Server.Address + ":" + strconv.Itoa(v.Server.Port)
}

var version = "1.0.0"

func main() {

	service.InitBaseFlags(version)
	logger, err := logging.GetLoggerFromCustomPath("../conf")
	if err != nil {
		color.Red("creating logger has failed [%s]", err.Error())
		os.Exit(0)
	}
	logger.Info("starting up...")
	applicationConfig := config.GetApplicationConfig("../conf")
	conf := AppConfig{}
	if err := applicationConfig.Unmarshal(&conf); err != nil {
		logger.Error("no config.yml found in the gateway, check the config.yml is located under config directory, terminating all the services now..")
		os.Exit(0)
	}
	fmt.Print("conf : ", conf)
	var router = service_zmq.GetRouter(conf.DeploymentMode)
	masterService := service_zmq.MasterService{}
	masterService.Boostrap(conf.ServicePort, logger)

	// let's start the base services
	baseServiceManager := base_modules.BaseServiceManager{Logger: logger, Router: router}
	fmt.Print("conf.BaseServices : ", conf.BaseServices)
	baseServiceManager.StartBaseModulesWithCustomConfig(conf.BaseServices, "../base_services")

	// start the base service auth and project
	var projectConfig []common.ProjectDatasourceConfig
	if err := project.InitStandaloneService(router, projectConfig); err != nil {
		logger.Info("Failed to init project service", logging.Error(err))
		os.Exit(0)
	}

	projectService := common.GetService("project").ServiceInterface.(common.ProjectInterface)
	projectConfig = projectService.GetProjectDatasourceInfo()

	if err := auth.InitStandaloneService(router, projectConfig); err != nil {
		logger.Info("Failed to init auth service", logging.Error(err))
		os.Exit(0)
	}

	for _, serviceName := range conf.CustomServices {

		if serviceName == ContentService {
			if err := content.InitStandaloneService(router, projectConfig); err != nil {
				logger.Info("Failed to init content service", logging.Error(err))
				os.Exit(0)
			}
		}

		if serviceName == Machines {
			if err := machines.InitStandaloneService(router, projectConfig); err != nil {
				logger.Info("Failed to init machine service", logging.Error(err))
				os.Exit(0)
			}
		}
		if serviceName == MouldsService {
			if err := moulds.InitStandaloneService(router, projectConfig); err != nil {
				logger.Info("Failed to init moulds service", logging.Error(err))
				os.Exit(0)
			}
		}
		if serviceName == MaintenanceService {
			if err := maintenance.InitStandaloneService(router, projectConfig); err != nil {
				logger.Info("Failed to init maintenance service", logging.Error(err))
				os.Exit(0)
			}
		}
		if serviceName == Notifications {
			if err := notifications.InitStandaloneService(router, projectConfig); err != nil {
				logger.Info("Failed to init notification service", logging.Error(err))
				os.Exit(0)
			}
		}
		if serviceName == ProductionOrderService {
			if err := production_order.InitStandaloneService(router, projectConfig); err != nil {
				logger.Info("Failed to init production order service", logging.Error(err))
				os.Exit(0)
			}
		}

		if serviceName == IoTService {
			if err := cmd.Bootstrap(router); err != nil {
				logger.Info("Failed to init IoT service", logging.Error(err))
				os.Exit(0)
			}
		}
		if serviceName == AnalyticsService {
			if err := analytics.InitStandaloneService(router, projectConfig); err != nil {
				logger.Info("Failed to init Analytics service", logging.Error(err))
				os.Exit(0)
			}
		}
		if serviceName == SchedulerService {
			if err := scheduler.InitStandaloneService(router, projectConfig); err != nil {
				logger.Info("Failed to init Scheduler service", logging.Error(err))
				os.Exit(0)
			}
		}
		if serviceName == ITService {
			if err := it.InitStandaloneService(router, projectConfig); err != nil {
				logger.Info("Failed to init IT service", logging.Error(err))
				os.Exit(0)
			}
		}
		if serviceName == FixedAssetService {
			if err := fixed_assets.InitStandaloneService(router, projectConfig); err != nil {
				logger.Info("Failed to init fixed assets service", logging.Error(err))
				os.Exit(0)
			}
		}
		if serviceName == FactoryService {
			if err := factory.InitStandaloneService(router, projectConfig); err != nil {
				logger.Info("Failed to init factory assets service", logging.Error(err))
				os.Exit(0)
			}
		}
		if serviceName == IncidentService {
			if err := incident.InitStandaloneService(router, projectConfig); err != nil {
				logger.Info("Failed to start incident service service", logging.Error(err))
				os.Exit(0)
			}
		}

		if serviceName == FacilityService {
			if err := facility.InitStandaloneService(router, projectConfig); err != nil {
				logger.Info("Failed to start facility service", logging.Error(err))
				os.Exit(0)
			}
		}

		if serviceName == ToolingService {
			if err := tooling.InitStandaloneService(router, projectConfig); err != nil {
				logger.Info("Failed to start tooling service", logging.Error(err))
				os.Exit(0)
			}
		}
		if serviceName == WorkHourAssignmentService {
			if err := work_hour_assignment.InitStandaloneService(router, projectConfig); err != nil {
				logger.Info("Failed to start Work Hour Assignment service", logging.Error(err))
				os.Exit(0)
			}
		}
		if serviceName == TicketsService {
			if err := tickets.InitStandaloneService(router, projectConfig); err != nil {
				logger.Info("Failed to start Tickets service", logging.Error(err))
				os.Exit(0)
			}
		}
		if serviceName == ManufacturingService {
			if err := manufacturing.InitStandaloneService(router, projectConfig); err != nil {
				logger.Info("Failed to start Manufacturing service", logging.Error(err))
				os.Exit(0)
			}
		}
		if serviceName == BatchManagementService {
			if err := batch_management.InitStandaloneService(router, projectConfig); err != nil {
				logger.Info("Failed to start Batch Management service", logging.Error(err))
				os.Exit(0)
			}
		}
		if serviceName == TraceabilityService {
			if err := traceability.InitStandaloneService(router, projectConfig); err != nil {
				logger.Info("Failed to start traceability service", logging.Error(err))
				os.Exit(0)
			}
		}
		if serviceName == QAService {
			if err := qa.InitStandaloneService(router, projectConfig); err != nil {
				logger.Info("Failed to start QA service", logging.Error(err))
				os.Exit(0)
			}
		}
		if serviceName == LabourManagementService {
			if err := labour_management.InitStandaloneService(router, projectConfig); err != nil {
				logger.Info("Failed to start Labour management service", logging.Error(err))
				os.Exit(0)
			}
		}
		if serviceName == MachineDownTimeService {
			if err := machine_downtime.Bootstrap(router, "../services/machine_downtime/conf"); err != nil {
				logger.Info("Failed to start Machine Downtime  service", logging.Error(err))
				os.Exit(0)
			}
		}
		if serviceName == SparePartInventoryService {
			if err := spare_part_inventory.Bootstrap(router, "../services/spare_part_inventory/conf"); err != nil {
				logger.Info("Failed to start Spare Part Inventory service", logging.Error(err))
				os.Exit(0)
			}
		}
		if serviceName == EnergyManagementService {
			if err := energy_management.Bootstrap(router, "../services/energy_management/conf"); err != nil {
				logger.Info("Failed to start energy management service", logging.Error(err))
				os.Exit(0)
			}
		}

	}
	masterService.LoadComponents()
	logger.Info("Gateway service is finally started after starting all the configured services ", logging.String("service_address", conf.getAddress()))
	if err := http.ListenAndServe(conf.getAddress(), router); err != nil {
		logger.Error("error starting HTTP server", logging.Error(err))
	}
}
