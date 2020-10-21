package nginx

func Save(conf *Config) (Caches, error) {
	caches, err := conf.Save()
	if err != nil {
		return nil, err
	}
	return caches, nil
}
