package config

// Flag: 				DevMdoe (bool)
// default: 		false
// description: Run in development mode
var DevModeFlag bool

// Flag: 				EnvMode (bool)
// default: 		false
// description: Print env before execution
var EnvModeFlag bool

// Flag: 				Confirm (bool)
// default: 		false
// description: Confirm action before execution
var ConfirmFlag bool

// Flag: 				Host (string)
// default: 		"""
// description: Set the host for the server. Overrides env var HOST
var HostFlag string

// Flag: 				ProtectedPortFlag (string)
// default: 		""
// description: Set the protected api port for the server. Overrides env var PORT
var ProtectedPortFlag string

// Flag: 				PublicPortFlag (string)
// default: 		""
// description: Set the public api port for the server. Overrides env var PORT
var PublicPortFlag string

// Flag: 				LogLevel (string)
// default: 		""
// description: Set the log level. Overrides env var LOG_LEVEL
var LogLevelFlag string

// Flag: 				NoBorder (bool)
// default: 		false
// description: Print tables without border styling
var NoBorderFlag bool

// Flag: 				StartWatcher (bool)
// default: 		false
// description: Start watcher daemon
var StartWatcherFlag bool

// Flag: 				pushEndpoint (string)
// default:
// description:
var PushEndpoint string

// Flag: 				samplingFraction (string)
// default:
// description:
var SamplingFraction string

// Flag: 				exporterType (string)
// default:
// description:
var ExporterType string
