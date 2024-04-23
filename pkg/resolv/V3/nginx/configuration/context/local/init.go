package local

func init() {
	err := registerContextBuilders()
	if err != nil {
		panic(err)
	}
	err = registerContextParseFuncs()
	if err != nil {
		panic(err)
	}
}
