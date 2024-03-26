package options

import (
	"fmt"
	"github.com/ClessLi/bifrost/internal/pkg/logger"
	"github.com/spf13/pflag"
	"go.uber.org/zap/zapcore"
	"path/filepath"
	"strings"
	"time"
)

const (
	flagInfoLevel             = "log.info-level"
	flagErrorLevel            = "log.error-level"
	flagDisableCaller         = "log.disable-caller"
	flagDisableStacktrace     = "log.disable-stacktrace"
	flagFormat                = "log.format"
	flagEnableColor           = "log.enable-color"
	flagInfoOutputPaths       = "log.info-output-paths"
	flagErrorOutputPaths      = "log.error-output-paths"
	flagInnerErrorOutputPaths = "log.inner-error-output-paths"
	flagDevelopment           = "log.development"
	flagName                  = "log.name"

	consoleFormat = "console"
	jsonFormat    = "json"

	defaultLogBaseDir       = "logs"
	defaultInfoLogFileName  = "bifrost.log"
	defaultErrorLogFileName = "bifrost_error.log"
)

type LoggerOptions struct {
	Name                  string   `json:"name"                     mapstructure:"name"`
	InfoOutputPaths       []string `json:"info-output-paths"        mapstructure:"info-output-paths"`
	ErrorOutputPaths      []string `json:"error-output-paths"       mapstructure:"error-output-paths"`
	InnerErrorOutputPaths []string `json:"inner-error-output-paths" mapstructure:"inner-error-output-paths"`
	InfoLevel             string   `json:"info-level"               mapstructure:"info-level"`
	ErrorLevel            string   `json:"error-level"              mapstucture:"error-level"`
	Format                string   `json:"format"                   mapstructure:"format"`
	DisableCaller         bool     `json:"disable-caller"           mapstructure:"disable-caller"`
	DisableStacktrace     bool     `json:"disable-stacktrace"       mapstructure:"disable-stacktrace"`
	EnableColor           bool     `json:"enable-color"             mapstructure:"enable-color"`
	Development           bool     `json:"development"              mapstructure:"development"`
}

func NewLoggerOptions() *LoggerOptions {
	opts := &LoggerOptions{
		InfoLevel:             zapcore.InfoLevel.String(),
		ErrorLevel:            zapcore.WarnLevel.String(),
		DisableCaller:         false,
		DisableStacktrace:     false,
		Format:                consoleFormat,
		EnableColor:           false,
		Development:           false,
		InfoOutputPaths:       []string{filepath.Join(defaultLogBaseDir, defaultInfoLogFileName)},
		ErrorOutputPaths:      []string{filepath.Join(defaultLogBaseDir, defaultErrorLogFileName)},
		InnerErrorOutputPaths: []string{"stderr"},
	}

	return opts
}

func (o *LoggerOptions) Validate() []error {
	var errs []error
	var zapLevel zapcore.Level
	// validate Info log
	if err := zapLevel.UnmarshalText([]byte(o.InfoLevel)); err != nil {
		errs = append(errs, err)
	}

	format := strings.ToLower(o.Format)
	if format != consoleFormat && format != jsonFormat {
		errs = append(errs, fmt.Errorf("not a valid log format: %q", o.Format))
	}

	// validate Error log
	if err := zapLevel.UnmarshalText([]byte(o.ErrorLevel)); err != nil {
		errs = append(errs, err)
	}

	return errs
}

func (o *LoggerOptions) AddFlags(fs *pflag.FlagSet) {
	// Add Info log flags
	fs.StringVar(&o.InfoLevel, flagInfoLevel, o.InfoLevel, "Minimum Info log output `LEVEL`.")
	fs.BoolVar(&o.DisableCaller, flagDisableCaller, o.DisableCaller, "Disable output of caller information in the log.")
	fs.BoolVar(&o.DisableStacktrace, flagDisableStacktrace,
		o.DisableStacktrace, "Disable the log to record a stack trace for all messages at or above panic level.")
	fs.StringVar(&o.Format, flagFormat, o.Format, "Log output `FORMAT`, support plain or json format.")
	fs.BoolVar(&o.EnableColor, flagEnableColor, o.EnableColor, "Enable output ansi colors in plain format logs.")
	fs.StringSliceVar(&o.InfoOutputPaths, flagInfoOutputPaths, o.InfoOutputPaths, "Output paths of Info log.")
	fs.StringSliceVar(&o.InnerErrorOutputPaths, flagInnerErrorOutputPaths, o.InnerErrorOutputPaths, "Inner Error output paths of log.")
	fs.BoolVar(
		&o.Development,
		flagDevelopment,
		o.Development,
		"Development puts the logger in development mode, which changes "+
			"the behavior of DPanicLevel and takes stacktraces more liberally.",
	)
	fs.StringVar(&o.Name, flagName, o.Name, "The name of the logger.")

	// Add Error log flags
	fs.StringVar(&o.ErrorLevel, flagErrorLevel, o.ErrorLevel, "Minimum Error log output `LEVEL`.")

	fs.StringSliceVar(&o.ErrorOutputPaths, flagErrorOutputPaths, o.ErrorOutputPaths, "Output paths of Error log.")
}

func (o *LoggerOptions) Complete() error {
	defaultLogSecondLevelDir := time.Now().Format("20060102")
	if o.InfoOutputPaths == nil || (len(o.InfoOutputPaths) == 1 && o.InfoOutputPaths[0] == "stdout") {
		o.InfoOutputPaths = append(o.InfoOutputPaths, filepath.Join(defaultLogBaseDir, defaultLogSecondLevelDir, defaultInfoLogFileName))
	}
	if o.ErrorOutputPaths == nil || (len(o.ErrorOutputPaths) == 1 && o.ErrorOutputPaths[0] == "stdout") {
		o.ErrorOutputPaths = []string{filepath.Join(defaultLogBaseDir, defaultLogSecondLevelDir, defaultErrorLogFileName)}
	}

	return nil
}

func (o *LoggerOptions) ApplyTo(conf *logger.Config) error {
	conf.InfoLogOpts.Format = o.Format
	conf.InfoLogOpts.Name = o.Name
	conf.InfoLogOpts.ErrorOutputPaths = o.InnerErrorOutputPaths
	conf.InfoLogOpts.Level = o.InfoLevel
	conf.InfoLogOpts.EnableColor = o.EnableColor
	conf.InfoLogOpts.Development = o.Development
	conf.InfoLogOpts.DisableCaller = o.DisableCaller
	conf.InfoLogOpts.DisableStacktrace = o.DisableStacktrace
	conf.InfoLogOpts.OutputPaths = o.InfoOutputPaths

	conf.ErrLogOpts.Format = o.Format
	conf.ErrLogOpts.ErrorOutputPaths = nil
	conf.ErrLogOpts.Level = o.ErrorLevel
	conf.ErrLogOpts.EnableColor = o.EnableColor
	conf.ErrLogOpts.Development = o.Development
	conf.ErrLogOpts.DisableCaller = o.DisableCaller
	conf.ErrLogOpts.DisableStacktrace = o.DisableStacktrace
	conf.ErrLogOpts.OutputPaths = o.ErrorOutputPaths

	return nil
}
