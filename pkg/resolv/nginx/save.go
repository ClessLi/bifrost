package nginx

func Save(conf *Config) error {
	err := conf.Save()
	if err != nil {
		return err
	}
	return nil
}
