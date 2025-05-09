package consts

const (
	MachineDownTimeMasterTable             = "machine_downtime_master"
	MachineDownTimeSettingTable            = "machine_downtime_setting"
	MachineDownTimeFaultTypeTable          = "machine_downtime_fault_type"
	MachineDownTimeFaultCodeTable          = "machine_downtime_fault_code"
	MachineDownTimeStatusTable             = "machine_downtime_status"
	MachineDownTimeEmailEscalationTable    = "machine_downtime_email_escalation"
	MachineDownTimeSignalHistoryTable      = "machine_downtime_signal_history"
	MachineDownTimeEmailTemplateTable      = "labour_management_email_template"
	MachineDownTimeEmailTemplateFieldTable = "labour_management_email_template_field"

	DowntimeStatus_Fault_Reportd             = 1
	DowntimeStatus_Fault_Under_Investigation = 2
	DowntimeStatus_Fault_Repaired            = 3
	DowntimeStatus_Fault_Cacelled            = 4

	ProjectID = "906d0fd569404c59956503985b330132"
)
