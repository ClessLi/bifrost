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

	err = registerJsonRegMatchers()
	if err != nil {
		panic(err)
	}
	err = registerJsonUnmarshalerBuilders()
	if err != nil {
		panic(err)
	}
}
