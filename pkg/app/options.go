package app

import (
	cliflag "github.com/marmotedu/component-base/pkg/cli/flag"
)

// CliOptions abstracts configuration options for reading parameters from the
// command line.
type CliOptions interface {
	Flags() (fss cliflag.NamedFlagSets)
	Validate() []error
}

// ConfigurableOptions abstracts configuration options for reading parameters
// from a configuration file.
type ConfigurableOptions interface {
	ApplyFlags() []error
}

// CompletableOptions abstracts configuration options which can completed.
type CompletableOptions interface {
	Complete() error
}

// PrintableOptions abstracts configuration options which can printed.
type PrintableOptions interface {
	String() string
}
