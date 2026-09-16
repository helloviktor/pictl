package ipc

// Command identifies an operation that can be requested from the pictl helper process.
type Command string

const (
	CommandCPUUsage         Command = "cpu_usage"
	CommandMemoryUsage      Command = "memory_usage"
	CommandDiskUsage        Command = "disk_usage"
	CommandCPUTemperature   Command = "cpu_temperature"
	CommandAvailableUpdates Command = "available_updates"
	CommandApplyUpdates     Command = "apply_updates"
	CommandRestart          Command = "restart"
	CommandShutdown         Command = "shutdown"
)
