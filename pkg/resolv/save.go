package resolv

func Save(conf *Config) error {
	err := conf.save()
	if err != nil {
		return err
	}
	return nil
}
