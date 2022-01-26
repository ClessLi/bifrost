package options

func (o *Options) Validate() []error {
	errors := make([]error, 0)

	errors = append(errors, o.GenericServerRunOptions.Validate()...)
	errors = append(errors, o.SecureServing.Validate()...)
	errors = append(errors, o.InsecureServing.Validate()...)
	errors = append(errors, o.RAOptions.Validate()...)
	errors = append(errors, o.WebServerConfigsOptions.Validate()...)
	errors = append(errors, o.MonitorOptions.Validate()...)
	errors = append(errors, o.WebServerLogWatcherOptions.Validate()...)
	errors = append(errors, o.Log.Validate()...)

	return errors
}
