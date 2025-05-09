package const_util

const (
	AssemblyMachineHelpSignalViewTable   = "assembly_machine_help_signal_view"
	AssemblyHelpSignalProcessedTimeTable = "assembly_help_signal_processed_time"

	SelectAssemblyHistoryMessage        = "SELECT body ,ts  FROM message WHERE ts > ? AND topic = 'machines/L3_AssemblyLine' ORDER BY ts ASC"
	SelectAssemblyMasterFromMessageFlag = "SELECT id FROM assembly_machine_master WHERE object_info->>'$.lineMappingMessageFlag' = ? AND object_info->>'$.helpButtonStationNo' = ?"
)
