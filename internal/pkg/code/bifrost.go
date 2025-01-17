package code

//go:generate codegen -type=int -fullname=Bifrost

// bifrost: configuration errors.
const (
	// ErrConfigurationTypeMismatch - 500: Configuration type mismatch.
	ErrConfigurationTypeMismatch int = iota + 110001

	// ErrSameConfigFingerprint - 500: Same config fingerprint.
	ErrSameConfigFingerprint

	// ErrSameConfigFingerprints - 500: Same config fingerprint between files and configuration.
	ErrSameConfigFingerprints

	// ErrConfigManagerIsRunning - 500: Config manager is running.
	ErrConfigManagerIsRunning

	// ErrConfigManagerIsNotRunning - 500: Config manager is not running.
	ErrConfigManagerIsNotRunning

	// ErrWebServerNotFound - 400: Web server not found.
	ErrWebServerNotFound

	// ErrConfigurationNotFound - 400: Web server configuration not found.
	ErrConfigurationNotFound

	// ErrParserNotFound - 500: Parser not found.
	ErrParserNotFound

	// ErrUnknownKeywordString - 500: Unknown keyword string.
	ErrUnknownKeywordString

	// ErrInvalidConfig - 500: Invalid parser.Config.
	ErrInvalidConfig

	// ErrParseFailed - 500: Config parse failed.
	ErrParseFailed

	// ErrV3ContextIndexOutOfRange - 500: Index of the Context's children is out of range.
	ErrV3ContextIndexOutOfRange

	// ErrV3NullContextPosition - 500: Null Context position.
	ErrV3NullContextPosition

	// ErrV3SetFatherContextFailed - 500: Set father Context failed.
	ErrV3SetFatherContextFailed

	// ErrV3OperationOnErrorContext - 500: Performing operations on Error Context.
	ErrV3OperationOnErrorContext

	// ErrV3InvalidContext - 500: Invalid Context.
	ErrV3InvalidContext

	// ErrV3InvalidOperation - 500: Invalid operation.
	ErrV3InvalidOperation

	// ErrV3ContextNotFound - 500: Queried context not found.
	ErrV3ContextNotFound

	// ErrV3ConversionToContextFailed - 500: Conversion to context failed.
	ErrV3ConversionToContextFailed
)

// bifrost: statistics errors.
const ()

// bifrost: metrics errors.
const (
	// ErrStopMonitoringTimeout - 500: Stop monitoring timeout.
	ErrStopMonitoringTimeout int = iota + 110201

	// ErrMonitoringServiceSuspension - 500: Monitoring service suspension.
	ErrMonitoringServiceSuspension

	// ErrMonitoringStarted - 500: Monitoring is already started.
	ErrMonitoringStarted
)

// bifrost: log watch errors.
const (
	// ErrLogsDirPath - 500: Logs dir is not exist or is not a directory.
	ErrLogsDirPath int = iota + 110301

	// ErrLogBufferIsNotExist - 500: Log buffer is not exist.
	ErrLogBufferIsNotExist

	// ErrLogBufferIsExist - 500: Log buffer is already exist.
	ErrLogBufferIsExist

	// ErrLogIsLocked - 500: Log is locked.
	ErrLogIsLocked

	// ErrLogIsUnlocked - 500: Log is unlocked.
	ErrLogIsUnlocked

	// ErrUnknownLockError - 500: Unknown lock error.
	ErrUnknownLockError
)
